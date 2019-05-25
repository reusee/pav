package parser

import (
	"sync/atomic"
	"unicode"
)

var nextClusterID int64

func Literal(s string) *Instruction {
	return RuneSeq([]rune(s))
}

func RuneSeq(runes []rune) *Instruction {
	if len(runes) == 0 {
		return nil
	}
	return &Instruction{
		Op:   OpRune,
		Rune: runes[0],
		Next: RuneSeq(runes[1:]),
	}
}

func RuneSet(runes ...rune) *Instruction {
	return &Instruction{
		Op:    OpRune,
		Runes: runes,
	}
}

func RuneRange(r1, r2 rune) *Instruction {
	return &Instruction{
		Op:        OpRune,
		RuneRange: [2]rune{r1, r2},
	}
}

func Rune(r rune) *Instruction {
	return &Instruction{
		Op:   OpRune,
		Rune: r,
	}
}

func AnyRune() *Instruction {
	return &Instruction{
		Op:        OpRune,
		RuneRange: [2]rune{0, unicode.MaxRune},
	}
}

func Inverse(inst *Instruction) *Instruction {
	i := *inst
	i.Inverse = true
	return &i
}

func Seq(instructions ...*Instruction) *Instruction {
	if len(instructions) == 0 {
		return nil
	}
	return &Instruction{
		Op:   OpCall,
		Inst: instructions[0],
		Next: Seq(instructions[1:]...),
	}
}

func Shortest(instructions ...*Instruction) *Instruction {
	if len(instructions) == 0 { // NOCOVER
		return nil
	}
	return &Instruction{
		Op:          OpClone,
		Insts:       instructions,
		ClusterID:   atomic.AddInt64(&nextClusterID, 1),
		ClusterType: ClusterShortest,
	}
}

func First(instructions ...*Instruction) *Instruction {
	return Shortest(instructions...)
}

func Longest(instructions ...*Instruction) *Instruction {
	if len(instructions) == 0 { // NOCOVER
		return nil
	}
	return &Instruction{
		Op:    OpClone,
		Insts: instructions,
	}
}

func Optional(inst *Instruction) *Instruction {
	return &Instruction{
		Op: OpClone,
		Insts: []*Instruction{
			// zero
			nil,
			// one
			{
				Op:   OpCall,
				Inst: inst,
			},
		},
	}
}

func ZeroOrMore(inst *Instruction) *Instruction {
	var cloneInst *Instruction
	cloneInst = &Instruction{
		Op: OpClone,
		Insts: []*Instruction{
			// zero
			nil,
			// more
			&Instruction{
				Op:   OpCall,
				Inst: inst,
				Next: &Instruction{
					Op:    OpIndirect,
					InstP: &cloneInst,
				},
			},
		},
	}
	return cloneInst
}

func OneOrMore(inst *Instruction) *Instruction {
	return Seq(
		// one
		inst,
		// more
		ZeroOrMore(inst),
	)
}

func Indirect(p **Instruction) *Instruction {
	return &Instruction{
		Op:    OpIndirect,
		InstP: p,
	}
}

func Named(name string) *Instruction {
	return &Instruction{
		Op:   OpCall,
		Name: name,
	}
}
