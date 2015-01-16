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
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer wm.Close()

	var windows []*Window
	for {
		select {
		case win := <-wm.MapWindow:
			pt("window mapped %v\n", win)
			win.Focus()
			n := len(windows)
			if n > 1 {
				win.SetGeometry(n*50, n*50, 500, 100)
			} else {
				win.SetPos(n*50, n*50)
				win.SetSize(500, 100)
			}
			windows = append(windows, win)
		case win := <-wm.UnmapWindow:
			pt("window unmap %v\n", win)
		case stroke := <-wm.Strokes:
			switch stroke.Sym {
			case Key_F1:
				exec.Command("sakura").Start()
			case Key_F2:
				windows[0].Above(nil)
				windows[0].Focus()
			case Key_F3:
				windows[0].Below(windows[1])
				windows[1].Focus()
			case Key_F4:
				wm.GetFocus().Opposite(nil)
			}
			pt("stroke %v\n", stroke)
		case <-testSigs:
			return
		}
	}
}
