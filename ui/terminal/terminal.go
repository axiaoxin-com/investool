// 终端操作界面

package terminal

import "github.com/rivo/tview"

// Run 启动终端操作界面
func Run() {
	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}
