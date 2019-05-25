package parser

import "reflect"

type JSONParser struct{}

func NewJSONParserVM(initInst *Instruction) *VM {
	o := new(JSONParser)
	v := reflect.ValueOf(o)
	t := reflect.TypeOf(o)
	vm := &VM{
		Routines: make(map[string]Routine),
	}
	for i := 0; i < v.NumMethod(); i++ {
		fn := v.Method(i).Interface()
		if fn, ok := fn.(func() *Instruction); ok {
			name := t.Method(i).Name
			vm.Routines[name] = Routine{
				Start: fn(),
			}
		}
	}
	vm.Threads = []*Thread{
		{
			PC: initInst,
		},
	}
	return vm
}

func (_ JSONParser) Lexical(str string) *Instruction {
	return Seq(
		Named("Blank"),
		Literal(str),
	)
}

func (_ JSONParser) Blank() *Instruction {
	return ZeroOrMore(
		RuneSet(' ', '\b', '\f', '\n', '\r', '\t'),
	)
}

func (j JSONParser) Object() *Instruction {
	return Seq(
		j.Lexical("{"),
		Optional(
			Seq(
				Named("Blank"),
				Named("String"),
				j.Lexical(":"),
				Named("Value"),
				ZeroOrMore(
					Seq(
						j.Lexical(","),
						Named("String"),
						j.Lexical(":"),
						Named("Value"),
					),
				),
			),
		),
		j.Lexical("}"),
	)
}

func (j JSONParser) Array() *Instruction {
	return Seq(
		j.Lexical("["),
		Optional(
			Seq(
				Named("Value"),
				ZeroOrMore(
					Seq(
						j.Lexical(","),
						Named("Value"),
					),
				),
			),
		),
		j.Lexical("]"),
	)
}

func (j JSONParser) Value() *Instruction {
	return Seq(
		Named("Blank"),
		Longest(
			Named("String"),
			Named("Number"),
			Named("Object"),
			Named("Array"),
			j.Lexical("true"),
			j.Lexical("false"),
			j.Lexical("null"),
		),
	)
}

func (j JSONParser) String() *Instruction {
	return Seq(
		j.Lexical(`"`),
		ZeroOrMore(
			Longest(
				Inverse(
					RuneSet('"', '\\'),
				),
				Literal(`\"`),
				Literal(`\\`),
				Literal(`\/`),
				Literal(`\b`),
				Literal(`\f`),
				Literal(`\n`),
				Literal(`\r`),
				Literal(`\t`),
				Seq(
					Literal(`\u`),
					RuneRange('0', '9'),
					RuneRange('0', '9'),
					RuneRange('0', '9'),
					RuneRange('0', '9'),
				),
			),
		),
		Literal(`"`),
	)
}

func (j JSONParser) Number() *Instruction {
	return Seq(
		Optional(
			j.Lexical("-"),
		),
		Longest(
			j.Lexical("0"),
			Seq(
				RuneRange('1', '9'),
				ZeroOrMore(
					RuneRange('0', '9'),
				),
			),
		),
		Optional(
			Seq(
				j.Lexical("."),
				OneOrMore(
					RuneRange('0', '9'),
				),
			),
		),
		Optional(
			Seq(
				Longest(
					j.Lexical("e"),
					j.Lexical("E"),
				),
				Optional(
					Longest(
						j.Lexical("+"),
						j.Lexical("-"),
					),
				),
				OneOrMore(
					RuneRange('0', '9'),
				),
			),
		),
	)
}
