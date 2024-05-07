<h1>dual-center-switch-system(dcss) 后端</h1>

## 介绍

双中心切换系统WEB后端

## 用法

### 开发模式启动

go run main.go

### 静态编译(需要安装upx和go)
>
 编译前需要将前端的`dist/*`文件 复制到`embed/dist/`文件夹下

#### Linux

```bash
bash build.sh
```

#### Window

```bash
build.bat
```

### 错误记录

- Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
solution:
MinGW-w64 - for 32 and 64 bit Windows Files
go env -w CGO_ENABLED=1
