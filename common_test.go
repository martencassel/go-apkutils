package apkutils

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"testing"

	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

func TestTarFile(t *testing.T) {
	t.Run("Test write a .tar.gz file", func(t *testing.T) {
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		w.WriteString("Hello World")
		w.Flush()
		n, result, err := TarGzip("test.tar.gz", buf.Bytes(), true)
		if err != nil {
			t.Fatal("Error writing tar file:", err)
		}
		dgst := digest.FromBytes(result)
		ioutil.WriteFile("test.tar.gz", result, 0644)
		assert.Equal(t, 11, n)
		assert.True(t, dgst.String() == "sha256:1c5db109432161a7e01e5dcaba49a3b011749317dd1d54925ccf65df482f10d9")
	})

}
