package parser

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// Rule ...
type Rule struct {
	Type     string      `json:"type"`
	Name     string      `json:"name"`
	Attr     string      `json:"attr"`
	Parent   interface{} `json:"parent"`
	Page     string      `json:"page"`
	Path     string      `json:"path"`
	Children []*Rule     `json:"children,omitempty"`
}

// Parser ...
type Parser struct {
	URL      string  `json:"url"`
	All      bool    `json:"all"`
	Rule     []*Rule `json:"rule"`
	Host     string  `json:"host"`
	Name     string  `json:"name"`
	Domen    string  `json:"domen"`
	Limit    string  `json:"limit"`
	PathType int     `json:"path_type"`
}

type fass struct {
	rule Parser
	res  [][]string
}

// New ...
func New(ctx context.Context) (context.Context, context.CancelFunc) {
	// create chrome instance
	return chromedp.NewContext(
		ctx,
		chromedp.WithLogf(log.Printf),
	)
}

func rules(rule string) Parser {
	var r Parser
	if err := json.Unmarshal([]byte(rule), &r); err != nil {
		log.Fatal(err)
	}
	return r
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func export(data [][]string) {

	file, err := os.Create("./result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}

}

func keys(data map[string]map[string]interface{}) []string {
	var keys []string
	for key := range data {
		for k := range data[key] {
			keys = append(keys, k)
		}
		break
	}
	return keys
}

func normalize(data map[string]map[string]interface{}) [][]string {
	var f [][]string
	keys := keys(data)
	for key := range data {
		var line []string
		for i := range keys {
			v := data[key][keys[i]]
			switch v.(type) {
			case *string:
				htmlString := removeTabs(*(v.(*string)))
				htmlString = removeTagLansana(htmlString, "script")
				htmlString = removeTagLansana(htmlString, "style")
				line = append(line, htmlString)
			case *[]*cdp.Node:
				log.Printf("%s %T %#v", keys[i], v, *(v.(*[]*cdp.Node)))
			}
		}
		f = append(f, line)
	}
	return f
}

// Run ...
func Run(ctx context.Context, ruleData string) [][]string {
	rule := rules(ruleData)

	data := make(map[string]map[string]interface{})
	content(ctx, rule.URL, rule.Rule, data)

	res := normalize(data)
	//log.Printf("data %#v", data)

	return res
}

// ExportCSV ...
func ExportCSV(res [][]string) {
	export(res)
}

func runTasks(ctx context.Context, url string) error {
	return chromedp.Run(ctx,
		chromedp.Tasks{
			chromedp.Navigate(url),
			chromedp.Sleep(2 * time.Second),
		},
	)
}

func tasks(rules []*Rule, data map[string]interface{}) chromedp.Tasks {
	var tasks chromedp.Tasks

	for _, r := range rules {

		if r.Type == "link" {
			var nodes []*cdp.Node
			data[r.Name] = &nodes
		} else {
			var x string
			data[r.Name] = &x
		}

		tasks = append(tasks, tasksData(r, data))
	}
	return tasks
}

func tasksData(rule *Rule, content map[string]interface{}) chromedp.QueryAction {
	var ok bool
	switch rule.Type {
	case "text":
		return chromedp.TextContent(rule.Path, content[rule.Name].(*string))
	case "attr":
		return chromedp.AttributeValue(rule.Path, rule.Attr, content[rule.Name].(*string), &ok)
	case "link":
		return chromedp.Nodes(rule.Path, content[rule.Name].(*[]*cdp.Node))
	case "html":
		return chromedp.OuterHTML(rule.Path, content[rule.Name].(*string))
	}
	return nil
}

func content(ctx context.Context, url string, rule []*Rule, data map[string]map[string]interface{}) {

	if len(rule) <= 0 {
		return
	}

	if err := runTasks(ctx, url); err != nil {
		log.Fatal(err)
	}
	if _, ok := data[url]; !ok {
		data[url] = make(map[string]interface{})
	}
	tasksData := tasks(rule, data[url])
	log.Printf("%#v", tasksData)

	err := chromedp.Run(ctx, tasksData)
	if err != nil {
		log.Fatal(err)
	}
	if _, ok := data[url]["link"]; ok {
		log.Printf("%#v", data[url]["link"].(*[]*cdp.Node))
		if len(rule[0].Children) > 0 {
			for _, node := range *(data[url]["link"].(*[]*cdp.Node)) {
				content(ctx, node.AttributeValue("href"), rule[0].Children, data)
			}

		}
	}

	return
}
