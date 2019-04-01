package gdu

type Entry interface {
	Kind() string
	GetSize() uint64
	GetName() string
}

type Entries []Entry

func (e Entries) Len() int {
	return len(e)
}

func (e Entries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Entries) Less(i, j int) bool {
	return e[i].GetSize() < e[j].GetSize()
}

type FileEntry struct {
	name string
	size uint64
}

func (*FileEntry) Kind() string {
	return "File"
}

func (n *FileEntry) GetSize() uint64 {
	return n.size
}

func (n *FileEntry) GetName() string {
	return n.name
}

func NewFileEntry(name string, size uint64) FileEntry {
	return FileEntry{name, size}
}
