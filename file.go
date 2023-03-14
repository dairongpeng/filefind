package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	pattern = "[\u4e00-\u9fa5]+"
	reg     = "[^\u4e00-\u9fa5]"
)

var (
	fileNames = make([]string, 0)
	golds     = make([]*rs, 0)
)

func readFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	reader := bufio.NewReader(file)

	fLine := 0
	results := &rs{
		Path:   path,
		Values: make([]*V, 0),
	}
	for {
		fLine++
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		// 去掉空白符
		l := strings.ReplaceAll(string(line), "\t", "")
		isMatch, err := regexp.MatchString(pattern, l)
		if err != nil {
			return err
		}

		if isMatch { // 匹配到中文，且不是注释
			if !strings.Contains(l, "//") && !strings.Contains(l, "/*") && !strings.Contains(l, "*/") {
				regExp := regexp.MustCompile(reg)
				str := regExp.ReplaceAllString(l, "")
				val := &V{
					Line: fLine,
					Val:  str,
				}
				results.Values = append(results.Values, val)
			}
		}
	}
	log.Println(color.CyanString("%s", StringResults(results)))
	if len(results.Values) > 0 {
		golds = append(golds, results)
	}
	return nil
}

type rs struct {
	Path   string `json:"path"`
	Values []*V   `json:"values"`
}

type V struct {
	Line int    `json:"line"`
	Val  string `json:"val"`
}

func StringResults(rs *rs) string {
	b, err := json.Marshal(rs)
	if err != nil {
		return fmt.Sprintf("%+v", rs)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", rs)
	}
	return out.String()
}
