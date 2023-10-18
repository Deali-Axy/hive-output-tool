# hive-output-tool

hive导出工具

## build

build on linux, run on linux

```bash
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build
```

build on windows, run on linux

```powershell
go env -w CGO_ENABLED=0
go env -w GOARCH=amd64
go env -w GOOS=linux
go build
```

## deploy

需要将 `.env` 和 `hive-output-tool` 可执行文件夹放在同一个目录

```bash
tar czvf hive-tools.tar.gz .
```

## updates

- 0.2 - 允许输入的SQL中包含换行
- 0.3 - 修复输入 q 无法退出的问题

