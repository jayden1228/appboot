package tool

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gogf/gf/os/gfile"
)

//检查某字符是否存在文件里
func CheckFileContainsChar(filename, s string) bool {
	data := gfile.GetContents(filename)
	if len(data) > 0 {
		return strings.LastIndex(data, s) > 0
	}
	return false
}

//拼接特殊字符串
func FormatField(field string, formats []string) string {
	if len(formats) == 0 {
		return ""
	}
	buf := bytes.Buffer{}
	for key := range formats {
		buf.WriteString(fmt.Sprintf(`%s:"%s" `, formats[key], field))
	}
	return "`" + strings.TrimRight(buf.String(), " ") + "`"
}

//添加注释
func AddToComment(s string, suff string) string {
	if strings.EqualFold(s, "") {
		return ""
	}
	return "// " + s + suff
}

//判断是否包存在某字符
func InArrayString(str string, arr []string) bool {
	for _, val := range arr {
		if val == str {
			return true
		}
	}
	return false
}

//检查字符串,去掉特殊字符
func CheckCharDoSpecial(s string, char byte, regs string) string {
	reg := regexp.MustCompile(regs)
	var result []string
	if arr := reg.FindAllString(s, -1); len(arr) > 0 {
		buf := bytes.Buffer{}
		for key, val := range arr {
			if val != string(char) {
				buf.WriteString(val)
			}
			if val == string(char) && buf.Len() > 0 {
				result = append(result, buf.String())
				buf.Reset()
			}
			//处理最后一批数据
			if buf.Len() > 0 && key == len(arr)-1 {
				result = append(result, buf.String())
			}
		}
	}
	return strings.Join(result, string(char))
}
func CheckCharDoSpecialArr(s string, char byte, reg string) []string {
	s = CheckCharDoSpecial(s, char, reg)
	return strings.Split(s, string(char))
}

// 添加``符号
func AddQuote(str string) string {
	return "`" + str + "`"
}

// 去掉 `符号
func CleanQuote(str string) string {
	return strings.Replace(str, "`", "", -1)
}

func Gofmt(path string) bool {
	if !ExecCommand("goimports", "-l", "-w", path) {
		if !ExecCommand("gofmt", "-l", "-w", path) {
			return ExecCommand("go", "fmt", path)
		}
	}
	return true
}

func ExecCommand(name string, args ...string) bool {
	cmd := exec.Command(name, args...)
	_, err := cmd.Output()
	if err != nil {
		return false
	}
	return true
}
