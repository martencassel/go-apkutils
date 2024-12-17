package apkutils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"path/filepath"
)

// ReadGzipHeader reads the header of a gzip file if found.
// There are 3 signature bytes that occur in a specific order.
func ReadGzipHeader(buf []byte) bool {
	if len(buf) < 3 {
		return false
	}
	if buf[0] != GzipID1 || buf[1] != GzipID2 || buf[2] != GzipDeflate {
		return false
	}
	return true
}

// TarGzip create a tar.gz file data from some source bytes.
// Optionally, you can specify not to write an EnfOfTar header.
// This function can be used to create signature.tar.gz files for signed APKINDEX files,
// and APKINDEX.unsigned.tar.gz of an APKINDEX file.
func TarGzip(filename string, b []byte, writeEOFTar bool) (int, []byte, error) {
	nRead := len(b)
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	tw.WriteHeader(&tar.Header{
		Name: filepath.Base(filename),
		Size: int64(nRead),
		Mode: 0600,
	})
	n, err := tw.Write(b)
	if err != nil {
		return 0, nil, err
	}
	if writeEOFTar {
		err = tw.Close()
		if err != nil {
			return 0, nil, err
		}
	}
	err = gz.Close()
	if err != nil {
		return 0, nil, err
	}
	ret := buf.Bytes()
	return n, ret, nil
}
