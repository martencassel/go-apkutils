package apkutils

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ApkIndex is an APKINDEX header and payload
type ApkIndex struct {
	PullChecksum string
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
