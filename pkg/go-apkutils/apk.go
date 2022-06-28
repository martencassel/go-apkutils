package apkutils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
)

// readApkFile reads the APK file and parses metadata from .PKGINFO and extracts the pull checksum and its size.
func readApkFile(f io.Reader) (*ApkFile, error) {
	var pkgInfo *PkgInfo
	var pullChecksum string
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	packageFileSize := len(buf.Bytes())
	bytes_ := buf.Bytes()
	var offsets []int
	for i, _ := range bytes_ {
		if ReadGzipHeader(bytes_[i:]) {
			offsets = append(offsets, i)
		}
	}
	// length := len(bytes_)
	//	block1 := bytes_[offsets[0]:offsets[1]]
	block2 := bytes_[offsets[1]:offsets[2]]
	//	block3 := bytes_[offsets[2]:length]
	// ioutil.WriteFile("block1.gz", block1, 0644)
	// ioutil.WriteFile("block2.gz", block2, 0644)
	// ioutil.WriteFile("block3.gz", block3, 0644)
	Sha1CheckSum := sha1.Sum(block2)
	str := base64.StdEncoding.EncodeToString(Sha1CheckSum[:20])
	pullChecksum = fmt.Sprintf("Q1%s", str)
	br := bytes.NewReader(block2)
	uncompressedStream, err := gzip.NewReader(br)
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if header.Name == ".PKGINFO" {
			bs, _ := ioutil.ReadAll(tarReader)
			_pkgInfo, err := readPkgInfo(bytes.NewBuffer(bs))
			if err != nil {
				panic(err)
			}
			pkgInfo = _pkgInfo
		}
	}
	defer uncompressedStream.Close()
	return &ApkFile{
		PkgFileSize:  packageFileSize,
		PullChecksum: pullChecksum,
		PkgInfo:      pkgInfo,
	}, nil
}
