package main

import (
	"container/list"
	"fmt"
	"log"
	"os/exec"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/reusee/wmutil"
)

var (
	pt = fmt.Printf
)

func main() {
	var wm *wmutil.Wm
	windows := list.New()
	mod := uint16(xproto.ModMaskControl)
	keyBindings := map[wmutil.Stroke]func(){
		wmutil.Stroke{mod, wmutil.Key_Return}: func() {
			exec.Command("sakura").Start()
		},
		wmutil.Stroke{mod, wmutil.Key_o}: func() {
			exec.Command("dmenu_run").Start()
		},
		wmutil.Stroke{mod, wmutil.Key_z}: func() {
			win := wm.PointingWindow()
			if win != nil {
				win.Destroy()
			}
		},
		wmutil.Stroke{mod, wmutil.Key_j}: func() {
			front := windows.Remove(windows.Front()).(*wmutil.Window)
			back := windows.Back().Value.(*wmutil.Window)
			front.Below(back)
			windows.PushBack(front)
		},
	}

	var strokes []wmutil.Stroke
	for stroke, _ := range keyBindings {
		strokes = append(strokes, stroke)
	}
	var err error
	wm, err = wmutil.New(&wmutil.Config{
		Strokes: strokes,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer wm.Close()

	screenWidth := int(wm.DefaultScreen.WidthInPixels)
	screenHeight := int(wm.DefaultScreen.HeightInPixels)
	for {
		select {
		case win := <-wm.Map:
			win.SetGeometry(0, 0, screenWidth, screenHeight)
			win.WarpPointer()
			windows.PushFront(win)
		case win := <-wm.Unmap:
			_ = win
		case stroke := <-wm.Stroke:
			if cb, ok := keyBindings[stroke]; ok {
				cb()
			}
		}
	}
}
