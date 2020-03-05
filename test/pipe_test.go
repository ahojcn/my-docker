package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

/**
管道测试
*/

// 栗子1：同步操作，先用管道的写端写数据，然后管道的读端读数据
func Test004(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		log.Fatalln("os.pipe error:", err)
	}
	defer reader.Close()
	defer writer.Close()

	_, err = writer.Write([]byte("pipe content"))
	if err != nil {
		log.Fatalln("writer.Write error:", err)
	}

	buf := make([]byte, 100)
	n, err := reader.Read(buf)
	if err != nil {
		log.Fatalln("reader.Read(buf) error:", err)
	}
	log.Fatalln("read:", string(buf[:n]))
}

// 栗子2：异步，启动两个 goroutine 一个写端写10次，管道读端读管道里面所有内容
func Test005(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		log.Fatalln("os.pipe error:", err)
	}

	go func() {
		for i := 0; i < 10; i++ {
			content := fmt.Sprintf("%s-%d\n", "pipe content", i)
			_, err := writer.WriteString(content)
			if err != nil {
				log.Fatalln("writer.Write error:", err)
			}
		}
		writer.Close()
	}()

	go func() {
		n, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatalln("reader.Read(buf) error:", err)
		}
		log.Printf("Read content:%q\n", n)
	}()

	for i := 0; i <= 100; i++ {
		time.Sleep(time.Second * 1)
	}
}
