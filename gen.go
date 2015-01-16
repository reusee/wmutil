// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
)

func main() {
	genSymDefs()
	exec.Command("gofmt -w .").Run()
}

func genSymDefs() {
	buf := new(bytes.Buffer)
	w := func(format string, args ...interface{}) {
		fmt.Fprintf(buf, format, args...)
	}
	w(`package wmutil

var (
`)
	content, err := ioutil.ReadFile("/usr/include/X11/keysymdef.h")
	if err != nil {
		log.Fatal(err)
	}
	defPattern := regexp.MustCompile(`(?m:^#define (?:XK_)([^\s]+)\s+(\S+))`)
	matches := defPattern.FindAllSubmatch(content, -1)
	for _, group := range matches {
		w("\tKey_%s uint32 = %s\n", group[1], group[2])
	}
	w(")\n")
	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("symdef.go", src, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
