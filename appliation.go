package fit

import (
	"flag"
	"net/http"
)

const (
	LOG_TAG = "fit->"
)

var (
	g_App  *appication = nil
	Modles []interface{}
)

type appication struct {
	Server *http.Server
}

func (app *appication) InitModels(models []interface{}) {
	Modles = models
}

func (app *appication) Init() (bool, error) {
	//read from configuration file and initialize application
	loglevel := flag.Int("loglevel", Verbose, `Log print level setting.
	Log level: Verbose(5) Debug(4) Info(3) Warn(2) Error(1) Assert(0) Nonelog(-1)`)
	logtag := flag.String("logtag", "", "Log print filter setting. The default is ''")
	confFile := flag.String("confile", "config/fit.conf", "Configuration file.")
	ok, err := Config().LoadConfig(*confFile)
	if ok == false {
		Logger().LogError(LOG_TAG, err.Error())
	}

	flag.Parse()

	Logger().SetLogLevel(*loglevel)
	Logger().SetLogTag(*logtag)

	return ok, err
}

//implement function at application level
func (app *appication) Start() (bool, error) {

	app.Server = &http.Server{
		Handler:        Router(),
		Addr:           Config().Port,
		ReadTimeout:    Config().ReadTimeout,
		WriteTimeout:   Config().WriteTimeout,
		MaxHeaderBytes: Config().MaxHeaderBytes,
	}
	Logger().LogInfo(LOG_TAG, "start to listen on port "+Config().Port)
	// 同步数据库
	//if  err := Engine().Sync2(Modles...); err!=nil{
	//	//fmt.Println("fail to sync database: ", err)
	//	Logger().LogError("fail to sync database: ", err)
	//} else {
	//	Logger().LogInfo("success to sync database。。。")
	//}
	SartOK = true //启动OK
	err := app.Server.ListenAndServe()

	if err != nil {
		Logger().LogError(LOG_TAG, err.Error())
		SartOK = false //启动失败
		return false, err
	}
	return true, nil
}

func (app *appication) Stop(wait bool) {
	app.Server.Shutdown(nil)
}

func App() *appication {
	if g_App == nil {
		g_App = &appication{}
	}

	return g_App
}
