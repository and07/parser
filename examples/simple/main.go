package main

import (
	"context"

	"github.com/and07/parser/pkg/parser"
)

func main() {

	ctx, cancel := parser.New(context.Background())
	defer cancel()

	res := parser.Run(ctx, ruleData1)
	parser.ExportCSV(res, "./result.csv")
}
