package mapFormat

import (
	"fmt"
	"math"
	"sort"
)

type BestdoriV2Note struct {
	BestdoriV2BasicNote
	Type        string                `json:"type"`
	BPM         float64               `json:"bpm,omitempty"`
	Connections []BestdoriV2BasicNote `json:"connections,omitempty"`
	Direction   string                `json:"direction,omitempty"`
	Width       int                   `json:"width,omitempty"`
}

type BestdoriV2BasicNote struct {
	Beat_  float64 `json:"beat,omitempty"`
	Lane_  float64 `json:"lane,omitempty"`
	Flick  bool    `json:"flick,omitempty"`
	Hidden bool    `json:"hidden,omitempty"`
}

func (note BestdoriV2Note) Beat() float64 {
	if len(note.Connections) == 0 {
		return note.Beat_
	} else {
		return note.Connections[0].Beat_
	}
}

func (note BestdoriV2Note) Lane() float64 {
	if len(note.Connections) == 0 {
		return note.Lane_
	} else {
		return note.Connections[0].Lane_
	}
}

type BestdoriV2Chart []BestdoriV2Note

func (formatChart BestdoriV2Chart) Len() int {
	return len(formatChart)
}

func (formatChart BestdoriV2Chart) Less(i, j int) bool {
	if formatChart[i].Beat() < formatChart[j].Beat() {
		if formatChart[i].Lane() < formatChart[j].Lane() {
			return true
		}
	}
	return false
}

func (formatChart BestdoriV2Chart) Swap(i, j int) {
	formatChart[i], formatChart[j] = formatChart[j], formatChart[i]
}

type ChartBestdoriV2 struct {
	chart BestdoriV2Chart
}

func fixLane(lane float64, noteHidden bool) (fix float64) {
	if !noteHidden {
		if lane < 0.0 {
			return 0.0
		} else if lane > 7.0 {
			return 7.0
		} else {
			return lane
		}
	} else {
		return lane
	}
}

func (formatChart BestdoriV2Chart) Decode() (Chart Chart) {
	SlideCounter := 0
	sort.Sort(formatChart)
	// 首先，我们先排序，然后将基本信息填上
	for _, formatNote := range formatChart {
		var note Note
		if formatNote.Beat() < 0 {
			// 过滤掉节拍异常的音符
			continue
		}
		if formatNote.Type == "Single" {
			// 检测到该音符是单点音符
			// 注入基本信息
			note = Note{
				Type:  NoteTypeSingle,
				Beat:  formatNote.Beat(),
				Lane:  fixLane(formatNote.Lane_, false),
				Flick: formatNote.Flick,
			}
			// 注入侧滑信息
			if formatNote.Direction == "Left" {
				note.Flick = true
				note.Direction = -formatNote.Width
			} else if formatNote.Direction == "Right" {
				note.Flick = true
				note.Direction = formatNote.Width
			}
			Chart = append(Chart, note)
		} else if formatNote.Type == "BPM" {
			// 检测到该音符是BPM音符
			// 注入基本信息
			note = Note{
				Type: NoteTypeBpm,
				BPM:  math.Abs(formatNote.BPM),
				Beat: fixLane(formatNote.Beat_, false),
			}
			Chart = append(Chart, note)
		} else if formatNote.Type == "Slide" || formatNote.Type == "Long" {
			// 检测到该音符为绿条
			// 检测connection字段中的信息
			connectionsCount := len(formatNote.Connections)
			if connectionsCount == 0 {
				// 长度为0 非法 跳过
				continue
			} else if connectionsCount == 1 {
				// 长度为1 退化为单点
				if formatNote.Connections[0].Beat_ < 0 {
					continue
				}
				// 注入基本信息
				note = Note{
					Type:  NoteTypeSingle,
					Beat:  formatNote.Connections[0].Beat_,
					Lane:  fixLane(formatNote.Connections[0].Lane_, false),
					Flick: formatNote.Connections[0].Flick,
				}
				Chart = append(Chart, note)
			} else {
				// 长度正常
				SlideCounter++
				//注入绿条首
				note = Note{
					Type:   NoteTypeSlide,
					Beat:   formatNote.Connections[0].Beat_,
					Lane:   fixLane(formatNote.Connections[0].Lane_, false),
					Pos:    0,
					Status: SlideStart,
				}
				Chart = append(Chart, note)
				// 注入绿条中间键、尾键
				for i := 1; i < connectionsCount-1; i++ {
					if formatNote.Connections[i].Beat_ < 0 {
						// 过滤掉节拍异常的音符
						continue
					}
					note = Note{
						Type:   NoteTypeSlide,
						Beat:   formatNote.Connections[i].Beat_,
						Lane:   fixLane(formatNote.Connections[i].Lane_, formatNote.Connections[i].Hidden),
						Pos:    SlideCounter,
						Hidden: formatNote.Connections[i].Hidden,
						Status: SlideEnd,
						Flick:  formatNote.Connections[i].Flick,
					}
					Chart[len(Chart)-1].Status = SlideTick
					Chart[len(Chart)-1].Flick = false
					Chart = append(Chart, note)
				}
			}
		}
	}

	currentBPM := 120.0
	offsetBeat := 0.0
	offsetTime := 0.0
	sort.Sort(Chart)
	for i := range Chart {
		Chart[i].Time = (Chart[i].Beat-offsetBeat)*(60.0/currentBPM) + offsetTime
		if Chart[i].Type == NoteTypeBpm {
			offsetTime = Chart[i].Time
			offsetBeat = Chart[i].Beat
			currentBPM = Chart[i].BPM
		}
	}
	return Chart
}

