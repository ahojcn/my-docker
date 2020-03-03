package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/command"
	_ "mydocker/nsenter"
	"os"
	"runtime"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = "implementation of mydocker"

	app.Commands = []cli.Command{
		command.RunCommand,
		command.InitCommand,
		command.CommitCommand,
		command.ListCommand,
		command.LogCommand,
		command.ExecCommand,
		command.StopCommand,
	}

	app.Before = func(context *cli.Context) error {
		// Log as JSON instead of the default ASCII formatter.
		log.AddHook(lineHook{
			Field:  "source",
			Skip:   0,
			levels: nil,
		})
		log.SetFormatter(&log.TextFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln("Run() error", err)
	}
}

// line number hook for log the call context,
type lineHook struct {
	Field string
	// skip为遍历调用栈开始的索引位置
	Skip   int
	levels []log.Level
}

// Levels implement levels
func (hook lineHook) Levels() []log.Level {
	allLevels := make([]log.Level, 10)
	allLevels = append(allLevels, log.InfoLevel)
	allLevels = append(allLevels, log.DebugLevel)
	allLevels = append(allLevels, log.ErrorLevel)
	allLevels = append(allLevels, log.FatalLevel)
	allLevels = append(allLevels, log.WarnLevel)
	allLevels = append(allLevels, log.PanicLevel)
	return allLevels
}

// Fire implement fire
func (hook lineHook) Fire(entry *log.Entry) error {
	entry.Data[hook.Field] = findCaller(hook.Skip)
	return nil
}

func findCaller(skip int) string {
	file := ""
	line := 0
	var pc uintptr
	// 遍历调用栈的最大索引为第11层.
	for i := 0; i < 11; i++ {
		file, line, pc = getCaller(skip + i)
		// 过滤掉所有logrus包，即可得到生成代码信息
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}

	fullFnName := runtime.FuncForPC(pc)

	fnName := ""
	if fullFnName != nil {
		fnNameStr := fullFnName.Name()
		// 取得函数名
		parts := strings.Split(fnNameStr, ".")
		fnName = parts[len(parts)-1]
	}

	return fmt.Sprintf("%s:%d:%s()", file, line, fnName)
}

func getCaller(skip int) (string, int, uintptr) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0, pc
	}
	n := 0

	// 获取包名
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line, pc
}
