package daemon

/*
#include "unistd.h"
*/
import "C"

//func Daemon() {
//	cmd := exec.Command(os.Args[0])
//	cmd.Stdin = nil
//	cmd.Stderr = nil
//	cmd.Stdout = nil
//	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
//
//	err := cmd.Start()
//	if err == nil {
//		cmd.Process.Release()
//		os.Exit(0)
//	}
//}

func Daemon() {
	C.daemon(1, 0)
}
