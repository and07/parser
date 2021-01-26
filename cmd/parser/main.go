package main

import (
	"context"

	"github.com/and07/parser/pkg/parser"
	"github.com/chromedp/chromedp"
)

func main() {

	// 禁用chrome headless
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := parser.New(allocCtx)
	defer cancel()

	res := parser.Run(ctx, ruleData)
	parser.ExportCSV(res, "./result.csv")
}
