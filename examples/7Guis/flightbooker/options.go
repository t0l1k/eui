package flightbooker

import (
	"github.com/t0l1k/eui"
)

type OptionButton struct {
	*eui.Button
	items    []string
	selected int // Index of selected item
	onSelect *eui.Signal[*OptionButton]
}

func NewOptionButton(items []string, initial int, fn func(*OptionButton)) *OptionButton {
	ob := &OptionButton{
		items:    items,
		selected: initial,
		onSelect: eui.NewSignal(func(a, b *OptionButton) bool { return false }),
	}
	ob.Button = eui.NewButton(ob.SelectedText(), func(b *eui.Button) { ob.showPopup() })
	ob.onSelect.Connect(fn)
	return ob
}

func (ob *OptionButton) SelectedIndex() int { return ob.selected }

func (ob *OptionButton) SelectedText() string {
	if ob.selected >= 0 && ob.selected < len(ob.items) {
		return ob.items[ob.selected]
	}
	return "Select Option"
}

func (ob *OptionButton) Select(index int) {
	if index >= 0 && index < len(ob.items) {
		ob.selected = index
		ob.updateText()
		if ob.onSelect != nil {
			ob.onSelect.Emit(ob)
		}
	}
}

func (ob *OptionButton) updateText() {
	txt := ""
	if ob.selected >= 0 && ob.selected < len(ob.items) {
		txt = ob.items[ob.selected]
	}
	ob.Button.SetText(txt)
}

func (ob *OptionButton) showPopup() {
	menu := newPopupMenu(ob)
	menu.SetViewType(eui.ViewModal)
	eui.GetUi().ShowModal(menu)
}

// --- PopupMenu: The modal list ---

type popupMenu struct {
	*eui.Container
	owner    *OptionButton
	resizeId int64
}

func newPopupMenu(owner *OptionButton) *popupMenu {
	m := &popupMenu{
		Container: eui.NewContainer(eui.NewVBoxLayout(1)),
		owner:     owner,
	}

	// Populate items
	for i, it := range owner.items {
		idx := i // capture loop variable
		text := it
		btn := eui.NewButton(text, func(bb *eui.Button) {
			owner.Select(idx)
			eui.GetUi().HideModal()
		})
		m.Add(btn)
	}

	// Initial positioning
	m.updateRect()

	// Listen for tick to keep position relative to owner button (handles resize and layout changes)
	m.resizeId = eui.GetUi().TickListener().Connect(func(ev eui.Event) {
		m.updateRect()
	})

	return m
}

func (m *popupMenu) Close() {
	if m.resizeId != 0 {
		eui.GetUi().TickListener().Disconnect(m.resizeId)
		m.resizeId = 0
	}
}

// Hit ensures the modal captures mouse events (blocking underlying UI)
func (m *popupMenu) Hit(pt eui.Point[int]) eui.Drawabler {
	if m.IsHidden() {
		return nil
	}
	return m
}

func (m *popupMenu) MouseUp(md eui.MouseData) {
	p := md.Pos()
	if !p.In(m.Rect()) {
		eui.GetUi().HideModal()
		return
	}
	// Propagate click to children (menu items)
	for _, ch := range m.Children() {
		if mh, ok := ch.(interface {
			Hit(eui.Point[int]) eui.Drawabler
		}); ok {
			if target := mh.Hit(p); target != nil {
				if btn, ok := ch.(*eui.Button); ok {
					btn.MouseUp(md) // Trigger selection
				}
				return
			}
		}
	}
}

func (m *popupMenu) updateRect() {
	w0, h0 := eui.GetUi().Size()
	itemH := 28 // Approximate height of a button item
	padding := 4

	// Calculate required size
	mw := m.owner.Button.Rect().W // Match button width at minimum
	if mw < 100 {
		mw = 100
	}
	mh := padding*2 + len(m.owner.items)*itemH

	// Target position relative to owner button
	tRect := m.owner.Button.Rect()
	var mx, my int

	if !tRect.IsEmpty() {
		mx = tRect.X
		my = tRect.Y + tRect.H // Below button

		// Check bottom overflow
		if my+mh > h0 {
			my = tRect.Y - mh // Show above if no space below
		}
		// Check right overflow
		if mx+mw > w0 {
			mx = w0 - mw - 5
		}
		if mx < 0 {
			mx = 5
		}
	} else {
		// Fallback center
		mx = (w0 - mw) / 2
		my = (h0 - mh) / 2
	}
	m.SetRect(eui.NewRect([]int{mx, my, mw, mh}))
}
