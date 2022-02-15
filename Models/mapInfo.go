package Models

type RegularType int

const (
	RegularTypeUnknown RegularType = iota
	RegularTypeRegular
	RegularTypeIrregular
)

type DifficultyDescription int

const (
	DifficultyLow    = -1 // 该项难度偏低
	DifficultyNormal = 0  // 该项难度正常
	DifficultyHigh   = 1  // 该项难度偏高
)

type BpmInfo struct {
	BPMLow  float64 `json:"bpm-low"`
	BPMHigh float64 `json:"bpm-high"`
	MainBPM float64 `json:"main-bpm"`
}

type IrregularInfo struct {
	Irregular     RegularType `json:"irregular"`      // 存在多压/交叉（出张）0 失败 1 标准 2 非标准
	IrregularInfo string      `json:"irregular-info"` // 无法分析的第一个错误情况
}

type NoteCount struct {
	Single         int `json:"single"`
	Flick          int `json:"flick"`
	SlideStart     int `json:"slideStart"`
	SlideTick      int `json:"slideTick"`
	SlideEnd       int `json:"slideEnd"`
	SlideFlick     int `json:"slideFlick"`
	SlideHidden    int `json:"slideHidden"`
	DirectionLeft  int `json:"direction-left"`
	DirectionRight int `json:"direction-right"`
}

type Distribution struct {
	Note []int `json:"note"`
	Hit  []int `json:"hit"`
}

// MapMetricsBasic 将会放入数据库存档的数据部分
type MapMetricsBasic struct {
	IrregularInfo
	TotalNote int     `json:"total_note"`
	TotalTime float64 `json:"total_time"`
	TotalNPS  float64 `json:"total_nps"`
	SPRhythm  bool    `json:"sp_rhythm"`
}

// MapMetricsStandard 基础部分，不要求正常谱面
type MapMetricsStandard struct {
	MapMetricsBasic

	BpmInfo
	TotalHitNote int     `json:"total-hit-note"`
	MaxScreenNPS float64 `json:"max-screen-nps"`
	TotalHPS     float64 `json:"total-hps"`

	NoteCount    NoteCount
	Distribution Distribution
}

// MapMetricsExtend 扩展部分，要求正常谱面，非正常时为nil
type MapMetricsExtend struct {
	LeftPercent       float64 `json:"left-percent"`
	MaxSpeed          float64 `json:"max-speed"`
	FingerMaxHPS      float64 `json:"finger-max-hps"`
	FlickNoteInterval float64 `json:"flick-note-interval"`
	NoteFlickInterval float64 `json:"note-flick-interval"`
}

// MapDifficultyStandard 基础部分，不要求正常谱面
type MapDifficultyStandard struct {
	TotalNPS            float64 `json:"total-nps"`
	TotalHPS            float64 `json:"total-hps"`
	MaxScreenNPS        float64 `json:"max-screen-nps"`
	Difficulty          float64 `json:"difficulty"`
	BlueWhiteDifficulty float64 `json:"blue-white-difficulty"`
}

// MapDifficultyExtend 扩展部分，要求正常谱面，非正常时为nil
type MapDifficultyExtend struct {
	MaxSpeed          DifficultyDescription `json:"max-speed"`
	FingerMaxHPS      DifficultyDescription `json:"finger-max-hps"`
	FlickNoteInterval DifficultyDescription `json:"flick-note-interval"`
	NoteFlickInterval DifficultyDescription `json:"note-flick-interval"`
}

type MapInfo struct {
	MapMetrics          *MapMetricsStandard    `json:"map-metrics"`
	MapMetricsExtend    *MapMetricsExtend      `json:"map-metrics-extend,omitempty"`
	MapDifficulty       *MapDifficultyStandard `json:"map-difficulty"`
	MapDifficultyExtend *MapDifficultyExtend   `json:"map-difficulty-extend,omitempty"`
}
