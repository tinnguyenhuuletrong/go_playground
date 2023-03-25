package grammar_play

import (
	"log"

	"github.com/alecthomas/repr"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

func Play_INI_Parser() {
	var (
		iniLexer = lexer.MustSimple([]lexer.SimpleRule{
			{`Ident`, `[a-zA-Z][a-zA-Z_\d]*`},
			{`String`, `"(?:\\.|[^"])*"`},
			{`Float`, `\d+(?:\.\d+)?`},
			{`Punct`, `[][=]`},
			{"comment", `[#;][^\n]*`},
			{"whitespace", `\s+`},
		})
		parser = participle.MustBuild[INI](
			participle.Lexer(iniLexer),
			participle.Unquote("String"),
			participle.Union[Value](String{}, Number{}, Boolean{}),
		)
	)

	inp := `
	a = "a"
	b = 123
	c = 3.14
	t = true
	f = false

	[server]
		addr = "127.0.01"
		port = 3001
	`

	ini, err := parser.ParseString("", inp)
	if err != nil {
		log.Fatalln(err)
		return
	}
	repr.Println(ini)
}

type INI struct {
	Properties []*Property `@@*`
	Sections   []*Section  `@@*`
}

type Section struct {
	Identifier string      `"[" @Ident "]"`
	Properties []*Property `@@*`
}

type Property struct {
	Key   string `@Ident "="`
	Value Value  `@@`
}

type Value interface{ value() }

type String struct {
	String string `@String`
}

func (String) value() {}

type Number struct {
	Number float64 `@Float`
}

func (Number) value() {}

type Boolean struct {
	Boolean bool `@("true" | "false")`
}

func (Boolean) value() {}
