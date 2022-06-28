package index

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadIndex(t *testing.T) {
	t.Run("Read a APKINDEX file", func(t *testing.T) {
		f, err := os.Open("../testdata/APKINDEX")
		if err != nil {
			t.Fatal("Error opening APKINDEX file:", err)
		}
		index, err := ReadApkIndex(f)
		if err != nil {
			t.Fatal("Error reading APKINDEX file:", err)
		}
		assert.Equal(t, len(index.Entries), 3)
	})
}
