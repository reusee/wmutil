package wmutil

import "github.com/BurntSushi/xgb/xproto"

func (w *Wm) internAtoms() error {
	atoms := []string{
		"WM_STATE",
	}
	for _, atom := range atoms {
		reply, err := xproto.InternAtom(w.Conn, false, uint16(len(atom)), atom).Reply()
		if err != nil {
			return err
		}
		w.Atoms[atom] = reply.Atom
	}
	return nil
}
