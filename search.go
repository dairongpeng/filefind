package main

import (
	"io/ioutil"
	"strings"
)

// DeepSearchFiles 获取文件夹下所有文件
func DeepSearchFiles(folder string) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			DeepSearchFiles(folder + "/" + file.Name())
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
