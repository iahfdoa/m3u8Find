package result

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"sync"
)

type Result struct {
	sync.RWMutex
	TagUserinfo map[string][]UsersInfo
}

func (r *Result) HasTags() bool {
	r.RLock()
	defer r.RUnlock()

	if len(r.TagUserinfo) <= 0 {
		return false
	} else {
		return true
	}
}

func (r *Result) HasTag(tag string) bool {
	r.RLock()
	defer r.RUnlock()
	_, ok := r.TagUserinfo[tag]
	if ok {
		return true
	} else {
		return false
	}
}
func (r *Result) AddTag(tag string, user []UsersInfo) {
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.TagUserinfo[tag]; !ok {
		r.TagUserinfo[tag] = []UsersInfo{}
	}
	r.TagUserinfo[tag] = user
}
func (r *Result) GetTag() chan *map[string][]UsersInfo {
	r.RLock()

	out := make(chan *map[string][]UsersInfo)

	go func() {
		defer close(out)
		defer r.RUnlock()

		for k, v := range r.TagUserinfo {
			m := make(map[string][]UsersInfo)
			m[k] = v
			out <- &m
		}
	}()

	return out
}

type UsersInfo struct {
	Id                   int    `json:"id" csv:"主播ID"`
	Username             string `json:"username" csv:"名字"`
	Country              string `json:"country" csv:"国家"`
	IsLive               bool   `json:"isLive" csv:"是否在直播"`
	IsOnline             bool   `json:"isOnline" csv:"是否在线"`
	Gender               string `json:"gender" csv:"性别"`
	PreviewUrlThumbBig   string `json:"previewUrlThumbBig" csv:"图片封面大"`
	PreviewUrlThumbSmall string `json:"previewUrlThumbSmall" csv:"图片封面小"`
	HlsPlaylist          string `json:"hlsPlaylist" csv:"直播地址"`
	ParentTag            string `csv:"分类"`
}

//type UsersInfo struct {
//	Id                   int    `json:"id" csv:"id"`
//	Username             string `json:"username" csv:"username"`
//	Country              string `json:"country" csv:"country"`
//	IsLive               bool   `json:"isLive" csv:"isLive"`
//	IsOnline             bool   `json:"isOnline" csv:"isOnline"`
//	Gender               string `json:"gender" csv:"gender"`
//	PreviewUrlThumbBig   string `json:"previewUrlThumbBig" csv:"previewUrlThumbBig"`
//	PreviewUrlThumbSmall string `json:"previewUrlThumbSmall" csv:"previewUrlThumbSmall"`
//	HlsPlaylist          string `json:"hlsPlaylist" csv:"hlsPlaylist"`
//}

func HasInArray(Id int, users []UsersInfo) (ok bool) {
	ok = false
	for _, user := range users {
		if user.Id == Id {
			ok = true
			break
		}
	}
	return
}

var NumberOfCsvFieldsErr = errors.New("导出的字段与csv标记不匹配")

func (u *UsersInfo) CSVHeaders() ([]string, error) {
	ty := reflect.TypeOf(*u)
	var headers []string
	for i := 0; i < ty.NumField(); i++ {
		headers = append(headers, ty.Field(i).Tag.Get("csv"))
	}
	if len(headers) != ty.NumField() {
		return nil, NumberOfCsvFieldsErr
	}
	return headers, nil
}

//
func (u *UsersInfo) CSVFields() ([]string, error) {
	var fields []string
	vl := reflect.ValueOf(*u)
	for i := 0; i < vl.NumField(); i++ {
		fields = append(fields, fmt.Sprint(vl.Field(i).Interface()))
	}
	if len(fields) != vl.NumField() {
		return nil, NumberOfCsvFieldsErr
	}
	return fields, nil
}
