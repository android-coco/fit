package fit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"reflect"
	"strconv"
	"time"
)

var (
	DEFAULT_COMMENT             = "#"
	DEFAULT_COMMENT_SEM         = ";"
	g_Config            *config = nil
)

type config struct {
	DocRoot           string
	RunMode           string
	LogLevel          int
	LogFilePath       string
	LogTag            string
	Port              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
	WriteTimeout      time.Duration

	SessionTimeout time.Duration
	SessionKey     string

	MaxHeaderBytes int
	MaxBytesReader int

	DriverName string
	HostName   string
	DBPort     string
	UserName   string
	Password   string
	DataBase   string


	MSDriverName string
	MSServer   string
	MSDBPort     string
	MSUserId   string
	MSPassword   string
	MSDataBase   string
}

func (conf *config) LoadConfig(fname string) (bool, error) {
	f, err := os.Open(fname)
	if err != nil {
		return false, err
	}

	defer f.Close()

	buf := bufio.NewReader(f)
	var lineNum int = 0
	for {
		lineNum++
		bytes, _, err := buf.ReadLine()

		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}

		line := strings.TrimSpace(string(bytes))
		switch {
		case len(line) == 0:
			continue
		case strings.HasPrefix(line, DEFAULT_COMMENT):
			continue
		case strings.HasPrefix(line, DEFAULT_COMMENT_SEM):
			continue
		default:
			if ok, err := conf.assignConfig(line, lineNum); ok == false {
				return ok, err
			}
		}
	}

	return true, nil
}

func (conf *config) assignConfig(line string, lineNum int) (bool, error) {
	optionVal := strings.SplitN(line, "=", 2)
	if len(optionVal) != 2 {
		return false, fmt.Errorf("parse  the content error : line %d , %s = ? ", lineNum, optionVal[0])
	}
	key := strings.TrimSpace(optionVal[0])
	value := strings.TrimSpace(optionVal[1])

	pv := reflect.ValueOf(g_Config)

	if pv.Kind() == reflect.Ptr && !pv.Elem().CanSet() {
		return false, fmt.Errorf("assign config error: reflect can not set")
	} else {
		pv = pv.Elem()
	}

	pf := pv.FieldByName(key)

	switch pf.Kind() {
	case reflect.String:
		pf.SetString(value)
	case reflect.Int:
		val, _ := fitInt(value)
		pf.SetInt(int64(val))
	case reflect.Int64:
		timeDur, _ := fitTimeDuration(value)
		durReflect := reflect.ValueOf(timeDur)
		pf.Set(durReflect)
	default:
		return false, fmt.Errorf("config does not contain this reflect.type: %s", pf.Kind())
	}

	return true, nil
}

func fitInt(value string) (int, error) {
	return strconv.Atoi(value)
}

func fitTimeDuration(value string) (time.Duration, error) {
	return time.ParseDuration(value)
}

func Config() *config {
	if g_Config == nil {
		g_Config = &config{
			DocRoot:           "/test/defalut",
			RunMode:           "Debug",
			LogLevel:          Verbose,
			Port:              ":80",
			ReadTimeout:       20,
			ReadHeaderTimeout: 20,
			IdleTimeout:       20,
			WriteTimeout:      20,
			MaxHeaderBytes:    1 << 12,
			MaxBytesReader:    1 << 20,

			SessionTimeout: 7200,
			SessionKey:     "fitsid",

			LogTag:      "",
			LogFilePath: "",

			DriverName: "mysql",
			HostName:   "127.0.0.1",
			DBPort:     "3306",
			UserName:   "root",
			Password:   "123456",
			DataBase:   "test",

			MSDriverName: "mysql",
			MSServer:   "127.0.0.1",
			MSDBPort:     "3306",
			MSUserId:   "root",
			MSPassword:   "123456",
			MSDataBase:   "test",
		}
	}

	return g_Config
}
