# 使用
```bash
go mod tidy
go build
# 使用代理
./m3u8Find -p http://localhost:1082
# 输出m3u8 文件
./m3u8Find -p http://localhost:1082 -om
# 输出m3u8到一个文件夹
./m3u8Find -p http://localhost:1082 -om -oa
# 输出csv
./m3u8Find -p http://localhost:1082 -oaa
# 查看更多帮助
./m3u8Find -h
```