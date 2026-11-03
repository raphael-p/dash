package tviewcomponents

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HideableFlex struct {
	*tview.Flex
	hidden bool
}

func NewHideableFlex() *HideableFlex {
	flex := tview.NewFlex()
	flex.Box = tview.NewBox()
	return &HideableFlex{flex, false}
}

func (hf *HideableFlex) Draw(screen tcell.Screen) {
	if hf.hidden {
		hf.Box.Draw(screen)
		screen.HideCursor()
		return
	}
	hf.Flex.Draw(screen)
}
