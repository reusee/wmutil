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
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer wm.Close()

	for {
		select {
		case win := <-wm.MapWindow:
			pt("window mapped %v\n", win)
		case win := <-wm.UnmapWindow:
			pt("window unmap %v\n", win)
		case stroke := <-wm.Strokes:
			exec.Command("sakura").Start()
			pt("stroke %v\n", stroke)
		case <-testSigs:
			return
		}
	}
}
