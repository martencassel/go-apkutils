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
	entries := make(map[string]*apkutils.IndexEntry)
	checksum := ""
	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)
		if s == "" {
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
				entry := &apkutils.IndexEntry{
					PullChecksum: lineData,
				}
				entries[lineData] = entry
				checksum = lineData
				break
			}
		case "P":
			{
				fmt.Println(lineData)
				record := entries[checksum]
				record.PackageName = lineData
				break
			}
		case "V":
			{
				record := entries[checksum]
				record.PackageVersion = lineData
				break
			}
		case "A":
			{
				record := entries[checksum]
				record.PackageArchitecture = lineData
				break
			}
		case "S":
			{
				record := entries[checksum]
				record.PackageSize = lineData
				break
			}
		case "I":
			{
				record := entries[checksum]
				record.PackageInstalledSize = lineData
				break
			}
		case "T":
			{
				record := entries[checksum]
				record.PackageDescription = lineData
				break
			}
		case "U":
			{
				record := entries[checksum]
				record.PackageUrl = lineData
				break
			}
		case "L":
			{
				record := entries[checksum]
				record.PackageLicense = lineData
				break
			}
		case "o":
			{
				record := entries[checksum]
				record.PackageOrigin = lineData
				break
			}
		case "m":
			{
				record := entries[checksum]
				record.PackageMaintainer = lineData
				break
			}
		case "t":
			{
				record := entries[checksum]
				record.BuildTimeStamp = lineData
				break
			}
		case "c":
			{
				record := entries[checksum]
				record.GitCommitAport = lineData
				break
			}
		case "D":
			{
				record := entries[checksum]
				record.PullDependencies = lineData
				break
			}
		case "p":
			{
				record := entries[checksum]
				record.PackageProvides = lineData
				break
			}
		}
	}
	v := make([]*apkutils.IndexEntry, 0, len(entries))
	for _, value := range entries {
		v = append(v, value)
	}
	return &apkutils.ApkIndex{Entries: v}
}
