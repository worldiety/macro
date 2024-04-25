package render

import "strings"

// A Dir contains other dirs and files.
type Dir struct {
	DirName  string
	MimeType string
	Files    []*File
	Dirs     []*Dir
}

func (n *Dir) Name() string {
	return n.DirName
}

func (n *Dir) Directory(name string) *Dir {
	for _, dir := range n.Dirs {
		if dir.Name() == name {
			return dir
		}
	}

	return nil
}

func (n *Dir) String() string {
	return n.StringIndent(0)
}

func (n *Dir) StringIndent(indent int) string {
	sb := &strings.Builder{}
	indentMe(sb, indent)
	sb.WriteString(n.DirName)
	sb.WriteString("[")
	sb.WriteString(n.MimeType)
	sb.WriteString("]:\n")
	for _, file := range n.Files {
		sb.WriteString(file.StringIndent(indent + 2))
	}

	for _, dir := range n.Dirs {
		sb.WriteString(dir.StringIndent(indent + 2))
	}

	return sb.String()
}

func indentMe(sb *strings.Builder, idn int) {
	for i := 0; i < idn; i++ {
		sb.WriteByte(' ')
	}
}
