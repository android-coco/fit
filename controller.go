package fit

import (
	"html/template"
	"net/http"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"errors"
)

type ControllerInterface interface {
	Get(w *Response, r *Request, p Params)
	Post(w *Response, r *Request, p Params)
	Head(w *Response, r *Request, p Params)
	Options(w *Response, r *Request, p Params)
	Put(w *Response, r *Request, p Params)
	Patch(w *Response, r *Request, p Params)
	Delete(w *Response, r *Request, p Params)
}

// rendering Json data (© jp)
type ConterllerRenderingJson interface {
	RenderingJsonAutomatically(result int, errMsg string)
	RenderingJson(result int, errMsg string, datas interface{})
}

//json result struct
type Result struct {
	Result   int         `json:"result"` // 错误码
	ErrorMsg string      `json:"errmsg"` // 错误描述
	Datas    interface{} `json:"datas"`  //数据
}

type Data map[interface{}]interface{}
type Controller struct {
	Data     Data //界面数据
	JsonData Result
	// template data
	TplName   string
	ViewPath  string
	TplPrefix string
	TplExt    string
}

// Get adds a request function to handle GET request.
func (c *Controller) Get(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

// Post adds a request function to handle post request.
func (c *Controller) Post(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

// Delete adds a request function to handle delete request.
func (c *Controller) Delete(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

// Put adds a request function to handle put request.
func (c *Controller) Put(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

// Head adds a request function to handle head request.
func (c *Controller) Head(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

// Patch adds a request function to handle patch request.
func (c *Controller) Patch(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

// Options adds a request function to handle options request.
func (c *Controller) Options(w *Response, r *Request, p Params) {
	http.Error(w.Writer(), "Method Not Allowed", 405)
}

// 重定向
func (c *Controller) Redirect(w *Response, r *Request, url string, code int) {
	http.Redirect(w.writer, r.Request, url, code)
}

func (c *Controller) LoadView(w *Response, tplname ...string) {
	var err error
	var t *template.Template
	var args []string
	for _, arg := range tplname {
		arg = "view/" + arg
		args = append(args, arg)
	}

	t, err = template.ParseFiles(args...) //从文件创建一个模板

	if err != nil {
		Logger().LogError("Fatal error ", err.Error())
	}

	err = t.Execute(w.Writer(), c.Data)

	if err != nil {
		Logger().LogError("Fatal error ", err.Error())
	}
}

/*Checking session and then load a special view */
func (c *Controller) LoadViewSafely(w *Response, r *Request, tplname ...string) (success bool) {
	session, err_s := GlobalManager().SessionStart(w, r)
	if err_s != nil || session == nil {
		// Session失效 重新登录
		Logger().LogError("JP ", "Session失效 重新登录"+err_s.Error())
		return false
	} else {
		userinfo := session.Get("UserInfo")
		if userinfo == nil {
			// 未登录
			Logger().LogDebug("JP ", "未登录")
			return false
		} else {
			// 已登录
			c.LoadView(w, tplname...)
			Logger().LogDebug("JP ", "已登录")
			return true
		}
	}
}

// json 输出函数
func (c *Controller) ResponseToJson(w *Response) {
	if c.JsonData.Datas == nil {
		c.JsonData.Datas = []interface{}{}
	}
	b, err := json.Marshal(c.JsonData)
	if err != nil {
		fmt.Fprint(w.Writer(), err.Error())
		return
	}
	fmt.Fprint(w.Writer(), string(b))
}

func (c *Controller) RenderingJson(result int, errMsg string, datas interface{}) {
	c.JsonData.Datas = datas
	c.JsonData.ErrorMsg = errMsg
	c.JsonData.Result = result
}

func (c *Controller) RenderingJsonAutomatically(result int, errMsg string) {
	c.RenderingJson(result, errMsg, make(map[string]interface{}, 0))
}

//reflect 接受较多的参数时使用
func (c *Controller) FitSetStruct(bean interface{}, r *Request) (err error) {
	// 如果不是指针 直接返回
	if reflect.TypeOf(bean).Kind() != reflect.Ptr {
		Logger().LogError("fit set struct err:", reflect.TypeOf(bean).Kind())
		return
	}
	rt := reflect.TypeOf(bean).Elem()
	rv := reflect.ValueOf(bean).Elem()


	for index := 0; index < rt.NumField(); index++ {
		sName := rt.Field(index).Name
		tags := splitTag(rt.Field(index).Tag.Get("fit"))

		if rv.Field(index).CanSet() {
			if len(tags) > 0 {
				if tags[0] == "-" { // 不映射
					continue
				}

				if len(tags) >= 2 {
					errors.New("目前未定义多重tag")
				}
				if tags[0] != "" {
					sName = tags[0]
				}
			}
			//fmt.Println("model struct kind:", rv.Field(index).Kind(), "------ tag name:", sName)
			switch rv.Field(index).Kind() {
			case reflect.Int:
				val := int64(r.FormIntValue(sName))
				rv.Field(index).SetInt(val)
			case reflect.Int64:
				rv.Field(index).SetInt(r.FormInt64Value(sName))
			case reflect.String:
				rv.Field(index).SetString(r.FormValue(sName))
			case reflect.Struct:
				//datetime, _ := time.ParseInLocation("2006-01-02 15:04:05", r.FormValue("datetime"), time.Local)

				//rv.Field(index).SetString(r.FormValue("datetime"))

				fmt.Println("f.field:", rv.Field(index))
			default:
				err = errors.New("undeclared reflect type")
				return
			}
		} else {
			err = errors.New("set failed")
			return
		}

	}
	return
}

func splitTag(tag string) (tags []string)  {
	tag = strings.TrimSpace(tag) // 去除空格多余
	tags = strings.Split(tag, " ")
	return
}

/*
		case "VAA01":
			// 病人ID
			//rt.Field(index).Tag.Get("fit")
			VAA01 := r.FormInt64Value("pid")
			nrl3Mod.VAA01 = VAA01
		case "BCK01":
			// 科室ID
			BCK01 := r.FormIntValue("did")
			nrl3Mod.BCK01 = BCK01
		case "BCE01A":
			// 护士ID
			nrl3Mod.BCE01A = r.FormValue("uid")
		case "BCE03A":
			// 护士名
			nrl3Mod.BCE03A = r.FormValue("username")
		case "DateTime":
			// 记录时间
			datetime, err4 := time.ParseInLocation("2006-01-02 15:04:05", r.FormValue("datetime"), time.Local)
			if err4 != nil {
				fit.Logger().LogError("NRL3 :", err4)
				errflag = true
			}
			nrl3Mod.DateTime = datetime

		if sName == "VAA01" {
			// 病人ID
			VAA01 := r.FormInt64Value("pid")
			nrl3Mod.VAA01 = VAA01
		} else {
			if rv.Field(index).CanSet() {
				//val := int(r.FormIntValue(sName))
				fmt.Println("model struct kind:", rv.Field(index).Kind())
				switch rv.Field(index).Kind() {
				case reflect.Int:
					rv.Field(index).SetInt(r.FormInt64Value(sName))
				case reflect.Int64:
					rv.Field(index).SetInt(r.FormInt64Value(sName))
				case reflect.String:
					rv.Field(index).SetString(r.FormValue(sName))
				}
			} else {
				fmt.Println("set failed")
			}
		}*/