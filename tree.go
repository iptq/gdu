package gdu

import "sync"

type TreeNode struct {
	Parent *TreeNode
	name   string
	size   uint64
	files  map[string]Entry
	mutex  sync.RWMutex
}

func NewNode(name string) TreeNode {
	return TreeNode{
		name:  name,
		size:  0,
		files: make(map[string]Entry),
	}
}

func (*TreeNode) Kind() string {
	return "Directory"
}

func (n *TreeNode) AddSize(by uint64) {
	n.mutex.RLock()
	n.size += by
	n.mutex.RUnlock()
	return
}

func (n *TreeNode) GetSize() (size uint64) {
	n.mutex.RLock()
	size = n.size
	n.mutex.RUnlock()
	return
}

func (n *TreeNode) GetName() string {
	return n.name
}

func (n *TreeNode) Get(key string) (value Entry) {
	n.mutex.RLock()
	value = n.files[key]
	n.mutex.RUnlock()
	return
}

func (n *TreeNode) GetEntries() Entries {
	n.mutex.RLock()
	result := make([]Entry, len(n.files))
	i := 0
	for _, entry := range n.files {
		result[i] = entry
		i += 1
	}
	n.mutex.RUnlock()
	return result
}

func (n *TreeNode) GetFiles() []string {
	n.mutex.RLock()
	result := make([]string, len(n.files))
	i := 0
	for name, _ := range n.files {
		result[i] = name
		i += 1
	}
	n.mutex.RUnlock()
	return result
}

func (n *TreeNode) Insert(key string, value Entry) {
	n.mutex.Lock()
	n.files[key] = value
	n.mutex.Unlock()
}
