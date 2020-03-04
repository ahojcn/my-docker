package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"syscall"
)

func main() {
	log.Infoln("pid:", os.Getpid())

	/* 5-2 cmd.Run()
		cmd := exec.Command("sh")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatalln("test exec error:", err)
		}
	*/

	// 5-3 使用 syscall.exec 程序替换
	command := "/bin/sh"
	if err := syscall.Exec(command, []string{command}, os.Environ()); err != nil {
		log.Fatalln("syscall.Exec error:", err)
	}
}
