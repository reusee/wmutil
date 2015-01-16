package wmutil

//go:generate go run gen.go

import (
	"log"
	"os"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Wm struct {
	Conn          *xgb.Conn
	Setup         *xproto.SetupInfo
	DefaultScreen *xproto.ScreenInfo
	DefaultRootId xproto.Window
	Windows       map[xproto.Window]*Window
	CodeToSyms    [][]uint32
	SymToCodes    map[uint32][]byte

	logger *log.Logger

	MapWindow   chan *Window
	UnmapWindow chan *Window
	Strokes     chan Stroke
}

type Window struct {
	Id                          xproto.Window
	Parent                      xproto.Window
	X, Y, Width, Height, Border int
	Mapped                      bool
}

type Config struct {
	Logger  *log.Logger
	Strokes []Stroke
}

type Stroke struct {
	Modifiers uint16
	Sym       uint32
}

func New(config *Config) (*Wm, error) {
	if config == nil {
		config = new(Config)
	}

	// connect
	conn, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}

	// infos
	setup := xproto.Setup(conn)
	defaultScreen := setup.DefaultScreen(conn)
	defaultRootId := defaultScreen.Root

	// grab control
	if err := xproto.ChangeWindowAttributesChecked(conn, defaultRootId, xproto.CwEventMask, []uint32{uint32(
		xproto.EventMaskSubstructureRedirect |
			xproto.EventMaskSubstructureNotify |
			xproto.EventMaskPropertyChange)}).Check(); err != nil {
		return nil, ef("another wm is running: %v", err)
	}

	// read keyboard mapping
	kmReply, err := xproto.GetKeyboardMapping(conn, 8, 248).Reply()
	if err != nil {
		return nil, ef("get keyboard mapping: %v", err)
	}
	keycodeToKeysyms := make([][]uint32, 256)
	keysymToKeycodes := make(map[uint32][]byte)
	for keycode := 8; keycode <= 255; keycode++ {
		start := (keycode - 8) * int(kmReply.KeysymsPerKeycode)
		keysyms := kmReply.Keysyms[start : start+int(kmReply.KeysymsPerKeycode)]
		for _, sym := range keysyms {
			if sym == 0 {
				continue
			}
			keysym := uint32(sym)
			keycodeToKeysyms[keycode] = append(keycodeToKeysyms[keycode], uint32(keysym))
			keysymToKeycodes[keysym] = append(keysymToKeycodes[keysym], byte(keycode))
		}
	}

	// get modifier mask
	mmReply, err := xproto.GetModifierMapping(conn).Reply()
	if err != nil {
		return nil, ef("get modifier mapping: %v", err)
	}
	symToModMask := func(sym uint32) uint16 {
		code := keysymToKeycodes[sym][0]
		index := 0
	loop:
		for ; index < 8; index++ {
			start := index * int(mmReply.KeycodesPerModifier)
			codes := mmReply.Keycodes[start : start+int(mmReply.KeycodesPerModifier)]
			for _, c := range codes {
				if byte(c) == code {
					break loop
				}
			}
		}
		return []uint16{
			xproto.ModMaskShift,
			xproto.ModMaskLock,
			xproto.ModMaskControl,
			xproto.ModMask1,
			xproto.ModMask2,
			xproto.ModMask3,
			xproto.ModMask4,
			xproto.ModMask5,
		}[index]
	}
	numlockModMask := symToModMask(Key_Num_Lock)

	// grab keys
	if err := xproto.UngrabKeyChecked(conn, xproto.GrabAny, defaultRootId, xproto.ModMaskAny).Check(); err != nil {
		return nil, ef("ungrab keys: %v", err)
	}
	ignoreModifiers := []uint16{
		0,
		xproto.ModMaskLock,
		numlockModMask,
		xproto.ModMaskLock | numlockModMask,
	}
	for _, stroke := range config.Strokes {
		for _, mod := range ignoreModifiers {
			keycodes := keysymToKeycodes[stroke.Sym]
			for _, code := range keycodes {
				if err := xproto.GrabKeyChecked(conn, true, defaultRootId, stroke.Modifiers|mod,
					xproto.Keycode(code), xproto.GrabModeAsync, xproto.GrabModeAsync).Check(); err != nil {
					return nil, ef("grab key: %v", err)
				}
			}
		}
	}

	wm := &Wm{
		Conn:          conn,
		Setup:         setup,
		DefaultScreen: defaultScreen,
		DefaultRootId: defaultRootId,
		Windows:       make(map[xproto.Window]*Window),
		MapWindow:     make(chan *Window),
		UnmapWindow:   make(chan *Window),
		Strokes:       make(chan Stroke),
		CodeToSyms:    keycodeToKeysyms,
		SymToCodes:    keysymToKeycodes,
	}
	if config.Logger == nil {
		wm.logger = log.New(os.Stdout, "==|>", log.Lmicroseconds)
	} else {
		wm.logger = config.Logger
	}
	go wm.loop()

	return wm, nil
}

