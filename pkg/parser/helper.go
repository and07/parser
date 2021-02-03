package parser

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"

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
