package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/and07/parser/pkg/parser"
	"github.com/chromedp/chromedp"
)

const requestTimeout = 5 * time.Second

func httpClient(requestTimeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: requestTimeout,
	}
}

func main() {

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	parser.Run(ctx, ruleData)
}
