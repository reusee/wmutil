package wmutil

import "github.com/BurntSushi/xgb/xproto"

func (w *Window) ReadLock(fn func()) {
	w.RLock()
	fn()
	w.RUnlock()
}

func (w *Window) WriteLock(fn func()) {
	w.Lock()
	fn()
	w.Unlock()
}

func (w *Window) SetPos(x, y int) {
	if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
		xproto.ConfigWindowX|xproto.ConfigWindowY, []uint32{uint32(x), uint32(y)}).Check(); err != nil {
		w.wm.pt("ERROR: set window position: %v\n", err)
	} else {
		w.WriteLock(func() {
			w.X, w.Y = x, y
		})
	}
}

func (w *Window) SetSize(width, height int) {
	if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
		xproto.ConfigWindowWidth|xproto.ConfigWindowHeight, []uint32{uint32(width), uint32(height)}).Check(); err != nil {
		w.wm.pt("ERROR: set window size: %v\n", err)
	} else {
		w.WriteLock(func() {
			w.Width, w.Height = width, height
		})
	}
}

func (w *Window) SetGeometry(x, y, width, height int) {
	if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
		xproto.ConfigWindowX|xproto.ConfigWindowY|xproto.ConfigWindowWidth|xproto.ConfigWindowHeight,
		[]uint32{uint32(x), uint32(y), uint32(width), uint32(height)}).Check(); err != nil {
		w.wm.pt("ERROR: set window geometry: %v\n", err)
	} else {
		w.WriteLock(func() {
			w.X, w.Y, w.Width, w.Height = x, y, width, height
		})
	}
}

func (w *Window) Above(sibling *Window) {
	if sibling != nil {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowSibling|xproto.ConfigWindowStackMode,
			[]uint32{uint32(sibling.Id), uint32(xproto.StackModeAbove)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v\n", err)
		}
	} else {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowStackMode, []uint32{uint32(xproto.StackModeAbove)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v\n", err)
		}
	}
}

func (w *Window) Below(sibling *Window) {
	if sibling != nil {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowSibling|xproto.ConfigWindowStackMode,
			[]uint32{uint32(sibling.Id), uint32(xproto.StackModeBelow)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v\n", err)
		}
	} else {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowStackMode, []uint32{uint32(xproto.StackModeBelow)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v\n", err)
		}
	}
}

func (w *Window) TopIf(sibling *Window) {
	if sibling != nil {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowSibling|xproto.ConfigWindowStackMode,
			[]uint32{uint32(sibling.Id), uint32(xproto.StackModeTopIf)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v\n", err)
		}
	} else {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowStackMode, []uint32{uint32(xproto.StackModeTopIf)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v\n", err)
		}
	}
}

func (w *Window) BottomIf(sibling *Window) {
	if sibling != nil {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowSibling|xproto.ConfigWindowStackMode,
			[]uint32{uint32(sibling.Id), uint32(xproto.StackModeBottomIf)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v\n", err)
		}
	} else {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowStackMode, []uint32{uint32(xproto.StackModeBottomIf)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v\n", err)
		}
	}
}

func (w *Window) Opposite(sibling *Window) {
	if sibling != nil {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowSibling|xproto.ConfigWindowStackMode,
			[]uint32{uint32(sibling.Id), uint32(xproto.StackModeOpposite)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v\n", err)
		}
	} else {
		if err := xproto.ConfigureWindowChecked(w.wm.Conn, w.Id,
			xproto.ConfigWindowStackMode, []uint32{uint32(xproto.StackModeOpposite)}).Check(); err != nil {
			w.wm.pt("ERROR: set window above: %v\n", err)
		}
	}
}

func (w *Window) Destroy() {
	if err := xproto.DestroyWindowChecked(w.wm.Conn, w.Id).Check(); err != nil {
		w.wm.pt("ERROR: destroy window: %v\n", err)
	}
}

func (w *Window) WarpPointer() {
	if err := xproto.WarpPointerChecked(w.wm.Conn, 0, w.Id, 0, 0, 0, 0, 0, 0).Check(); err != nil {
		w.wm.pt("ERROR: warp pointer: %v\n", err)
	}
}

func (wm *Wm) PointingWindow() *Window {
	reply, err := xproto.QueryPointer(wm.Conn, wm.DefaultRootId).Reply()
	if err != nil {
		wm.pt("ERROR: query pointer: %v\n", err)
	}
	win, ok := wm.Windows[reply.Child]
	if !ok {
		return nil
	}
	return win
}

func (w *Wm) FocusPointerRoot() {
	if err := xproto.SetInputFocusChecked(w.Conn, 0, xproto.InputFocusPointerRoot, 0).Check(); err != nil {
		w.pt("ERROR: set focus to pointer root %v", err)
	}
}
