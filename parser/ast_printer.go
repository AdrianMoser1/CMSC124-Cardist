package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (p *AstPrinter) Print(expr Expr) string {
	if expr == nil {
		return ""
	}
	return expr.Accept(p)
}

func (p *AstPrinter) VisitLiteral(expr *LiteralExpr) string {
	if expr.Value == nil {
		return "nil"
	}
	switch v := expr.Value.(type) {
	case float64:
		s := strconv.FormatFloat(v, 'f', -1, 64)
		if !strings.Contains(s, ".") {
			s += ".0"
		}
		return s
	case bool:
		return fmt.Sprintf("%t", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (p *AstPrinter) VisitUnary(expr *UnaryExpr) string {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) VisitBinary(expr *BinaryExpr) string {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGrouping(expr *GroupingExpr) string {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(p))
	}

	builder.WriteString(")")
	return builder.String()
}
