// internal/log/log.go
package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// 日志级别
const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	DisableLevel
)

// 日志级别名称
var levelNames = []string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

// Logger 表示日志记录器
type Logger struct {
	mu       sync.Mutex
	prefix   string      // 日志前缀
	level    int         // 当前日志级别
	output   io.Writer   // 输出位置
	logger   *log.Logger // 标准库日志记录器
	colorful bool        // 是否启用颜色
}

var (
	// 默认日志记录器
	defaultLogger = &Logger{
		prefix:   "geerpc",
		level:    InfoLevel,
		output:   os.Stdout,
		colorful: true,
	}
)

func init() {
	defaultLogger.logger = log.New(defaultLogger.output, "", log.LstdFlags)
}

// SetLevel 设置日志级别
func SetLevel(level int) {
	if level < DebugLevel || level > DisableLevel {
		return
	}
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	defaultLogger.level = level
}

// SetPrefix 设置日志前缀
func SetPrefix(prefix string) {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	defaultLogger.prefix = prefix
}

// SetOutput 设置日志输出位置
func SetOutput(w io.Writer) {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	defaultLogger.output = w
	defaultLogger.logger = log.New(w, "", log.LstdFlags)
}

// Debug 输出调试级别日志
func Debug(format string, v ...interface{}) {
	defaultLogger.logf(DebugLevel, format, v...)
}

// Info 输出信息级别日志
func Info(format string, v ...interface{}) {
	defaultLogger.logf(InfoLevel, format, v...)
}

// Warn 输出警告级别日志
func Warn(format string, v ...interface{}) {
	defaultLogger.logf(WarnLevel, format, v...)
}

// Error 输出错误级别日志
func Error(format string, v ...interface{}) {
	defaultLogger.logf(ErrorLevel, format, v...)
}

// Fatal 输出致命错误日志并退出程序
func Fatal(format string, v ...interface{}) {
	defaultLogger.logf(FatalLevel, format, v...)
	os.Exit(1)
}

// logf 实际的日志记录函数
func (l *Logger) logf(level int, format string, v ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// 获取调用者信息
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	file = filepath.Base(file)

	// 构建日志消息
	msg := fmt.Sprintf(format, v...)
	timeStr := time.Now().Format("2006/01/02 15:04:05")
	logStr := fmt.Sprintf("%s [%s] %s:%d %s: %s",
		timeStr, levelNames[level], file, line, l.prefix, msg)

	// 输出日志
	fmt.Fprintln(l.output, logStr)

	// 如果是致命错误，同时记录堆栈信息
	if level == FatalLevel {
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, true)
		fmt.Fprintf(l.output, "=== Stack trace ===\n%s\n", buf[:n])
	}
}
