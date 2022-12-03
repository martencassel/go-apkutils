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
		w.WriteString("hello-world-1")
		n1, result, err := TarGzip("hello1", buf.Bytes(), true)
		assert.NoError(t, err)
		if err != nil {
			t.Fatal("Error writing tar file:", err)
		}
		buf = bytes.Buffer{}
		w = bufio.NewWriter(&buf)
		w.WriteString("hello-world-2")
		n2, result, err := TarGzip("hello2", buf.Bytes(), false)
		if err != nil {
			t.Fatal("Error writing tar file:", err)
		}
		dgst := digest.FromBytes(result)
		ioutil.WriteFile("test.tar.gz", result, 0644)
		assert.Equal(t, 11, n1)
		assert.Equal(t, 11, n2)
		assert.True(t, dgst.String() == "sha256:9e0e1c95ae9cafec545573dd6827fcff2e6587fea6aa9eacfae6d530e59a4150")
	})

}
