package pav

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGoLexer(t *testing.T) {
	filepath.Walk(filepath.Join(runtime.GOROOT(), "src"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		runes := []rune(string(content))

		pt("%s\n", path)
		vm := NewVMFromObject(new(GoLexer), Named("Program"))
		for i, r := range runes {
			//pt("%d %d\n", i, len(vm.Threads))
			res := vm.Step(r)
			if i == len(runes)-1 {
				if len(res.Matched) != 1 {
					t.Fatal("not fully parsed")
				}
			}
			if len(vm.Threads) == 0 {
				stop := i + 64
				if stop > len(runes) {
					stop = len(runes)
				}
				pt("%q\n", string(runes[i:stop]))
				t.Fatalf("stop at %d\n", i)
			}
		}

		return nil
	})
}
