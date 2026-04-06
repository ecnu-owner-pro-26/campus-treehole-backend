package model

const (
	TagHappy    = 1 << iota // 1
	TagSad                  // 2
	TagDaydream             // 4
	TagAnxious              // 8
	TagTired                // 16
	TagCalm
	TagCrack
	TagStudy
	TagSlack
	TagWork
	TagSport
	TagEat // 32
)

// 标签名称到位掩码的映射
var TagNameToBit = map[string]int64{
	"开心": TagHappy,
	"难过": TagSad,
	"发呆": TagDaydream,
	"焦虑": TagAnxious,
	"疲惫": TagTired,
	"平静": TagCalm,
	"裂开": TagCrack,
	"学习": TagStudy,
	"摸鱼": TagSlack,
	"搬砖": TagWork,
	"运动": TagSport,
	"干饭": TagEat,
}

// 位掩码到标签名称的映射（用于响应）
var TagBitToName = map[int64]string{
	TagHappy:    "开心",
	TagSad:      "难过",
	TagDaydream: "发呆",
	TagAnxious:  "焦虑",
	TagTired:    "疲惫",
	TagCalm:     "平静",
	TagCrack:    "裂开",
	TagStudy:    "学习",
	TagSlack:    "摸鱼",
	TagWork:     "搬砖",
	TagSport:    "运动",
	TagEat:      "干饭",
}
