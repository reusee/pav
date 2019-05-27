package pav

import (
	"bytes"
	"encoding/json"
	"testing"
)

func eq(t *testing.T, args ...interface{}) {
	t.Helper()
	if len(args)%2 != 0 {
		t.Fatal("must be even number of args")
	}
	type Result struct {
		J1    []byte
		J2    []byte
		Equal bool
	}
	var results []Result
	for i := 0; i < len(args); i += 2 {
		j1, err := json.Marshal(args[i])
		if err != nil {
			t.Fatal(err)
		}
		j2, err := json.Marshal(args[i+1])
		if err != nil {
			t.Fatal(err)
		}
		results = append(results, Result{
			J1:    j1,
			J2:    j2,
			Equal: bytes.Equal(j1, j2),
		})
	}
	fatal := false
	for i, res := range results {
		if !res.Equal {
			fatal = true
			pt(
				"pair %d not equal:\n%s\n------\n%s\n",
				i+1,
				res.J1,
				res.J2,
			)
		}
	}
	if fatal {
		t.Fatal()
	}
}

func match(vm *VM, input string) bool {
	var res StepResult
	for _, r := range []rune(input) {
		res = vm.Step(r)
	}
	return len(res.Matched) >= 1
}
