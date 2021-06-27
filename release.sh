#!/bin/bash
echo "Realeasing ..."
goreleaser release --rm-dist
qshell batchdelete blogpostpics -y -i /tmp/qiniu-prefetch-del-list.txt
