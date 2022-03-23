// +build go1.16

package zstatic

import (
	"embed"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
)

func EmbedMustBytes(fs embed.FS, name string) ([]byte, error) {
	file, err := fs.ReadFile(name)
	if err == nil {
		return file, nil
	}
	return MustBytes(name)
}

func EmbedBytes(fs embed.FS, name string) []byte {
	file, _ := EmbedMustBytes(fs, name)
	return file
}

func EmbedMustString(fs embed.FS, name string) (string, error) {
	file, err := fs.ReadFile(name)
	if err == nil {
		return zstring.Bytes2String(file), nil
	}
	return MustString(name)
}

func EmbedString(fs embed.FS, name string) string {
	file, _ := EmbedMustString(fs, name)
	return file
}

func NewEmbedFileserver(fs embed.FS, dir string, fn ...func(ctype string, content []byte, err error)) func(c *znet.Context) {
	const defFile = "index.html"

	isCb := len(fn) > 0
	return func(c *znet.Context) {
		name := c.GetParam("file")
		if name == "" {
			name = defFile
		}
		content, err := EmbedMustBytes(fs, filepath.Join(dir, name))
		ctype := mime.TypeByExtension(filepath.Ext(name))
		if isCb {
			fn[0](ctype, content, err)
			return
		}
		if err != nil {
			c.String(404, err.Error())
			return
		}
		c.SetContentType(ctype)
		c.Byte(http.StatusOK, content)
	}
}
