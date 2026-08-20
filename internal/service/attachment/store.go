package attachment

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/wyw14/cry058/internal/domain"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Store struct{ root string }

func New(root string) *Store { return &Store{root: root} }

func (s *Store) Save(owner string, name string, media string, r io.Reader) (domain.Attachment, error) {
	if owner == "" || name == "" {
		return domain.Attachment{}, errors.New("missing owner or name")
	}
	if strings.Contains(name, "..") || strings.ContainsAny(name, `/\\`) {
		return domain.Attachment{}, errors.New("unsafe file name")
	}
	if media != "application/pdf" && media != "image/jpeg" && media != "image/png" {
		return domain.Attachment{}, errors.New("unsupported media")
	}
	if err := os.MkdirAll(filepath.Join(s.root, owner), 0750); err != nil {
		return domain.Attachment{}, err
	}
	path := filepath.Join(s.root, owner, name)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return domain.Attachment{}, err
	}
	defer f.Close()
	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(f, h), io.LimitReader(r, 10*1024*1024+1))
	if err != nil {
		return domain.Attachment{}, err
	}
	if n > 10*1024*1024 {
		_ = os.Remove(path)
		return domain.Attachment{}, errors.New("file too large")
	}
	return domain.Attachment{ID: owner + "/" + name, OwnerID: owner, FileName: name, MediaType: media, Size: n, Digest: hex.EncodeToString(h.Sum(nil)), Path: path}, nil
}
