package ui

import (
	"github.com/rivo/tview"
)

type GlossariesDialog struct {
	tview.Modal

	glossariesDropDown *tview.DropDown
}

func newGlossariesPage(ui *UI) *GlossariesDialog {
	page := &GlossariesDialog{}

	page.glossariesDropDown = tview.NewDropDown()

	return page
}
