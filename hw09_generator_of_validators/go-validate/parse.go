package main

import (
	"go/parser"
	"go/token"
)

func ParseFile(file string) err {

	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, "", nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Продолжить тут
	// Разобраться с генерацией

}