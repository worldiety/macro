package ast

import (
	"go/token"
	"strconv"
)

// TokenKind determines which kind of token, like INT, FLOAT, IMAG, CHAR, or STRING is meant.
type TokenKind int

const (
	TokenInt    TokenKind = TokenKind(token.INT)
	TokenFloat  TokenKind = TokenKind(token.FLOAT)
	TokenImag   TokenKind = TokenKind(token.IMAG)
	TokenChar   TokenKind = TokenKind(token.CHAR)
	TokenString TokenKind = TokenKind(token.STRING)
	TokenIdent  TokenKind = TokenKind(token.IDENT)
)

// A BasicLit represents a literal of a basic type.
type BasicLit struct {
	Kind TokenKind
	Val  string // the actual literal string, strings and chars must include the according escapes.
	Obj
}

func NewStrLit(v string) *BasicLit {
	return NewBasicLit(TokenString, strconv.Quote(v))
}

func NewIntLit(i int) *BasicLit {
	return NewBasicLit(TokenString, strconv.Itoa(i))
}

func NewInt64Lit(i int64) *BasicLit {
	return NewBasicLit(TokenString, strconv.FormatInt(i, 10))
}

func NewBoolLit(b bool) *BasicLit {
	return NewBasicLit(TokenIdent, strconv.FormatBool(b))
}

func NewIdentLit(ident string) *BasicLit {
	return NewBasicLit(TokenIdent, ident)
}

func NewBasicLit(kind TokenKind, value string) *BasicLit {
	return &BasicLit{
		Kind: kind,
		Val:  value,
	}
}

func (n *BasicLit) exprNode() {

}
