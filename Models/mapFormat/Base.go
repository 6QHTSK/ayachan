package mapFormat

type NoteType int

const (
	NoteTypeBpm NoteType = iota
	NoteTypeSingle
	NoteTypeSlide
)

type SlideStatus int

const (
	SlideTick SlideStatus = iota
	SlideStart
	SlideEnd
)

type Note struct {
	Type      NoteType    // 音符类型
	BPM       float64     // BPM信息
	Beat      float64     // 节拍数
	Time      float64     // 判定时间
	Lane      float64     // 轨道号
	Direction int         // 左右滑建，普通滑键的大小：-3~-1 左滑键；0 普通滑建；1~3右滑建
	Pos       int         // type为NoteType_Slide时，所属的绿条编号
	Status    SlideStatus // 是否为最后一个SlideTick
	Flick     bool        // 是否为粉键/滑键
	Hidden    bool        // 是否为隐藏音符
}

type Chart []Note

func (chart Chart) Len() int {
	return len(chart)
}

func (chart Chart) Less(i, j int) bool {
	if chart[i].Beat < chart[j].Beat {
		if chart[i].Lane < chart[j].Lane {
			return true
		}
	}
	return false
}

func (chart Chart) Swap(i, j int) {
	chart[i], chart[j] = chart[j], chart[i]
}
