package parser

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/chromedp/chromedp"
)

type Logger interface {
	// all levels + Prin
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

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
	Limit    int     `json:"limit"`
	PathType int     `json:"path_type"`
}

// Crawler ...
type Crawler struct {
	parser    Parser
	headsline []string
	res       chan [][]string
	logger    Logger
	w         io.Writer
}

type contextKey struct{}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

// CrawlerOption is a context option.
type CrawlerOption = func(*Crawler)

// WithLogger ...
func WithLogger(log Logger) CrawlerOption {
	return func(c *Crawler) { c.logger = log }
}

// WithConfigs ...
func WithConfigs(rulePath string) CrawlerOption {
	return func(c *Crawler) {
		var conf []byte
		var err error
		if isInputFromPipe() {
			conf = readConfig(os.Stdin)
		} else {
			if conf, err = readConfigFromFile(rulePath); err != nil {
				panic(err)
			}
		}

		c.parser = rules(conf)
	}
}

// WithWriter ...
func WithWriter(outPath string) CrawlerOption {
	return func(c *Crawler) {
		var writer io.Writer
		if outPath != "" {
			writer, _ = fileWriter(outPath)
		} else {
			writer = os.Stdout
		}
		c.w = writer
	}
}

// New ...
func New(ctx context.Context, opts ...CrawlerOption) (context.Context, context.CancelFunc) {

	c := &Crawler{
		res: make(chan [][]string, 100),
	}

	for _, o := range opts {
		o(c)
	}

	for i := 0; i < 15; i++ {
		if c.logger != nil {
			c.logger.Debugf("----worker req %d----", i)
		}

		go func(cnt int) {
			for req := range c.res {
				if err := c.write(req); err != nil {
					log.Println(err)
					if c.logger != nil {
						c.logger.Errorln(err.Error())
					}
				}
			}
		}(i)
	}

	headsline := make([]string, 0)
	head(c.parser.Rule, &headsline)
	c.headsline = headsline

	ctx = context.WithValue(ctx, contextKey{}, c)
	// create chrome instance
	return chromedp.NewContext(
		ctx,
		chromedp.WithLogf(log.Printf),
	)
}

// FromContext extracts the Context data stored inside a context.Context.
func FromContext(ctx context.Context) *Crawler {
	c, _ := ctx.Value(contextKey{}).(*Crawler)
	return c
}

// Run ...
func Run(ctx context.Context) (context.Context, error) {

	c := FromContext(ctx)

	c.res <- headliner(c.headsline)

	urlBase := c.parser.URL

	if c.parser.Limit > 1 {
		for i := 1; i <= c.parser.Limit; i++ {
			data := make(map[string]map[string]interface{})
			url := fmt.Sprintf(urlBase, i) //TODO
			if c.logger != nil {
				c.logger.Debugf("url: %s", url)
			}
			if err := content(ctx, url, c.parser.Rule, data); err != nil {
				return ctx, err
			}
			//c.res <- normalize(c.headsline, data)

		}

		return ctx, nil
	}

	data := make(map[string]map[string]interface{})
	if err := content(ctx, urlBase, c.parser.Rule, data); err != nil {
		return ctx, err
	}

	//c.res <- normalize(c.headsline, data)

	return ctx, nil
}

func (c *Crawler) write(res [][]string) error {
	return write(res, c.w)
}
