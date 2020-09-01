package log

import (
	"bytes"
	"fmt"
	beegolog "gameserver/utils/log/beego"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"io/ioutil"
	"runtime"
	"sync"
)

var beego *beegolog.BeeLogger
var once sync.Once

func LogBeego() *beegolog.BeeLogger {
	once.Do(func() {
		// 设置配置文件
		logConf := make(map[string]interface{})
		logConf["filename"] = viper.GetString("general.logsDir")
		logConf["maxlines"] = 50000
		logConf["maxsize"] = 10240000
		logConf["maxdays"] = 1

		confStr, err := jsoniter.Marshal(logConf)
		if err != nil {
			fmt.Println("marshal failed,err:", err)
			return
		}

		beego = beegolog.NewLogger()
		beego.SetLogger("file", string(confStr)) // 设置日志记录方式：本地文件记录
		beego.SetLevel(beegolog.LevelDebug)      // 设置日志写入缓冲区的等级
		//beego.SetLogger(beegolog.AdapterConsole)

		if runtime.GOOS == "linux" {
			// LINUX系统
			beego.SetLevel(beegolog.LevelWarning) // 设置日志写入缓冲区的等级
			beego.Async(1000)
		} else if runtime.GOOS == "windows" {
			beego.SetLevel(beegolog.LevelDebug) // 设置日志写入缓冲区的等级
			beego.SetLogger(beegolog.AdapterConsole)
		}
		beego.EnableFuncCallDepth(true) // 输出log时能显示输出文件名和行号（非必须）

	})
	return beego
}

// Emergency logs a message at emergency level.
func Emergency(f interface{}, v ...interface{}) {
	LogBeego().Emergency(beegolog.FormatLog(f, v...))
}

// Alert logs a message at alert level.
func Alert(f interface{}, v ...interface{}) {
	LogBeego().Alert(beegolog.FormatLog(f, v...))
}

// Critical logs a message at critical level.
func Critical(f interface{}, v ...interface{}) {
	LogBeego().Critical(beegolog.FormatLog(f, v...))
}

// Error logs a message at error level.
func Error(f interface{}, v ...interface{}) {
	LogBeego().Error(beegolog.FormatLog(f, v...))
}

// Warning logs a message at warning level.
func Warning(f interface{}, v ...interface{}) {
	LogBeego().Warn(beegolog.FormatLog(f, v...))
}

// Notice logs a message at notice level.
func Notice(f interface{}, v ...interface{}) {
	LogBeego().Notice(beegolog.FormatLog(f, v...))
}

// Info compatibility alias for Warning()
func Info(f interface{}, v ...interface{}) {
	LogBeego().Info(beegolog.FormatLog(f, v...))
}

// Debug logs a message at debug level.
func Debug(f interface{}, v ...interface{}) {
	LogBeego().Debug(beegolog.FormatLog(f, v...))
}

func Stack(f interface{}, v ...interface{}) {
	//logging.print(errorLog, stack("%v", err...))

	LogBeego().Debug(beegolog.FormatLog(f, stack(v...)))
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
)

// This package exposes some handy traceback functionality buried in the runtime.
//
// It can also be used to provide context to errors reducing the temptation to
// panic carelessly, just to get stack information.
//
// The theory is that most errors that are created with the fmt.Errorf
// style are likely to be rare, but require more context to debug
// properly. The additional cost of computing a stack trace is
// therefore negligible.
type StackError interface {
	Error() string
	StackTrace() string
}

type stackError struct {
	err        error
	stackTrace string
}

func (e stackError) Error() string {
	return fmt.Sprintf("%v\n%v", e.err, e.stackTrace)
}

func (e stackError) StackTrace() string {
	return e.stackTrace
}

func stack(args ...interface{}) error {
	stack := ""
	// See if any arg is already embedding a stack - no need to
	// recompute something expensive and make the message unreadable.
	for _, arg := range args {
		if stackErr, ok := arg.(stackError); ok {
			stack = stackErr.stackTrace
			break
		}
	}

	if stack == "" {
		// magic 5 trims off just enough stack data to be clear
		stack = stack1(3)
	}

	return stackError{fmt.Errorf("%v", args...), stack}
}

func stack1(calldepth int) string {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := calldepth; ; i++ { // Caller we care about is the user, 2 frames up
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		line-- // in stack trace, lines are 1-indexed but our array is 0-indexed
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.String()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.Trim(lines[n], " \t")
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
