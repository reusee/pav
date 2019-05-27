package pav

import (
	"strings"
	"testing"
)

func TestSingleRune(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op:   OpRune,
					Rune: 'a',
				},
			},
		},
	}

	res := vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 1,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
}

func TestMultipleRune(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op:    OpRune,
					Runes: []rune{'a', 'b'},
					Next: &Instruction{
						Op:   OpRune,
						Rune: 'a',
						Next: &Instruction{
							Op:        OpRune,
							RuneRange: [2]rune{'a', 'z'},
							Next: &Instruction{
								Op:   OpRune,
								Rune: 'a',
							},
						},
					},
				},
			},
		},
	}

	res := vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 1,
	)
	res = vm.Step('b')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)

}

func TestMultipleRuneFail(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op:   OpRune,
					Rune: 'a',
					Next: &Instruction{
						Op:   OpRune,
						Rune: 'a',
						Next: &Instruction{
							Op:   OpRune,
							Rune: 'a',
							Next: &Instruction{
								Op:   OpRune,
								Rune: 'a',
							},
						},
					},
				},
			},
		},
	}

	res := vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
	res = vm.Step('b')
	eq(t,
		len(res.Failed), 1,
		len(res.Matched), 0,
	)

}

func TestJump(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op:   OpRune,
					Rune: 'a',
					Next: &Instruction{
						Op: OpJump,
						Inst: &Instruction{
							Op:   OpRune,
							Rune: 'a',
							Next: &Instruction{
								Op: OpJump,
								Inst: &Instruction{
									Op:   OpRune,
									Rune: 'a',
									Next: &Instruction{
										Op: OpJump,
										Inst: &Instruction{
											Op:   OpRune,
											Rune: 'a',
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	res := vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 1,
	)
	res = vm.Step('b')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)

}

func TestJump2(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op: OpJump,
					Inst: &Instruction{
						Op:   OpRune,
						Rune: 'a',
						Next: &Instruction{
							Op: OpJump,
							Inst: &Instruction{
								Op:   OpRune,
								Rune: 'a',
								Next: &Instruction{
									Op: OpJump,
									Inst: &Instruction{
										Op:   OpRune,
										Rune: 'a',
										Next: &Instruction{
											Op: OpJump,
											Inst: &Instruction{
												Op:   OpRune,
												Rune: 'a',
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	res := vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 1,
	)
	res = vm.Step('b')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)

}

func TestCall(t *testing.T) {
	vm := &VM{
		Routines: map[string]Routine{
			"MatchA": Routine{
				Start: &Instruction{
					Op:   OpRune,
					Rune: 'a',
				},
			},
		},
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op:   OpCall,
					Name: "MatchA",
				},
			},
		},
	}

	res := vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 1,
	)

}

func TestCallLoop(t *testing.T) {
	vm := &VM{
		Routines: map[string]Routine{
			"MatchA": Routine{
				Start: &Instruction{
					Op:   OpRune,
					Rune: 'a',
					Next: &Instruction{
						Op:   OpCall,
						Name: "MatchA",
					},
				},
			},
		},
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op:   OpCall,
					Name: "MatchA",
				},
			},
		},
	}
	for i := 0; i < 32; i++ {
		res := vm.Step('a')
		eq(t,
			len(res.Failed), 0,
			len(res.Matched), 0,
		)
	}
	res := vm.Step('b')
	eq(t,
		len(res.Failed), 1,
		len(res.Matched), 0,
	)
}

func TestCall2(t *testing.T) {
	vm := &VM{
		Routines: map[string]Routine{
			"MatchA": Routine{
				Start: &Instruction{
					Op:   OpRune,
					Rune: 'a',
				},
			},
		},
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op:   OpCall,
					Name: "MatchA",
					Next: &Instruction{
						Op:   OpCall,
						Name: "MatchA",
						Next: &Instruction{
							Op:   OpCall,
							Name: "MatchA",
							Next: &Instruction{
								Op:   OpCall,
								Name: "MatchA",
							},
						},
					},
				},
			},
		},
	}

	res := vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 1,
	)
	res = vm.Step('b')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
	)

}

func TestNestedCall(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op: OpCall,
					Inst: &Instruction{
						Op: OpCall,
						Inst: &Instruction{
							Op: OpCall,
							Inst: &Instruction{
								Op: OpCall,
								Inst: &Instruction{
									Op: OpCall,
									Inst: &Instruction{
										Op: OpCall,
										Inst: &Instruction{
											Op: OpCall,
											Inst: &Instruction{
												Op:   OpRune,
												Rune: 'a',
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	res := vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 1,
		len(vm.Threads), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(res.Failed), 0,
		len(res.Matched), 0,
		len(vm.Threads), 0,
	)
}

func TestClone(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op: OpClone,
					Insts: []*Instruction{
						{
							Op:   OpRune,
							Rune: 'a',
						},
						{
							Op:   OpRune,
							Rune: 'b',
						},
					},
				},
			},
		},
	}
	res := vm.Step('a')
	eq(t,
		len(res.Failed), 1,
		len(res.Matched), 1,
	)
}

func TestClone2(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op: OpClone,
					Insts: []*Instruction{
						{
							Op:   OpRune,
							Rune: 'a',
						},
						{
							Op:   OpRune,
							Rune: 'b',
						},
					},
					Next: &Instruction{
						Op: OpClone,
						Insts: []*Instruction{
							{
								Op:   OpRune,
								Rune: 'c',
							},
							{
								Op:   OpRune,
								Rune: 'd',
							},
						},
					},
				},
			},
		},
	}

	res := vm.Step('b')
	eq(t,
		len(vm.Threads), 2, // c | d
		len(res.Matched), 0,
		len(res.Failed), 1, // a
	)
	res = vm.Step('d')
	eq(t,
		len(vm.Threads), 0,
		len(res.Matched), 1, // d
		len(res.Failed), 1, // c
	)
}

