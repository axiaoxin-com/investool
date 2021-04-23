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
	"path/filepath"
	"strings"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/staticfiles"
	"github.com/golang/freetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
)

// ExportPic 导出股票名称+代码图片，一张图片最多 50 个，超过则导出多张图片
func (e Exportor) ExportPic(ctx context.Context, filename string) (result []byte, err error) {
	height := 3200
	width := 650

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
	fontSize := float64(15)
	fc := freetype.NewContext()
	fc.SetDst(img)
	fc.SetFont(ffont)
	fc.SetClip(img.Bounds())
	fc.SetFontSize(fontSize)
	fc.SetSrc(fgColor)
	fc.SetDPI(300)
	fc.SetHinting(font.HintingNone)

	// 按分组写入股票名称+代码到不同图片
	for i, stocks := range e.Stocks.ChunkedBySize(50) {
		for j, stock := range stocks {
			pt := freetype.Pt(40, (j+1)*int(fc.PointToFixed(fontSize)>>6)+40)
			line := fmt.Sprintf("%d.%s    %s", i+1, stock.Name, stock.Code)
			_, err = fc.DrawString(line, pt)
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
		oriFilename := filename
		if i > 0 {
			dir, base := filepath.Split(filename)
			ext := filepath.Ext(base)
			fn := strings.TrimSuffix(base, ext)

			base = fmt.Sprint(fn, "_", i, ext)
			filename = filepath.Join(dir, base)
		}
		err = ioutil.WriteFile(filename, result, 0666)
		err = errors.Wrap(err, "write file error")
		filename = oriFilename
	}
	return
}
