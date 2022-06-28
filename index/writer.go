package index

import (
	"io"

	apkutils "github.com/martencassel/go-apkutils"
)

// Writer provides sequential writing of APKINDEX index entries from APK files
type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) WriteIndexEntry(apk *apkutils.ApkFile) {
	s := apk.ToIndexEntry()
	io.WriteString(w.w, s)
}

func (w *Writer) Write(entry *apkutils.IndexEntry) {
	io.WriteString(w.w, entry.String())
}

func (w *Writer) Close() {
	io.WriteString(w.w, "\n")
}
