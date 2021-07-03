#!/bin/bash
go mod tidy && \
go test -race ./... && \
rm -rf ./statics/js/app.min.*.js && \
bash ./misc/scripts/app.min.js.sh && \
bash ./misc/scripts/bumpversion.sh && \
echo "cdn预取..."
qshell cdnprefetch -i /tmp/qiniu-prefetch-refresh-list.txt && \
echo "cdn刷新..."
qshell cdnrefresh -i /tmp/qiniu-prefetch-refresh-list.txt && \
echo "cdn清理..."
qshell batchdelete blogpostpics -y -i /tmp/qiniu-del-list.txt && \
echo "release..."
goreleaser release --rm-dist  && \
echo "push..."
git push && git push --tags -f
