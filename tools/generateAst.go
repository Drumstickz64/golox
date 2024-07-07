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

	content += fmt.Sprintf("package %s\n", packageName)
	content += "\n"
	content += "import \"github.com/Drumstickz64/golox/token\"\n"
	content += "\n"

	content += defineVisitor(baseName, kinds)

	content += fmt.Sprintf("type %s interface {\n", baseName)
	content += fmt.Sprintf("	Accept(%sVisitor) any\n", baseName)
	content += "}\n"
	content += "\n"

	for _, kind := range kinds {
		content += defineType(kind, baseName)
	}

	pth := path.Join(outputDir, packageName+".go")
	if err := os.WriteFile(pth, []byte(content), 1); err != nil {
		return err
	}

	return formatFile(pth)
}

func defineVisitor(baseName string, kinds []string) string {
	content := ""
	content += fmt.Sprintf("type %sVisitor interface {\n", baseName)
	for _, kind := range kinds {
		kindName := strings.TrimSpace(strings.Split(kind, ":")[0])
		content += fmt.Sprintf("	Visit%s(exp *%s) any\n", kindName, kindName)
	}
	content += "}\n"

	content += "\n"

	return content
}

func defineType(kind, baseName string) string {
	content := ""
	kindName := strings.TrimSpace(strings.Split(kind, ":")[0])
	fields := strings.TrimSpace(strings.Split(kind, ":")[1])

	content += fmt.Sprintf("type %s struct {\n", kindName)
	for _, field := range strings.Split(fields, ", ") {
		content += "\t" + field + "\n"
	}
	content += "}\n"

	content += "\n"

	selfName := strings.ToLower(kindName[0:1])
	content += fmt.Sprintf("func (%s *%s) Accept(visitor %sVisitor) any {\n", selfName, kindName, baseName)
	content += fmt.Sprintf("	return visitor.Visit%s(%s)\n", kindName, selfName)
	content += "}\n"

	content += "\n"

	return content
}

func formatFile(pth string) error {
	return exec.Command("go", "fmt", pth).Run()
}
