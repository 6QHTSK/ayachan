package pkg

import "math"

type HandType int

const (
	UnknownHand HandType = iota
	LeftHand
	RightHand
	TryLeftHand
	TryRightHand
)

type ParsedNote struct {
	Note
	Hand         HandType
	NotePrevious *ParsedNote
	NoteAfter    *ParsedNote
}

type ParsedChart []ParsedNote

func (chart ParsedChart) Len() int {
	return len(chart)
}

func (chart ParsedChart) Less(i, j int) bool {
	if chart[i].Beat < chart[j].Beat {
		if chart[i].Lane < chart[j].Lane {
			return true
		}
	}
	return false
}

func (chart ParsedChart) Swap(i, j int) {
	chart[i], chart[j] = chart[j], chart[i]
}

func (chart Chart) InitParseChart() (ParsedChart ParsedChart) {
	for _, note := range chart {
		if note.Type == NoteTypeBpm || note.Type == NoteTypeSlide && note.Hidden {
			continue
		}
		ParsedChart = append(ParsedChart, ParsedNote{
			Note:         note,
			Hand:         UnknownHand,
			NotePrevious: nil,
			NoteAfter:    nil,
		})
	}
	return ParsedChart
}

func (note ParsedNote) GetIntervalFront() (interval float64) {
	if note.NotePrevious == nil {
		return math.Inf(1)
	}
	return math.Max(0.025, note.Time-note.NotePrevious.Time)
}

func (note ParsedNote) GetIntervalBack() (interval float64) {
	if note.NoteAfter == nil {
		return math.Inf(1)
	}
	return math.Max(0.025, note.NoteAfter.Time-note.Time)
}

func (note ParsedNote) GetGapFront() (gap float64) {
	if note.NotePrevious == nil {
		return 0
	}
	return note.Lane - note.NotePrevious.Lane
}
