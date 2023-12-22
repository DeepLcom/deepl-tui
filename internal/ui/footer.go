package ui

import (
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-shellwords"
	"github.com/rivo/tview"
)

func (ui *UI) setupFooter() {
	cmdline := tview.NewInputField()

	cmdline.
		SetFieldStyle(
			tcell.StyleDefault.
				Background(tview.Styles.PrimitiveBackgroundColor).
				Foreground(tview.Styles.PrimaryTextColor),
		).
		SetLabelStyle(
			tcell.StyleDefault,
		).
		SetBorder(true)

	cmdline.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyBacktab:
			// ignore backtab, else it would finish editing
			return nil
		case tcell.KeyTab:
			// remap key, else it would finish editing
			cmdline.Autocomplete()
			return nil
		}

		return event
	})

	cmdline.SetDoneFunc(func(key tcell.Key) {
		var err error
		defer func() {
			cmdline.
				SetLabel("").
				SetDisabled(true)

			if err != nil {
				cmdline.SetText(err.Error())
			}

			ui.SetFocus(ui.pages)
		}()

		text := cmdline.GetText()
		cmdline.SetText("")

		if key != tcell.KeyEnter {
			return
		}

		args, err := shellwords.Parse(text)
		if err != nil {
			return
		}

		if len(args) < 1 {
			return
		} else if len(args) > 1 {
			err = errors.New("invalid command")
			return
		}

		switch args[0] {
		case "translate", "glossaries":
			ui.switchToPage(args[0])
		case "size":
			_, _, w, h := ui.translatePage.GetInnerRect()
			err = errors.New(fmt.Sprintf("(%d, %d)", w, h))
		default:
			err = errors.New("invalid command")
		}
	})

	ui.footer = cmdline
}
