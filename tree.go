package gdu

import "sync"

// TreeNode describes a single node in the output file tree
type TreeNode struct {
	Parent *TreeNode
	name   string
	size   uint64
	files  map[string]Entry
	mutex  sync.RWMutex
}

// NewNode is a constructor for TreeNode
func NewNode(name string) TreeNode {
	return TreeNode{
		name:  name,
		size:  0,
		files: make(map[string]Entry),
	}
}

// Kind returns the kind of entry (directory)
func (*TreeNode) Kind() string {
	return "Directory"
}

// AddSize adds to the file size
func (n *TreeNode) AddSize(by uint64) {
	n.mutex.RLock()
	n.size += by
	n.mutex.RUnlock()
	return
}

// GetSize gets the file size of the entire node
func (n *TreeNode) GetSize() (size uint64) {
	n.mutex.RLock()
	size = n.size
	n.mutex.RUnlock()
	return
}

// GetName returns the name of the directory, appended with "/"
func (n *TreeNode) GetName() string {
	return n.name
}

// Get retrieves the entry for a particular file
func (n *TreeNode) Get(key string) (value Entry) {
	n.mutex.RLock()
	value = n.files[key]
	n.mutex.RUnlock()
	return
}

// GetEntries returns the children as a list of Entrys
func (n *TreeNode) GetEntries() Entries {
	n.mutex.RLock()
	result := make([]Entry, len(n.files))
	i := 0
	for _, entry := range n.files {
		result[i] = entry
		i++
	}
	n.mutex.RUnlock()
	return result
}

// GetFiles returns the children as a list of filenames
func (n *TreeNode) GetFiles() []string {
	n.mutex.RLock()
	result := make([]string, len(n.files))
	i := 0
	for name := range n.files {
		result[i] = name
		i++
	}
	n.mutex.RUnlock()
	return result
}

// Insert inserts a new entry into the tree
func (n *TreeNode) Insert(key string, value Entry) {
	n.mutex.Lock()
	n.files[key] = value
	n.mutex.Unlock()
}
