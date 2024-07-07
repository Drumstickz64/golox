package expr

import "github.com/Drumstickz64/golox/token"

type ExprVisitor interface {
	VisitTernary(exp *Ternary) any
	VisitBinary(exp *Binary) any
	VisitGrouping(exp *Grouping) any
	VisitLiteral(exp *Literal) any
	VisitUnary(exp *Unary) any
}

type Expr interface {
	Accept(ExprVisitor) any
}

type Ternary struct {
	Left   Expr
	Middle Expr
	Right  Expr
}

func (t *Ternary) Accept(visitor ExprVisitor) any {
	return visitor.VisitTernary(t)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b *Binary) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) any {
	return visitor.VisitGrouping(g)
}

type Literal struct {
	Value any
}

func (l *Literal) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteral(l)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u *Unary) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnary(u)
}
