package fit

import (
	"log"
	"os"
	"time"
)

type logFile struct {
	loggerBase
	FileFolder  string //log file folder
	LogFileName string //log file name
	FileHandle  *os.File
}

func newLoggerToFile(tag string, level int) *logFile {

	currentDir, err := os.Getwd()
	if err != nil {
		return nil
	}

	lfileFolder := currentDir + string(os.PathSeparator) + "log"

	err = os.MkdirAll(lfileFolder, os.ModePerm)
	if err != nil {
		return nil
	}

	lfileName := lfileFolder + string(os.PathSeparator) + time.Now().Format("2006-01-02") + ".txt"
	fHandle, err := os.OpenFile(lfileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)

	if err != nil {
		return nil
	}

	return &logFile{
		loggerBase: loggerBase{LogLevel: level,
			Log: log.New(fHandle, tag, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
			Tag: tag,
		},
		FileFolder:  lfileFolder,
		LogFileName: lfileName,
		FileHandle:  fHandle,
	}

}

func (lFile *logFile) Term() {
	lFile.FileHandle.Close()
}
