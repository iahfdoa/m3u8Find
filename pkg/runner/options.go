package runner

import (
	"fmt"
	"github.com/iahfdoa/m3u8Find/pkg/utils"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"os"
	"path/filepath"
)

type UserOptions struct {
	// Input 配置

	// Config 配置
	MainUrl         string // MainUrl https://zh.superchat.live/
	Proxy           string // 代理
	ProxyAuth       string // 代理认证（username:password）
	Rate            int    // 每秒请求次数
	ReconnectionNum int    // 重连次数
	Timeout         int    // 连接超时时间
	GroupRole       string // 分组规则
	AutoLimit       bool   // 自动校验
	Uniq            string // Uniq 字母、数字混合16位随机数

	// Output 输出
	Output              string // 输出位置
	OutputM3u8          bool   // OutputM3u8 是否输出M3u8信息
	OutputAllUserOnFile bool   // OutputAllUserOnFile 是否输出所有主播到一个m3u8
	OutputAllAnchors    bool   // OutputAllAnchors  是否输出个主播的详细信息

	// ApiRoutes 路由规则
	ModelsRoute      string // ModelsRoute 模式路由
	TagsRoute        string // TagRoute 标签路由
	UsersRoute       string // UserRoute 用户路由
	ModelsCountRoute string // ModelsCountRoute 模式路由直播总和

	// V2ApiRoutes V2路由规则
	ConfigRouteV2      string // ConfigRoute 配置路由
	ModelsRouteV2      string // ModelsRouteV2 模式路由 V2
	ModelsCountRouteV2 string // ModelsCountRouteV2 模式路由直播总和 V2
	TagsRouteV2        string // TagsRouteV2 标签路由 V2
	UsersRouteV2       string // UsersRouteV2 用户路由 V2

	// ModelsRoutes 参数 设定
	PrimaryTag string // PrimaryTag 主标签 girls,couples,men,trans
	Limit      int    // Limit 模块的最大请求
	ParentTag  string // ParentTag 父标签 tagLanguageUkrainian(乌克兰),tagLanguageChinese(中国),autoTagNew(新的),autoTagVr,fetishes,autoTagRecordablePublic

	// debug
	HealthCheck bool // 检查网络是否正常,代理是否可用.
	Debug       bool // 调试模式

}

func (userOptions *UserOptions) DefaultUserOptions() {
	if userOptions.Output == "" {
		userOptions.Output = filepath.Join(".", "output")
		userOptions.Output, _ = filepath.Abs(userOptions.Output)
	}
	userOptions.Uniq = utils.RandStr(16)
	userOptions.MainUrl = MainUrl
	userOptions.ModelsRoute = ModelsRoute
	userOptions.TagsRoute = TagsRoute
	userOptions.UsersRoute = UsersRoute
	userOptions.ModelsCountRoute = ModelsCountRoute
	userOptions.ConfigRouteV2 = ConfigRouteV2
	userOptions.ModelsRouteV2 = ModelsRouteV2
	userOptions.ModelsCountRoute = ModelsCountRouteV2
	userOptions.TagsRouteV2 = TagsRouteV2
	userOptions.UsersRouteV2 = UsersRouteV2

}
func ParseUserOptions() *UserOptions {
	userOptions := &UserOptions{}
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`m3u8Find 是一块爬虫工具`)
	// // Input 配置
	//	Proxy     string // 代理
	//	ProxyAuth string // 代理认证（username:password）
	flagSet.CreateGroup(`Input`, `用户输入`,
		flagSet.StringVarP(&userOptions.Proxy, "proxy", "p", "", "代理"),
		flagSet.StringVarP(&userOptions.ProxyAuth, "proxy-auth", "pa", "", "代理认证（username:password）"),
	)
	// 	// Config 配置
	//	MainUrl   string // MainUrl https://zh.superchat.live/
	//	Uniq      string // Uniq 字母、数字混合16位随机数
	flagSet.CreateGroup(`Config`, `配置`,
		flagSet.BoolVarP(&userOptions.AutoLimit, "auto-limit", "al", AutoLimit, "自动校验Limit"),
		flagSet.IntVarP(&userOptions.Rate, "rate", "r", Rate, "每秒请求次数"),
		flagSet.IntVarP(&userOptions.ReconnectionNum, "recon", "rn", ReconnectionNum, "重连次数"),
		flagSet.IntVarP(&userOptions.Timeout, "timeout", "t", Timeout, "连接超时时间"),
		flagSet.StringVarP(&userOptions.GroupRole, "group-role", "gr", GroupRole, ReturnU("分组规则", &GroupRoles)),
	)

	// 	// 	// Output 输出
	//	Output           string // 输出位置
	//	OutputM3u8       bool   // OutputM3u8 是否输出M3u8信息
	//	OutputAllAnchors bool   // OutputAllAnchors  是否输出个主播的详细信息
	flagSet.CreateGroup("Output", `输出`,
		flagSet.StringVarP(&userOptions.Output, "output", "o", Output, "输出路径"),
		flagSet.BoolVarP(&userOptions.OutputM3u8, "output-m3u8", "om", OutputM3u8, "是否输出M3u8信息"),
		flagSet.BoolVarP(&userOptions.OutputAllUserOnFile, "outpit-all", "oa", OutputAllUserOnFile, "是否输出主播到一个文件夹"),
		flagSet.BoolVarP(&userOptions.OutputAllAnchors, "output-all-anchors", "oaa", OutputAllAnchors, "是否输出个主播的详细信息"),
	)
	flagSet.CreateGroup("ModelsRoutes", `ModelsRoutes 路由参数设定`,
		flagSet.StringVarP(&userOptions.PrimaryTag, "primary-tag", "pt", PrimaryTag, ReturnU("主标签", &PrimaryTags)),
		flagSet.IntVar(&userOptions.Limit, "limit", Limit, `模块的最大请求 (主播个数)`),
	)
	flagSet.CreateGroup("debug", "debug",
		flagSet.BoolVarP(&userOptions.HealthCheck, "hc", "health-check", false, "运行诊断检查"),
		flagSet.BoolVarP(&userOptions.Debug, "debug", "d", false, "调试模式"),
	)
	_ = flagSet.Parse()
	var err error
	if userOptions.Output != "" {
		userOptions.Output, _ = filepath.Abs(userOptions.Output)
	}

	userOptions.DefaultUserOptions()
	if userOptions.HealthCheck {
		gologger.Print().Msgf("%s\n", userOptions.DoHealthCheck())
		os.Exit(0)
	}
	if userOptions.Debug {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
	} else {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelWarning)
	}
	err = userOptions.validateOptions()
	if err != nil {
		gologger.Fatal().Msgf("程序正在退出: %s\n", err)
		os.Exit(0)
	}
	if userOptions.Proxy == "" {
		gologger.Warning().Msg("连接需要科学上网")
	} else {
		gologger.Info().Msg(fmt.Sprintf("proxy setting ok : %v", userOptions.Proxy))

	}
	return userOptions
}
