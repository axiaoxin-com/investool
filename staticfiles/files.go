// Package staticfiles embed 静态文件
package staticfiles

import "embed"

// FS 静态文件资源
//go:embed font.ttf
var FS embed.FS
