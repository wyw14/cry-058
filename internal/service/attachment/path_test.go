package attachment

import (
	"bytes"
	"testing"
)

func TestAttachmentStoreRejectsTraversal(t *testing.T) {
	s := New(t.TempDir())
	if _, e := s.Save("owner", "../escape.txt", "image/png", bytes.NewBufferString("x")); e == nil {
		t.Fatal("unsafe attachment accepted")
	}
}
