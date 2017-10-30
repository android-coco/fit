package fit

import (
	"crypto/rand"
	r "math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/tealeg/xlsx"
	"errors"
	"fmt"
	"io/ioutil"
)

var alphaNum = []byte(`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`)

// RandomCreateBytes generate random []byte by specify chars.
func RandomCreateBytes(n int, alphabets ...byte) []byte {
	if len(alphabets) == 0 {
		alphabets = alphaNum
	}
	var bytes = make([]byte, n)
	var randBy bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randBy = true
	}
	for i, b := range bytes {
		if randBy {
			bytes[i] = alphabets[r.Intn(len(alphabets))]
		} else {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return bytes
}
func CheckError(err error) {
	if err != nil {
		Logger().LogError("Fatal error ", err.Error())
		os.Exit(1)
	}
}
func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GetParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Logger().LogError("File Path:", err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//key 标题,vlue数据
func ExportExcel(r *Response, titles []string, data [][]string, filename string) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return err
	}
	if titles == nil || data == nil {
		return errors.New("参数不能为空！")
	}
	if len(titles) != len(data[0]) {
		return errors.New("数据格式错误！")
	}
	//加入标题
	if titles != nil && len(titles) > 0 {
		row = sheet.AddRow()
		for _, k := range titles {
			cell = row.AddCell()
			cell.Value = k
		}
	}
	//加入内容
	if data != nil && len(data) > 0 {
		for _, k := range data {
			row = sheet.AddRow()
			for _, k1 := range k {
				cell = row.AddCell()
				cell.Value = k1
			}
		}
	}
	err = file.Save(filename + ".xlsx") //生成临时文件
	if err != nil {
		return err
	}
	file1, err1 := os.Open("./" + filename + ".xlsx") //打开临时文件
	defer os.Remove("./" + filename + ".xlsx")        //删除临时文件
	defer file1.Close()                               //关闭文件
	if err1 != nil {
		return err1
	}
	b, _ := ioutil.ReadAll(file1)
	r.Writer().Header().Set("Accept-Ranges", "bytes")
	r.Writer().Header().Set("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	r.Writer().Header().Set("Cache-Control", "max-age=0")
	r.Writer().Header().Set("Cache-Control", "max-age=1")
	r.Writer().Header().Set("Pragma", "no-cache")
	r.Writer().Header().Set("Expires", "0")
	r.Writer().Header().Set("Content-Disposition", "attachment;filename="+fmt.Sprintf("%s", filename+".xlsx")) //文件名称
	r.Writer().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	_, err = r.Writer().Write(b) //输出到浏览器
	if err != nil {
		return err
	}
	return nil
}
