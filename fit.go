package fit

var SartOK bool //标记是否启动成功
func Start() bool {
	initConfig, _ := App().Init()
	isStart, _ := App().Start()
	start := isStart && initConfig
	App().RegisterMime()
	return start
}

func Stop() bool {
	App().Stop(true)
	return true
}

