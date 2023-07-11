package parser

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

func removeTagLansana(s, tag string) string {
	startingScriptTag := "<" + tag
	endingScriptTag := "</" + tag + ">"

	var script string

	for {
		startingScriptTagIndex := strings.Index(s, startingScriptTag)
		endingScriptTagIndex := strings.Index(s, endingScriptTag)

		if startingScriptTagIndex > -1 && endingScriptTagIndex > -1 {
			script = s[startingScriptTagIndex : endingScriptTagIndex+len(endingScriptTag)]
			s = strings.Replace(s, script, "", 1)
			continue
		}

		break
	}

	return s
}

func removeHTMLTag(htmlString string) string {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		log.Fatal(err)
	}
	removeTag("script", doc)
	removeTag("style", doc)
	buf := bytes.NewBuffer([]byte{})
	if err := html.Render(buf, doc); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

func removeTag(tag string, n *html.Node) {
	// if note is script tag
	if n.Type == html.ElementNode && n.Data == tag {
		n.Parent.RemoveChild(n)
		return // script tag is gone...
	}
	// traverse DOM
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		removeTag(tag, c)
	}
}

func removeTabs(htmlString string) string {
	//htmlString = strings.ReplaceAll(htmlString, "\"", "'")
	htmlString = strings.ReplaceAll(htmlString, "\r\n", "")
	htmlString = strings.ReplaceAll(htmlString, "\t", "")
	htmlString = strings.ReplaceAll(htmlString, "\n", "")
	return htmlString
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getFile(filePath string) (*os.File, error) {
	if filePath == "" {
		return nil, errors.New("please input a file")
	}
	if !fileExists(filePath) {
		return nil, errors.New("the file provided does not exist")
	}

	file, e := os.Open(filepath.Clean(filePath))
	if e != nil {
		return nil, errors.Wrapf(e, "unable to read the file %s", filePath)
	}
	return file, nil
}

func head(rule []*Rule, headline *[]string) {
	if len(rule) == 0 {
		return
	}
	h := make([]string, 0)
	for _, r := range rule {
		h = append(h, r.Name)
		if len(r.Children) > 0 {
			head(r.Children, &h)
		}
	}

	*headline = append(*headline, h...)

}

// Conf ...
func Conf(ctx context.Context, r io.Reader) (context.Context, error) {

	ctx = context.WithValue(ctx, "rule", readConfig(r))
	return ctx, nil
}

func readConfigFromFile(fileName string) ([]byte, error) {
	file, e := getFile(fileName)
	if e != nil {
		return nil, e
	}
	defer file.Close()

	return readConfig(file), nil
}

func readConfig(r io.Reader) []byte {
	byteValue, _ := ioutil.ReadAll(r)
	return byteValue
}

// RuleConfig ...
func RuleConfig(ctx context.Context, fileName string) (context.Context, error) {
	file, e := getFile(fileName)
	if e != nil {
		return ctx, e
	}
	defer file.Close()

	return Conf(ctx, file)
}

// ExportCSV ...
func ExportCSV(ctx context.Context, path string) error {
	return export(ctx, path)
}

// Output ...
func Output(ctx context.Context, w io.Writer) error {
	res := ctx.Value("res").([][]string)

	return write(res, w)
}

func write(res [][]string, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	for _, value := range res {
		err := writer.Write(value)
		if err != nil {
			return err
		}
	}
	return nil
}

func runTasks(ctx context.Context, url string) error {
	log.Printf("chromedp.Run url %s", url)

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

func content(ctx context.Context, url string, rule []*Rule, data map[string]map[string]interface{}) error {

	if len(rule) <= 0 {
		return nil
	}

	if err := runTasks(ctx, url); err != nil {
		log.Printf("runTasks url %s %s", url, err)
		return err
	}
	if _, ok := data[url]; !ok {
		data[url] = make(map[string]interface{})
	}
	tasksData := tasks(rule, data[url])
	log.Printf("tasksData url %s %#v", url, tasksData)

	for _, t := range tasksData {
		log.Println(t)
		c, cancel := context.WithTimeout(ctx, 30*time.Second)
		err := chromedp.Run(c, t)
		if err != nil {
			log.Printf("chromedp.Run err %s ", err)
			cancel()
			continue
		}
		cancel()
	}

	log.Println("-------")
	c := FromContext(ctx)
	if c != nil {
		dataNormalize := normalize(c.headsline, data)
		log.Printf("dataNormalize %#v", dataNormalize)
		if len(dataNormalize) > 0 && len(dataNormalize[0]) > 0 {
			c.res <- dataNormalize
			delete(data, url)
		}
	}

	if _, ok := data[url]["link"]; ok {
		defer log.Printf("url %#v", data[url]["link"].(*[]*cdp.Node))
		if len(rule[0].Children) > 0 {
			for _, node := range *(data[url]["link"].(*[]*cdp.Node)) {
				url := node.AttributeValue("href")
				if err := content(ctx, url, rule[0].Children, data); err != nil {
					log.Println(err)
				}

			}
		}
	}

	return nil
}

func rules(rule []byte) Parser {
	var r Parser
	if err := json.Unmarshal(rule, &r); err != nil {
		log.Fatal(err)
	}
	return r
}

func fileWriter(path string) (*os.File, error) {
	return os.Create(path)
}

func export(ctx context.Context, path string) error {

	file, err := fileWriter(path)
	if err != nil {
		return err
	}

	defer file.Close()

	return Output(ctx, file)

}

func headliner(headsline []string) [][]string {
	var f [][]string
	var headline []string
	for i := range headsline {
		if headsline[i] == "link" {
			continue
		}

		headline = append(headline, headsline[i])

	}
	f = append(f, headline)
	return f
}

func normalize(headsline []string, data map[string]map[string]interface{}) [][]string {
	var f [][]string
	var headline []string
	for i := range headsline {
		if headsline[i] == "link" {
			continue
		}
		for key := range data {
			if _, ok := data[key][headsline[i]]; !ok {
				headline = append(headline, headsline[i])
			}
		}
	}

	for key := range data {
		var line []string
		for i := range headline {
			v := data[key][headline[i]]
			switch v.(type) {
			case *string:
				htmlString := removeTabs(*(v.(*string)))
				htmlString = removeTagLansana(htmlString, "script")
				htmlString = removeTagLansana(htmlString, "style")
				line = append(line, htmlString)
			case *[]*cdp.Node:
				log.Printf("%s %T %#v", headline[i], v, *(v.(*[]*cdp.Node)))
			}
		}
		if len(line) > 0 {
			f = append(f, line)
		}

	}
	return f
}
