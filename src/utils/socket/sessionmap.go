// 通过 read 和 dirty 两个字段将读写分离，读的数据存在只读字段 read 上，将最新写入的数据则存在 dirty 字段上
// 读取时会先查询 read，不存在再查询 dirty，写入时则只写入 dirty
// 读取 read 并不需要加锁，而读或写 dirty 都需要加锁
// 另外有 misses 字段来统计 read 被穿透的次数（被穿透指需要读 dirty 的情况），超过一定次数则将 dirty 数据同步到 read 上
// 对于删除数据则直接通过标记来延迟删除
package socket

import (
	"gameserver/utils/socket/server"
	"sync"
	"sync/atomic"
	"unsafe"
)

type SessionMap struct {
	mu     sync.Mutex                // 加锁作用，保护 dirty 字段
	read   atomic.Value              // 只读的数据，实际数据类型为 readOnly
	dirty  map[uint]*entrySessionMap // 最新写入的数据
	misses int                       // 计数器，每次需要读 dirty 则 +1统计 read 被穿透的次数（被穿透指需要读 dirty 的情况），超过一定次数则将 dirty 数据同步到 read 上
}

type readOnlySessionMap struct {
	m       map[uint]*entrySessionMap
	amended bool // 表示 dirty 里存在 read 里没有的 key，通过该字段决定是否加锁读 dirty
}

var expungedSessionMap = unsafe.Pointer(new(server.Session))

type entrySessionMap struct {
	//p == nil: 键值已经被删除，且 m.dirty == nil
	//p == expunged: 键值已经被删除，但 m.dirty!=nil 且 m.dirty 不存在该键值（expunged 实际是空接口指针）
	//除以上情况，则键值对存在，存在于 m.read.m 中，如果 m.dirty!=nil 则也存在于 m.dirty
	p unsafe.Pointer // 等同于 *interface{}
}

func newEntrySessionMap(i server.Session) *entrySessionMap {
	return &entrySessionMap{p: unsafe.Pointer(&i)}
}

func (m *SessionMap) Load(key uint) (value server.Session, ok bool) {
	read, _ := m.read.Load().(readOnlySessionMap) // 首先尝试从 read 中读取 readOnly 对象
	e, ok := read.m[key]
	if !ok && read.amended { // 如果不存在则尝试从 dirty 中获取
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlySessionMap) // 由于上面 read 获取没有加锁，为了安全再检查一次
		e, ok = read.m[key]
		if !ok && read.amended { // 确实不存在则从 dirty 获取
			e, ok = m.dirty[key]
			m.missLocked() // 调用 miss 的逻辑
		}
		m.mu.Unlock()
	}
	if !ok {
		return value, false
	}
	return e.load()
}

func (e *entrySessionMap) load() (value server.Session, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedSessionMap {
		return value, false
	}
	return *(*server.Session)(p), true
}

func (m *SessionMap) Store(key uint, value server.Session) {
	read, _ := m.read.Load().(readOnlySessionMap)
	if e, ok := read.m[key]; ok && e.tryStore(&value) { // 如果 read 里存在，则尝试存到 entry 里
		return
	}
	// 如果上一步没执行成功，则要分情况处理
	m.mu.Lock()
	read, _ = m.read.Load().(readOnlySessionMap)
	if e, ok := read.m[key]; ok { // 和 Load 一样，重新从 read 获取一次
		if e.unexpungeLocked() { // 情况 1：read 里存在
			m.dirty[key] = e // 如果 p == expunged，则需要先将 entry 赋值给 dirty（因为 expunged 数据不会留在 dirty）
		}
		e.storeLocked(&value) // 用值更新 entry
	} else if e, ok := m.dirty[key]; ok {
		e.storeLocked(&value) // 情况 2：read 里不存在，但 dirty 里存在，则用值更新 entry
	} else {
		if !read.amended { // 情况 3：read 和 dirty 里都不存在
			m.dirtyLocked()                                            // 如果 amended == false，则调用 dirtyLocked 将 read 拷贝到 dirty（除了被标记删除的数据）
			m.read.Store(readOnlySessionMap{m: read.m, amended: true}) // 然后将 amended 改为 true
		}
		m.dirty[key] = newEntrySessionMap(value) // 然后将 amended 改为 true
	}
	m.mu.Unlock()
}

