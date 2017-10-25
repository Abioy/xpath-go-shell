package main

import (
	"encoding/json"
	"fmt"
	//"launchpad.net/xmlpath"
	xmlpath "gopkg.in/xmlpath.v2"
	"os"
)

type Target map[string]*xmlpath.Path

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

func Usage(errs ...string) {
	rc := 0
	var msg string
	for _, err := range errs {
		msg += "[ERROR] "
		msg += err
		msg += "\n"
		rc = 1
	}
	msg += usage()
	fmt.Fprint(os.Stderr, msg)
	os.Exit(rc)
}

func main() {
	args := os.Args
	if args == nil || len(args) < 3 {
		Usage("len(args) < 3")
	}
	pattern := args[1]
	targets_str := args[2]
	path, err := xmlpath.Compile(pattern)
	if err != nil {
		Usage("<PATH> illegal")
	}
	var data map[string]string
	if json.Unmarshal([]byte(targets_str), &data) != nil {
		Usage("<TARGET> illegal")
	}
	target := Target{}
	for k, v := range data {
		target[k], err = xmlpath.Compile(v)
		if err != nil {
			Usage("<TARGET> path illegal")
		}
	}

	file := os.Stdin
	root, err := xmlpath.ParseHTML(file)
	if err != nil {
		Usage("parse os.Stdin failed")
	}

	it := path.Iter(root)
	for it.Next() {
		result := map[string]string{}
		for k, v := range target {
			sub_it := v.Iter(it.Node())
			result[k] = ""
			if sub_it.Next() {
				result[k] = sub_it.Node().String()
			}
		}
		d, _ := json.Marshal(result)
		fmt.Println(string(d))
	}
}
