package chartFormat

import (
	"ayachanV2/Models"
	"ayachanV2/Models/mapFormat"
)

type DiffType int

const (
	Diff_Easy DiffType = iota
	Diff_Normal
	Diff_Hard
	Diff_Expert
	Diff_SpecialC
)

type Author struct {
	Username string `json:"username"` // 谱面作者用户名
	Nickname string `json:"nickname"` // 谱面作者昵称
}

type BestdoriChartItem struct {
	// id = ChartID * 10 + Official ? 9 : Diff json不
	ChartID int      `json:"id"`      // Bestdori的谱面ID
	Title   string   `json:"title"`   // 谱面的标题
	Artists string   `json:"artists"` // 谱面的艺术家
	Author  Author   `json:"author"`
	Diff    DiffType `json:"diff"`
	Level   int      `json:"level"`
	SongUrl struct {
		Cover string `json:"cover"`
		Audio string `json:"audio"`
	} `json:"song_url"` // 谱面资源
	Official            bool                      `json:"official"`        //是否为官谱
	Likes               int                       `json:"likes,omitempty"` // 喜爱数
	PostTime            uint64                    `json:"time,omitempty"`  // 时间戳
	LastUpdateTime      uint64                    `json:"last_update"`
	Chart               mapFormat.BestdoriV2Chart `json:"chart,omitempty"` // 谱面
	Models.MapInfoBasic                           // 一些易得的谱面数据 用于某些项
}

type CharterRankItem struct {
	Rank   int // 排名
	Author Author
	Count  int // 计数
}

type SongRankItem struct {
	Rank    int      // 排名
	ChartID int      `json:"id"`      // Bestdori的谱面ID
	Title   string   `json:"title"`   // 谱面的标题
	Artists string   `json:"artists"` // 谱面的艺术家
	Author  Author   `json:"author"`
	Diff    DiffType `json:"diff"`
	Level   int      `json:"level"`
	Likes   int      `json:"likes,omitempty"` // 喜爱数
}

type CharterSelfInfoBasic struct {
	Charter         Author
	TotalPost       int
	TotalLike       int
	TotalNote       int
	TotalTime       float64
	AverageDiff     float64
	AverageDiffRank float64
}
