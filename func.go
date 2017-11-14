package main

import (
	"fmt"
	"strings"
)

// Func represents information about Function declaration
// needede for this tool. It's a stipped down version of
// ast.FuncDecl basically.
type Func struct {
	Name         string
	ParamCount   int
	Params       []Var
	ReturnsCount int
	Returns      []Var
}

// String implements Stringer interface for Func.
func (f *Func) String() string {
	return fmt.Sprintf("func %s(%s) %s {…}", f.Name, f.ParamsNames(), f.ReturnsNames())
}

// ParamsNames returns human-readable comma separated string for function parameters.
func (f *Func) ParamsNames() string {
	var ret []string
	for _, arg := range f.Params {
		ret = append(ret, arg.String())
	}
	return strings.Join(ret, ", ")
}

// ReturnsNames returns human-readable comma separated string for function returns.
func (f *Func) ReturnsNames() string {
	var ret []string
	for _, arg := range f.Returns {
		ret = append(ret, arg.String())
	}
	return strings.Join(ret, ", ")
}
