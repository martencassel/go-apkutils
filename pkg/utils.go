package apkutils

func ReadGzipHeader(buf []byte) bool {
	if buf[0] != gzipID1 || buf[1] != gzipID2 || buf[2] != gzipDeflate {
		return false
	}
	return true
}
