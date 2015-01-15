package wmutil

import (
	"os"
	"os/exec"
	"os/signal"
	"testing"
)

var testSigs chan os.Signal

func TestMain(m *testing.M) {
	testSigs = make(chan os.Signal)
	signal.Notify(testSigs, os.Interrupt, os.Kill)
	os.Exit(m.Run())
}

func TestConnect(t *testing.T) {
	wm, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer wm.Close()

	exec.Command("xfce4-terminal").Start()

	<-testSigs
}