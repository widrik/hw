package main

import (
	"io"
	"log"
	"os"
	"strings"

	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"path/filepath"
)

type GeneratatedFile struct {
	Name    string
	Structs []Struct
}

type Struct struct {
	Name    string
	VarName string
	Fields  []Field
}

type Field struct {
	Names      []string
	Type       FieldType
	Validators []Validator
}

type FieldType struct {
	VarType string
	Type    string
}

type Validator struct {
	Type  string
	Value interface{}
}

func Generate(fileName string) (error) {
	astFile, err := parseFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	resultFileName := strings.ReplaceAll(fileName, filepath.Ext(fileName), "_validator.go")
	resultFile, err := os.Create(resultFileName)
	defer func() {
		err := resultFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		return err
	}

	structsData, err := CreateStructsData(*astFile)
	if err != nil {
		return err
	}

	template, err := CreateTemplate(structsData)
	if err != nil {
		log.Fatal(err)
	}

	resultedTemplate, err := format.Source(template.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.WriteString(resultFile, string(resultedTemplate))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func parseFile(fileName string) (file *ast.File, err error) {
	astFile, err := parser.ParseFile(token.NewFileSet(), fileName, nil, 0)
	if err != nil {
		return
	}

	return astFile, err
}
