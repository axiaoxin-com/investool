// Package statics embed 静态文件
package statics

import "embed"

// Files 静态文件资源
//go:embed files/*
var Files embed.FS
