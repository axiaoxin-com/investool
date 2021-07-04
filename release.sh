#!/bin/bash
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
