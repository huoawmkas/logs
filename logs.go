package logs

// 日志
import (
	"fmt"
	"os"
	"runtime"
	"time"
)

const (
	LOG_ERROR = iota
	LOG_WARING
	LOG_INFO
	LOG_DEBUG
)

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

var log *mylog

// 初始化
func init() {
	log = new(mylog)

	// default setting
	log.dir = "../log"
	log.file = "logs"
	log.level = LOG_INFO
}

func Init(dir string, file string, level int, savefile bool, prefix ...string) {
	log.setLevel(level)
	if len(prefix) != 0 {
		log.setPrefix(fmt.Sprintf("[%v] ", prefix[0]))
	}
	if savefile {
		log.setDir(dir)
		log.setFile(file)
		log.setSavefile(savefile)
		log.log = make(chan string, 100)
		go log.run()
	}
}

func Error(err ...interface{}) {
	if LOG_ERROR > log.level {
		return
	}
	log.write(LOG_ERROR, fmt.Sprint(err...))
}

func Waring(war ...interface{}) {
	if LOG_WARING > log.level {
		return
	}
	log.write(LOG_WARING, fmt.Sprint(war...))
}

func Info(info ...interface{}) {
	if LOG_INFO > log.level {
		return
	}
	log.write(LOG_INFO, fmt.Sprint(info...))
}

func Debug(deb ...interface{}) {
	if LOG_DEBUG > log.level {
		return
	}
	log.write(LOG_DEBUG, fmt.Sprint(deb...))
}

func Errorf(format string, v ...interface{}) {
	if LOG_ERROR > log.level {
		return
	}
	log.write(LOG_ERROR, fmt.Sprintf(format, v...))
}

func Waringf(format string, v ...interface{}) {
	if LOG_WARING > log.level {
		return
	}
	log.write(LOG_WARING, fmt.Sprintf(format, v...))
}

func Infof(format string, v ...interface{}) {
	if LOG_INFO > log.level {
		return
	}
	log.write(LOG_INFO, fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...interface{}) {
	if LOG_DEBUG > log.level {
		return
	}
	log.write(LOG_DEBUG, fmt.Sprintf(format, v...))
}

/*
 * 日志执行函数
 */
type mylog struct {
	log      chan string // 日志chan
	dir      string      // 日志存放目录
	file     string      // 日志文件名
	savefile bool        // 是否保存到文件
	level    int         // 日志级别
	prefix   string      // 日志前缀
}

func (l *mylog) setDir(dir string) {
	l.dir = dir
}

func (l *mylog) setFile(file string) {
	l.file = file
}

func (l *mylog) setSavefile(b bool) {
	l.savefile = b
}

func (l *mylog) setLevel(level int) {
	l.level = level
}

func (l *mylog) setPrefix(prefix string) {
	l.prefix = prefix
}

func (l *mylog) getLevelString(level int) string {
	switch level {
	case LOG_ERROR:
		return "ERRO"
	case LOG_WARING:
		return "WARI"
	case LOG_INFO:
		return "INFO"
	case LOG_DEBUG:
		return "DEBU"
	}

	return "unknown"
}

func (l *mylog) write(level int, str string) {
	levelColor := getColorByLevel(level)
	// 输出日志
	pc, _, line, _ := runtime.Caller(2)
	p := runtime.FuncForPC(pc)
	t := time.Now().Local().Format("2006/01/02 15:04:05.999999")
	str = fmt.Sprintf("%s%-26s \u001B[%dm[%s]\u001B[0m %s(%d): %s\n", l.prefix, t, levelColor, l.getLevelString(level), p.Name(), line, str)
	// 输出到控制台
	if false == l.savefile {
		fmt.Print(str)
		return
	}

	// 输出到文件
	l.log <- str
}

func getColorByLevel(level int) int {
	switch level {
	case LOG_DEBUG:
		return colorGray
	case LOG_WARING:
		return colorYellow
	case LOG_ERROR:
		return colorRed
	default:
		return colorBlue
	}
}

func (l *mylog) run() {

	var last, now time.Time

	var fp *os.File

	for {
		str := <-l.log

		// 获取时间
		now = time.Now()

		if (now.Day() != last.Day()) || (fp == nil) { // 跨天或之前文件打开失败,则重新创建文件
			if fp != nil {
				fp.Close()
			}

			// 判断文件夹是否存在
			_, err := os.Stat(l.dir)
			if nil != err {
				os.MkdirAll(l.dir, os.ModePerm)
			}

			path := fmt.Sprintf("%s/%s-%04d-%02d-%02d.log", l.dir, l.file,
				now.Year(), now.Month(), now.Day())
			fp, _ = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
		}

		if fp != nil {
			fp.WriteString(str)
		}
		last = now
	}
}
