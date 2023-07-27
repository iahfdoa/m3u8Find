package scan

import (
	"encoding/json"
	"fmt"
	"github.com/iahfdoa/m3u8Find/pkg/result"
	"github.com/iahfdoa/m3u8Find/pkg/utils"
	"github.com/projectdiscovery/gologger"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

type Scanner struct {
	client      *http.Client
	ScanResults *result.Result
	primaryTag  string
	ModelsRoute string
	Phase       Phase
	debug       bool
	retries     int
	rate        int
}

func (s *Scanner) GetModelsRouteResult(tag string, limit int) ([]result.UsersInfo, bool) {
	client := s.client
	primaryTag := s.primaryTag
	parentTag := tag
	payload := s.ModelsRoute + "?primaryTag=" + primaryTag + "&parentTag=" + parentTag + "&limit=" + strconv.Itoa(limit)
	request, _ := http.NewRequest("GET", payload, nil)
	dial, err := utils.ReconDial(client, request, 1, s.retries)
	if err != nil {
		gologger.Debug().Msgf("%s\n", err.Error())
		return nil, false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			gologger.Debug().Msg(err.Error())
		}
	}(dial.Body)

	type modelsResult struct {
		Models []map[string]interface{} `json:"models"`
	}
	l := &modelsResult{
		Models: nil,
	}
	//var l modelsResult
	dialStr, _ := ioutil.ReadAll(dial.Body)
	//fmt.Println(string(dialStr))
	err = json.Unmarshal(dialStr, &l)
	if err != nil {
		gologger.Debug().Msgf("反序列化失败-> %s", err.Error())
		return nil, false
	}
	var users []result.UsersInfo

	for _, sliceModels := range l.Models {
		user := result.UsersInfo{
			PreviewUrlThumbBig:   fmt.Sprintf("%v", sliceModels["previewUrlThumbBig"]),
			PreviewUrlThumbSmall: sliceModels["previewUrlThumbSmall"].(string),
			Username:             sliceModels["username"].(string),
			Id:                   int(sliceModels["id"].(float64)),
			IsLive:               sliceModels["isLive"].(bool),
			IsOnline:             sliceModels["isOnline"].(bool),
			Gender:               sliceModels["gender"].(string),
			HlsPlaylist:          utils.RemoveHlsPlaylistP(sliceModels["hlsPlaylist"].(string)),
			Country:              sliceModels["country"].(string),
			ParentTag:            parentTag,
		}
		users = append(users, user)

	}
	gologger.Debug().Msgf("请求成功 %s ,获得 %d 个主播信息", payload, len(users))
	return users, true
}

type UsersInfo struct {
	Id                   int    `json:"id" csv:"id"`
	Username             string `json:"username" csv:"username"`
	Country              string `json:"country" csv:"country"`
	IsLive               bool   `json:"isLive" csv:"isLive"`
	IsOnline             bool   `json:"isOnline" csv:"isOnline"`
	Gender               string `json:"gender" csv:"gender"`
	PreviewUrlThumbBig   string `json:"previewUrlThumbBig" csv:"previewUrlThumbBig"`
	PreviewUrlThumbSmall string `json:"previewUrlThumbSmall" csv:"previewUrlThumbSmall"`
	HlsPlaylist          string `json:"hlsPlaylist" csv:"hlsPlaylist"`
}

func NewScanner(options *Options) *Scanner {

	debug := options.Debug
	retries := options.Retries
	rate := options.Rate
	modelsRoute := options.ModelsRoute
	scanResult := &result.Result{
		RWMutex:     sync.RWMutex{},
		TagUserinfo: make(map[string][]result.UsersInfo),
	}
	return &Scanner{
		client:      options.Client,
		ScanResults: scanResult,
		primaryTag:  options.PrimaryTag,
		ModelsRoute: modelsRoute,
		Phase:       Phase{},
		debug:       debug,
		retries:     retries,
		rate:        rate,
	}
}

const (
	Init State = iota
	HostDiscovery
	Scan
	Done
	Guard
)

type State int
type Phase struct {
	sync.RWMutex
	State
}

func (phase *Phase) Is(state State) bool {
	phase.RLock()
	defer phase.RUnlock()

	return phase.State == state
}

func (phase *Phase) Set(state State) {
	phase.Lock()
	defer phase.Unlock()

	phase.State = state
}
