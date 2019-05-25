package parser

import (
	"fmt"
)

type VM struct {
	Routines map[string]Routine
	Threads  []*Thread
	step     int
}

type Thread struct {
	Stack     []Frame
	PC        *Instruction
	Match     bool
	instStats []instStat
}

type instStat struct {
	Inst    *Instruction
	Counter int
	Bound   int
}

type Frame struct {
	Return      *Instruction
	ClusterID   int64
	ClusterType ClusterType
}

type Routine struct {
	Start *Instruction
}

type Instruction struct {
	Next *Instruction
	Op   Op

	// OpCall
	Name        string
	ClusterID   int64
	ClusterType ClusterType

	// OpJump, OpCall
	Inst *Instruction

	// OpClone
	Insts []*Instruction

	// OpRune
	Rune      rune
	Runes     []rune
	RuneRange [2]rune
	Inverse   bool

	// OpIndirect
	InstP **Instruction
}

type Op uint8

const (
	OpRune Op = iota + 1
	OpCall
	OpJump
	OpClone
	OpReturn
	OpIndirect
)

type ClusterType uint8

const (
	ClusterShortest ClusterType = iota + 1
	ClusterBound
)

type StepResult struct {
	Matched []*Thread
	Failed  []*Thread
}

func (v *VM) PrepareToFeed(thread *Thread) {
	for {

		// implicit return
		if thread.PC == nil {
			thread.PC = &Instruction{
				Op: OpReturn,
			}
		}

		// ready to feed
		if thread.PC.Op == OpRune {
			return
		}

		//TODO restart with larger bound
		// update counter
		added := false
		for i, c := range thread.instStats {
			if c.Inst == thread.PC {
				if c.Counter >= c.Bound {
					v.kill(thread)
					return
				}
				thread.instStats[i].Counter++
				added = true
				break
			}
		}
		if !added {
			thread.instStats = append(thread.instStats, instStat{
				Inst:    thread.PC,
				Counter: 0,
				Bound:   64,
			})
		}

		switch thread.PC.Op {

		case OpCall:
			thread.Stack = append(thread.Stack, Frame{
				Return:      thread.PC.Next,
				ClusterID:   thread.PC.ClusterID,
				ClusterType: thread.PC.ClusterType,
			})
			if thread.PC.Inst != nil {
				thread.PC = thread.PC.Inst
			} else if thread.PC.Name != "" {
				thread.PC = v.Routines[thread.PC.Name].Start
			} else { // NOCOVER
				panic(fmt.Errorf("bad instruction: %+v", thread.PC))
			}

		case OpJump:
			thread.PC = thread.PC.Inst

		case OpClone:
			inst := thread.PC
			for i, start := range thread.PC.Insts {
				var t *Thread
				if i == 0 {
					// use current thread
					t = thread
				} else {
					// create new thread
					stack := make([]Frame, len(thread.Stack))
					copy(stack, thread.Stack)
					counters := make([]instStat, len(thread.instStats))
					copy(counters, thread.instStats)
					t = &Thread{
						Stack:     stack,
						instStats: counters,
					}
					v.Threads = append(v.Threads, t)
				}
				// set pc
				if start == nil {
					t.PC = inst.Next
				} else {
					t.PC = &Instruction{
						Op:          OpCall,
						Inst:        start,
						ClusterID:   inst.ClusterID,
						ClusterType: inst.ClusterType,
						Next:        inst.Next,
					}
				}
			}

		case OpReturn:
			if len(thread.Stack) > 0 {
				v.UnwindStack(thread)
			} else {
				thread.PC = nil
				return // no more frames
			}

		case OpIndirect:
			thread.PC = *thread.PC.InstP

		}

	}

}

func (v *VM) UnwindStack(thread *Thread) {
	frame := thread.Stack[len(thread.Stack)-1]
	thread.PC = frame.Return
	thread.Stack = thread.Stack[:len(thread.Stack)-1]

	// clustered frames
	if frame.ClusterID > 0 {
		switch frame.ClusterType {

		case ClusterShortest:
			// unwind threads containing frames in the same cluster
			if thread.Match {
			loop_thread:
				for _, t := range v.Threads {
					if t == thread {
						continue
					}
					for _, f := range t.Stack {
						if f.ClusterID == frame.ClusterID {
							v.kill(t)
							continue loop_thread
						}
					}
				}
			}

		}
	}

}

func (v *VM) Step(input rune) (
	result StepResult,
) {

	for i := 0; i < len(v.Threads); i++ {
		v.PrepareToFeed(v.Threads[i])
	}

	for _, thread := range v.Threads {
		// feed rune
		if thread.PC != nil {
			if thread.PC.Op != OpRune { // NOCOVER
				panic("bad code path")
			}

			inst := thread.PC

			if len(inst.Runes) > 0 {
				// runes
				thread.Match = false
				for _, r := range inst.Runes {
					if input == r {
						thread.Match = true
						break
					}
				}
			} else if inst.RuneRange[0] != inst.RuneRange[1] {
				// rune range
				thread.Match = input >= inst.RuneRange[0] &&
					input <= inst.RuneRange[1]
			} else {
				// single rune
				thread.Match = input == inst.Rune
			}
			if inst.Inverse {
				thread.Match = !thread.Match
			}

			if thread.Match {
				thread.PC = thread.PC.Next
			} else {
				v.kill(thread)
			}

		}

		thread.instStats = thread.instStats[:0]
	}

	v.step++

	for i := 0; i < len(v.Threads); i++ {
		v.PrepareToFeed(v.Threads[i])
	}

	// purge stopped threads
	i := 0
	for i < len(v.Threads) {
		thread := v.Threads[i]
		if thread.PC == nil {
			if thread.Match {
				result.Matched = append(result.Matched, v.Threads[i])
			} else {
				result.Failed = append(result.Failed, v.Threads[i])
			}
			v.Threads = append(v.Threads[:i], v.Threads[i+1:]...)
			continue
		}
		i++
	}

	return
}

func (v *VM) kill(t *Thread) {
	for len(t.Stack) > 0 {
		v.UnwindStack(t)
	}
	t.PC = nil
	t.Match = false
}

func (v *VM) dumpThreads() { // NOCOVER
	pt("---- %d threads ----\n", len(v.Threads))
	for _, thread := range v.Threads { // NOCOVER
		pt("%+v\n", thread.PC)
	}
	pt("---- ----\n") // NOCOVER
}
