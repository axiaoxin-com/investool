#!/bin/bash
go mod tidy && \
go test -race ./... && \
rm -rf ./statics/js/app.min.*.js && \
bash ./misc/scripts/app.min.js.sh && \
bash ./misc/scripts/bumpversion.sh && \
qshell cdnprefetch -i /tmp/qiniu-prefetch-refresh-list.txt && \
qshell cdnrefresh -i /tmp/qiniu-prefetch-refresh-list.txt && \
qshell batchdelete blogpostpics -y -i /tmp/qiniu-del-list.txt && \
goreleaser release --rm-dist  && \
git push && git push --tags -f
