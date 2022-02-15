package ChartFormat

import (
	"github.com/6QHTSK/ayachan/Models"
)

type AuthorID int

const (
	Author_Aya AuthorID = iota
	Author_Aya6QHTSK
	Author_6QHTSK
	Author_Other
)

type Chart struct {
	Submitted bool   // 是否已经提交到BanGround社区
	ChartID   int    // 已提交：BanGround社区谱面ID,未提交：Bestdori社区Ex/SpID
	Title     string // 谱面的标题
	Artists   string // 谱面的艺术家
	Author    struct {
		AuthorID   AuthorID // 谱面作者用户名
		AuthorName string   // Author_Other 其他作者时使用的名义
	}
	Diffs []struct {
		BestdoriChartID        int // Bestdori的谱面ID 0 就是未发表
		Diff                   DiffType
		Level                  int
		Models.MapMetricsBasic // 一些易得的谱面数据
	}
	SongUrl struct {
		Cover string
		Audio string
	} // 谱面资源
	Rel string // 访问网址
}
