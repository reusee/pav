package parser

import (
	"testing"
)

func TestJSONEmptyObject(t *testing.T) {
	for _, input := range []string{
		"{}",
		" {}",
		" {}",
		" { }",
		" \n{ \t}",
	} {
		vm := NewJSONParserVM(&Instruction{
			Op:   OpCall,
			Name: "Object",
		})
		eq(t,
			match(vm, input), true,
		)
	}
}

func TestJSONEmptyArray(t *testing.T) {
	for _, input := range []string{
		"[]",
		" []",
		" []",
		" [ ]",
		" \n[ \t]",
	} {
		vm := NewJSONParserVM(&Instruction{
			Op:   OpCall,
			Name: "Array",
		})
		eq(t,
			match(vm, input), true,
		)
	}
}

func TestJSONString(t *testing.T) {
	for _, input := range []string{
		`""`,
		`"ab"`,
		`"a"`,
		`"\n\""`,
		`"\n\"\\"`,
		`"\n\"\\\u1234"`,
	} {
		vm := NewJSONParserVM(&Instruction{
			Op:   OpCall,
			Name: "String",
		})
		eq(t,
			match(vm, input), true,
		)
	}
}

func TestJSONNumber(t *testing.T) {
	for _, input := range []string{
		"0",
		"-0",
		"42",
		"-42",
		"-42.1",
		"-999.99999",
		"-999.0",
		"-999.0",
		"1e10",
		"9e-20",
		"-999.01e10",
		"-999.01e+10",
		"-999.01e-10",
		"-999.01E-10",
		"999.01E-10",
	} {
		vm := NewJSONParserVM(&Instruction{
			Op:   OpCall,
			Name: "Number",
		})
		eq(t,
			match(vm, input), true,
		)
	}
}

func TestJSONValue(t *testing.T) {
	for _, input := range []string{
		`"foo"`,
		`42`,
		`true`,
		`false`,
		`null`,
	} {
		vm := NewJSONParserVM(&Instruction{
			Op:   OpCall,
			Name: "Value",
		})
		eq(t,
			match(vm, input), true,
		)
	}
}

func TestJSONArray(t *testing.T) {
	for _, input := range []string{
		`[42]`,
		`[42, 42]`,
		`[42, 42, 42]`,
		`["foo"]`,
		`["foo", "bar"]`,
		`[true]`,
		`[{}]`,
	} {
		vm := NewJSONParserVM(&Instruction{
			Op:   OpCall,
			Name: "Array",
		})
		eq(t,
			match(vm, input), true,
		)
	}
}

func TestJSONObject(t *testing.T) {
	for _, input := range []string{
		`{"foo": "bar"}`,
		`{"foo": "bar", "bar": []}`,
		`{"foo": "bar", "bar": [], "baz": 42}`,
	} {
		vm := NewJSONParserVM(&Instruction{
			Op:   OpCall,
			Name: "Object",
		})
		eq(t,
			match(vm, input), true,
		)
	}
}
