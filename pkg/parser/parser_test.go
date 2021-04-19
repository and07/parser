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
	"limit":1,
	"path_type":0
 }
`

var rule1 string = `
{
	"rule":[
	   {
		  "name":"link",
		  "children":[
			 {
				"name":"title1"
			 },
			 {
				"name":"test"
			 },
			 {
				"name":"img1"
			 },
			 {
				"name":"cat"
			 }
		  ]
	   }
	]
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

	rule := rules([]byte(rule))

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

func TestHead(t *testing.T) {
	rule := rules([]byte(rule))
	heads := make([]string, 0)
	head(rule.Rule, &heads)
	if len(heads) != 2 {
		t.Error(heads)
	}

	rule1 := rules([]byte(rule1))
	heads1 := make([]string, 0)
	head(rule1.Rule, &heads1)
	if len(heads1) != 5 {
		t.Error(heads1)
	}
}
