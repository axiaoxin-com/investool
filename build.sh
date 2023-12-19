CGO_ENABLED=0 go build -ldflags "-X github.com/axiaoxin-com/investool/version.Version=`TZ=Asia/Shanghai date +%y%m%d%H%M`"
