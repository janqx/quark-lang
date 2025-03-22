package parser

type SourceFile struct {
	Filename string
	Ext      string
	Size     int   // number of file bytes
	Lines    []int // number of columns per line
}

type SourceFileSet struct {
	Files []*SourceFile
}
