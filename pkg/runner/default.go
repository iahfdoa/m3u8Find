package runner

import "strings"

var PrimaryTags = []string{`girls`, `couples`, `trans`, `men`}
var ParentTags = []string{
	"tagLanguageUkrainian",
	"tagLanguageChinese",
	"autoTagNew",
	"autoTagVr",
	"fetishes",
	"autoTagRecordablePublic",
	"ageTeen",
	"ageYoung",
	"ageMilf",
	"ageMature",
	"ageOld",
	"ethnicityMiddleEastern",
	"ethnicityAsian",
	"ethnicityEbony",
	"ethnicityIndian",
	"ethnicityLatino",
	"ethnicityWhite",
	"bodyTypePetite",
	"bodyTypeAthletic",
	"bodyTypeMedium",
	"bodyTypeCurvy",
	"bodyTypeBBW",
	"bodyTypeBig",
	"hairColorBlonde",
	"hairColorBlack",
	"hairColorRed",
	"hairColorColorful",
	"privatePriceEight",
	"privatePriceSixteenToTwentyFour",
	"privatePriceThirtyTwoSixty",
	"privatePriceNinetyPlus",
	"autoTagBestPrivates",
	"autoTagRecordablePrivate",
	"autoTagSpy",
	"autoTagP2P",
}
var GroupRoles = []string{`gender`, `country`}

const (
	Version             = "0.0.1"
	MainHost            = `zh.superchat.live`
	MainUrl             = `https://` + MainHost
	ReconnectionNum     = 5
	Timeout             = 3
	GroupRole           = `gender`
	AutoLimit           = false
	Rate                = 150
	OutputAllUserOnFile = false

	v1               = `/api/front`
	ModelsRoute      = v1 + `/models`
	ModelsCountRoute = ModelsRoute + `/count`
	TagsRoute        = ModelsRoute + `/liveTags`
	UsersRoute       = ModelsRoute + `/username`

	v2                 = `/api/front/v2`
	ConfigRouteV2      = v2 + `/config/data`
	ModelsRouteV2      = v2 + `/models`
	ModelsCountRouteV2 = ModelsRouteV2 + `/count`
	TagsRouteV2        = ModelsRouteV2 + `/liveTags`
	UsersRouteV2       = ModelsRouteV2 + `/username`
	Output             = `output`
	OutputM3u8         = false
	OutputAllAnchors   = false
	PrimaryTag         = `girls`
	Limit              = 120
)

func ReturnU(res string, tags *[]string) string {
	res += " "
	for _, tag := range *tags {
		res += tag + ","
	}
	res = strings.TrimSuffix(res, ",")
	return res
}
