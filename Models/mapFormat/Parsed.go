package mapFormat

type HandType int

const (
	UnknownHand HandType = iota
	LeftHand
	RightHand
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
		ParsedChart = append(ParsedChart, ParsedNote{
			Note:         note,
			Hand:         UnknownHand,
			NotePrevious: nil,
			NoteAfter:    nil,
		})
	}
	return ParsedChart
}
