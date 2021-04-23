// 导出股票名称+代码图片

package exportor

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/staticfiles"
	"github.com/golang/freetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
)

// ExportPic 导出股票名称+代码图片
func (e Exportor) ExportPic(ctx context.Context, filename string) (result []byte, err error) {
	width := 500
	height := 36 * len(e.Stocks)

	leftTop := image.Point{0, 0}
	rightBottom := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{leftTop, rightBottom})

	bgColor, fgColor := image.White, image.Black
	draw.Draw(img, img.Bounds(), bgColor, image.ZP, draw.Src)
	// set font
	fontBytes, err := staticfiles.FS.ReadFile("font.ttf")
	if err != nil {
		err = errors.Wrap(err, "read embed file error")
		return
	}
	ffont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		err = errors.Wrap(err, "parse font error")
		return
	}
	fontSize := float64(18)
	fc := freetype.NewContext()
	fc.SetDst(img)
	fc.SetFont(ffont)
	fc.SetClip(img.Bounds())
	fc.SetFontSize(fontSize)
	fc.SetSrc(fgColor)
	fc.SetDPI(144)
	fc.SetHinting(font.HintingFull)

	// 写入股票名称+代码
	for i, stock := range e.Stocks {
		pt := freetype.Pt(50, (i+1)*int(fc.PointToFixed(fontSize)>>6))
		line := fmt.Sprintf("%d.%s    %s", i+1, stock.Name, stock.Code)
		_, err := fc.DrawString(line, pt)
		if err != nil {
			logging.Errorf(ctx, "draw %s error: %s", line, err.Error())
			continue
		}
	}

	// 生成图片
	picbuff := new(bytes.Buffer)
	err = png.Encode(picbuff, img)
	if err != nil {
		err = errors.Wrap(err, "png encode error")
		return
	}
	result = picbuff.Bytes()
	err = ioutil.WriteFile(filename, result, 0666)
	err = errors.Wrap(err, "write file error")
	return
}
