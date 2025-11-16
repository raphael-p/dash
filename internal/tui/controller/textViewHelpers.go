package controller

import "github.com/gdamore/tcell/v2"

type inputCapture struct {
	runeHotkeys  []rune
	eventHotkeys []tcell.Key
	action       func(r rune)
}

type InputCaptureFactory struct {
	inputCaptures []inputCapture
}

func NewInputCaptureFactory(allowNumericInput bool) *InputCaptureFactory {
	return &InputCaptureFactory{[]inputCapture{}}
}

func (icf *InputCaptureFactory) registerInputCature(ic inputCapture) {
	icf.inputCaptures = append(icf.inputCaptures, ic)
}

func (icf *InputCaptureFactory) make() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		e := event.Key()
		r := event.Rune()

		// mouse scroll
		if e == tcell.KeyDown || e == tcell.KeyUp {
			return event
		}

		for _, ic := range icf.inputCaptures {
			for _, eh := range ic.eventHotkeys {
				if eh == e {
					ic.action(r)
					return nil
				}
			}
			for _, rh := range ic.runeHotkeys {
				if rh == r {
					ic.action(r)
					return nil
				}
			}
		}
		return nil
	}
}