func (w *Wm) Close() {
	xproto.ChangeWindowAttributes(w.Conn, w.DefaultRootId, xproto.CwEventMask, []uint32{uint32(
		xproto.EventMaskNoEvent)})
	w.Conn.Close()
}

func (w *Wm) pt(format string, args ...interface{}) {
	w.logger.Printf(format, args...)
}

func (w *Wm) loop() {
	for {
		ev, xerr := w.Conn.WaitForEvent()
		if ev == nil && xerr == nil {
			w.logger.Fatal("bad conn")
		}

		if xerr != nil {
			w.pt("ERROR: %v\n", xerr)
		}

		if ev != nil {
			w.pt("EVENT: %T %v\n", ev, ev)
			switch ev := ev.(type) {

			case xproto.CreateNotifyEvent:
				if ev.OverrideRedirect { // do not manage override-redirect windows
					continue
				}
				win := &Window{
					Id:     ev.Window,
					Parent: ev.Parent,
					X:      int(ev.X),
					Y:      int(ev.Y),
					Width:  int(ev.Width),
					Height: int(ev.Height),
					Border: int(ev.BorderWidth),
				}
				w.Windows[win.Id] = win
				w.pt("managed window %v\n", win.Id)

			case xproto.ConfigureRequestEvent:
				if win, ok := w.Windows[ev.Window]; ok && win.Mapped { // skip managed mapped window request
					notifyEv := xproto.ConfigureNotifyEvent{ // not moving or resizing
						Event:  win.Id,
						Window: win.Id,
						X:      int16(win.X),
						Y:      int16(win.Y),
						Width:  uint16(win.Width),
						Height: uint16(win.Height),
					}
					xproto.SendEvent(w.Conn, false, win.Id, xproto.EventMaskStructureNotify, string(notifyEv.Bytes()))
				} else { // configure as requested
					var vals []uint32
					flags := ev.ValueMask
					if xproto.ConfigWindowX&flags > 0 {
						vals = append(vals, uint32(ev.X))
					}
					if xproto.ConfigWindowY&flags > 0 {
						vals = append(vals, uint32(ev.Y))
					}
					if xproto.ConfigWindowWidth&flags > 0 {
						vals = append(vals, uint32(ev.Width))
					}
					if xproto.ConfigWindowHeight&flags > 0 {
						vals = append(vals, uint32(ev.Height))
					}
					if xproto.ConfigWindowBorderWidth&flags > 0 {
						vals = append(vals, 0) // do not set border width
					}
					if xproto.ConfigWindowSibling&flags > 0 {
						vals = append(vals, uint32(ev.Sibling))
					}
					if xproto.ConfigWindowStackMode&flags > 0 {
						vals = append(vals, uint32(ev.StackMode))
					}
					xproto.ConfigureWindow(w.Conn, ev.Window, flags, vals)
				}

			case xproto.MapRequestEvent:
				xproto.MapWindow(w.Conn, ev.Window)
				if win, ok := w.Windows[ev.Window]; ok {
					win.Mapped = true
					w.MapWindow <- win
				}
			case xproto.MapNotifyEvent:

			case xproto.UnmapNotifyEvent:
				if win, ok := w.Windows[ev.Window]; ok {
					win.Mapped = false
					w.UnmapWindow <- win
				}

			case xproto.DestroyNotifyEvent:
				delete(w.Windows, ev.Window)

			case xproto.KeyPressEvent:
				w.Strokes <- Stroke{
					Modifiers: ev.State,
					Sym:       w.CodeToSyms[ev.Detail][0],
				}
			case xproto.KeyReleaseEvent:

			}
		}
	}
}
