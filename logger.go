package fit

import (
	"container/list"
	"os"
	"sync"
)

var (
	g_LoggerManager *loggerManager = nil
)

const (
	Verbose = 5
	Debug   = 4
	Info    = 3
	Warn    = 2
	Error   = 1
	Assert  = 0
	Nonelog = -1
)

type LoggerInterface interface {
	LogVerbose(tag string, msg ...interface{})
	LogDebug(tag string, msg ...interface{})
	LogInfo(tag string, msg ...interface{})
	LogWarn(tag string, msg ...interface{})
	LogError(tag string, msg ...interface{})
	LogAssert(tag string, msg ...interface{})
	SetLogLevel(level int)
	SetLogTag(tag string)
	Term()
}

type loggerManager struct {
	loggerSet *list.List
	logLevel  int
	logTag    string
	locker    *sync.RWMutex
}

func (log *loggerManager) forEach() {
	defer func() {
		log.locker.Unlock()
	}()

	log.locker.Lock()
	for e := log.loggerSet.Front(); e != nil; e = e.Next() {
		//fmt.Print(e.Value.(int))
		//e.Value.Term()
	}
}

func (log *loggerManager) LogVerbose(tag string, msg ...interface{}) {
	defer func() {
		log.locker.Unlock()
	}()

	log.locker.Lock()
	for e := log.loggerSet.Front(); e != nil; e = e.Next() {
		if log, ok := e.Value.(LoggerInterface); ok {
			log.LogVerbose(tag, msg)
		}
	}
}

func (log *loggerManager) LogDebug(tag string, msg ...interface{}) {
	defer func() {
		log.locker.Unlock()
	}()

	log.locker.Lock()
	for e := log.loggerSet.Front(); e != nil; e = e.Next() {
		if log, ok := e.Value.(LoggerInterface); ok {
			log.LogDebug(tag, msg)
		}
	}
}
func (log *loggerManager) LogInfo(tag string, msg ...interface{}) {
	defer func() {
		log.locker.Unlock()
	}()

	log.locker.Lock()
	for e := log.loggerSet.Front(); e != nil; e = e.Next() {
		if log, ok := e.Value.(LoggerInterface); ok {
			log.LogInfo(tag, msg)
		}
	}
}
func (log *loggerManager) LogWarn(tag string, msg ...interface{}) {
	defer func() {
		log.locker.Unlock()
	}()

	log.locker.Lock()
	for e := log.loggerSet.Front(); e != nil; e = e.Next() {
		if log, ok := e.Value.(LoggerInterface); ok {
			log.LogWarn(tag, msg)
		}
	}
}
func (log *loggerManager) LogError(tag string, msg ...interface{}) {
	defer func() {
		log.locker.Unlock()
	}()

	log.locker.Lock()
	for e := log.loggerSet.Front(); e != nil; e = e.Next() {
		if log, ok := e.Value.(LoggerInterface); ok {
			log.LogError(tag, msg)
		}
	}
}
func (log *loggerManager) LogAssert(tag string, msg ...interface{}) {
	defer func() {
		log.locker.Unlock()
	}()

	log.locker.Lock()
	for e := log.loggerSet.Front(); e != nil; e = e.Next() {
		if log, ok := e.Value.(LoggerInterface); ok {
			log.LogAssert(tag, msg)
		}
	}
}

func (log *loggerManager) SetLogTag(tag string) {
	defer func() {
		log.locker.Unlock()
	}()

	log.locker.Lock()
	for e := log.loggerSet.Front(); e != nil; e = e.Next() {
		if log, ok := e.Value.(LoggerInterface); ok {
			log.SetLogTag(tag)
		}
	}
}

func (log *loggerManager) SetLogLevel(level int) {
	log.logLevel = level

	defer func() {
		log.locker.Unlock()
	}()

	log.locker.Lock()
	for e := log.loggerSet.Front(); e != nil; e = e.Next() {
		if log, ok := e.Value.(LoggerInterface); ok {
			log.SetLogLevel(level)
		}
	}
}

func (log *loggerManager) GetLogLevel() int {
	return log.logLevel
}

func (log *loggerManager) RegistLogger(loggerInstance LoggerInterface) {
	defer func() {
		log.locker.Unlock()
	}()

	log.locker.Lock()
	log.loggerSet.PushBack(loggerInstance)
}

func (log *loggerManager) UnRegistLogger(loggerInstance LoggerInterface) {
	defer func() {
		log.locker.Unlock()
	}()

	log.locker.Lock()
	for e := log.loggerSet.Front(); e != nil; e = e.Next() {
		//fmt.Print(e.Value.(int))
		//e.Value.Term()
	}
}

func Logger() *loggerManager {
	if g_LoggerManager == nil {
		g_LoggerManager = &loggerManager{
			loggerSet: list.New(),
			locker:    new(sync.RWMutex),
		}
		g_LoggerManager.RegistLogger(newLoggerBase(os.Stdout, Config().LogTag, Config().LogLevel))
		g_LoggerManager.RegistLogger(newLoggerToFile(Config().LogTag, Config().LogLevel))
	}
	return g_LoggerManager
}
