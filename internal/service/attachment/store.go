package attachment

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wyw14/cry058/internal/domain"
)

const maxAttachmentSize int64 = 10 * 1024 * 1024

type Store struct{ root string }

func New(root string) *Store {
	if root == "" {
		root = "."
	}
	return &Store{root: filepath.Clean(root)}
}

func (s *Store) Save(owner, name, media string, r io.Reader) (domain.Attachment, error) {
	dir, path, err := s.resolvePath(owner, name)
	if err != nil {
		return domain.Attachment{}, err
	}
	if !allowedMedia(media) {
		return domain.Attachment{}, errors.New("unsupported media")
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return domain.Attachment{}, err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return domain.Attachment{}, err
	}
	h := sha256.New()
	n, copyErr := io.Copy(io.MultiWriter(f, h), io.LimitReader(r, maxAttachmentSize+1))
	// Close before any removal: Windows refuses to delete a still-open file.
	closeErr := f.Close()
	if copyErr != nil {
		_ = os.Remove(path)
		return domain.Attachment{}, copyErr
	}
	if closeErr != nil {
		_ = os.Remove(path)
		return domain.Attachment{}, closeErr
	}
	if n > maxAttachmentSize {
		_ = os.Remove(path)
		return domain.Attachment{}, errors.New("file too large")
	}
	return domain.Attachment{
		ID:        owner + "/" + name,
		OwnerID:   owner,
		FileName:  name,
		MediaType: media,
		Size:      n,
		Digest:    hex.EncodeToString(h.Sum(nil)),
		Path:      path,
	}, nil
}

// resolvePath validates owner and name, then resolves the on-disk path while
// guaranteeing it stays inside the owner's directory. Filenames that could
// traverse out (../, separators, drive letters, absolute paths) are rejected
// rather than rewritten, and the final path is checked against the owner
// directory as defence in depth.
func (s *Store) resolvePath(owner, name string) (dir, path string, err error) {
	if err = safeSegment("owner", owner); err != nil {
		return "", "", err
	}
	if err = safeSegment("name", name); err != nil {
		return "", "", err
	}
	dir = filepath.Clean(filepath.Join(s.root, owner))
	if !within(s.root, dir) {
		return "", "", errors.New("owner escapes attachment root")
	}
	path = filepath.Clean(filepath.Join(dir, name))
	if !within(dir, path) {
		return "", "", errors.New("attachment path escapes owner directory")
	}
	return dir, path, nil
}

// safeSegment rejects any identifier that could traverse or anchor to an
// absolute path: empties, "." and "..", path separators (both Unix and
// Windows), and Windows drive letters.
func safeSegment(field, value string) error {
	if value == "" {
		return errors.New("missing " + field)
	}
	if value == "." || value == ".." || strings.ContainsAny(value, `/\`) {
		return errors.New(field + " must be a single path component")
	}
	if len(value) >= 2 && value[1] == ':' {
		return errors.New(field + " must not contain a drive letter")
	}
	return nil
}

// within reports whether path is contained in base after cleaning.
func within(base, path string) bool {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)
	return rel != ".." && !strings.HasPrefix(rel, "../")
}

func allowedMedia(media string) bool {
	switch media {
	case "application/pdf", "image/jpeg", "image/png":
		return true
	}
	return false
}
