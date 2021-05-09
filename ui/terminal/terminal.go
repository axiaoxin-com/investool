// 终端操作界面

package terminal

import "github.com/rivo/tview"

// Run 启动终端操作界面
func Run() error {
	gridLayout := GridLayout()
	app := tview.NewApplication()
	if err := app.SetRoot(gridLayout, true).Run(); err != nil {
		return err
	}
	return nil
}

func GridLayout() *tview.Grid {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	filterForm := FilterForm()
	main := newPrimitive("Main content")
	checkerOptionsForm := CheckerOptionsForm()

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("X-STOCK"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (filterForm and side bar are hidden).
	grid.AddItem(filterForm, 0, 0, 0, 0, 0, 0, false).
		AddItem(main, 1, 0, 1, 3, 0, 0, false).
		AddItem(checkerOptionsForm, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(filterForm, 1, 0, 1, 1, 0, 100, false).
		AddItem(main, 1, 1, 1, 1, 0, 100, false).
		AddItem(checkerOptionsForm, 1, 2, 1, 1, 0, 100, false)
	return grid
}

// FilterForm 筛选条件表单
func FilterForm() *tview.Form {
	form := tview.NewForm().
		AddInputField("ROE(%)", "8.0", 20, nil, nil).
		AddButton("重置", nil)

	form.SetBorder(true).SetTitle("筛选条件").SetTitleAlign(tview.AlignLeft)
	return form
}

// CheckerOptionsForm 检测条件表单
func CheckerOptionsForm() *tview.Form {
	form := tview.NewForm().
		AddInputField("最低 ROE(%)", "8.0", 20, nil, nil).
		AddButton("重置", nil)

	form.SetBorder(true).SetTitle("检测条件").SetTitleAlign(tview.AlignLeft)
	return form
}