func typeConvert(typeNum NoteType) (typeString string) {
	if typeNum == NoteTypeBpm {
		return "BPM"
	} else if typeNum == NoteTypeSingle {
		return "Single"
	} else if typeNum == NoteTypeSlide {
		return "Slide"
	}
	return ""
}

func directionConvert(DirectionValue int) (DirectionString string, Width int) {
	if DirectionValue == 0 {
		return "", 0
	} else if DirectionValue < 0 {
		return "Left", -DirectionValue
	} else {
		return "Right", DirectionValue
	}
}

func (chart Chart) EncodeBestdoriV2() (formatChart BestdoriV2Chart, err error) {
	SlideSuccessFlag := false
	for i, note := range chart {
		if note.Type == NoteTypeSlide && note.Status == SlideStart {
			var basicNoteList []BestdoriV2BasicNote
			for j := i; j < len(chart); j++ {
				if note.Type == NoteTypeSlide && note.Pos == chart[j].Pos {
					var tick = chart[j]
					basicNote := BestdoriV2BasicNote{
						Beat_:  tick.Beat,
						Lane_:  tick.Lane,
						Flick:  tick.Flick,
						Hidden: tick.Hidden,
					}
					basicNoteList = append(basicNoteList, basicNote)
					if tick.Status == SlideEnd {
						formatChart = append(formatChart, BestdoriV2Note{
							Type:        "Slide",
							Connections: basicNoteList,
						})
						basicNoteList = []BestdoriV2BasicNote{}
						SlideSuccessFlag = true
					}
				}
			}
			// 未查找到绿条尾
			if !SlideSuccessFlag {
				return formatChart, fmt.Errorf("EncodeBestdoriV2:找不到绿条尾")
			} else {
				SlideSuccessFlag = false
			}
		} else {
			formatNote := BestdoriV2Note{
				Type: typeConvert(note.Type),
				BPM:  note.BPM,
			}
			formatNote.Direction, formatNote.Width = directionConvert(note.Direction)
			formatNote.BestdoriV2BasicNote = BestdoriV2BasicNote{
				Beat_: note.Beat,
				Lane_: note.Lane,
				Flick: note.Flick && note.Direction == 0,
			}
			formatChart = append(formatChart, formatNote)
		}
	}
	return formatChart, nil
}
