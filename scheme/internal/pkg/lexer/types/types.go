package types

const (
	KwAnd    = "and"
	KwOr     = "or"
	KwXor    = "xor"
	KwNot    = "not"
	KwCons   = "cons"
	KwCdr    = "cdr"
	KwEq     = "="
	KwLt     = "<"
	KwGt     = ">"
	KwGtEq   = ">="
	KwLtEq   = "<="
	KwLambda = "lambda"
	KwIf     = "if"
)
const (
	LiteralTrue  = "#t"
	LiteralFalse = "#f"
)

type RuneClassifier func(rune) bool