func TestReturn(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op: OpClone,
					Insts: []*Instruction{
						{
							Op:   OpRune,
							Rune: 'a',
							Next: &Instruction{
								Op: OpReturn,
							},
						},
						{
							Op:   OpRune,
							Rune: 'b',
						},
					},
				},
			},
		},
	}
	res := vm.Step('a')
	eq(t,
		len(res.Failed), 1,
		len(res.Matched), 1,
	)
}

func TestReturn2(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op: OpReturn,
				},
			},
		},
	}
	res := vm.Step('a')
	eq(t,
		len(res.Failed), 1,
		len(res.Matched), 0,
		len(vm.Threads), 0,
	)
}

func TestShortest(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op: OpClone,
					Insts: []*Instruction{
						{
							Op:   OpRune,
							Rune: 'a',
						},
						{
							Op:   OpRune,
							Rune: 'a',
							Next: &Instruction{
								Op:   OpRune,
								Rune: 'a',
							},
						},
					},
					ClusterID:   1,
					ClusterType: ClusterShortest,
				},
			},
		},
	}

	res := vm.Step('a')
	eq(t,
		len(vm.Threads), 0,
		len(res.Matched), 1,
		len(res.Failed), 1,
	)
}

func TestShortest2(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op:   OpRune,
					Rune: 'a',
					Next: &Instruction{
						Op: OpClone,
						Insts: []*Instruction{
							{
								Op:   OpRune,
								Rune: 'a',
							},
							{
								Op:   OpRune,
								Rune: 'a',
								Next: &Instruction{
									Op:   OpRune,
									Rune: 'a',
								},
							},
						},
						ClusterID:   1,
						ClusterType: ClusterShortest,
						Next: &Instruction{
							Op:   OpRune,
							Rune: 'b',
						},
					},
				},
			},
		},
	}

	res := vm.Step('a')
	eq(t,
		len(vm.Threads), 2, // a | aa
		len(res.Matched), 0,
		len(res.Failed), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(vm.Threads), 1,
		len(res.Matched), 0,
		len(res.Failed), 1, // aa failed
	)
	res = vm.Step('b')
	eq(t,
		len(res.Matched), 1,
	)

}

func TestLongest(t *testing.T) {
	vm := &VM{
		Threads: []*Thread{
			{
				PC: &Instruction{
					Op:   OpRune,
					Rune: 'a',
					Next: &Instruction{
						Op: OpClone,
						Insts: []*Instruction{
							{
								Op:   OpRune,
								Rune: 'a',
							},
							{
								Op:   OpRune,
								Rune: 'a',
								Next: &Instruction{
									Op:   OpRune,
									Rune: 'a',
								},
							},
						},
						Next: &Instruction{
							Op:   OpRune,
							Rune: 'b',
						},
					},
				},
			},
		},
	}

	res := vm.Step('a')
	eq(t,
		len(vm.Threads), 2, // a | aa
		len(res.Matched), 0,
		len(res.Failed), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(vm.Threads), 2,
		len(res.Matched), 0,
		len(res.Failed), 0,
	)
	res = vm.Step('a')
	eq(t,
		len(vm.Threads), 1,
		len(res.Matched), 0,
		len(res.Failed), 1, // aab failed
	)
	res = vm.Step('b')
	eq(t,
		len(vm.Threads), 0,
		len(res.Matched), 1,
		len(res.Failed), 0,
	)

}

func TestLeftRecursion(t *testing.T) {
	var A *Instruction
	A = Longest(
		Rune('a'),
		Seq(
			Indirect(&A),
			Rune('b'),
		),
	)
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					A,
					Literal("AAA"),
				),
			},
		},
	}

	runes := []rune("abbbAAA")
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

func TestLeftRecursion2(t *testing.T) {
	//TODO

	var A *Instruction
	A = Longest(
		Rune('a'),
		Seq(
			Indirect(&A),
			Rune('b'),
		),
	)
	vm := &VM{
		Threads: []*Thread{
			{
				PC: Seq(
					A,
					Literal("AAA"),
				),
			},
		},
	}

	runes := []rune("a" + strings.Repeat("b", 128) + "AAA")
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
