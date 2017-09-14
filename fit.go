package fit



func Start() bool {
	initConfig, _ := App().Init()
	isStart, _ := App().Start()
	start :=  isStart && initConfig
	return start
}

func Stop() bool {
	App().Stop(true)
	return true
}
