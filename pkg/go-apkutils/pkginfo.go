package apkutils

import (
	"bufio"
	"bytes"
	"strings"
)

type PkgInfo struct {
	PkgName       string
	PkgVer        string
	PkgDesc       string
	PkgUrl        string
	PkgBuildDate  string
	PkgPackager   string
	PkgSize       string
	PkgArch       string
	PkgOrigin     string
	PkgCommit     string
	PkgMaintainer string
	PkgLicense    string
	PkgProvides   []string
	PkgDepends    []string
	PkgDataHash   string
}

func readPkgInfo(buf *bytes.Buffer) (*PkgInfo, error) {
	pkgInfo := &PkgInfo{}
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
