package parser

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	"github.com/chromedp/chromedp"
)

var rule string = `
{
	"url":"index.html",
	"all":false,
	"rule":[
			 {
				"type":"text",
				"name":"H1",
				"attr":null,
				"page":"index.html",
				"path":"//html/body/h1[1]"
			 },
			 {
				"type":"html",
				"name":"P",
				"attr":null,
				"page":"index.html",
				"path":"//html/body/p[1]"
			 }
	],
	"host":"index.html",
	"name":"",
	"domen":"index.html",
	"limit":"1",
	"path_type":0
 }
`

func TestExtractData(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("could not get working directory: %v", err))
	}
	testdataDir := "file://" + path.Join(wd, "testdata")

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	rule := rules(rule)

	data := make(map[string]map[string]interface{})
	content(ctx, testdataDir+"/"+rule.URL, rule.Rule, data)
	for key := range data {
		if "Заголовок страницы" != *(data[key]["H1"].(*string)) {
			t.Error("H1 ERROR ", *(data[key]["H1"].(*string)))
		}
		if "<p>Основной текст.</p>" != *(data[key]["P"].(*string)) {
			t.Error("P ERROR ", *(data[key]["P"].(*string)))
		}
	}

}
