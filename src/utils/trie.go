package utils
import (
"unicode/utf8"
)

/*
脏字过滤库
*/

type Trie struct {
	Root *TrieNode
}

type TrieNode struct {
	Children map[rune]*TrieNode
	End      bool
}

func NewTrie() Trie {
	var r Trie
	r.Root = NewTrieNode()
	return r
}

func NewTrieNode() *TrieNode {
	n := new(TrieNode)
	n.Children = make(map[rune]*TrieNode)
	return n
}

func (this *Trie) Inster(txt string) {
	if len(txt) < 1 {
		return
	}
	node := this.Root
	key := []rune(txt)
	for i := 0; i < len(key); i++ {
		if _, exists := node.Children[key[i]]; !exists {
			node.Children[key[i]] = NewTrieNode()
		}
		node = node.Children[key[i]]
	}

	node.End = true
}

func (this *Trie) HasDirty(txt string) bool {
	if len(txt) < 1 {
		return false
	}
	node := this.Root
	key := []rune(txt)
	var chars []rune = nil
	slen := len(key)
	for i := 0; i < slen; i++ {
		if _, exists := node.Children[key[i]]; exists {
			node = node.Children[key[i]]
			for j := i + 1; j < slen; j++ {
				if _, exists := node.Children[key[j]]; exists {
					node = node.Children[key[j]]
					if node.End == true {
						if chars == nil {
							chars = key
						}
						for t := i; t <= j; t++ {
							return true
						}
						i = j
						node = this.Root
						break
					}
				}
			}
			node = this.Root
		}
	}
	return false
}

func (this *Trie) Replace(txt string) string {
	if len(txt) < 1 {
		return txt
	}
	node := this.Root
	key := []rune(txt)
	var chars []rune = nil
	slen := len(key)
	for i := 0; i < slen; i++ {
		if _, exists := node.Children[key[i]]; exists {
			node = node.Children[key[i]]
			for j := i + 1; j < slen; j++ {
				if _, exists := node.Children[key[j]]; exists {
					node = node.Children[key[j]]
					if node.End == true {
						if chars == nil {
							chars = key
						}
						for t := i; t <= j; t++ {
							c, _ := utf8.DecodeRuneInString("*")
							chars[t] = c
						}
						i = j
						node = this.Root
						break
					}
				}
			}
			node = this.Root
		}
	}
	if chars == nil {
		return txt
	} else {
		return string(chars)
	}
}

