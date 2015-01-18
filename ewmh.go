package wmutil

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func (w *Wm) setSupported() error {
	atoms := []xproto.Atom{
	//TODO
	}
	buf := make([]byte, len(atoms)*4)
	for i, atom := range atoms {
		xgb.Put32(buf[i*4:], uint32(atom))
	}
	err := xproto.ChangePropertyChecked(w.Conn, xproto.PropModeReplace, w.DefaultRootId,
		w.Atom("_NET_SUPPORTED"), xproto.AtomAtom, 32, uint32(len(atoms)), buf).Check()
	return err
}
