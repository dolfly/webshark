package sharkd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var sockpath string

func init() {
	sharkd, err := lookupsharkd()
	if err != nil {
		log.Fatal(err)
	}
	sockpath = filepath.Join(os.TempDir(), fmt.Sprintf("sharkd.%d.sock", os.Getpid()))
	scmd := exec.Command(sharkd, "unix:"+sockpath)
	err = scmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
func lookupsharkd() (string, error) {
	sharkdpath, err := exec.LookPath("sharkd")
	if err != nil {
		fmt.Printf("lookup sharkd error: %s", err.Error())
		return "", err
	}
	return sharkdpath, nil
}
