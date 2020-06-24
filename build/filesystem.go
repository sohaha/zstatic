package build

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sohaha/zlsgo/zfile"
)

var cwd string

func init() {
	cwd, _ = os.Getwd()
}

// CompressFile reads the given file and converts it to a gzip compressed hex string
func CompressFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	var byteBuffer bytes.Buffer
	writer := gzip.NewWriter(&byteBuffer)
	_, _ = writer.Write(data)
	_ = writer.Close()
	return byteBuffer.Bytes(), nil
}

// FindGoFiles finds all go files recursively from the given directory
func FindGoFiles(directory string) ([]string, error) {
	result := make([]string, 0)
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			basePath := filepath.Base(path)
			if strings.HasPrefix(basePath, ".") {
				if zfile.DirExist(basePath) {
					return filepath.SkipDir
				}
				return nil
			}
			goFilePath := filepath.Ext(path)
			if goFilePath == ".go" {
				isMewnFile := strings.HasSuffix(path, "____tmp.go")
				if !isMewnFile {
					result = append(result, path)
				}
			}
			return nil
		})
	return result, err
}

// DecompressHex decompresses the gzip/hex encoded data
func DecompressHex(hexdata []byte) ([]byte, error) {
	datareader := bytes.NewReader(hexdata)
	gzipReader, err := gzip.NewReader(datareader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	return ioutil.ReadAll(gzipReader)
}

// DecompressHexString decompresses the gzip/hex encoded data
func DecompressHexString(hexdata string) ([]byte, error) {
	data, err := hex.DecodeString(hexdata)
	if err != nil {
		panic(err)
	}
	datareader := bytes.NewReader(data)

	gzipReader, err := gzip.NewReader(datareader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	return ioutil.ReadAll(gzipReader)
}

func HasMewnReference(filename string) (bool, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return false, err
	}
	for _, imprt := range node.Imports {
		if imprt.Path.Value == `"github.com/sohaha/zzz/lib/static"` {
			return true, nil
		}
	}
	return false, nil
}

func GetMewnFiles(args []string, ignoreErrors bool) (mewnFiles []string, err error) {
	var goFiles []string
	if len(args) > 0 {
		for _, inputFile := range args {
			inputFile, err = filepath.Abs(inputFile)
			if err != nil && !ignoreErrors {
				return
			}
			inputFile = filepath.ToSlash(inputFile)
			goFiles = append(goFiles, inputFile)
		}
	} else {
		goFiles, err = FindGoFiles(cwd)
		if err != nil && !ignoreErrors {
			return
		}
	}

	var isReferenced bool
	for _, goFile := range goFiles {
		isReferenced, err = HasMewnReference(goFile)
		if err != nil && !ignoreErrors {
			return
		}
		if isReferenced {
			mewnFiles = append(mewnFiles, goFile)
		}
	}

	return
}
