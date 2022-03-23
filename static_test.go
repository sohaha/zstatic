package zstatic

import (
	"testing"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/zfile"
)

func TestStatic(t *testing.T) {
	tt := zlsgo.NewTest(t)

	tt.EqualExit(true, zfile.FileExist("README.md"))

	g, _ := Group("./")
	s, err := g.MustString("README.md")
	tt.EqualNil(err)
	t.Log(len(s))

}
