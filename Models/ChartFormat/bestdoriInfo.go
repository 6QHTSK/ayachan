package ChartFormat

import (
	"ayachan/Models"
	"ayachan/Models/MapFormat"
	"time"
)

type DiffType int

const (
	Diff_Easy DiffType = iota
	Diff_Normal
	Diff_Hard
	Diff_Expert
	Diff_Special
)

type Author struct {
	Username string `json:"username"` // 谱面作者用户名
	Nickname string `json:"nickname"` // 谱面作者昵称
}

type BestdoriChartItem struct {
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
	Likes               int                       `json:"likes,omitempty"` // 喜爱数
	PostTime            uint64                    `json:"time,omitempty"`  // 时间戳
	LastUpdateTime      time.Time                 `json:"last_update"`
	Content             string                    `json:"content"`
	Chart               MapFormat.BestdoriV2Chart `json:"chart,omitempty"` // 谱面
	Models.MapInfoBasic                           // 一些易得的谱面数据 用于某些项
}

type BestdoriChartUpdateItem struct {
	ChartID  int    `json:"chart_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Diff     int    `json:"diff"`
	Level    int    `json:"chart_level"`
	Likes    int    `json:"likes"`
}

type CharterRankItem struct {
	Author Author
	Count  int // 计数
}

type SongRankItem struct {
	ChartID int      `json:"id"`      // Bestdori的谱面ID
	Title   string   `json:"title"`   // 谱面的标题
	Artists string   `json:"artists"` // 谱面的艺术家
	Author  Author   `json:"author"`
	Diff    DiffType `json:"diff"`
	Level   int      `json:"level"`
	Likes   int      `json:"likes,omitempty"` // 喜爱数
}

type CharterSelfInfoBasic struct {
	Charter     Author
	TotalPost   int
	TotalLike   int
	TotalNote   int
	TotalTime   float64
	AverageDiff float64
}
