// +build go1.16

package zstatic

import (
	"embed"
)

func EmbedMustBytes(fs embed.FS, name string) ([]byte, error) {
	file, err := fs.ReadFile(name)
	if err == nil {
		return file, nil
	}
	return MustBytes(name)
}
