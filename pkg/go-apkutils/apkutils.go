package apkutils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// ApkIndex is an APKINDEX header and payload
type ApkIndex struct {
	PullChecksum string
	f            io.Reader
}

// ReadIndex reads the APKINDEX file.
func ReadIndex(f io.Reader) *apkIndex {
	index, err := readApkIndex(f)
	if err != nil {
		panic(err)
	}
	return index
}

type ApkFile struct {
	f            io.Reader
	PullChecksum string
	PkgInfo      *PkgInfo
	PkgFileSize  int
}

func (apkFile *ApkFile) ToIndexEntry() string {
	return fmt.Sprintf("C:%s\nP:%s\nV:%s\nA:%s\nS:%s\nI:%s\nT:%s\nU:%s\nL:%s\no:%s\nm:%s\nt:%s\nc:%s\nD:%s\np:%s\n\n",
		apkFile.PullChecksum,                           // C
		apkFile.PkgInfo.PkgName,                        // P
		apkFile.PkgInfo.PkgVer,                         // V
		apkFile.PkgInfo.PkgArch,                        // A
		strconv.Itoa(apkFile.PkgFileSize),              // S
		apkFile.PkgInfo.PkgSize,                        // I
		apkFile.PkgInfo.PkgDesc,                        // T
		apkFile.PkgInfo.PkgUrl,                         // U
		apkFile.PkgInfo.PkgLicense,                     // L
		apkFile.PkgInfo.PkgOrigin,                      // o
		apkFile.PkgInfo.PkgMaintainer,                  // m
		apkFile.PkgInfo.PkgBuildDate,                   // t
		apkFile.PkgInfo.PkgCommit,                      // c
		strings.Join(apkFile.PkgInfo.PkgDepends, " "),  // D
		strings.Join(apkFile.PkgInfo.PkgProvides, " ")) // p
}

// ReadApkFile reads the APK file.
func ReadApkFile(f io.Reader) (*ApkFile, error) {
	apk, err := readApkFile(f)
	if err != nil {
		return nil, err
	}
	return apk, nil
}

func createTarArchive(filename string, buf io.Writer) error {
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	err := addToArchive(tw, filename)
	if err != nil {
		return err
	}
	return nil
}

func addToArchive(tw *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	header.Name = filename
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}
	return nil
}
