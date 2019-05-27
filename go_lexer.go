package pav

type GoLexer struct{}

func (_ GoLexer) Program() *Instruction {
	return ZeroOrMore(
		Longest(
			Named("Comment"),
			Named("Token"),
			Named("Blank"),
		),
	)
}

func (_ GoLexer) Comment() *Instruction {
	return Longest(
		Named("LineComment"),
		Named("GeneralComment"),
	)
}

func (_ GoLexer) Token() *Instruction {
	return First(
		Named("Keyword"),
		Named("Identifier"),
		Named("OperatorAndPunctuation"),
		Named("Literal"),
	)
}

func (_ GoLexer) Blank() *Instruction {
	return RuneSet(
		' ',
		'\b',
		'\t',
		'\r',
		'\n',
	)
}

//TODO semicolon

func (_ GoLexer) Identifier() *Instruction {
	return Seq(
		Named("Letter"),
		ZeroOrMore(
			Longest(
				Named("Letter"),
				RuneCategory("Nd"),
			),
		),
	)
}

func (_ GoLexer) Letter() *Instruction {
	return Longest(
		RuneCategory("L"),
		Rune('_'),
	)
}

func (_ GoLexer) Keyword() *Instruction {
	return Longest(
		Literal("break"),
		Literal("default"),
		Literal("func"),
		Literal("interface"),
		Literal("select"),
		Literal("case"),
		Literal("defer"),
		Literal("go"),
		Literal("map"),
		Literal("struct"),
		Literal("chan"),
		Literal("else"),
		Literal("goto"),
		Literal("package"),
		Literal("switch"),
		Literal("const"),
		Literal("fallthrough"),
		Literal("if"),
		Literal("range"),
		Literal("type"),
		Literal("continue"),
		Literal("for"),
		Literal("import"),
		Literal("return"),
		Literal("var"),
	)
}

func (_ GoLexer) OperatorAndPunctuation() *Instruction {
	return Longest(
		Literal("+"),
		Literal("&"),
		Literal("+="),
		Literal("&="),
		Literal("&&"),
		Literal("=="),
		Literal("!="),
		Literal("("),
		Literal(")"),
		Literal("-"),
		Literal("|"),
		Literal("-="),
		Literal("|="),
		Literal("||"),
		Literal("<"),
		Literal("<="),
		Literal("["),
		Literal("]"),
		Literal("*"),
		Literal("^"),
		Literal("*="),
		Literal("^="),
		Literal("<-"),
		Literal(">"),
		Literal(">="),
		Literal("{"),
		Literal("}"),
		Literal("/"),
		Literal("<<"),
		Literal("/="),
		Literal("<<="),
		Literal("++"),
		Literal("="),
		Literal(":="),
		Literal(","),
		Literal(";"),
		Literal("%"),
		Literal(">>"),
		Literal("%="),
		Literal(">>="),
		Literal("--"),
		Literal("!"),
		Literal("..."),
		Literal("."),
		Literal(":"),
		Literal("&^"),
		Literal("&^="),
	)
}

func (_ GoLexer) Literal() *Instruction {
	return First(
		Named("IntegerLiteral"),
		Named("FloatLiteral"),
		Named("ImaginaryLiteral"),
		Named("RuneLiteral"),
		Named("StringLiteral"),
	)
}

func (_ GoLexer) IntegerLiteral() *Instruction {
	return First(
		Named("DecimalLiteral"),
		Named("OctalLiteral"),
		Named("HexLiteral"),
	)
}

func (_ GoLexer) DecimalLiteral() *Instruction {
	return Seq(
		RuneRange('1', '9'),
		ZeroOrMore(
			RuneRange('0', '9'),
		),
	)
}

func (_ GoLexer) FloatLiteral() *Instruction {
	return Longest(
		Seq(
			Named("Decimals"),
			Literal("."),
			Optional(Named("Decimals")),
			Optional(Named("Exponent")),
		),
		Seq(
			Named("Decimals"),
			Named("Exponent"),
		),
		Seq(
			Literal("."),
			Named("Decimals"),
			Optional(Named("Exponent")),
		),
	)
}

func (_ GoLexer) Decimals() *Instruction {
	return Seq(
		RuneRange('0', '9'),
		ZeroOrMore(
			RuneRange('0', '9'),
		),
	)
}

func (_ GoLexer) ImaginaryLiteral() *Instruction {
	return Seq(
		Longest(
			Named("Decimals"),
			Named("FloatLiteral"),
		),
		Literal("i"),
	)
}

func (_ GoLexer) RuneLiteral() *Instruction {
	return Seq(
		Literal(`'`),
		First(
			Named("UnicodeValue"),
			Named("ByteValue"),
		),
		Literal(`'`),
	)
}

func (_ GoLexer) StringLiteral() *Instruction {
	return First(
		Named("RawStringLiteral"),
		Named("InterpretedStringLiteral"),
	)
}

func (_ GoLexer) RawStringLiteral() *Instruction {
	return Seq(
		Literal("`"),
		ZeroOrMore(
			First(
				Named("UnicodeChar"),
				Rune('\n'),
			),
		),
		Literal("`"),
	)
}

func (_ GoLexer) OctalLiteral() *Instruction {
	return Seq(
		Rune('0'),
		ZeroOrMore(
			RuneRange('0', '7'),
		),
	)
}

func (_ GoLexer) HexLiteral() *Instruction {
	return Seq(
		Rune('0'),
		First(
			Rune('x'),
			Rune('X'),
		),
		Named("HexDigit"),
		ZeroOrMore(
			Named("HexDigit"),
		),
	)
}

func (_ GoLexer) InterpretedStringLiteral() *Instruction {
	return Seq(
		Rune('"'),
		ZeroOrMore(
			Longest(
				Named("UnicodeValue"),
				Named("ByteValue"),
			),
		),
		Rune('"'),
	)
}

func (_ GoLexer) LineComment() *Instruction {
	return Seq(
		Literal("//"),
		ZeroOrMore(
			RuneInverse(
				Rune('\n'),
			),
		),
	)
}

func (_ GoLexer) GeneralComment() *Instruction {
	return Seq(
		Literal("/*"),
		ZeroOrMore(AnyRune()),
		Literal("*/"),
	)
}

func (_ GoLexer) UnicodeValue() *Instruction {
	return Longest(
		RuneInverse(Rune('\n')),
		Seq(
			Literal(`\u`),
			Named("HexDigit"),
			Named("HexDigit"),
			Named("HexDigit"),
			Named("HexDigit"),
		),
		Seq(
			Literal(`\U`),
			Named("HexDigit"),
			Named("HexDigit"),
			Named("HexDigit"),
			Named("HexDigit"),
		),
		First(
			Literal(`\a`),
			Literal(`\b`),
			Literal(`\f`),
			Literal(`\n`),
			Literal(`\r`),
			Literal(`\t`),
			Literal(`\v`),
			Literal(`\\`),
			Literal(`\'`),
			Literal(`\"`),
		),
	)
}

func (_ GoLexer) ByteValue() *Instruction {
	return First(
		Seq(
			Rune('\\'),
			RuneRange('0', '7'),
			RuneRange('0', '7'),
			RuneRange('0', '7'),
		),
		Seq(
			Literal(`\x`),
			Named("HexDigit"),
			Named("HexDigit"),
		),
	)
}

func (_ GoLexer) HexDigit() *Instruction {
	return First(
		RuneRange('0', '9'),
		RuneRange('a', 'f'),
		RuneRange('A', 'F'),
	)
}
