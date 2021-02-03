# PARSER

This is a basic parser


## use WINDOWS
``` sh
./bin/parser -r=$(pwd)\\cmd\\parser\\rule.json -c=$HOME\\parser.csv

./bin/parser -r=$(pwd)\\cmd\\parser\\rule.json > out.csv

cat $(pwd)\\cmd\\parser\\rule.json | ./bin/parser > $HOME\\test.csv
```

## use MAC
``` sh
./bin/parser -r=$(pwd)/cmd/parser/rule.json -c=$HOME/parser.csv

./bin/parser -r=$(pwd)/cmd/parser/rule.json > out.csv

cat $(pwd)/cmd/parser/rule.json | ./bin/parser > $HOME/test.csv
```

``` go

package main

import (
	"context"

	"github.com/and07/parser/pkg/parser"
)


var rule string = `
{
	"url":"http://example.com/index.html",
	"all":false,
	"rule":[
			 {
				"type":"text",
				"name":"H1",
				"attr":null,
				"page":"http://example.com/index.html",
				"path":"//html/body/h1[1]"
			 },
			 {
				"type":"html",
				"name":"P",
				"attr":null,
				"page":"http://example.com/index.html",
				"path":"//html/body/p[1]"
			 }
	],
	"host":"http://example.com/index.html",
	"name":"",
	"domen":"example.com",
	"limit":"1",
	"path_type":0
 }
`

func main() {

	ctx, cancel := parser.New(context.Background())
	defer cancel()

	res := parser.Run(ctx, rule)
	parser.ExportCSV(res)
}
```