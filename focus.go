package eui

type FocusManager struct{ hovered, focused Drawabler }

func (fm *FocusManager) SetHovered(w Drawabler) {
	if fm.hovered != w {
		if fm.hovered != nil {
			fm.hovered.SetState(StateNormal)
		}
	}
	fm.hovered = w
	if fm.hovered != nil {
		fm.hovered.SetState(StateHover)
	}
}
func (fm *FocusManager) SetFocused(w Drawabler) {
	if fm.focused == w {
		return
	}
	if fm.focused != nil {
		fm.focused.SetState(StateNormal)
	}
	fm.hovered = nil
	fm.focused = w
	if w != nil {
		w.SetState(StateFocused)
	}
}
func (fm *FocusManager) Blur() {
	if fm.focused != nil {
		if fm.focused.State().IsFocused() {
			fm.focused.SetState(StateNormal)
		}
		fm.focused = nil
	}
	if fm.hovered != nil {
		if fm.hovered.State().IsHovered() {
			fm.hovered.SetState(StateNormal)
		}
		fm.hovered = nil
	}
}
func (fm *FocusManager) Hovered() Drawabler { return fm.hovered }
func (fm *FocusManager) Focused() Drawabler { return fm.focused }
