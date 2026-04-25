package candy_ast

import "candy/candy_token"

type Parameter struct {
	Token    candy_token.Token
	Name     *Identifier
	TypeName Expression
	Default  Expression
}
