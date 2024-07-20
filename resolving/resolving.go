package resolving

import (
	"fmt"
	"os"

	"github.com/Drumstickz64/golox/ast"
	"github.com/Drumstickz64/golox/errors"
	"github.com/Drumstickz64/golox/interpreting"
	"github.com/Drumstickz64/golox/token"
)

type functionType int

const (
	FUNCTION_TYPE_NONE functionType = iota
	FUNCTION_TYPE_FUNCTION
	FUNCTION_TYPE_METHOD
)

type classType int

const (
	CLASS_TYPE_NONE classType = iota
	CLASS_TYPE_CLASS
)

type Resolver struct {
	interpreter  *interpreting.Interpreter
	scopes       []map[string]bool
	currFunction functionType
	currClass    classType
	hadError     bool
}

func NewResolver(interpreter *interpreting.Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
		scopes:      []map[string]bool{},
	}
}

func (r *Resolver) Resolve(statements []ast.Stmt) bool {
	r.resolveBlock(statements)
	return r.hadError
}

func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) (any, error) {
	r.beginScope()
	defer r.endScope()
	r.resolveBlock(stmt.Statements)

	return nil, nil
}

func (r *Resolver) VisitClassStmt(stmt *ast.ClassStmt) (any, error) {
	enclosingClass := r.currClass
	r.currClass = CLASS_TYPE_CLASS
	defer func() { r.currClass = enclosingClass }()

	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.beginScope()
	defer r.endScope()
	r.scopes[len(r.scopes)-1]["this"] = true

	for _, method := range stmt.Methods {
		declaration := FUNCTION_TYPE_METHOD
		r.resolveFunction(method, declaration)
	}

	return nil, nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.ExpressionStmt) (any, error) {
	r.resolveExpr(stmt.Expression)
	return nil, nil
}

func (r *Resolver) VisitVarStmt(statement *ast.VarStmt) (any, error) {
	r.declare(statement.Name)
	if statement.Initializer != nil {
		r.resolveExpr(statement.Initializer)
	}
	r.define(statement.Name)

	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt *ast.WhileStmt) (any, error) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil, nil
}

func (r *Resolver) VisitVariableExpr(expr *ast.VariableExpr) (any, error) {
	if len(r.scopes) > 0 {
		scope := r.scopes[len(r.scopes)-1]
		isDefined, exists := scope[expr.Name.Lexeme]
		if exists && !isDefined {
			r.reportError(expr.Name, "can't read local variable in its own initializer")
			return nil, nil
		}
	}

	r.resolveLocal(expr, expr.Name)

	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(statement *ast.FunctionStmt) (any, error) {
	r.declare(statement.Name)
	r.define(statement.Name)

	r.resolveFunction(statement, FUNCTION_TYPE_FUNCTION)

	return nil, nil
}

func (r *Resolver) VisitIfStmt(stmt *ast.IfStmt) (any, error) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}

	return nil, nil
}

func (r *Resolver) VisitPrintStmt(stmt *ast.PrintStmt) (any, error) {
	r.resolveExpr(stmt.Expression)
	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt *ast.ReturnStmt) (any, error) {
	if r.currFunction == FUNCTION_TYPE_NONE {
		r.reportError(stmt.Keyword, "can't return from top-level code")
	}

	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}

	return nil, nil
}

func (r *Resolver) VisitAssignmentExpr(expr *ast.AssignmentExpr) (any, error) {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)

	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr *ast.BinaryExpr) (any, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr *ast.CallExpr) (any, error) {
	r.resolveExpr(expr.Callee)

	for _, arg := range expr.Arguments {
		r.resolveExpr(arg)
	}

	return nil, nil
}

func (r *Resolver) VisitGetExpr(expr *ast.GetExpr) (any, error) {
	r.resolveExpr(expr.Object)
	return nil, nil
}

func (r *Resolver) VisitSetExpr(expr *ast.SetExpr) (any, error) {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)
	return nil, nil
}

func (r *Resolver) VisitThisExpr(expr *ast.ThisExpr) (any, error) {
	if r.currClass == CLASS_TYPE_NONE {
		r.reportError(expr.Keyword, "can't use 'this' outside of a class")
		return nil, nil
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *ast.GroupingExpr) (any, error) {
	r.resolveExpr(expr.Expression)
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr *ast.LiteralExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr *ast.LogicalExpr) (any, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) resolveBlock(statements []ast.Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) resolveStmt(statement ast.Stmt) {
	statement.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveFunction(funStmt *ast.FunctionStmt, funType functionType) {
	enclosingType := r.currFunction
	r.currFunction = funType
	defer func() { r.currFunction = enclosingType }()

	r.beginScope()
	defer r.endScope()

	for _, param := range funStmt.Parameters {
		r.declare(param)
		r.define(param)
	}

	r.resolveBlock(funStmt.Body)
}

func (r *Resolver) declare(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	_, alreadyDefined := scope[name.Lexeme]
	if alreadyDefined {
		r.reportError(name, "there is already a variable with that name in this scope")
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		scope := r.scopes[i]
		_, hasName := scope[name.Lexeme]
		if hasName {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) reportError(tok token.Token, msg any) {
	r.hadError = true

	if tok.Kind == token.EOF {
		r.report(errors.NewBuildtimeError(tok.Line, tok.Column, " at end", msg))
	}

	r.report(errors.NewBuildtimeError(tok.Line, tok.Column, " at '"+tok.Lexeme+"'", msg))
}

func (r *Resolver) report(err error) {
	fmt.Fprintln(os.Stderr, err)
}
