package main

import (
	"errors"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	// folder 指定工作目录
	folder string
	// export 指定导出目录
	export bool
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[FILE-FIND] ")
}

func main() {
	rootCmd := cobra.Command{
		Use: "filectl",
	}
	rootCmd.PersistentFlags().StringVarP(&folder, "folder", "f", "", "need find files from folder")

	scanCmd := &cobra.Command{
		Use:     "scan",
		Aliases: []string{"run", "find"},
		Short:   "Run file find run cmd",
		Long:    "read and filter all files from folder",
		// Args:    cobra.NoArgs,
		Version: "v1",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Print(color.CyanString("read source folder: %v", folder))
			if len(folder) == 0 {
				return errors.New("folder is nil Use --help for help")
			}

			DeepSearchFiles(folder)
			if err := do(); err != nil {
				return err
			}

			if export {
				if err := exportToJson(); err != nil {
					return err
				}
			}
			return nil
		},
	}

	// PersistentFlags和Flags的区别是PersistentFlags定义的选项可以被传递到子命令中
	scanCmd.Flags().BoolVarP(&export, "export", "e", false, "if export result to json file")

	rootCmd.AddCommand(scanCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
