# README

Windows下编译Linux平台64位可执行程序：

```
SET CGO_ENABLED=0 // 跨平台编译默认禁用
SET GOOS=linux
SET GOARCH=amd64 // 默认当前架构
go build
```

