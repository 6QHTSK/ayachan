package bestdori

import (
	"fmt"
	"github.com/6QHTSK/ayachan/internal/pkg/generalModels"
	"github.com/6QHTSK/ayachan/internal/server/database"
	"math"
)

type pair struct {
	Low       interface{}
	High      interface{}
	Valid     bool
	HighValid bool
}

type elementBool struct {
	Value bool
	Valid bool
}

func (p *pair) set(low interface{}, high interface{}) {
	p.Low = low
	p.High = high
	p.Valid = true
	p.HighValid = true
}

func (p *pair) setLow(low interface{}) {
	p.Low = low
	p.Valid = true
	p.HighValid = false
}

type SearchBody struct {
	queryString string
	page        int64 //page start at 0
	limit       int64
	level       pair
	diff        pair
	SP          elementBool
	Regular     elementBool
	Time        pair
	NPS         pair
}

func NewSearch(queryString string, page int64, limit int64) SearchBody {
	return SearchBody{
		queryString: queryString,
		page:        page,
		limit:       limit,
	}
}

func (s *SearchBody) FilterLevel(levelLow int, levelHigh int) *SearchBody {
	s.level.set(levelLow, levelHigh)
	return s
}

func (s *SearchBody) FilterDiff(diffLow int, diffHigh int) *SearchBody {
	s.diff.set(diffLow, diffHigh)
	return s
}

func (s *SearchBody) FilterSP(isSp bool) *SearchBody {
	s.SP.Value = isSp
	s.SP.Valid = true
	return s
}

func (s *SearchBody) FilterIrregular(isRegular bool) *SearchBody {
	s.Regular.Value = isRegular
	s.Regular.Valid = true
	return s
}

func (s *SearchBody) FilterTime(timeLow float64, timeHigh float64) *SearchBody {
	s.Time.set(timeLow, timeHigh)
	return s
}

func (s *SearchBody) FilterNPS(NPSLow float64, NPSHigh float64) *SearchBody {
	s.NPS.set(NPSLow, NPSHigh)
	return s
}

func (s *SearchBody) FilterTimeLow(timeLow float64) *SearchBody {
	s.Time.setLow(timeLow)
	return s
}

func (s *SearchBody) FilterNPSLow(NPSLow float64) *SearchBody {
	s.NPS.setLow(NPSLow)
	return s
}

func (s *SearchBody) Filter() (filter []string) {
	if s.level.Valid {
		filter = append(filter, fmt.Sprintf("chart_level >= %d", s.level.Low), fmt.Sprintf("chart_level <= %d", s.level.High))
	}
	if s.diff.Valid {
		filter = append(filter, fmt.Sprintf("diff >= %d", s.diff.Low), fmt.Sprintf("diff <= %d", s.diff.High))
	}
	if s.Time.Valid {
		filter = append(filter, fmt.Sprintf("total_time >= %f", s.Time.Low))
		if s.Time.HighValid {
			filter = append(filter, fmt.Sprintf("total_time <= %f", s.Time.High))
		}
	}
	if s.NPS.Valid {
		filter = append(filter, fmt.Sprintf("total_nps >= %f", s.NPS.Low))
		if s.NPS.HighValid {
			filter = append(filter, fmt.Sprintf("total_nps <= %f", s.NPS.High))
		}
	}
	if s.SP.Valid {
		if s.SP.Value {
			filter = append(filter, "sp_rhythm = true")
		} else {
			filter = append(filter, "sp_rhythm != true")
		}
	}
	if s.Regular.Valid {
		if s.Regular.Value {
			filter = append(filter, "irregular = 1")
		} else {
			filter = append(filter, "irregular != 1")
		}
	}
	return filter
}

func (s *SearchBody) Search() (documents []generalModels.BestdoriChartItem, totalCount int64, totalPage int64, err error) {
	documents, totalCount, err = database.Query(s.queryString, s.page, s.limit, s.Filter())
	totalPage = int64(math.Ceil(float64(totalPage) / float64(s.limit)))
	return documents, totalCount, totalPage, err
}
