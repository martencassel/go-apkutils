package index

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	apkutils "github.com/martencassel/go-apkutils"
)

// ReadApkIndex reads an APKINDEX file to a ApkIndex struct.
func ReadApkIndex(f io.Reader) (*apkutils.ApkIndex, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	// If the file is gzipped tar, we need to extract APKINDEX from the tar.
	if apkutils.ReadGzipHeader(buf.Bytes()) {
		gr := bytes.NewReader(buf.Bytes())
		uncompressedStream, err := gzip.NewReader(gr)
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
			if header.Name == "APKINDEX" {
				bs, _ := ioutil.ReadAll(tarReader)
				index := parseApkIndex(bytes.NewBuffer(bs))
				return index, nil
			}
		}
	}
	index := parseApkIndex(buf)
	return index, nil
}

// parseApkIndex parses an APKINDEX file into a ApkIndex struct.
func parseApkIndex(buf *bytes.Buffer) *apkutils.ApkIndex {
	scanner := bufio.NewScanner(strings.NewReader(buf.String()))
	var indexEntries []apkutils.IndexEntry
	var curEntry apkutils.IndexEntry
	pkgSet := make(map[string]int)
	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)
		if s == "" {
			if pkgSet[curEntry.PullChecksum] == 1 {
				continue
			}
			pkgSet[curEntry.PullChecksum] += 1
			indexEntries = append(indexEntries, curEntry)
			continue
		}
		if line[1] != ':' {
			continue
		}
		lineTag := fmt.Sprintf("%c", line[0])
		lineData := line[2:]
		switch lineTag {
		case "C":
			{
				curEntry = apkutils.IndexEntry{
					PullChecksum: lineData,
				}
				break
			}
		case "P":
			{
				curEntry.PackageName = lineData
				break
			}
		case "V":
			{
				curEntry.PackageVersion = lineData
				break
			}
		case "A":
			{
				curEntry.PackageArchitecture = lineData
				break
			}
		case "S":
			{
				curEntry.PackageSize = lineData
				break
			}
		case "I":
			{
				curEntry.PackageInstalledSize = lineData
				break
			}
		case "T":
			{
				curEntry.PackageDescription = lineData
				break
			}
		case "U":
			{
				curEntry.PackageUrl = lineData
				break
			}
		case "L":
			{
				curEntry.PackageLicense = lineData
				break
			}
		case "o":
			{
				curEntry.PackageOrigin = lineData
				break
			}
		case "m":
			{
				curEntry.PackageMaintainer = lineData
				break
			}
		case "t":
			{
				curEntry.BuildTimeStamp = lineData
				break
			}
		case "c":
			{
				curEntry.GitCommitAport = lineData
				break
			}
		case "D":
			{
				curEntry.PullDependencies = lineData
				break
			}
		case "p":
			{
				curEntry.PackageProvides = lineData
				break
			}
		}
	}
	return &apkutils.ApkIndex{Entries: indexEntries}
}
