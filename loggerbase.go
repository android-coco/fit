package fit

import (
	"fmt"
	"io"
	"log"
)

type loggerBase struct {
	LoggerInterface
	Log      *log.Logger
	LogLevel int
	Tag      string //if the tag matched, then output
}

func newLoggerBase(out io.Writer, tag string, level int) *loggerBase {
	return &loggerBase{
		LogLevel: level,
		Log:      log.New(out, tag, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		Tag:      tag,
	}
}

func (lBase *loggerBase) LogVerbose(tag string, msg ...interface{}) {
	if lBase.LogLevel >= Verbose {
		if len(lBase.Tag) != 0 {
			if lBase.Tag != tag {
				return
			}
		}
		lBase.Log.SetPrefix("[Verbose]"+tag)
		lBase.Log.Output(3, fmt.Sprintln(msg...))
	}
}
func (lBase *loggerBase) LogDebug(tag string, msg ...interface{}) {
	if lBase.LogLevel >= Debug {
		if len(lBase.Tag) != 0 {
			if lBase.Tag != tag {
				return
			}
		}
		lBase.Log.SetPrefix("[Debug]"+tag)
		lBase.Log.Output(3, fmt.Sprintln(msg...))
	}
}
func (lBase *loggerBase) LogInfo(tag string, msg ...interface{}) {
	if lBase.LogLevel >= Info {
		if len(lBase.Tag) != 0 {
			if lBase.Tag != tag {
				return
			}
		}
		lBase.Log.SetPrefix("[Info]"+tag)
		lBase.Log.Output(3, fmt.Sprintln(msg...))
	}
}
func (lBase *loggerBase) LogWarn(tag string, msg ...interface{}) {
	if lBase.LogLevel >= Warn {
		if len(lBase.Tag) != 0 {
			if lBase.Tag != tag {
				return
			}
		}
		lBase.Log.SetPrefix("[Warn]"+tag)
		lBase.Log.Output(3, fmt.Sprintln(msg...))
	}
}
func (lBase *loggerBase) LogError(tag string, msg ...interface{}) {
	if lBase.LogLevel >= Error {
		if len(lBase.Tag) != 0 {
			if lBase.Tag != tag {
				return
			}
		}
		lBase.Log.SetPrefix("[Error]"+tag)
		lBase.Log.Output(3, fmt.Sprintln(msg...))
	}
}
func (lBase *loggerBase) LogAssert(tag string, msg ...interface{}) {
	if lBase.LogLevel >= Assert {
		if len(lBase.Tag) != 0 {
			if lBase.Tag != tag {
				return
			}
		}
		lBase.Log.SetPrefix("[Assert]"+tag)
		lBase.Log.Output(3, fmt.Sprintln(msg...))
	}
}
func (lBase *loggerBase) SetLogLevel(level int) {
	lBase.LogLevel = level
}
func (lBase *loggerBase) SetLogTag(tag string) {
	lBase.Tag = tag
}
func (lBase *loggerBase) Term() {
}
