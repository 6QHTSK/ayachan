package pkg

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
	Beat_  *float64 `json:"beat,omitempty"`
	Lane_  *float64 `json:"lane,omitempty"`
	Flick  bool     `json:"flick,omitempty"`
	Hidden bool     `json:"hidden,omitempty"`
}

func (note BestdoriV2Note) Beat() (value float64) {
	if len(note.Connections) == 0 {
		if note.Beat_ == nil {
			return value
		}
		value = *note.Beat_
	} else {
		if note.Connections[0].Beat_ == nil {
			return value
		}
		value = *note.Connections[0].Beat_
	}
	return value
}

func (note BestdoriV2Note) Lane() (value float64) {
	if len(note.Connections) == 0 {
		if note.Lane_ == nil {
			return value
		}
		value = *note.Lane_
	} else {
		if note.Connections[0].Lane_ == nil {
			return value
		}
		value = *note.Connections[0].Lane_
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
				return false, fmt.Errorf("????????????????????????????????????")
			}
			if formatNote.Width < 0 || formatNote.Width > 3 {
				return false, fmt.Errorf("??????????????????")
			}
			fallthrough
		case "Single":
			if formatNote.Lane_ == nil {
				return false, fmt.Errorf("??????????????????Lane??????")
			}
			// ???Beat,?????????Lane??????Decode??????
			if len(formatNote.Connections) != 0 {
				return false, fmt.Errorf("?????????????????????Connections??????")
			}
			fallthrough
		case "BPM":
			if formatNote.Beat_ == nil {
				return false, fmt.Errorf("?????????/BPM?????????Beat??????")
			}
			if len(formatNote.Connections) != 0 {
				return false, fmt.Errorf("BPM???????????????Connections??????")
			}
			// BPM??????????????????0.0Beat???BPM????????????Decode????????????
		case "Long":
			fallthrough
		case "Slide":
			// ???????????????????????????????????????Decode????????????
			for _, formatTick := range formatNote.Connections {
				if formatTick.Beat_ == nil || formatTick.Lane_ == nil {
					return false, fmt.Errorf("???Slide/Long?????????????????????Beat/Lane??????")
				}
			}
		default:
			// ????????????????????????Decode????????????
		}
	}
	return true, nil
}

func (formatChart BestdoriV2Chart) Decode() (Chart Chart) {
	SlideCounter := 0
	sort.Sort(formatChart)
	FirstBPMBeat := math.Inf(1)

	// ??????BPM???????????????0
	for _, formatNote := range formatChart {
		if formatNote.Type == "BPM" {
			FirstBPMBeat = formatNote.Beat()
			break
		}
	}

	// ??????????????????????????????????????????????????????
	for _, formatNote := range formatChart {
		var note Note

		if formatNote.Beat() < FirstBPMBeat {
			continue //????????????????????????BPM???????????????????????????
		}

		if formatNote.Type == "Single" || formatNote.Type == "Directional" {
			// ?????????????????????????????????
			// ??????????????????
			note = Note{
				Type:  NoteTypeSingle,
				Beat:  formatNote.Beat() - FirstBPMBeat,
				Lane:  fixLane(formatNote.Lane(), false),
				Flick: formatNote.Flick,
			}
			// ??????????????????
			if formatNote.Direction == "Left" {
				note.Flick = true
				note.Direction = -formatNote.Width
			} else if formatNote.Direction == "Right" {
				note.Flick = true
				note.Direction = formatNote.Width
			}
			Chart = append(Chart, note)
		} else if formatNote.Type == "BPM" {
			// ?????????????????????BPM??????
			// ??????????????????
			note = Note{
				Type: NoteTypeBpm,
				BPM:  math.Abs(formatNote.BPM),
				Beat: formatNote.Beat() - FirstBPMBeat,
			}
			Chart = append(Chart, note)
		} else if formatNote.Type == "Slide" || formatNote.Type == "Long" {
			// ???????????????????????????
			// ??????connection??????????????????
			connectionsCount := len(formatNote.Connections)
			if connectionsCount == 0 {
				// ?????????0 ?????? ??????
				continue
			} else if connectionsCount == 1 {
				// ?????????1 ???????????????
				note = Note{
					Type:  NoteTypeSingle,
					Beat:  formatNote.Beat() - FirstBPMBeat,
					Lane:  fixLane(formatNote.Lane(), false),
					Flick: formatNote.Connections[0].Flick,
				}
				Chart = append(Chart, note)
			} else {
				// ????????????
				SlideCounter++
				//???????????????
				note = Note{
					Type:   NoteTypeSlide,
					Beat:   formatNote.Beat() - FirstBPMBeat,
					Lane:   fixLane(formatNote.Lane(), false),
					Pos:    SlideCounter,
					Status: SlideStart,
				}
				Chart = append(Chart, note)
				// ??????????????????????????????
				for i := 1; i < connectionsCount; i++ {
					note = Note{
						Type:   NoteTypeSlide,
						Beat:   *formatNote.Connections[i].Beat_ - FirstBPMBeat,
						Lane:   fixLane(*formatNote.Connections[i].Lane_, formatNote.Connections[i].Hidden),
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
	// ?????????????????????
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

func (formatChart *BestdoriV2Chart) Encode(chart Chart) error {
	*formatChart = []BestdoriV2Note{}
	SlideSuccessFlag := false
	for i, note := range chart {
		if note.Type == NoteTypeSlide && note.Status == SlideStart {
			var basicNoteList []BestdoriV2BasicNote
			for j := i; j < len(chart); j++ {
				if note.Type == NoteTypeSlide && note.Pos == chart[j].Pos {
					var tick = chart[j]
					basicNote := BestdoriV2BasicNote{
						Beat_:  &tick.Beat,
						Lane_:  &tick.Lane,
						Flick:  tick.Flick,
						Hidden: tick.Hidden,
					}
					basicNoteList = append(basicNoteList, basicNote)
					if tick.Status == SlideEnd {
						*formatChart = append(*formatChart, BestdoriV2Note{
							Type:        "Slide",
							Connections: basicNoteList,
						})
						basicNoteList = []BestdoriV2BasicNote{}
						SlideSuccessFlag = true
					}
				}
			}
			// ?????????????????????
			if !SlideSuccessFlag {
				return fmt.Errorf("Encode:??????????????????")
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
				Beat_: &note.Beat,
				Lane_: &note.Lane,
				Flick: note.Flick && note.Direction == 0,
			}
			*formatChart = append(*formatChart, formatNote)
		}
	}
	return nil
}
