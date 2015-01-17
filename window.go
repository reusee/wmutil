package wmutil

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

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
	atomDeleteWindow := w.wm.Atoms["WM_DELETE_WINDOW"]
	for _, atom := range w.Protocols {
		if atom == atomDeleteWindow {
			// send message
			msg := xproto.ClientMessageEvent{
				Format: 32,
				Window: w.Id,
				Type:   w.wm.Atoms["WM_PROTOCOLS"],
				Data: xproto.ClientMessageDataUnionData32New([]uint32{
					uint32(atomDeleteWindow),
					0, 0, 0, 0, // must be 20-bytes long
				}),
			}
			if err := xproto.SendEventChecked(w.wm.Conn, false, w.Id, xproto.EventMaskNoEvent, string(msg.Bytes())).Check(); err != nil {
				w.wm.pt("ERROR: send client message: %v\n", err)
			}
			return
		}
	}
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

func (w *Window) GetStrsProperty(atom xproto.Atom) (ret []string) {
	reply, err := xproto.GetProperty(w.wm.Conn, false, w.Id, atom, xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		w.wm.pt("ERROR: get window property: %v\n", err)
		return
	}
	start := 0
	for i, c := range reply.Value {
		if c == 0 {
			ret = append(ret, string(reply.Value[start:i]))
			start = i + 1
		}
	}
	if start < int(reply.ValueLen) {
		ret = append(ret, string(reply.Value[start:]))
	}
	return
}

func (w *Window) GetWindowIdProperty(atom xproto.Atom) xproto.Window {
	reply, err := xproto.GetProperty(w.wm.Conn, false, w.Id, atom, xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		w.wm.pt("ERROR: get window property: %v\n", err)
		return 0
	}
	if len(reply.Value) == 0 {
		return 0
	}
	return xproto.Window(xgb.Get32(reply.Value))
}

func (w *Window) GetAtomsProperty(atom xproto.Atom) (ret []xproto.Atom) {
	reply, err := xproto.GetProperty(w.wm.Conn, false, w.Id, atom, xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		w.wm.pt("ERROR: get window property: %v\n", err)
		return
	}
	for i := uint32(0); i < reply.ValueLen; i += 4 {
		ret = append(ret, xproto.Atom(xgb.Get32(reply.Value[i*4:])))
	}
	return
}

func (w *Window) ChangeInt32sProperty(atom, what xproto.Atom, ints ...uint32) {
	buf := make([]byte, len(ints)*4)
	for i, integer := range ints {
		xgb.Put32(buf[i*4:], integer)
	}
	err := xproto.ChangePropertyChecked(w.wm.Conn, xproto.PropModeReplace, w.Id, atom, what,
		32, uint32(len(buf)/4), buf).Check()
	if err != nil {
		w.wm.pt("ERROR: change window ints property: %v\n", err)
	}
}
