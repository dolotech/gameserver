package utils

import (
	"encoding/json"
	"gameserver/utils/log"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

var josonConfigOnce sync.Once
var jsonConfig *JsonConfig

func Json() *JsonConfig {
	josonConfigOnce.Do(func() {
		jsonConfig = &JsonConfig{
			resource:      make(map[string][]byte, 50),
			refleMap:      make(map[string]reflect.Type, 50),
			unsafePointer: make(map[string]interface{}, 50),
		}
	})
	return jsonConfig
}

type JsonConfig struct {
	resource      map[string][]byte
	refleMap      map[string]reflect.Type
	unsafePointer map[string]interface{}
}

func (this *JsonConfig) Fill() {
	for name, v := range this.refleMap {
		inst := reflect.New(v)
		if data, ok := this.resource[name]; ok {
			err := json.Unmarshal(data, inst.Interface())
			if err != nil {
				log.Error(err)
				continue
			}
			pointer := this.unsafePointer[name]
			instPointer := unsafe.Pointer(inst.Pointer())
			p := unsafe.Pointer(reflect.ValueOf(pointer).Pointer())
			atomic.SwapPointer((*unsafe.Pointer)(p), *(*unsafe.Pointer)(instPointer))
		}
	}
}

func (this *JsonConfig) LoadFile(pathDir string) {
	files := []string{}
	traverseDir(pathDir, &files)
	for _, table1 := range files {
		fileSuffix := path.Ext(table1)
		if fileSuffix == ".json" {
			this.parserIWant(table1)
		}
	}
}

func (this *JsonConfig) RegistJson(name string, i interface{}) {
	this.refleMap[name] = reflect.TypeOf(i).Elem()
	this.unsafePointer[name] = i
}

//解析json文件并检测需要字段用于计数
func (this *JsonConfig) parserIWant(filePath string) {
	_, e := os.Stat(filePath)
	if e != nil {
		log.Error("File Open Error", filePath)
		return
	}
	data, err := ioutil.ReadFile(filePath)
	if err == io.EOF {
		log.Error("File Is Commplete")
		return
	} else if err != nil {
		log.Error(err)
		return
	}
	if len(data) == 0 {
		log.Error("len(data) ==0 ")
		return
	}
	arr := strings.Split(path.Base(filePath), string(os.PathSeparator))
	name := arr[len(arr)-1]
	name = strings.Replace(name, path.Ext(filePath), "", -1)
	this.resource[name] = data
}

//获取指定目录下的所有文件和目录
func traverseDir(dirPath string, files *[]string) (err error) {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Info(err, dirPath)
		return err
	}
	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			//traverseDir(dirPth+PthSep+fi.Name(), files)
		} else {
			//log.Info(dirPath+PthSep+fi.Name())
			*files = append(*files, dirPath+PthSep+fi.Name())
		}
	}
	return nil
}
