package wmutil

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

	logger *log.Logger

	NewWindow chan *Window
	Strokes   chan Stroke
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
	Key       byte
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
	if err := xproto.UngrabKeyChecked(conn, xproto.GrabAny, defaultRootId, xproto.ModMaskAny).Check(); err != nil {
		return nil, ef("ungrab keys: %v", err)
	}
	ignoreModifiers := []uint16{
		0,
		xproto.ModMaskLock,
		xproto.ModMask2,
		xproto.ModMaskLock | xproto.ModMask2,
	}
	for _, stroke := range config.Strokes {
		for _, mod := range ignoreModifiers {
			if err := xproto.GrabKeyChecked(conn, true, defaultRootId, stroke.Modifiers|mod,
				xproto.Keycode(stroke.Key), xproto.GrabModeAsync, xproto.GrabModeAsync).Check(); err != nil {
				return nil, ef("grab key: %v", err)
			}
		}
	}

	wm := &Wm{
		Conn:          conn,
		Setup:         setup,
		DefaultScreen: defaultScreen,
		DefaultRootId: defaultRootId,
		Windows:       make(map[xproto.Window]*Window),
		NewWindow:     make(chan *Window),
		Strokes:       make(chan Stroke),
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
	//TODO notify and wait for loop to exit
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
			switch ev := ev.(type) {

			case xproto.CreateNotifyEvent:
				if ev.OverrideRedirect {
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
				w.addWindow(win)

			case xproto.ConfigureRequestEvent:
				win, ok := w.Windows[ev.Window]
				if !ok { // not managed window, add to Windows
					win = &Window{
						Id:     ev.Window,
						Parent: ev.Parent,
						X:      int(ev.X),
						Y:      int(ev.Y),
						Width:  int(ev.Width),
						Height: int(ev.Height),
						Border: int(ev.BorderWidth),
					}
					w.addWindow(win)
				}
				// not moving or resizing, layout will do these
				notifyEv := xproto.ConfigureNotifyEvent{
					Event:  win.Id,
					Window: win.Id,
					X:      int16(win.X),
					Y:      int16(win.Y),
					Width:  uint16(win.Width),
					Height: uint16(win.Height),
				}
				xproto.SendEvent(w.Conn, false, win.Id, xproto.EventMaskStructureNotify, string(notifyEv.Bytes()))

			case xproto.MapRequestEvent:
				win, ok := w.Windows[ev.Window]
				if !ok { // skip not managed window. may be override-redirect
					continue
				}
				if win.Mapped {
					continue
				}
				xproto.MapWindow(w.Conn, win.Id)

			case xproto.MapNotifyEvent:

			case xproto.KeyPressEvent:
				w.Strokes <- Stroke{
					Modifiers: ev.State,
					Key:       byte(ev.Detail),
				}
			case xproto.KeyReleaseEvent:

			default:
				w.pt("NOT HANDLED EVENT: %T %v\n", ev, ev)
			}
		}
	}
}

func (w *Wm) addWindow(win *Window) {
	w.pt("added window %v\n", win.Id)
	w.Windows[win.Id] = win
	w.NewWindow <- win
}
