package main

import (
	"log"
	"time"

	"github.com/itsbth/glox/ast"
	"github.com/itsbth/glox/parser"
	"github.com/itsbth/glox/scanner"
)

func main() {
	log.Printf("Hello, World!")
	scan := scanner.NewScanner("var a = /* this /* is */ nested */ 2 + 3;")
	go scan.Scan()
	go func() {
		// this is not a workable solution
		for err := range scan.Errors {
			log.Println(err.Error())
		}
	}()
	toks := make([]scanner.Token, 0)
	for token := range scan.Tokens {
		log.Printf("token: %s\n", token.String())
		toks = append(toks, token)
	}
	time.Sleep(2 * time.Second)
	node := ast.NewBinary(ast.NewIdentifier("a"), scanner.T_PLUS, ast.NewIdentifier("b"))
	log.Printf("node: %s", node.String())
	log.Println(node.Accept(PrettyPrinter{}))
	parser := parser.NewParser(toks)
	log.Printf("parsed: %s", parser.Parse())
}
