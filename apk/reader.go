package apk

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/martencassel/go-apkutils"
)

// ReadApk reads an ApkFile from the reader into an ApkFile struct.
func ReadApk(r io.Reader) (*apkutils.ApkFile, error) {
	apkFile, err := ReadApkFile(r)
	if err != nil {
		return nil, err
	}
	return apkFile, nil
}

// ReadApkFile reads the APK file and parses metadata from .PKGINFO and extracts the pull checksum and its size.
func ReadApkFile(f io.Reader) (*apkutils.ApkFile, error) {
	var pkgInfo *apkutils.PkgInfo
	// https://stackoverflow.com/questions/38837679/alpine-apk-package-repositories-how-are-the-checksums-calculated
	var pullChecksum string
	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, fmt.Errorf("empty file")
	}
	packageFileSize := len(buf.Bytes())
	bytes_ := buf.Bytes()
	var offsets []int
	for i, _ := range bytes_ {
		if apkutils.ReadGzipHeader(bytes_[i:]) {
			offsets = append(offsets, i)
		}
	}
	if len(offsets) < 3 {
		return nil, fmt.Errorf("invalid apk file")
	}
	block2 := bytes_[offsets[1]:offsets[2]]
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
	return &apkutils.ApkFile{
		PkgFileSize:  packageFileSize,
		PullChecksum: pullChecksum,
		PkgInfo:      pkgInfo,
	}, nil
}

// readPkgInfo reads the .PKGINFO file and parses metadata from it.
func readPkgInfo(buf *bytes.Buffer) (*apkutils.PkgInfo, error) {
	pkgInfo := &apkutils.PkgInfo{}
	scanner := bufio.NewScanner(strings.NewReader(buf.String()))
	var depends []string
	var provides []string
	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)
		index := strings.Index(line, "=")
		if s == "" {
			continue
		}
		if index == -1 {
			continue
		}
		part1 := line[:index]
		part1 = strings.TrimPrefix(part1, " ")
		part1 = strings.TrimSuffix(part1, " ")
		part2 := line[index+1:]
		part2 = strings.TrimPrefix(part2, " ")
		part2 = strings.TrimSuffix(part2, " ")
		if index == -1 {
			continue
		}
		switch part1 {
		case "pkgname":
			pkgInfo.PkgName = part2
		case "pkgver":
			pkgInfo.PkgVer = part2
		case "pkgdesc":
			pkgInfo.PkgDesc = part2
		case "url":
			pkgInfo.PkgUrl = part2
		case "builddate":
			pkgInfo.PkgBuildDate = part2
		case "packager":
			pkgInfo.PkgPackager = part2
		case "size":
			pkgInfo.PkgSize = part2
		case "arch":
			pkgInfo.PkgArch = part2
		case "origin":
			pkgInfo.PkgOrigin = part2
		case "commit":
			pkgInfo.PkgCommit = part2
		case "maintainer":
			pkgInfo.PkgMaintainer = part2
		case "license":
			pkgInfo.PkgLicense = part2
		case "provides":
			provides = append(provides, part2)
		case "depend":
			depends = append(depends, part2)
		}
	}
	pkgInfo.PkgDepends = depends
	pkgInfo.PkgProvides = provides
	return pkgInfo, nil
}
