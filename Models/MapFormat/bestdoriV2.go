package MapFormat

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
	Beat_  interface{} `json:"beat,omitempty"` //Omitempty need to be interface,Beat_ may be 0.0
	Lane_  interface{} `json:"lane,omitempty"`
	Flick  bool        `json:"flick,omitempty"`
	Hidden bool        `json:"hidden,omitempty"`
}

func (note BestdoriV2Note) Beat() (value float64) {
	var ok bool
	if len(note.Connections) == 0 {
		value, ok = note.Beat_.(float64)
	} else {
		value, ok = note.Connections[0].Beat_.(float64)
	}
	if !ok {
		return 0
	}
	return value
}

func (note BestdoriV2Note) Lane() (value float64) {
	var ok bool
	if len(note.Connections) == 0 {
		value, ok = note.Lane_.(float64)
	} else {
		value, ok = note.Connections[0].Lane_.(float64)
	}
	if !ok {
		return 0
	}
	return value
}

type BestdoriV2Chart []BestdoriV2Note

func (formatChart BestdoriV2Chart) Len() int {
	return len(formatChart)
}

func (formatChart BestdoriV2Chart) Less(i, j int) bool {
	if formatChart[i].Beat() == formatChart[j].Beat() {
		return formatChart[i].Lane() < formatChart[j].Lane()
	}
	return formatChart[i].Beat() < formatChart[j].Beat()
}

func (formatChart BestdoriV2Chart) Swap(i, j int) {
	formatChart[i], formatChart[j] = formatChart[j], formatChart[i]
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

func (formatChart BestdoriV2Chart) MapCheck() (result bool, err error) {
	for _, formatNote := range formatChart {
		switch formatNote.Type {
		case "Directional":
			if formatNote.Direction != "Left" && formatNote.Direction != "Right" {
				return false, fmt.Errorf("无法识别侧划音符的标识符")
			}
			if formatNote.Width < 0 || formatNote.Width > 3 {
				return false, fmt.Errorf("侧划音符超限")
			}
			fallthrough
		case "Single":
			currentBeat, ok := formatNote.Beat_.(float64)
			if !ok || currentBeat < 0.0 {
				return false, fmt.Errorf("无法解析某个音符的节拍/节拍数小于0")
			}
			_, ok = formatNote.Lane_.(float64) // 音符轨道可以被修正，在后续Decode部分修正
			if !ok {
				return false, fmt.Errorf("无法解析某个音符的轨道")
			}
			if len(formatNote.Connections) != 0 {
				return false, fmt.Errorf("单键错误的拥有Connections字段")
			}
		case "Long":
			fallthrough
		case "Slide":
			// 绿条长度可以被修正，在后续Decode部分修正
			for _, formatTick := range formatNote.Connections {
				currentBeat, ok := formatTick.Beat_.(float64)
				if !ok || currentBeat < 0.0 {
					return false, fmt.Errorf("无法解析某个绿条节点的节拍/节拍数小于0")
				}
				_, ok = formatTick.Lane_.(float64)
				if !ok {
					return false, fmt.Errorf("无法解析某个绿条节点的轨道")
				}
			}
		case "BPM":
			_, ok := formatNote.Beat_.(float64)
			if !ok {
				return false, fmt.Errorf("无法解析BPM节点的节拍")
			}
			if len(formatNote.Connections) != 0 {
				return false, fmt.Errorf("BPM错误的拥有Connections字段")
			}
			// BPM的正负，不在0.0Beat的BPM音符会在Decode部分修正
		default:
			// 不知道的音符会在Decode部分扔掉
		}
	}
	return true, nil
}

func (formatChart BestdoriV2Chart) Decode() (Chart Chart) {
	SlideCounter := 0
	sort.Sort(formatChart)
	FirstBPMBeat := math.Inf(1)

	// 首个BPM节拍校正至0
	for _, formatNote := range formatChart {
		if formatNote.Type == "BPM" {
			FirstBPMBeat = formatNote.Beat()
			break
		}
	}

	// 首先，我们先排序，然后将基本信息填上
	for _, formatNote := range formatChart {
		var note Note

		if formatNote.Beat() < FirstBPMBeat {
			continue //忽略所有在第一个BPM音符出现之前的音符
		}

		if formatNote.Type == "Single" || formatNote.Type == "Directional" {
			// 检测到该音符是单点音符
			// 注入基本信息
			note = Note{
				Type:  NoteTypeSingle,
				Beat:  formatNote.Beat() - FirstBPMBeat,
				Lane:  fixLane(formatNote.Lane(), false),
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
				Beat: formatNote.Beat() - FirstBPMBeat,
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
				note = Note{
					Type:  NoteTypeSingle,
					Beat:  formatNote.Beat() - FirstBPMBeat,
					Lane:  fixLane(formatNote.Lane(), false),
					Flick: formatNote.Connections[0].Flick,
				}
				Chart = append(Chart, note)
			} else {
				// 长度正常
				SlideCounter++
				//注入绿条首
				note = Note{
					Type:   NoteTypeSlide,
					Beat:   formatNote.Beat() - FirstBPMBeat,
					Lane:   fixLane(formatNote.Lane(), false),
					Pos:    SlideCounter,
					Status: SlideStart,
				}
				Chart = append(Chart, note)
				// 注入绿条中间键、尾键
				for i := 1; i < connectionsCount; i++ {
					note = Note{
						Type:   NoteTypeSlide,
						Beat:   formatNote.Connections[i].Beat_.(float64) - FirstBPMBeat,
						Lane:   fixLane(formatNote.Connections[i].Lane_.(float64), formatNote.Connections[i].Hidden),
						Pos:    SlideCounter,
						Hidden: formatNote.Connections[i].Hidden,
						Status: SlideEnd,
						Flick:  formatNote.Connections[i].Flick,
					}
					if i != 1 {
						Chart[len(Chart)-1].Status = SlideTick
						Chart[len(Chart)-1].Flick = false
					}
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
	// 空谱面特殊处理
	if Chart.Len() == 0 {
		Chart = append(Chart, Note{
			Type:      NoteTypeBpm,
			BPM:       120,
			Beat:      0,
			Time:      0,
			Lane:      0,
			Direction: 0,
			Pos:       0,
			Status:    0,
			Flick:     false,
			Hidden:    false,
		})
	}
	return Chart
}

func typeConvert(note Note) (typeString string) {
	if note.Type == NoteTypeBpm {
		return "BPM"
	} else if note.Type == NoteTypeSingle {
		if note.Direction != 0 {
			return "Directional"
		}
		return "Single"
	} else if note.Type == NoteTypeSlide {
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
				Type: typeConvert(note),
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
