package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"os"
)

const (
	pattern = "[\u4e00-\u9fa5]+"
	reg     = "[^\u4e00-\u9fa5]"
)

var (
	source    string
	export    bool
	fileNames = make([]string, 0)
	golds     = make([]*rs, 0)
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[FILE-FIND] ")
}

func main() {
	root := cobra.Command{
		Use:  "filefind",
		Long: "read and filter all files from folder",
	}

	root.PersistentFlags().StringVarP(&source, "source", "s", "", "need find files from folder")
	root.PersistentFlags().BoolVarP(&export, "export", "e", false, "if export result to json file")

	root.RunE = func(cmd *cobra.Command, args []string) error {
		log.Print(color.CyanString("read source folder: %v", source))
		if len(source) == 0 {
			return errors.New("folder is nil Use --help for help")
		}

		DeepSearchFiles(source)
		if err := do(); err != nil {
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
