package zstatic

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sohaha/zstatic/build"
)

type FS struct {
	g *build.FileGroup
}

var _ http.FileSystem = (*FS)(nil)

func NewFS(pattern string) (*FS, error) {
	g, err := Group(pattern)
	if err != nil {
		return nil, err
	}

	return &FS{g: g}, nil
}

func (f *FS) String(name string) string {
	return f.g.String(name)
}

func (f *FS) Open(name string) (http.File, error) {
	b, err := f.g.MustBytes(name)
	if err != nil {
		return nil, err
	}
	return newByteFile(b), nil
}

type byteFile struct {
	content []byte
	pos     int
}

func (f *byteFile) Close() error {
	return nil
}

func (f *byteFile) Read(p []byte) (int, error) {
	if f.pos >= len(f.content) {
		return 0, io.EOF
	}
	n := copy(p, f.content[f.pos:])
	f.pos += n
	return n, nil
}

func (f *byteFile) Seek(offset int64, whence int) (int64, error) {
	var newPos int
	switch whence {
	case io.SeekStart:
		newPos = int(offset)
	case io.SeekCurrent:
		newPos = f.pos + int(offset)
	case io.SeekEnd:
		newPos = len(f.content) + int(offset)
	}
	if newPos < 0 || newPos >= len(f.content) {
		return 0, errors.New("invalid seek position")
	}
	f.pos = newPos
	return int64(f.pos), nil
}

func (f *byteFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *byteFile) Stat() (os.FileInfo, error) {
	return &byteFileInfo{len(f.content)}, nil
}

type byteFileInfo struct {
	size int
}

func (f *byteFileInfo) Name() string {
	return ""
}

func (f *byteFileInfo) Size() int64 {
	return int64(f.size)
}

func (f *byteFileInfo) Mode() os.FileMode {
	return 0444
}

func (f *byteFileInfo) ModTime() time.Time {
	return time.Time{}
}

func (f *byteFileInfo) IsDir() bool {
	return false
}

func (f *byteFileInfo) Sys() interface{} {
	return nil
}

func newByteFile(content []byte) http.File {
	return &byteFile{content: content}
}
