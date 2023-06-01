package Interpreter

import (
	"fmt"

	"github.com/AnshVM/golox/Ast"
	"github.com/AnshVM/golox/Environment"
	"github.com/AnshVM/golox/Error"
	"github.com/AnshVM/golox/Parser"
	"github.com/AnshVM/golox/Tokens"
)

type Interpreter struct {
	Env *Environment.Environment
}

func NewInterpreter(env *Environment.Environment) *Interpreter {
	env.Define("clock", Clock())
	return &Interpreter{Env: env}
}

func (i *Interpreter) Interpret(stmts []Parser.Stmt) error {
	for _, stmt := range stmts {
		err := i.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) Exec(stmt Parser.Stmt) error {
	switch s := stmt.(type) {
	case Ast.ExpressionStmt:
		return i.ExecExpressionStmt(&s)
	case Ast.PrintStmt:
		return i.ExecPrintStmt(&s)
	case Ast.VarStmt:
		return i.ExecVarStmt(&s)
	case Ast.BlockStmt:
		return i.ExecBlockStmt(&s)
	case Ast.IfStmt:
		return i.ExecIfStmt(&s)
	case *Ast.WhileStmt:
		return i.ExecWhileStmt(s)
	case *Ast.Function:
		return i.ExecFuncStmt(s)
	}
	return nil
}

func (i *Interpreter) ExecReturnStmt(stmt *Ast.Return) (any, error) {
	var returnVal any = nil
	var err error
	if stmt.Value != nil {
		returnVal, err = i.Eval(stmt.Value)
	}
	return returnVal, err
}

func (i *Interpreter) ExecFuncStmt(stmt *Ast.Function) error {
	callable := CreateFunction(stmt)
	i.Env.Define(stmt.Name.Lexeme, callable)
	return nil
}

func (i *Interpreter) ExecWhileStmt(stmt *Ast.WhileStmt) error {
	for {
		condition, err := i.Eval(stmt.Condition)
		if err != nil {
			return err
		}
		if isTruthy(condition) {
			err := i.Exec(stmt.Body)
			if err != nil {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func (i *Interpreter) ExecIfStmt(stmt *Ast.IfStmt) error {
	condition, err := i.Eval(stmt.Condition)
	if err != nil {
		return err
	}
	if isTruthy(condition) {
		return i.Exec(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.Exec(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) ExecExpressionStmt(stmt *Ast.ExpressionStmt) error {
	_, err := i.Eval(stmt.Expression)
	return err
}

func (i *Interpreter) ExecPrintStmt(stmt *Ast.PrintStmt) error {
	result, err := i.Eval(stmt.Expression)
	if err == nil {
		fmt.Printf("%v\n", result)
	}
	return err
}

func (i *Interpreter) ExecVarStmt(stmt *Ast.VarStmt) error {
	if stmt.Initializer != nil {
		value, err := i.Eval(stmt.Initializer)
		if err == nil {
			i.Env.Define(stmt.Name.Lexeme, value)
		}
		return err
	}
	i.Env.Define(stmt.Name.Lexeme, nil)
	return nil
}

func (i *Interpreter) ExecBlockStmt(stmt *Ast.BlockStmt) error {
	_, err := i.executeBlock(stmt.Statements, &Environment.Environment{Enclosing: i.Env, Values: map[string]any{}})
	return err
}

func (i *Interpreter) executeBlock(statements []Parser.Stmt, env *Environment.Environment) (any, error) {
	prev := i.Env
	defer func() {
		i.Env = prev
	}()
	i.Env = env
	var err error
	for _, stmt := range statements {

		if stmt, ok := stmt.(*Ast.Return); ok {
			return i.ExecReturnStmt(stmt)
		}

		err = i.Exec(stmt)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) Eval(expr Parser.Expr) (any, error) {
	switch e := expr.(type) {
	case *Ast.LiteralExpr:
		return i.EvalLiteral(e), nil
	case *Ast.BinaryExpr:
		return i.EvalBinary(e)
	case *Ast.GroupingExpr:
		return i.EvalGrouping(e)
	case *Ast.UnaryExpr:
		return i.EvalUnary(e)
	case *Ast.ConditionalExpr:
		return i.EvalConditional(e)
	case *Ast.VariableExpr:
		return i.EvalVariable(e)
	case *Ast.AssignExpr:
		return i.EvalAssign(e)
	case *Ast.LogicalExpr:
		return i.EvalLogical(e)
	case *Ast.Call:
		return i.EvalCall(e)
	}
	return nil, Error.ErrRuntimeError
}

func (i *Interpreter) EvalCall(expr *Ast.Call) (any, error) {

	callee, err := i.Eval(expr.Callee)
	if err != nil {
		return nil, err
	}

	evaluatedArgs := []any{}
	for _, argExpr := range expr.Arguments {
		evalArg, err := i.Eval(argExpr)
		if err != nil {
			return nil, err
		}
		evaluatedArgs = append(evaluatedArgs, evalArg)
	}
	function, ok := callee.(*LoxCallable)
	if !ok {
		Error.ReportRuntimeError(expr.Paren, "Expression is not callable.")
		return nil, Error.ErrRuntimeError
	}
	if len(evaluatedArgs) != int(function.Arity()) {
		Error.ReportRuntimeError(expr.Paren, fmt.Sprintf("Expected %d arguments, got %d", function.Arity(), len(evaluatedArgs)))
		return nil, Error.ErrRuntimeError
	}
	return function.Call(i, evaluatedArgs)
}

func (i *Interpreter) EvalLogical(expr *Ast.LogicalExpr) (any, error) {
	left, err := i.Eval(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == Tokens.OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}
	return i.Eval(expr.Right)

}

func (i *Interpreter) EvalAssign(expr *Ast.AssignExpr) (any, error) {
	value, err := i.Eval(expr.Value)
	if err != nil {
		return nil, err
	}
	err = i.Env.Assign(expr.Name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) EvalLiteral(expr *Ast.LiteralExpr) any {
	return expr.Value
}

func (i *Interpreter) EvalGrouping(expr *Ast.GroupingExpr) (any, error) {
	evaluated, err := i.Eval(expr.Expression)
	if err != nil {
		return nil, err
	}
	return evaluated, nil
}

func (i *Interpreter) EvalUnary(expr *Ast.UnaryExpr) (any, error) {
	right, err := i.Eval(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case Tokens.MINUS:
		rightVal, err := checkNumberOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}
		return -rightVal, nil

	case Tokens.BANG:
		return !isTruthy(right), nil
	}
	return nil, Error.ErrRuntimeError
}

func (i *Interpreter) EvalVariable(expr *Ast.VariableExpr) (any, error) {
	return i.Env.Get(expr.Name)
}

func (i *Interpreter) EvalConditional(conditional *Ast.ConditionalExpr) (any, error) {
	cond, err := i.Eval(conditional.Condition)
	if err != nil {
		return nil, Error.ErrRuntimeError
	}
	if cond == true {
		return i.Eval(conditional.Then)
	} else {
		return i.Eval(conditional.Else)
	}
}

func (i *Interpreter) EvalBinary(binary *Ast.BinaryExpr) (any, error) {
	switch binary.Operator.Type {
	case Tokens.MINUS:
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left - right), nil

	case Tokens.SLASH:
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		if right == 0 {
			Error.ReportRuntimeError(binary.Operator, "Cannot divide by zero")
			return nil, Error.ErrRuntimeError
		}
		return (left / right), nil

	case Tokens.STAR:
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left * right), nil

	case Tokens.PLUS:
		left, right, err := i.EvalBinaryOperandsAny(binary)
		if err != nil {
			return nil, err
		}
		if isFloat32(right) && isFloat32(left) {
			return (left.(float32) + right.(float32)), nil
		}
		if isString(right) && isString(left) {
			return (left.(string) + right.(string)), nil
		}
		// if isString(right) && isFloat32(left) {
		// 	val, _ := left.(float32)
		// 	return (strconv.FormatFloat(float64(val), 'f', 0, 32) + right.(string)), nil
		// }
		// if isString(left) && isFloat32(right) {
		// 	val, _ := right.(float32)
		// 	return (left.(string) + strconv.FormatFloat(float64(val), 'f', 0, 32)), nil
		// }
		Error.ReportRuntimeError(binary.Operator, "Operands must strings or numbers")
		return nil, Error.ErrRuntimeError

	case Tokens.GREATER:
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left > right), nil

	case Tokens.GREATER_EQUAL:
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left >= right), nil

	case Tokens.LESS:
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left < right), nil

	case Tokens.LESS_EQAUL:
		left, right, err := i.EvalBinaryOperandsNumber(binary)
		if err != nil {
			return nil, err
		}
		return (left <= right), nil

	case Tokens.EQUAL_EQUAL:
		left, right, err := i.EvalBinaryOperandsAny(binary)
		if err != nil {
			return nil, err
		}
		return isEqual(left, right), nil

	case Tokens.BANG_EQUAL:
		left, right, err := i.EvalBinaryOperandsAny(binary)
		if err != nil {
			return nil, err
		}
		return isEqual(left, right), nil
	}

	// unreachable
	return nil, nil
}

func (i *Interpreter) EvalBinaryOperandsAny(binary *Ast.BinaryExpr) (any, any, error) {
	right, err := i.Eval(binary.Right)
	if err != nil {
		return nil, nil, err
	}
	left, err := i.Eval(binary.Left)
	if err != nil {
		return nil, nil, err
	}
	return left, right, nil
}

func (i *Interpreter) EvalBinaryOperandsNumber(binary *Ast.BinaryExpr) (float32, float32, error) {
	evalRight, err := i.Eval(binary.Right)
	if err != nil {
		return 0, 0, err
	}
	evalLeft, err := i.Eval(binary.Left)
	if err != nil {
		return 0, 0, err
	}
	right, err := checkNumberOperand(binary.Operator, evalRight)
	if err != nil {
		return 0, 0, err
	}
	left, err := checkNumberOperand(binary.Operator, evalLeft)
	if err != nil {
		return 0, 0, err
	}
	return left, right, nil
}

func checkNumberOperand(operator *Tokens.Token, right any) (float32, error) {
	if val, ok := right.(float32); ok {
		return val, nil
	} else {
		Error.ReportRuntimeError(operator, "Operand must be a number")
		return 0, Error.ErrRuntimeError
	}
}

// Handles the special case that one operand is boolean, and other is not
func isEqual(a any, b any) bool {
	if isBool(a) && !isBool(b) {
		b = isTruthy(b)
	}
	if isBool(b) && !isBool(a) {
		a = isTruthy(a)
	}
	return a == b
}

func isBool(x any) bool {
	return x == true || x == false
}

func isFloat32(val any) bool {
	_, ok := val.(float32)
	return ok
}

func isString(val any) bool {
	_, ok := val.(string)
	return ok
}

func isTruthy(val any) bool {
	if val == nil {
		return false
	}
	if boolVal, ok := val.(bool); ok {
		return boolVal
	}
	return true

}
