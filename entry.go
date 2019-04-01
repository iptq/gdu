package gdu

// Entry describes any entry in the final output table
type Entry interface {
	Kind() string
	GetSize() uint64
	GetName() string
}

// Entries is just a list of the type Entry, used to implement sort.Interface
type Entries []Entry

func (e Entries) Len() int {
	return len(e)
}

func (e Entries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Entries) Less(i, j int) bool {
	if e[i].Kind() == "Directory" && e[j].Kind() == "File" {
		return true
	} else if e[i].Kind() == "File" && e[j].Kind() == "Directory" {
		return false
	} else {
		return e[i].GetSize() > e[j].GetSize()
	}
}

// FileEntry describes a regular file in the final output
type FileEntry struct {
	name string
	size uint64
}

// Kind returns the kind of entry (file)
func (*FileEntry) Kind() string {
	return "File"
}

// GetSize returns the size of the file
func (n *FileEntry) GetSize() uint64 {
	return n.size
}

// GetName returns the name of the file
func (n *FileEntry) GetName() string {
	return n.name
}

// NewFileEntry is the constructor for a FileEntry
func NewFileEntry(name string, size uint64) FileEntry {
	return FileEntry{name, size}
}
