package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

type SyncWriter struct {
	m      sync.Mutex
	Writer io.Writer
}
type Live struct {
	LiveTagDetails map[string]map[string]int `json:"liveTagDetails"`
}

func (w *SyncWriter) Write(b []byte) (n int, err error) {
	w.m.Lock()
	defer w.m.Unlock()
	return w.Writer.Write(b)
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func RandStr(n int) string {
	cb := make(chan byte)
	var b []byte
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	go func() {
		for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
			if remain == 0 {
				cache, remain = src.Int63(), letterIdMax
			}
			if idx := int(cache & letterIdMask); idx < len(letters) {
				cb <- letters[idx]
				i--
			}
			cache >>= letterIdBits
			remain--
		}
		close(cb)
	}()

	for c := range cb {
		b = append(b, c)
	}
	return string(b)
}

func StrInArray(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	//index的取值：[0,len(str_array)]
	if index < len(strArray) && strArray[index] == target { //需要注意此处的判断，先判断 &&左侧的条件，如果不满足则结束此处判断，不会再进行右侧的判断
		return true
	}
	return false
}

func RemoveHlsPlaylistP(p string) string {
	compile := regexp.MustCompile(`^(.*/)(\d+)([_\dp]{0,6})(\.m3u8)$`)
	group := compile.FindAllStringSubmatch(p, -1)

	if len(group) < 1 && len(group[0]) < 5 {

		return p
	}
	var np string
	if group[0][3] == "_1080p" {
		np = p
	} else {
		np = group[0][1] + group[0][2] + group[0][4]
	}
	return np
}

// RegGender
// p  需要 找寻的字符串
// reg 正则匹配
// n 返回 正则中的数组第几位 若 n<=0  则 默认返回第0位
func RegGender(p string, reg string, n int) string {
	compile := regexp.MustCompile(reg)
	if compile == nil {
		return ""
	}
	np := compile.FindAllStringSubmatch(p, -1)
	if n > 0 {
		return np[0][n]
	} else {
		return np[0][0]
	}
}
func WriteAppend(data, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("文件打开失败 -> %s", err.Error())
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(data + "\n")

	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	return nil

}
func WriteFile(data string, filename string) error {
	f, err := os.Create(filename)
	//datas := stringsutil.SplitAny(data, "\n")
	defer f.Close()
	if err != nil {
		if os.IsPermission(err) {
			return errors.Wrap(err, "permission error")
		}
		return errors.New("打开文件失败: " + err.Error())
	}
	dataReader := bytes.NewReader([]byte(data))
	//writer := io.Writer()
	_, err = io.Copy(f, dataReader)
	if err != nil {
		return errors.New("写入失败: " + err.Error())
	}
	return nil

}

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]

}

// ReconDial 多次请求
func ReconDial(client *http.Client, req *http.Request, i int, max int) (*http.Response, error) {
	get, err := client.Do(req)
	if err != nil {
		if i < max {
			i++
			get, err = ReconDial(client, req, i, max)
		} else {
			err = errors.New(fmt.Sprintf("链接失败超过%v次-> `%v` ", max, req.URL.String()))
		}
	}
	return get, err
}
