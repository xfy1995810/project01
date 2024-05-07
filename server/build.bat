rem set GOOS=linux
set GOOS=linux
go build -trimpath -ldflags "-w"
upx.exe -9 -k  dcss
rm -f dcss.*