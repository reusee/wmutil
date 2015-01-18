package wmutil

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"
)

func (w *Wm) Atom(name string) xproto.Atom {
	if atom, ok := w.stringToAtom[name]; ok {
		return atom
	}
	reply, err := xproto.InternAtom(w.Conn, false, uint16(len(name)), name).Reply()
	if err != nil {
		log.Fatalf("atom intern error %v", err)
	}
	w.stringToAtom[name] = reply.Atom
	return reply.Atom
}

func (w *Wm) AtomName(atom xproto.Atom) string {
	if name, ok := w.atomToString[atom]; ok {
		return name
	}
	reply, err := xproto.GetAtomName(w.Conn, atom).Reply()
	if err != nil {
		log.Fatalf("get atom name error %v", err)
	}
	w.atomToString[atom] = reply.Name
	return reply.Name
}
