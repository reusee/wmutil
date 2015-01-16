package wmutil

import "github.com/BurntSushi/xgb/xproto"

func (w *Window) SetPos(x, y int) {
	if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
		xproto.ConfigWindowX|xproto.ConfigWindowY, []uint32{uint32(x), uint32(y)}).Check(); err != nil {
		w.wm.pt("ERROR: set window position: %v", err)
	}
}

func (w *Window) SetSize(width, height int) {
	if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
		xproto.ConfigWindowWidth|xproto.ConfigWindowHeight, []uint32{uint32(width), uint32(height)}).Check(); err != nil {
		w.wm.pt("ERROR: set window size: %v", err)
	}
}

func (w *Window) SetGeometry(x, y, width, height int) {
	if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
		xproto.ConfigWindowX|xproto.ConfigWindowY|xproto.ConfigWindowWidth|xproto.ConfigWindowHeight,
		[]uint32{uint32(x), uint32(y), uint32(width), uint32(height)}).Check(); err != nil {
		w.wm.pt("ERROR: set window geometry: %v", err)
	}
}

func (w *Window) Above(sibling *Window) {
	if sibling != nil {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowSibling|xproto.ConfigWindowStackMode,
			[]uint32{uint32(sibling.Id), uint32(xproto.StackModeAbove)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v", err)
		}
	} else {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowStackMode, []uint32{uint32(xproto.StackModeAbove)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v", err)
		}
	}
}

func (w *Window) Below(sibling *Window) {
	if sibling != nil {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowSibling|xproto.ConfigWindowStackMode,
			[]uint32{uint32(sibling.Id), uint32(xproto.StackModeBelow)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v", err)
		}
	} else {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowStackMode, []uint32{uint32(xproto.StackModeBelow)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v", err)
		}
	}
}

func (w *Window) TopIf(sibling *Window) {
	if sibling != nil {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowSibling|xproto.ConfigWindowStackMode,
			[]uint32{uint32(sibling.Id), uint32(xproto.StackModeTopIf)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v", err)
		}
	} else {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowStackMode, []uint32{uint32(xproto.StackModeTopIf)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v", err)
		}
	}
}

func (w *Window) BottomIf(sibling *Window) {
	if sibling != nil {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowSibling|xproto.ConfigWindowStackMode,
			[]uint32{uint32(sibling.Id), uint32(xproto.StackModeBottomIf)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v", err)
		}
	} else {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowStackMode, []uint32{uint32(xproto.StackModeBottomIf)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v", err)
		}
	}
}

func (w *Window) Opposite(sibling *Window) {
	if sibling != nil {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowSibling|xproto.ConfigWindowStackMode,
			[]uint32{uint32(sibling.Id), uint32(xproto.StackModeOpposite)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v", err)
		}
	} else {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowStackMode, []uint32{uint32(xproto.StackModeOpposite)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v", err)
		}
	}
}

func (w *Window) Destroy() {
	if err := xproto.DestroyWindowChecked(w.wm.Conn, w.Id).Check(); err != nil {
		w.wm.pt("ERROR: destroy window: %v", err)
	}
}
