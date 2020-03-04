package main

import (
	log "github.com/Sirupsen/logrus"
	"mydocker/command"
	"os/exec"
)

func main() {
	test01()
}

func test01() {
	copy("/home/ahojcn/test-mydocker/busybox", "/home/ahojcn/test-mydocker/tmp")
	copy("/home/ahojcn/test-mydocker/busybox/bin/top", "/home/ahojcn/test-mydocker/tmp")
}

func copy(src, dst string) {
	exist, _ := command.PathExists(src)
	if !exist {
		log.Errorln("src not exists:", src)
		return
	}
	exist, _ = command.PathExists(dst)
	if !exist {
		log.Errorln("dst not exists:", dst)
		return
	}
	if _, err := exec.Command("cp", "-r", src, dst).CombinedOutput(); err != nil {
		log.Errorf("cp -r %s %s error: %v", src, dst, err)
		return
	}
}
