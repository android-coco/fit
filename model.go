package fit


import (
	_ "github.com/go-sql-driver/mysql"//mysql
	_"github.com/denisenkom/go-mssqldb"//sqlserver
	_ "github.com/mattn/go-oci8"//orcale
	"github.com/go-xorm/xorm"
	"time"
	"github.com/go-xorm/core"
)

var (
	engine *xorm.Engine
	msEngine *xorm.Engine
)

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"`+time.Time(j).Format("2006-01-02 15:04:05")+`"`), nil
}
func MySqlEngine() *xorm.Engine  {
	if engine == nil {
		var err error
		engine, err = initEngine()
		engine.TZLocation = time.UTC
		engine.DatabaseTZ = time.UTC
		//SnakeMapper 支持struct为驼峰式命名，表结构为下划线命名之间的转换，这个是默认的Maper
		//映射同名设置默认
		engine.SetMapper(core.SameMapper{})
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
		//SnakeMapper 支持struct为驼峰式命名，表结构为下划线命名之间的转换，这个是默认的Maper
		//映射同名设置默认
		msEngine.SetMapper(core.SameMapper{})
		Logger().LogInfo("","dataSourceName", dataSourceName)
		//msEngine, err = xorm.NewEngine("mssql", "server=192.168.0.130;port=1433;database=test;user id=sa;password=youhao;")
		if err != nil {
			Logger().LogError("new sql server engine failed:", err)
			return nil
		}
		if err := msEngine.Ping(); err != nil {
			Logger().LogError("ms engine ping", err)
			return msEngine
		}
	}
	return msEngine
}

func OracleEngine() *xorm.Engine  {
	if msEngine == nil {
		var err error
		msEngine, err = xorm.NewEngine(Config().OracleDriverName, Config().OracleConnUrl)
		//"song/123456@192.168.0.105:1521/ORCL"
		Logger().LogError("new orcale engine failed:", msEngine,err)
		//SnakeMapper 支持struct为驼峰式命名，表结构为下划线命名之间的转换，这个是默认的Maper
		//映射同名设置默认
		msEngine.SetMapper(core.SameMapper{})
		if err != nil {
			Logger().LogError("new orcale engine failed:", err)
			return nil
		}
		if err := msEngine.Ping(); err != nil {
			Logger().LogError("ms engine ping", err)
			return msEngine
		}
	}
	return msEngine
}