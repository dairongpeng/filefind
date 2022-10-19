package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	pattern = "[\u4e00-\u9fa5]+"
	reg =  "[^\u4e00-\u9fa5]"
)

var (
	folder    string
	export    bool
	fileNames = make([]string, 0)
	golds     = make([]*rs, 0)
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[FILE-UTIL] ")
}

func main() {
	root := cobra.Command{
		Use:  "read-file",
		Long: "read all files from folder",
	}

	root.PersistentFlags().StringVar(&folder, "folder", "", "need find files from folder")
	root.PersistentFlags().BoolVar(&export, "export", false, "if export result to json file")

	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	root.RunE = func(cmd *cobra.Command, args []string) error {
		log.Print(color.CyanString("folder: %v", folder))
		if len(folder) == 0 {
			return errors.New("folder is nil Use --help for help")
		}
		GetFiles(folder)
		if er := do(); err != nil {
			return err
		}

		if export {
			if err := exportToJson(); err != nil {
				return err
			}
		}
		return nil
	}

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

// GetFiles 获取文件夹下所有文件
func GetFiles(folder string) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			GetFiles(folder + "/" + file.Name())
		} else {
			if strings.HasSuffix(file.Name(), ".go") { // 只查go源文件
				fileNames = append(fileNames, folder+"/"+file.Name())
			}
		}
	}
}

func do() error {
	if len(fileNames) == 0 {
		return nil
	}

	for _, fileName := range fileNames {
		err := readFile(fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

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

func exportToJson() error {
	if len(golds) == 0 {
		return nil
	}

	file, err := os.Create("result.json")
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(golds)
	if err != nil {
		return err
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
r
