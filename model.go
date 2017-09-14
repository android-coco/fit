package fit


import (
	_ "github.com/go-sql-driver/mysql"
	_"github.com/denisenkom/go-mssqldb"
	"github.com/go-xorm/xorm"
	"time"
)

var (
	engine *xorm.Engine
	msEngine *xorm.Engine
)
type JsonTime time.Time
func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"`+time.Time(j).Format("2006-01-02 15:04:05")+`"`), nil
}
func Engine() *xorm.Engine  {
	if engine == nil {
		var err error
		engine, err = initEngine()
		engine.TZLocation = time.UTC
		engine.DatabaseTZ = time.UTC
		if err != nil {
			Logger().LogError("fail to create engine: %v", err)
		}
	}
	return engine

}

func initEngine() (*xorm.Engine, error)  {
	var dataSource  string = Config().UserName + ":" + Config().Password + "@tcp(" + Config().HostName + ":" + Config().DBPort + ")/" + Config().DataBase + "?charset=utf8"
	//db, err := xorm.NewEngine(Config().DriverName, "root:123456@tcp(127.0.0.1:3307)/test?charset=utf8")
	db, err := xorm.NewEngine(Config().DriverName, dataSource)
	Logger().LogInfo("", "data source",dataSource)
	if err == nil {
		return db, err
	}
	return nil, err
}

func SQLServerEngine() *xorm.Engine  {
	if msEngine == nil {
		var err error
		var dataSourceName string  = "server=" + Config().MSServer + ";port=" + Config().MSDBPort + ";database=" + Config().MSDataBase + ";user id=" + Config().MSUserId + ";password=" + Config().MSPassword
		msEngine, err = xorm.NewEngine(Config().MSDriverName, dataSourceName)
		Logger().LogInfo("","dataSourceName", dataSourceName)
		//msEngine, err = xorm.NewEngine("mssql", "server=192.168.0.130;port=1433;database=test;user id=sa;password=youhao;")
		if err != nil {
			Logger().LogError("new sql server engine failed:", err)
			return nil
		}
		if err := msEngine.Ping(); err != nil {
			Logger().LogInfo("ms engine ping", err)
			return msEngine
		}
	}
	return msEngine
}