package runner

import (
	"encoding/json"
)

type LiveTagGroups struct {
	Alias string   `json:"alias"`
	Tags  []string `json:"tags"`
}

type LiveTagsInfo struct {
	LiveTagGroups  []*LiveTagGroups          `json:"liveTagGroups"`
	LiveTagDetails map[string]map[string]int `json:"liveTagDetails"`
}
type UsersInfo struct {
	Id                   int    `json:"id" csv:"主播ID"`
	Username             string `json:"username" csv:"名字"`
	Country              string `json:"country" csv:"国家"`
	IsLive               bool   `json:"isLive" csv:"是否在直播"`
	IsOnline             bool   `json:"isOnline" csv:"是否在线"`
	Gender               string `json:"gender" csv:"性别"`
	PreviewUrlThumbBig   string `json:"previewUrlThumbBig" csv:"图片封面"`
	PreviewUrlThumbSmall string `json:"previewUrlThumbSmall"`
	HlsPlaylist          string `json:"hlsPlaylist" csv:"直播地址"`
}

func (u *UsersInfo) JSON() ([]byte, error) {
	return json.Marshal(u)
}
