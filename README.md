## FileFind

FileFind是一款文件查找工具，你可以指定文件查找目录及相应的正则，FileFind会查找该目录下的所有文件，按照您提供的正则，进行报告输出。
目前一个使用场景是查找项目中使用到的中文编码（非注释）。

### 编译命令行工具
- 编译为二进制 `go build -o filectl .`，注意这里已经不是单文件程序了，直接使用`go run main.go` 会报错 
- 查看命令行帮助 `./filectl -h`
```shell
Usage:
  filectl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  scan        Run file find run cmd

Flags:
  -f, --folder string   need find files from folder
  -h, --help            help for filectl

Use "filectl [command] --help" for more information about a command.
```
- 查看子命令帮助`./filectl  scan -h`
```shell
read and filter all files from folder

Usage:
  filectl scan [flags]

Aliases:
  scan, run, find

Flags:
  -e, --export    if export result to json file
  -h, --help      help for scan
  -v, --version   version for scan

Global Flags:
  -f, --folder string   need find files from folder
```
### 查找指定文件夹
```shell
 ~/workspace/go-workspace/filefind/ [master+*] ./filectl -f ./example scan -e true
[FILE-FIND] read source folder: ./example
[FILE-FIND] {
    "path": "./example/main.go",
    "values": [
        {
            "line": 9,
            "val": "您好这里我使用了中文"
        },
        {
            "line": 11,
            "val": "这里不确定能不能被工具发现呢"
        }
    ]
}
 ~/workspace/go-workspace/filefind/ [master+*] 
```
### 结果输出
```shell
 ~/workspace/go-workspace/filefind/ [master+*] tree                     
.
├── README.md
├── example
│   └── main.go
├── export.go
├── file.go
├── filectl
├── go.mod
├── go.sum
├── main.go
├── result.json
└── search.go

1 directory, 10 files
 ~/workspace/go-workspace/filefind/ [master+*] cat result.json                   
[{"path":"./example/main.go","values":[{"line":9,"val":"您好这里我使用了中文"},{"line":11,"val":"这里不确定能不能被工具发现呢"}]}]
 ~/workspace/go-workspace/filefind/ [master+*] 
```