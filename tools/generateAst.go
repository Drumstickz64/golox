// generating the AST programatically is probably useless in Go,
// but I'm going to do it anyway because I want to

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/Drumstickz64/golox/errors"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: generateAst <output_directory>")
		os.Exit(64)
	}

	outputDir := os.Args[1]
	err := defineAst(outputDir, "Expr", []string{
		"Binary   : Left Expr, Operator token.Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal  : Value any",
		"Unary    : Operator token.Token, Right Expr",
	})

	if err != nil {
		errors.LogCliError("error while generating AST: "+err.Error(), 65)
	}
}

func defineAst(outputDir, baseName string, kinds []string) error {
	packageName := strings.ToLower(baseName)
	content := ""

	content += "package ast\n"
	content += "\n"
	content += "import \"github.com/Drumstickz64/golox/token\"\n"
	content += "\n"

	content += defineVisitor(baseName, kinds)

	content += fmt.Sprintf("type %s interface {\n", baseName)
	content += fmt.Sprintf("	Accept(%sVisitor) (any, error)\n", baseName)
	content += "}\n"
	content += "\n"

	for _, kind := range kinds {
		content += defineType(kind, baseName)
	}

	pth := path.Join(outputDir, packageName+".go")
	if err := os.WriteFile(pth, []byte(content), 0777); err != nil {
		return err
	}

	return formatFile(pth)
}

func defineVisitor(baseName string, kinds []string) string {
	content := ""
	content += fmt.Sprintf("type %sVisitor interface {\n", baseName)
	for _, kind := range kinds {
		kindName := strings.TrimSpace(strings.Split(kind, ":")[0])
		itemName := kindName + baseName
		content += fmt.Sprintf("	Visit%s(*%s) (any, error)\n", itemName, itemName)
	}
	content += "}\n"

	content += "\n"

	return content
}

func defineType(kind, baseName string) string {
	content := ""
	kindName := strings.TrimSpace(strings.Split(kind, ":")[0])
	itemName := kindName + baseName
	fields := strings.TrimSpace(strings.Split(kind, ":")[1])

	content += fmt.Sprintf("type %s struct {\n", itemName)
	for _, field := range strings.Split(fields, ", ") {
		content += "\t" + field + "\n"
	}
	content += "}\n"

	content += "\n"

	selfName := strings.ToLower(kindName[0:1])
	content += fmt.Sprintf("func (%s *%s) Accept(visitor %sVisitor) (any, error) {\n", selfName, itemName, baseName)
	content += fmt.Sprintf("	return visitor.Visit%s(%s)\n", itemName, selfName)
	content += "}\n"

	content += "\n"

	return content
}

func formatFile(pth string) error {
	return exec.Command("go", "fmt", pth).Run()
}
