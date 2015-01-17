package wmutil

import "github.com/BurntSushi/xgb/xproto"

func (w *Wm) internAtoms() error {
	atoms := []string{
		"WM_STATE",
		"WM_PROTOCOLS",
		"WM_DELETE_WINDOW",
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

var atomNames = make(map[xproto.Atom]string)

func (w *Wm) AtomName(atom xproto.Atom) string {
	if name, ok := atomNames[atom]; ok {
		return name
	}
	reply, err := xproto.GetAtomName(w.Conn, atom).Reply()
	if err != nil {
		w.pt("ERROR: get atom name: %v\n", err)
	}
	atomNames[atom] = reply.Name
	return reply.Name
}
