package zstatic

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zstatic/build"
)

// mainAssetDirectory stores all the assets
var mainAssetDirectory = build.NewAssetDirectory()
var rootFileGroup *build.FileGroup
var err error

func init() {
	rootFileGroup, err = mainAssetDirectory.NewFileGroup(".")
	if err != nil {
		zlog.Fatal(err)
	}
}

// String gets the asset value by name
func String(name string) string {
	return rootFileGroup.String(name)
}

// MustString gets the asset value by name
func MustString(name string) (string, error) {
	return rootFileGroup.MustString(name)
}

// Bytes gets the asset value by name
func Bytes(name string) []byte {
	return rootFileGroup.Bytes(name)
}

// MustBytes gets the asset value by name
func MustBytes(name string) ([]byte, error) {
	return rootFileGroup.MustBytes(name)
}

// AddAsset adds the given asset to the root context
func AddAsset(groupName, name, value string) {
	fileGroup := mainAssetDirectory.GetGroup(groupName)
	if fileGroup == nil {
		fileGroup, err = mainAssetDirectory.NewFileGroup(groupName)
		if err != nil {
			zlog.Fatal(err)
		}
	}
	_ = fileGroup.AddAsset(name, value)
}

// AddByteAsset adds the given asset to the root context
func AddByteAsset(groupName, name string, value []byte) {
	fileGroup := mainAssetDirectory.GetGroup(groupName)
	if fileGroup == nil {
		fileGroup, err = mainAssetDirectory.NewFileGroup(groupName)
		if err != nil {
			zlog.Fatal(err)
		}
	}
	_ = fileGroup.AddByteAsset(name, value)
}

// Entries returns the file entries as a slice of filenames
func Entries() []string {
	return rootFileGroup.Entries()
}

// Reset clears the file entries
func Reset() {
	rootFileGroup.Reset()
}

// All All
func All() map[string][]byte {
	return rootFileGroup.All()
}

// Group holds a group of assets
func Group(name string) (result *build.FileGroup, err error) {
	result = mainAssetDirectory.GetGroup(name)
	if result == nil {
		result, err = mainAssetDirectory.NewFileGroup(name)
	}
	return
}

func NewFileserver(dir string, handle ...func(c *znet.Context, name string, content []byte, err error) bool) func(c *znet.Context) {
	f, _ := NewFileserverAndGroup(dir, handle...)
	return f
}

func NewFileserverAndGroup(dir string, handle ...func(c *znet.Context, name string, content []byte, err error) bool) (func(c *znet.Context), *build.FileGroup) {
	const defFile = "index.html"
	f, _ := Group(dir)
	return func(c *znet.Context) {
		name := c.GetParam("file")
		if name == "" {
			name = defFile
		} else {
			name = strings.TrimPrefix(name, "/")
		}
		content, err := f.MustBytes(name)
		mime := zfile.GetMimeType(name, content)
		c.SetContentType(mime)
		if len(handle) > 0 && handle[0] != nil {
			if handle[0](c, name, content, err) {
				return
			}
		}
		if err != nil {
			c.String(404, err.Error())
			return
		}
		c.Byte(http.StatusOK, content)
	}, f
}

func LoadTemplate(pattern string) (t *template.Template, err error) {
	t = template.New("")

	var templateData *build.FileGroup
	templateData, err = Group(pattern)
	if err != nil {
		return nil, err
	}

	all := templateData.All()

	if len(all) == 0 {
		files, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			name := filepath.Base(file)
			bytes, _ := zfile.ReadFile(file)
			t, err = t.New(name).Parse(zstring.Bytes2String(bytes))
			if err != nil {
				return nil, err
			}
		}
	}

	for file := range all {
		name := filepath.Base(file)
		t, err = t.New(name).Parse(templateData.String(file))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
