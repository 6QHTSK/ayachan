package DatabaseModel

import (
	"github.com/6QHTSK/ayachan/Models"
	"github.com/6QHTSK/ayachan/Models/ChartFormat"
	"time"
)

// BestdoriFanMadeView also a map for meiliSearch item
type BestdoriFanMadeView struct {
	ChartID    int       `db:"chartID" json:"chart_id"`
	Title      string    `db:"title" json:"title"`
	Artists    string    `db:"artists" json:"artists"`
	Username   string    `db:"author" json:"username"`
	Nickname   string    `db:"nickname" json:"nickname"`
	Diff       int       `db:"diff" json:"diff"`
	ChartLevel int       `db:"chartLevel" json:"chart_level"`
	CoverURL   string    `db:"coverURL" json:"cover_url"`
	SongURL    string    `db:"songURL" json:"song_url"`
	Likes      int       `db:"likes" json:"likes"`
	PostTime   uint64    `db:"postTime" json:"post_time"`
	LastUpdate time.Time `db:"lastUpdate" json:"last_update"`
	TotalNote  int       `db:"totalNote" json:"total_note"`
	TotalTime  float64   `db:"totalTime" json:"total_time"`
	TotalNPS   float64   `db:"totalNPS" json:"total_nps"`
	SPRhythm   bool      `db:"SPRhythm" json:"sp_rhythm"`
	Irregular  int       `db:"irregular" json:"irregular"`
	Content    string    `db:"content" json:"content"`
	Version    int       `db:"version" json:"-"`
}

type BestdoriAuthorList struct {
	UserName string `db:"username"`
	NickName string `db:"nickname"`
}

type BestdoriFanMade struct {
	ChartID    int       `db:"chartID"`
	Title      string    `db:"title"`
	Artists    string    `db:"artists"`
	Author     string    `db:"author"`
	Diff       int       `db:"diff"`
	ChartLevel int       `db:"chartLevel"`
	CoverURL   string    `db:"coverURL"`
	SongURL    string    `db:"songURL"`
	Likes      string    `db:"likes"`
	PostTime   uint64    `db:"postTime"`
	LastUpdate time.Time `db:"lastUpdate"`
	Content    string    `db:"content"`
}

type BestdoriFanMadeMetrics struct {
	ChartID   int     `db:"chartID"`
	TotalNote int     `db:"totalNote"`
	TotalTime float64 `db:"totalTime"`
	//TotalNPS use TotalNote / TotalTime
	SPRhythm  bool `db:"SPRhythm"`
	Irregular bool `db:"irregular"`
	Version   int  `db:"version"`
}

type BestdoriOverAllInfo struct {
	ChartCount  int     `db:"chartCount" json:"chart-count"`
	Latest      int     `db:"latest" json:"latest"`
	TotalNote   int     `db:"totalNote" json:"total-note"`
	TotalTime   float64 `db:"totalTime" json:"total-time"`
	AuthorCount int     `db:"authorCount" json:"author-count"`
}

func (d BestdoriFanMadeView) ToBestdoriChart() ChartFormat.BestdoriChartItem {
	return ChartFormat.BestdoriChartItem{
		ChartID: d.ChartID,
		Title:   d.Title,
		Artists: d.Artists,
		Author: ChartFormat.Author{
			Username: d.Username,
			Nickname: d.Nickname,
		},
		Diff:  ChartFormat.DiffType(d.Diff),
		Level: d.ChartLevel,
		SongUrl: struct {
			Cover string `json:"cover"`
			Audio string `json:"audio"`
		}{
			Cover: d.CoverURL,
			Audio: d.SongURL,
		},
		Likes:          d.Likes,
		PostTime:       d.PostTime,
		LastUpdateTime: d.LastUpdate,
		Content:        d.Content,
		MapMetricsBasic: Models.MapMetricsBasic{
			IrregularInfo: Models.IrregularInfo{
				Irregular: Models.RegularType(d.Irregular),
			},
			TotalNote: d.TotalNote,
			TotalTime: d.TotalTime,
			TotalNPS:  d.TotalNPS,
			SPRhythm:  d.SPRhythm,
		},
	}
}
