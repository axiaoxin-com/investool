#!/bin/bash
go mod tidy && \
go test -race ./... && \
bash ./misc/scripts/app.min.js.sh && \
bash ./misc/scripts/bumpversion.sh && \
goreleaser release --rm-dist  && \
qshell batchdelete blogpostpics -y -i /tmp/qiniu-del-list.txt
