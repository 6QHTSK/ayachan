package generalModels

import (
	"github.com/6QHTSK/ayachan/pkg"
	"time"
)

type Author struct {
	Username string `json:"username"` // 谱面作者用户名
	Nickname string `json:"nickname"` // 谱面作者昵称
}

type BestdoriChartItem struct {
	ChartID int    `json:"id"`      // Bestdori的谱面ID
	Title   string `json:"title"`   // 谱面的标题
	Artists string `json:"artists"` // 谱面的艺术家
	Author  Author `json:"author"`
	Diff    int    `json:"diff"` // 0-4 Easy-Special
	Level   int    `json:"level"`
	SongUrl struct {
		Cover string `json:"cover"`
		Audio string `json:"audio"`
	} `json:"song_url"` // 谱面资源
	Likes               int                 `json:"likes,omitempty"` // 喜爱数
	PostTime            uint64              `json:"time,omitempty"`  // 时间戳
	LastUpdateTime      time.Time           `json:"last_update"`
	Content             string              `json:"content"`
	Chart               pkg.BestdoriV2Chart `json:"chart,omitempty"` // 谱面
	pkg.MapMetricsBasic                     // 一些易得的谱面数据 用于某些项
}

type BestdoriChartUpdateItem struct {
	ChartID  int    `json:"chart_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Diff     int    `json:"diff"`
	Level    int    `json:"chart_level"`
	Likes    int    `json:"likes"`
}
