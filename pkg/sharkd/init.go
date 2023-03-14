package sharkd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var sockpath string

func SockPath() string {
	return sockpath
}

func Start(ctx context.Context) {
	sharkd, err := lookupsharkd()
	if err != nil {
		log.Fatal(err)
	}
	sockpath = filepath.Join(os.TempDir(), fmt.Sprintf("sharkd.%d.sock", os.Getpid()))

	fmt.Println(sockpath, "11")
	scmd := exec.CommandContext(ctx, sharkd, "unix:"+sockpath)
	scmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err = scmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	var finish = make(chan struct{})
	defer close(finish)
	go func() {
		select {
		case <-ctx.Done():
			syscall.Kill(-scmd.Process.Pid, syscall.SIGKILL)
		case <-finish:
		}
	}()

	if err := scmd.Wait(); err != nil {
		fmt.Printf("sharkd process %d exit with error: %s", scmd.Process.Pid, err)
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
