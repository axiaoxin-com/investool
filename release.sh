#!/bin/bash
go mod tidy && \
go test -race ./... && \
rm -rf ./statics/js/app.min.*.js && \
bash ./misc/scripts/app.min.js.sh && \
bash ./misc/scripts/bumpversion.sh && \
echo "release..."
goreleaser release --rm-dist  && \
echo "push..."
git push && git push --tags -f
