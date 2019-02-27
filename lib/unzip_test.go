package unzip4win

import (
	"path/filepath"
	"testing"
)

func TestIsLookLikeZipFile(t *testing.T) {
	t.Run("is zip file", func(t *testing.T) {
		path := filepath.Join(".", "hoge", "fuga.zip")
		actual := IsLookLikeZipFile(path)
		if !actual {
			t.Error("expected: true, but actual is false")
		}
	})
	t.Run("is not zip file", func(t *testing.T) {
		path := filepath.Join(".", "hoge", "fuga.zip.txt")
		actual := IsLookLikeZipFile(path)
		if actual {
			t.Error("expected: false, but actual is true")
		}
	})
}
