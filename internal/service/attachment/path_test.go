package attachment

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAttachmentStoreRejectsTraversal(t *testing.T) {
	s := New(t.TempDir())
	if _, e := s.Save("owner", "../escape.txt", "image/png", bytes.NewBufferString("x")); e == nil {
		t.Fatal("unsafe attachment accepted")
	}
}

func TestSaveRejectsUnsafeNames(t *testing.T) {
	cases := []string{
		"..",
		".",
		"./sneak.txt",
		"sub/../escape.txt",
		"../../etc/passwd",
		"/etc/passwd",
		"a/b.png",
		"a\\b.png",
		"win\\evil.txt",
		"web/../../../../oob.png",
	}
	s := New(t.TempDir())
	for _, name := range cases {
		if a, e := s.Save("owner", name, "image/png", bytes.NewBufferString("x")); e == nil {
			t.Errorf("name %q was accepted: %+v", name, a)
		}
	}
}

func TestSaveRejectsUnsafeOwners(t *testing.T) {
	cases := []string{
		"",
		"..",
		"../admin",
		"a/b",
		"a\\b",
		"/root",
	}
	s := New(t.TempDir())
	for _, owner := range cases {
		if a, e := s.Save(owner, "ok.png", "image/png", bytes.NewBufferString("x")); e == nil {
			t.Errorf("owner %q was accepted: %+v", owner, a)
		}
	}
}

func TestSaveRejectsUnsupportedMedia(t *testing.T) {
	s := New(t.TempDir())
	if _, e := s.Save("owner", "ok.png", "application/x-msdownload", bytes.NewBufferString("x")); e == nil {
		t.Fatal("unsafe media accepted")
	}
	if _, e := s.Save("owner", "ok.png", "", bytes.NewBufferString("x")); e == nil {
		t.Fatal("empty media accepted")
	}
}

func TestSaveRejectsOversize(t *testing.T) {
	s := New(t.TempDir())
	big := bytes.Repeat([]byte{'x'}, 10*1024*1024+1)
	if _, e := s.Save("owner", "big.png", "image/png", bytes.NewReader(big)); e == nil {
		t.Fatal("oversize attachment accepted")
	}
	// nothing leaked to disk
	if _, e := os.Stat(filepath.Join(s.root, "owner", "big.png")); !os.IsNotExist(e) {
		t.Fatalf("oversize file left on disk: %v", e)
	}
}

func TestSaveHonoursHappyPath(t *testing.T) {
	root := t.TempDir()
	s := New(root)
	body := []byte("hello")
	a, e := s.Save("owner", "doc.png", "image/png", bytes.NewReader(body))
	if e != nil {
		t.Fatalf("save failed: %v", e)
	}
	if a.OwnerID != "owner" || a.FileName != "doc.png" || a.MediaType != "image/png" {
		t.Fatalf("bad metadata: %+v", a)
	}
	if a.Size != int64(len(body)) {
		t.Fatalf("size mismatch: got %d want %d", a.Size, len(body))
	}
	// path must live inside the owner directory and the file must be readable
	if !strings.HasPrefix(filepath.ToSlash(a.Path), filepath.ToSlash(filepath.Join(root, "owner"))+"/") {
		t.Fatalf("path escapes owner dir: %q", a.Path)
	}
	got, e := os.ReadFile(a.Path)
	if e != nil {
		t.Fatalf("read back failed: %v", e)
	}
	if !bytes.Equal(got, body) {
		t.Fatalf("content mismatch: got %q want %q", got, body)
	}
	// digest must match the content
	want := sha256.Sum256(body)
	if a.Digest != hex.EncodeToString(want[:]) {
		t.Fatalf("digest mismatch: got %s want %s", a.Digest, hex.EncodeToString(want[:]))
	}
}

func TestSaveRejectsDuplicateAndKeepsFirst(t *testing.T) {
	s := New(t.TempDir())
	body := []byte("first")
	if _, e := s.Save("owner", "doc.png", "image/png", bytes.NewReader(body)); e != nil {
		t.Fatalf("first save failed: %v", e)
	}
	if _, e := s.Save("owner", "doc.png", "image/png", bytes.NewReader([]byte("second"))); e == nil {
		t.Fatal("duplicate save accepted")
	}
	got, e := os.ReadFile(filepath.Join(s.root, "owner", "doc.png"))
	if e != nil {
		t.Fatalf("read failed: %v", e)
	}
	if !bytes.Equal(got, body) {
		t.Fatalf("first content clobbered: got %q want %q", got, body)
	}
}

func TestSaveCleansUpOnCopyError(t *testing.T) {
	s := New(t.TempDir())
	// a reader that errors after some bytes should leave no partial file
	_, e := s.Save("owner", "doc.png", "image/png", &boomReader{})
	if e == nil {
		t.Fatal("expected copy error")
	}
	if _, st := os.Stat(filepath.Join(s.root, "owner", "doc.png")); !os.IsNotExist(st) {
		t.Fatalf("partial file left on disk")
	}
}

type boomReader struct{ n int }

func (*boomReader) Read(p []byte) (int, error) {
	return 0, errors.New("boom")
}

var _ io.Reader = (*boomReader)(nil)
