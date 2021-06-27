#!/bin/bash
go mod tidy && \
go test -race ./... && \
bash ./misc/scripts/app.min.js.sh && \
bash ./misc/scripts/bumpversion.sh && \
goreleaser release --rm-dist  && \
git push --tags && \
qshell batchdelete blogpostpics -y -i /tmp/qiniu-prefetch-del-list.txt
