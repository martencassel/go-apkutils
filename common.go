package apkutils

func ReadGzipHeader(buf []byte) bool {
	if buf[0] != GzipID1 || buf[1] != GzipID2 || buf[2] != GzipDeflate {
		return false
	}
	return true
}
