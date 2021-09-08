#!/bin/bash
# 依赖：
# npm install uglify-js -g
# go get -u github.com/swaggo/swag/cmd/swag
# go get github.com/x-motemen/gobump/cmd/gobump
# qshell https://github.com/qiniu/qshell/releases/download/v2.6.2/qshell-v2.6.2-linux-amd64.tar.gz
# qshell account <Your AccessKey> <Your SecretKey> <Your Name>
go mod tidy && \
go test -race ./... && \
rm -rf ./statics/js/app.min.*.js && \
bash ./misc/scripts/app.min.js.sh && \
bash ./misc/scripts/bumpversion.sh && \
echo "push..."
git push && git push --tags -f && \
echo "release..."
goreleaser release --rm-dist
echo "cdn清理..."
qshell batchdelete blogpostpics -y -i /tmp/qiniu-del-list.txt
