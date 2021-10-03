package chartFormat

import (
	"ayachanV2/Models"
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
	Username string // 谱面作者用户名
	NickName string // 谱面作者昵称
}

type BestdoriChartItem struct {
	chartID int    // Bestdori的谱面ID
	Title   string // 谱面的标题
	Artists string // 谱面的艺术家
	Author  Author
	Diff    DiffType
	Level   int
	SongUrl struct {
		Cover string
		Audio string
	} // 谱面资源
	Likes               int    // 喜爱数
	time                uint32 // 时间戳
	Models.MapInfoBasic        // 一些易得的谱面数据 用于某些项
}

type CharterRankItem struct {
	Rank   int // 排名
	Author Author
	Count  int // 计数
}

type SongRankItem struct {
	Rank int // 排名
	BestdoriChartItem
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
