package wmutil

import (
	"os"
	"os/exec"
	"os/signal"
	"testing"

	"github.com/BurntSushi/xgb/xproto"
)

var testSigs chan os.Signal

func TestMain(m *testing.M) {
	testSigs = make(chan os.Signal)
	signal.Notify(testSigs, os.Interrupt, os.Kill)
	os.Exit(m.Run())
}

func TestConnect(t *testing.T) {
	wm, err := New(&Config{
		Strokes: []Stroke{
			{xproto.ModMaskControl, Key_F},
			{0, Key_F1},
			{0, Key_F2},
			{0, Key_F3},
			{0, Key_F4},
			{0, Key_F5},
			{0, Key_F6},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer wm.Close()

	var windows []*Window
	for {
		select {
		case win := <-wm.Map:
			pt("window %v mapped instance %s class %s\n", win, win.Instance, win.Class)
			n := len(windows)
			if n > 1 {
				win.SetGeometry(n*50, n*50, 500, 100)
			} else {
				win.SetPos(n*50, n*50)
				win.SetSize(500, 100)
			}
			windows = append(windows, win)
			win.WarpPointer()
		case win := <-wm.Unmap:
			pt("window unmap %v\n", win)
		case stroke := <-wm.Stroke:
			pt("stroke %v\n", stroke)
			switch stroke.Sym {
			case Key_F1:
				exec.Command("sakura").Start()
			case Key_F2:
				windows[0].Above(nil)
			case Key_F3:
				windows[0].Below(windows[1])
			case Key_F4:
				win := wm.PointingWindow()
				if win != nil {
					pt("destroy %v\n", win.Id)
					win.Destroy()
				}
			case Key_F5:
			case Key_F6:
			}
		case win := <-wm.NameChanged:
			win.ReadLock(func() {
				pt("window name: %v\n", win.Name)
			})
		case win := <-wm.IconChanged:
			win.ReadLock(func() {
				pt("window icon: %v\n", win.Icon)
			})
		case <-wm.Resize:
		case <-testSigs:
			return
		}
	}
}
