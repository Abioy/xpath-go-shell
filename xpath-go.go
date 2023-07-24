package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
)

type TargetV2 map[string]*xpath.Expr

func usage() string {
	var help string
	help += "Usage:  xpath-go <PATH> <TARGET>\n"
	help += "    Query html from stdin via xpath expression and output in json.\n"
	help += "\n"
	help += "Arguments:\n"
	help += "    PATH   : expression to match.\n"
	help += "    TARGET : raw json string of `key`/`value` pairs. `value` should be relative path expression from leaf node matched above.\n"
	help += "\n"
	help += "Examples:\n"
	help += "    cat test.html | xpath-go \"//div[@class=\\\"seckill-timer\\\"]\" \"{\\\"id\\\":\\\"./@id\\\"}\"\n"
	help += "\n"
	return help
}

func Usage(errs ...any) {
	fmt.Fprint(os.Stderr, "[ERROR] ")
	fmt.Fprint(os.Stderr, errs...)
	fmt.Fprintf(os.Stderr, "\n%s", usage())
	os.Exit(1)
}

func main() {
	args := os.Args
	if args == nil || len(args) < 3 {
		Usage("len(args) < 3")
	}
	pattern := args[1]
	targets_str := args[2]

	var err error
	// path, err := xmlpath.Compile(pattern)
	expr, err := xpath.Compile(pattern)
	if err != nil {
		Usage("<PATH> illegal")
	}
	var data map[string]string
	if json.Unmarshal([]byte(targets_str), &data) != nil {
		Usage("<TARGET> illegal json")
	}
	target := TargetV2{}
	for k, v := range data {
		target[k], err = xpath.Compile(v)
		if err != nil {
			Usage("<TARGET> path illegal: ", k, " => ", v)
		}
	}

	file := os.Stdin
	doc, err := htmlquery.Parse(file)
	if err != nil {
		Usage("parse os.Stdin failed")
	}

	navigator := htmlquery.CreateXPathNavigator(doc)
	t := expr.Select(navigator)
	for t.MoveNext() {
		curr := t.Current().Copy()
		nav := curr.(*htmlquery.NodeNavigator)
		result := map[string]string{}
		for k, v := range target {
			result[k] = ""
			switch x := v.Evaluate(nav).(type) {
			case bool, float64, string:
				result[k] = fmt.Sprintf("%v", x)
			case *xpath.NodeIterator:
				if x.MoveNext() {
					result[k] = x.Current().(*htmlquery.NodeNavigator).Value()
				}
			default:
				panic(fmt.Sprintf("%t", x))
			}
		}
		d, _ := json.Marshal(result)
		fmt.Println(string(d))
	}
}
