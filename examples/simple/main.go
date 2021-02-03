package main

import (
	"context"
	"log"

	"github.com/and07/parser/pkg/parser"
)

func main() {

	ctx, cancel := parser.New(context.Background())
	defer cancel()
	var err error
	if ctx, err = parser.RuleConfig(ctx, "./rule.json"); err != nil {
		log.Printf("ERROR parser.Rule %s", err)
	}
	parser.Run(ctx)
	parser.ExportCSV(ctx, "./result.csv")
}
