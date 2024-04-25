package render

import "strings"

// A File represents anything like a source code file, or xml, or json, or an image.
type File struct {
	FileName string
	MimeType string
	Buf      []byte
	Error    error
}

func (n *File) Name() string {
	return n.FileName
}

func (n *File) String() string {
	return n.StringIndent(0)
}

func (n *File) StringIndent(indent int) string {
	sb := &strings.Builder{}
	indentMe(sb, indent)
	sb.WriteString(n.FileName)
	sb.WriteString("[")
	sb.WriteString(n.MimeType)
	sb.WriteString("]:\n")
	if n.Error != nil {
		sb.WriteString(string(n.Buf))
		sb.WriteString("\n\n")
		sb.WriteString(n.Error.Error())
	} else {
		for _, s := range strings.Split(string(n.Buf), "\n") {
			indentMe(sb, indent+2)
			sb.WriteString(s)
			sb.WriteByte('\n')
		}

	}

	return sb.String()
}
