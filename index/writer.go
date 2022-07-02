package index

import (
	"io"

	apkutils "github.com/martencassel/go-apkutils"
)

// Writer provides sequential writing of APKINDEX index entries from APK files
type Writer struct {
	w io.Writer
}

// NewWriter create a new writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// WriteApk writes a APK file to an APKINDEX file, ie. a APKINDEX record for the APK.
func (w *Writer) WriteApk(apk *apkutils.ApkFile) {
	s := apk.ToIndexEntry()
	io.WriteString(w.w, s)
}

// Write writes a APKINDEX record to the underlying writer.
func (w *Writer) Write(entry *apkutils.IndexEntry) {
	data := entry.String()
	io.WriteString(w.w, data)
}

// Close closes the underlying writer, by writing a \n line.
func (w *Writer) Close() {
	io.WriteString(w.w, "\n")
}