func (e *entrySessionMap) tryStore(i *server.Session) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expungedSessionMap {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entrySessionMap) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expungedSessionMap, nil)
}
func (e *entrySessionMap) storeLocked(i *server.Session) {
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *SessionMap) LoadOrStore(key uint, value server.Session) (actual server.Session, loaded bool) {
	read, _ := m.read.Load().(readOnlySessionMap)
	if e, ok := read.m[key]; ok {
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, loaded
		}
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlySessionMap)
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		actual, loaded, _ = e.tryLoadOrStore(value)
	} else if e, ok := m.dirty[key]; ok {
		actual, loaded, _ = e.tryLoadOrStore(value)
		m.missLocked()
	} else {
		if !read.amended {
			m.dirtyLocked()
			m.read.Store(readOnlySessionMap{m: read.m, amended: true})
		}
		m.dirty[key] = newEntrySessionMap(value)
		actual, loaded = value, false
	}
	m.mu.Unlock()

	return actual, loaded
}

func (e *entrySessionMap) tryLoadOrStore(i server.Session) (actual server.Session, loaded, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == expungedSessionMap {
		return actual, false, false
	}
	if p != nil {
		return *(*server.Session)(p), true, true
	}

	ic := i
	for {
		if atomic.CompareAndSwapPointer(&e.p, nil, unsafe.Pointer(&ic)) {
			return i, false, true
		}
		p = atomic.LoadPointer(&e.p)
		if p == expungedSessionMap {
			return actual, false, false
		}
		if p != nil {
			return *(*server.Session)(p), true, true
		}
	}
}

func (m *SessionMap) Delete(key uint) {
	read, _ := m.read.Load().(readOnlySessionMap) // 获取逻辑和 Load 类似，read 不存在则查询 dirty
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlySessionMap)
		e, ok = read.m[key]
		if !ok && read.amended {
			delete(m.dirty, key)
		}
		m.mu.Unlock()
	}
	// 查询到 entry 后执行删除
	if ok {
		// 将 entry.p 标记为 nil，数据并没有实际删除
		// 真正删除数据并被被置为 expunged，是在 Store 的 tryExpungeLocked 中
		e.delete()
	}
}

func (e *entrySessionMap) delete() (hadValue bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expungedSessionMap {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return true
		}
	}
}

func (m *SessionMap) Range(f func(key uint, value server.Session) bool) {
	read, _ := m.read.Load().(readOnlySessionMap)
	if read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlySessionMap)
		if read.amended {
			read = readOnlySessionMap{m: m.dirty}
			m.read.Store(read)
			m.dirty = nil
			m.misses = 0
		}
		m.mu.Unlock()
	}

	for k, e := range read.m {
		v, ok := e.load()
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}

func (m *SessionMap) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnlySessionMap{m: m.dirty}) // 当 miss 积累过多，会将 dirty 存入 read，然后 将 amended = false，且 m.dirty = nil
	m.dirty = nil
	m.misses = 0
}

func (m *SessionMap) dirtyLocked() {
	if m.dirty != nil {
		return
	}
	read, _ := m.read.Load().(readOnlySessionMap)
	m.dirty = make(map[uint]*entrySessionMap, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() { // 判断 entry 是否被删除，否则就存到 dirty 中
			m.dirty[k] = e
		}
	}
}

func (e *entrySessionMap) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expungedSessionMap) { // 如果有 p == nil（即键值对被 delete），则会在这个时机被置为 expunged
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expungedSessionMap
}
