package ast

import (
	"github.com/Drumstickz64/golox/token"
)

type StmtVisitor interface {
	VisitExpressionStmt(*ExpressionStmt) (any, error)
	VisitPrintStmt(*PrintStmt) (any, error)
	VisitVarStmt(*VarStmt) (any, error)
}

type Stmt interface {
	Accept(StmtVisitor) (any, error)
}

type ExpressionStmt struct {
	Expression Expr
}

func (e *ExpressionStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitExpressionStmt(e)
}

type PrintStmt struct {
	Expression Expr
}

func (p *PrintStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitPrintStmt(p)
}

type VarStmt struct {
	Name        token.Token
	Initializer Expr
}

func (v *VarStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitVarStmt(v)
}
