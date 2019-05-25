package parser

import "testing"

func TestRuneSeq(t *testing.T) {
	runes := []rune("abcdefg")
	vm := &VM{
		Threads: []*Thread{
			{
				PC: RuneSeq(runes),
			},
		},
	}
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 0,
			)
		} else {
			eq(t,
				len(res.Matched), 0,
				len(res.Failed), 0,
			)
		}
	}
}

func TestLiteral(t *testing.T) {
	str := "abcdefg"
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Literal(str),
			},
		},
	}
	runes := []rune(str)
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 0,
			)
		} else {
			eq(t,
				len(res.Matched), 0,
				len(res.Failed), 0,
			)
		}
	}
}

func TestSeq(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("abc"),
					Literal("def"),
					Literal("ghi"),
				),
			},
		},
	}
	runes := []rune("abcdefghi")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 0,
			)
		} else {
			eq(t,
				len(res.Matched), 0,
				len(res.Failed), 0,
			)
		}
	}
}

func TestSeq2(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Seq(
						Literal("a"),
					),
					Seq(
						Literal("b"),
					),
					Seq(
						Literal("c"),
					),
				),
			},
		},
	}
	runes := []rune("abc")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 0,
			)
		} else {
			eq(t,
				len(res.Matched), 0,
				len(res.Failed), 0,
			)
		}
	}
}

func TestSeq3(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("a"),
					Seq(
						Literal("b"),
					),
					Seq(
						Literal("c"),
					),
				),
			},
		},
	}
	runes := []rune("abc")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 0,
			)
		} else {
			eq(t,
				len(res.Matched), 0,
				len(res.Failed), 0,
			)
		}
	}
}

func TestShortestComb(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("ab"),
					Shortest(
						Literal("aaa"),
						Literal("a"),
					),
					Literal("c"),
				),
			},
		},
	}
	runes := []rune("abac")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 0,
			)
		} else {
			eq(t,
				len(res.Matched), 0,
			)
		}
	}
}

func TestLongestComb(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("ab"),
					Longest(
						Literal("aaa"),
						Literal("a"),
					),
					Literal("c"),
				),
			},
		},
	}
	runes := []rune("abaaac")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 0,
			)
		} else {
			eq(t,
				len(res.Matched), 0,
			)
		}
	}
}

func TestOptional(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("ab"),
					Optional(
						Literal("x"),
					),
					Literal("c"),
				),
			},
		},
	}
	runes := []rune("abc")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1, // abc
				len(res.Failed), 1, // abxc
			)
		} else {
			eq(t,
				len(res.Matched), 0,
			)
		}
	}
}

func TestOptional2(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("ab"),
					Optional(
						Literal("x"),
					),
					Literal("c"),
				),
			},
		},
	}
	runes := []rune("abxc")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1, // abc
				len(res.Failed), 0, // abxc
			)
		} else {
			eq(t,
				len(res.Matched), 0,
			)
		}
	}
}

func TestOptional3(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("ab"),
					Optional(
						Optional(
							Literal("x"),
						),
					),
					Literal("c"),
				),
			},
		},
	}
	runes := []rune("abxc")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1, // abc
				len(res.Failed), 0, // abxc
			)
		} else {
			eq(t,
				len(res.Matched), 0,
			)
		}
	}
}

func TestZeroOrMore(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("a"),
					ZeroOrMore(Literal("a")),
					Literal("b"),
				),
			},
		},
	}
	runes := []rune("ab")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 1, // aab
			)
		} else {
			eq(t,
				len(res.Matched), 0,
			)
		}
	}
}

func TestZeroOrMore2(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("a"),
					ZeroOrMore(Literal("a")),
					Literal("b"),
				),
			},
		},
	}
	runes := []rune("aab")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 1, // aa*
			)
		} else {
			eq(t,
				len(res.Matched), 0,
			)
		}
	}
}

func TestZeroOrMore3(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("a"),
					ZeroOrMore(Literal("a")),
					Literal("b"),
				),
			},
		},
	}
	runes := []rune("aaaaaaaaaaaaaab")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 1, // aa*
			)
		} else {
			eq(t,
				len(res.Matched), 0,
			)
		}
	}
}

func TestZeroOrMore4(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("a"),
					ZeroOrMore(
						Literal("xy"),
					),
					Literal("b"),
				),
			},
		},
	}
	runes := []rune("axyxyxyxyxyxyb")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 1, // aa*
			)
		} else {
			eq(t,
				len(res.Matched), 0,
			)
		}
	}
}

func TestOneOrMore(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("a"),
					OneOrMore(
						Literal("xy"),
					),
					Literal("b"),
				),
			},
		},
	}
	runes := []rune("axyxyxyxyxyxyb")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 1, // aa*
			)
		} else {
			eq(t,
				len(res.Matched), 0,
			)
		}
	}
}

func TestOneOrMore2(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("a"),
					OneOrMore(
						Literal("xy"),
					),
					Literal("b"),
				),
			},
		},
	}
	runes := []rune("axyb")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 1, // aa*
			)
		} else {
			eq(t,
				len(res.Matched), 0,
			)
		}
	}
}

func TestOneOrMore3(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					Literal("a"),
					OneOrMore(
						Literal("xy"),
					),
					Literal("b"),
				),
			},
		},
	}
	runes := []rune("axyxyb")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
				len(res.Failed), 1, // aa*
			)
		} else {
			eq(t,
				len(res.Matched), 0,
			)
		}
	}
}

func TestFirst(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					First(
						Literal("a"),
					),
					First(
						Literal("c"),
						Literal("d"),
					),
				),
			},
		},
	}
	runes := []rune("ac")
	for i, r := range runes {
		res := vm.Step(r)
		if i == len(runes)-1 { // last
			eq(t,
				len(res.Matched), 1,
			)
		} else {
			eq(t,
				len(res.Matched), 0,
			)
		}
	}
}
