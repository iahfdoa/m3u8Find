package runner

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/iahfdoa/m3u8Find/pkg/result"
	"github.com/iahfdoa/m3u8Find/pkg/scan"
	"github.com/iahfdoa/m3u8Find/pkg/utils"
	"github.com/projectdiscovery/fileutil"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/ratelimit"
	"github.com/projectdiscovery/stringsutil"
	"github.com/remeh/sizedwaitgroup"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Runner struct {
	options    *UserOptions
	wgScan     sizedwaitgroup.SizedWaitGroup
	limiter    *ratelimit.Limiter
	client     *http.Client
	ParentTags *map[string]map[string]int
	scanner    *scan.Scanner
}

func (runner *Runner) Run() {
	runner.wgScan = sizedwaitgroup.New(runner.options.Rate)
	runner.limiter = ratelimit.New(context.Background(), uint(runner.options.Rate), time.Second)
	runner.scanner.Phase.Set(scan.Scan)
	for name, v := range *runner.ParentTags {
		runner.limiter.Take()
		runner.wgScan.Add()
		limit := v["modelsLive"]
		parentTag := name
		go runner.handle(limit, parentTag)

	}
	runner.wgScan.Wait()
	runner.scanner.Phase.Set(scan.Done)
	runner.handleOutput(runner.scanner.ScanResults)

}
func (runner *Runner) handleOutput(scanResults *result.Result) {
	var (
		output string
	)
	output = runner.options.Output
	m3u8 := runner.options.OutputM3u8
	csv := runner.options.OutputAllAnchors
	if !fileutil.FolderExists(output) {
		mkdirErr := os.MkdirAll(output, 0700)
		if mkdirErr != nil {
			gologger.Error().Msgf("Could not create output folder %s: %s\n", output, mkdirErr)
			return
		}
	}
	groupRole := runner.options.GroupRole
	filename := filepath.Join(output, "/"+runner.options.PrimaryTag)
	outputAllUserOnFile := runner.options.OutputAllUserOnFile
	var users []result.UsersInfo
	if m3u8 && scanResults.HasTags() {
		fa := filename + ".m3u8"
		for userinfo := range scanResults.GetTag() {
			for k, user := range *userinfo {
				for _, u := range user {
					if !result.HasInArray(u.Id, users) {
						users = append(users, u)
					}
				}
				if !outputAllUserOnFile {
					f := filename + "_" + k + ".m3u8"
					err := outputM3u8Users(f, groupRole, user)
					if err != nil {
						gologger.Debug().Msgf("文件写入失败 %s-> %s", f, err.Error())
						continue
					} else {
						gologger.Info().Msgf("文件写入成功 -> %s", f)
					}
				}

			}
		}
		if outputAllUserOnFile {
			err := outputM3u8AllUsers(fa, &users)
			if err != nil {
				gologger.Debug().Msgf("文件写入失败 %s-> %s", filename, err.Error())
				return
			} else {
				gologger.Info().Msgf("文件写入成功 -> %s", fa)
			}

		}

	}
	if csv && scanResults.HasTags() {
		for userinfo := range scanResults.GetTag() {
			for _, user := range *userinfo {
				for _, u := range user {
					if !result.HasInArray(u.Id, users) {
						users = append(users, u)
					}
				}
			}
		}
		f := filename + ".csv"
		err := outputCsvAllUsers(f, users)
		if err != nil {
			gologger.Debug().Msgf("文件写入失败 %s-> %s", f, err.Error())

		} else {
			gologger.Info().Msgf("文件写入成功 -> %s", f)
		}

	}

	if !m3u8 && !csv && scanResults.HasTags() {
		var l int
		for userinfo := range scanResults.GetTag() {
			for k, user := range *userinfo {
				l += len(user)
				gologger.Debug().Msgf("收集 %s 类型 -> %d 个", k, len(user))
			}
		}
		gologger.Info().Msgf("收集总主播 %d 个\n", l)
		gologger.Info().Msgf("查看详情请使用 -om 或 -oa 或 -oaa 参数\n")
	}

}
func (runner *Runner) handle(limit int, tag string) {
	defer runner.wgScan.Done()
	if runner.scanner.ScanResults.HasTag(tag) {
		return
	}
	runner.limiter.Take()
	userinfo, ok := runner.scanner.GetModelsRouteResult(tag, limit)
	if ok {
		runner.scanner.ScanResults.AddTag(tag, userinfo)
	}
}
func NewTag(client *http.Client, options *UserOptions) (*map[string]map[string]int, error) {

	primaryTag := options.PrimaryTag
	tagsRoute := options.TagsRoute
	mainUrl := options.MainUrl
	m := make(map[string]map[string]int)
	for _, v := range ParentTags {
		t := make(map[string]int)
		t["modelsLive"] = options.Limit
		m[v] = t
	}
	payload := mainUrl + tagsRoute + "?primaryTag" + primaryTag
	request, err := http.NewRequest("GET", payload, nil)

	if err != nil {
		return nil, err
	}
	liveGet, err := utils.ReconDial(client, request, 1, options.ReconnectionNum)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			gologger.Debug().Msg(err.Error())
		}
	}(liveGet.Body)
	var l utils.Live
	liveStrJson, _ := ioutil.ReadAll(liveGet.Body)
	err = json.Unmarshal(liveStrJson, &l)
	if err != nil {
		gologger.Debug().Msg("请求父类Tag失败")
		return nil, err
	}
	for k, v := range l.LiveTagDetails {
		value, ok := m[k]
		if !ok {
			m[k] = v
		} else if options.AutoLimit {
			m[k]["modelsLive"] = value["modelsLive"]
		}
	}
	return &m, nil
}
func NewRunner(options *UserOptions) (*Runner, error) {
	runner := &Runner{options: options}
	client, err := NewClient(options)
	if err != nil {
		return nil, err
	}

	scanner := scan.NewScanner(&scan.Options{
		Timeout:     time.Duration(options.Timeout),
		Retries:     options.ReconnectionNum,
		Rate:        options.Rate,
		Debug:       options.Debug,
		Client:      client,
		PrimaryTag:  options.PrimaryTag,
		ModelsRoute: options.MainUrl + options.ModelsRoute,
	})
	runner.scanner = scanner
	runner.client = client
	runner.ParentTags, err = NewTag(client, options)
	if err != nil {
		return nil, err
	}
	return runner, nil

}
func NewClient(options *UserOptions) (*http.Client, error) {
	t := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if options.Proxy != "" {
		proxyURL, err := url.Parse(options.Proxy)
		if err != nil {
			return nil, err
		}
		if options.ProxyAuth != "" {
			proxyURL.User = url.UserPassword(stringsutil.SplitAny(options.ProxyAuth)[0], stringsutil.SplitAny(options.ProxyAuth)[1])
		}
		t.Proxy = http.ProxyURL(proxyURL)
	}
	return &http.Client{Transport: t, Timeout: time.Duration(options.Timeout) * time.Second}, nil

}
