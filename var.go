package main

import (
	"fmt"
	"go/ast"
)

// Var represents single var for function parameters
// and returns.
// TODO(divan): word "var" probably is incorrect here,
// but I can't come up with better at the moment.:w
type Var struct {
	Name string
	Type Type
}

// String implements Stringer interface for Var.
func (t Var) String() string {
	return fmt.Sprintf("%s %s", t.Name, t.Type)
}

// Type defines supported types by this tool.
type Type string

const (
	TypeCChar   Type = "*C.char"
	TypeCInt    Type = "C.int"
	TypeUnknown Type = "N/A"
)

// TypeFromExpr attemtps to convert ast.Expr into
// Type. Incompattible or incorrect expression is ignored.
//
// It looks for SelectorExpr to catch two cases for now:
// "C.int" and "*C.char"
func TypeFromExpr(expr ast.Expr) Type {
	parseSelectorExpr := func(e *ast.SelectorExpr) (x, sel string) {
		if e == nil {
			return
		}
		if s, ok := e.X.(*ast.Ident); ok {
			return s.Name, e.Sel.Name
		}
		return
	}

	var sel *ast.SelectorExpr
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		sel = e
	case *ast.StarExpr:
		if v, ok := e.X.(*ast.SelectorExpr); ok {
			sel = v
		}
	}

	x, s := parseSelectorExpr(sel)
	if x == "C" {
		switch s {
		case "char":
			return TypeCChar
		case "int":
			return TypeCInt
		}
	}
	return TypeUnknown
}

// String implements Stringer interface for Type.
func (t Type) String() string {
	return string(t)
}
