package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s file.go\n", os.Args[0])
		os.Exit(1)
	}

	fd, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "library.go", fd, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// Inspect the AST and find all exported funcs
	var funcs []*Func
	ast.Inspect(f, func(node ast.Node) bool {
		if f, ok := node.(*ast.FuncDecl); ok {
			fn, ok := parseExportedFunc(f)
			if !ok {
				return true
			}

			// do something with fn
			funcs = append(funcs, fn)
		}
		return true
	})

	if len(funcs) == 0 {
		log.Fatal("CGO exported functions not found. Nothing to generate.")
	}

	out, err := GenerateOutput(funcs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}

// parseExportedFunc parses FuncDecl into our own Func with small subset
// of func information we need. It also filters out CGO-unexported functions.
func parseExportedFunc(f *ast.FuncDecl) (*Func, bool) {
	if exported := hasExportComment(f); !exported {
		return nil, false
	}

	fn := &Func{
		Name:         f.Name.Name,
		ParamCount:   f.Type.Params.NumFields(),
		Params:       listToVars(f.Type.Params),
		ReturnsCount: f.Type.Results.NumFields(),
		Returns:      listToVars(f.Type.Results),
	}
	return fn, true
}

// hasExportComment checks that FuncDecl has a special "//export" comment
// on the last line. See https://golang.org/cmd/cgo/#hdr-C_references_to_Go for
// more details.
func hasExportComment(f *ast.FuncDecl) bool {
	if f.Doc == nil {
		return false
	}

	lastComment := f.Doc.List[len(f.Doc.List)-1].Text
	return strings.HasPrefix(lastComment, "//export")
}

// listToVars converts ast.FieldList into slice of Var.
//
// ast.FieldList is used to describe function parameters and returns,
// so we take only information we need about them.
func listToVars(fields *ast.FieldList) []Var {
	var ret []Var
	if fields == nil {
		return ret
	}

	for _, field := range fields.List {
		// unnamed returns case
		if len(field.Names) < 1 {
			v := Var{
				Name: "",
				Type: TypeFromExpr(field.Type),
			}
			ret = append(ret, v)
			continue
		}

		// for declaration in form "name1, name2 string",
		// field.Names contains many entries — splitting it
		for _, name := range field.Names {
			v := Var{
				Name: name.Name,
				Type: TypeFromExpr(field.Type),
			}
			ret = append(ret, v)
		}
	}
	return ret
}
