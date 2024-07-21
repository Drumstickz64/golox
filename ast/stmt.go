package ast

import (
	"github.com/Drumstickz64/golox/token"
)

type StmtVisitor interface {
	VisitBlockStmt(*BlockStmt) (any, error)
	VisitClassStmt(*ClassStmt) (any, error)
	VisitExpressionStmt(*ExpressionStmt) (any, error)
	VisitWhileStmt(*WhileStmt) (any, error)
	VisitIfStmt(*IfStmt) (any, error)
	VisitPrintStmt(*PrintStmt) (any, error)
	VisitReturnStmt(*ReturnStmt) (any, error)
	VisitVarStmt(*VarStmt) (any, error)
	VisitFunctionStmt(*FunctionStmt) (any, error)
}

type Stmt interface {
	Accept(StmtVisitor) (any, error)
}

type BlockStmt struct {
	Statements []Stmt
}

func (b *BlockStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitBlockStmt(b)
}

type ClassStmt struct {
	Name       token.Token
	SuperClass *VariableExpr
	Methods    []*FunctionStmt
}

func (c *ClassStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitClassStmt(c)
}

type ExpressionStmt struct {
	Expression Expr
}

func (e *ExpressionStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitExpressionStmt(e)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (w *WhileStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitWhileStmt(w)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i *IfStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitIfStmt(i)
}

type PrintStmt struct {
	Expression Expr
}

func (p *PrintStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitPrintStmt(p)
}

type ReturnStmt struct {
	Keyword token.Token
	Value   Expr
}

func (r *ReturnStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitReturnStmt(r)
}

type VarStmt struct {
	Name        token.Token
	Initializer Expr
}

func (v *VarStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitVarStmt(v)
}

type FunctionStmt struct {
	Name       token.Token
	Parameters []token.Token
	Body       []Stmt
}

func (f *FunctionStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitFunctionStmt(f)
}
