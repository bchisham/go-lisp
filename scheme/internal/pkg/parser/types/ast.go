package types

type Type string

const (
	Bool               Type = "bool"
	Char               Type = "char"
	Float              Type = "float"
	Int                Type = "int"
	List               Type = "list"
	Pair               Type = "pair"
	Nil                Type = "nil"
	Lambda             Type = "lambda"
	Map                Type = "map"
	String             Type = "string"
	Identifier         Type = "identifier"
	Void               Type = "void"
	Quot               Type = "quot"
	RelationalOperator Type = "relationalOperator"
	ArithmeticOperator Type = "arithmeticOperator"
)
