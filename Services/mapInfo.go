package Services

import (
	"ayachanV2/Models"
	"ayachanV2/Models/chartFormat"
	"ayachanV2/Models/mapFormat"
	"math"
)

// basicInfoGetter 获得最基础的谱面信息(除了Irregular项)、Hit总数、HPS
func basicInfoGetter(Map mapFormat.Chart) (info Models.MapInfoBasic, TotalHitCount int, TotalHPS float64, err error) {
	var BPMList map[float64]float64
	var firstNoteTime float64
	BPMList = make(map[float64]float64)
	noteFlag := true // 检查前置区间内是否无note
	info.BPMLow = math.MaxFloat64
	info.BPMHigh = -1.0
	currentBPM := 120.0
	currentBPMStartTime := 0.0
	MainBPMTime := -1.0
	for _, note := range Map {
		switch note.Type {
		case mapFormat.NoteTypeSingle:
			if noteFlag {
				firstNoteTime = note.Time
				noteFlag = false
			}
			info.TotalNote++
			TotalHitCount++
			if note.Flick && note.Direction != 0 {
				info.SPRhythm = true
			}
		case mapFormat.NoteTypeSlide:
			if noteFlag {
				firstNoteTime = note.Time
				noteFlag = false
			}
			if !note.Hidden {
				info.TotalNote++
				if note.Status == mapFormat.SlideStart {
					TotalHitCount++
				}
			} else {
				info.SPRhythm = true
			}
		case mapFormat.NoteTypeBpm:
			if !noteFlag {
				info.BPMLow = math.Min(info.BPMLow, currentBPM)
				info.BPMHigh = math.Max(info.BPMHigh, currentBPM)
				BPMList[currentBPM] += note.Time - currentBPMStartTime
				if BPMList[currentBPM] > MainBPMTime {
					info.MainBPM = currentBPM
					MainBPMTime = BPMList[currentBPM]
				}
			}
			currentBPM = note.BPM
			currentBPMStartTime = note.Time
		}
	}
	// Append最后一个BPM数据
	info.BPMLow = math.Min(info.BPMLow, currentBPM)
	info.BPMHigh = math.Max(info.BPMHigh, currentBPM)
	BPMList[currentBPM] += Map[len(Map)-1].Time - currentBPMStartTime
	if BPMList[currentBPM] > MainBPMTime {
		info.MainBPM = currentBPM
		MainBPMTime = BPMList[currentBPM]
	}
	info.TotalTime = math.Max(20.0, Map[len(Map)-1].Time-firstNoteTime)
	info.TotalNPS = float64(info.TotalNote) / info.TotalTime
	TotalHPS = float64(TotalHitCount) / info.TotalTime
	return info, TotalHitCount, TotalHPS, nil
}

// counter 谱面计数器，谱面Note计数
func counter(Map mapFormat.Chart) (NoteCount Models.NoteCount) {
	for _, note := range Map {
		switch note.Type {
		case mapFormat.NoteTypeSingle:
			if note.Flick {
				if note.Direction > 0 {
					NoteCount.Direction.Total++
					NoteCount.Direction.Right++
				} else if note.Direction < 0 {
					NoteCount.Direction.Total++
					NoteCount.Direction.Left++
				} else {
					NoteCount.Flick++
				}
			} else {
				NoteCount.Single++
			}
		case mapFormat.NoteTypeSlide:
			switch note.Status {
			case mapFormat.SlideStart:
				NoteCount.SlideStart++
			case mapFormat.SlideTick:
				if note.Hidden {
					NoteCount.SlideHidden++
				} else {
					NoteCount.SlideTick++
				}
			case mapFormat.SlideEnd:
				if note.Flick {
					NoteCount.SlideFlick++
				} else {
					NoteCount.SlideEnd++
				}
			}
		}
	}
	return NoteCount
}

//distribution 谱面分布计数器，MaxScreenNPS
func distribution(Map mapFormat.Chart, totalTime float64) (MaxScreenNPS float64, Distribution Models.Distribution) {
	Distribution.Note = make([]int, int(math.Ceil(totalTime+0.01)))
	Distribution.Hit = make([]int, int(math.Ceil(totalTime+0.01)))
	for _, note := range Map {
		switch note.Type {
		case mapFormat.NoteTypeSingle:
			Distribution.Note[int(math.Floor(note.Time))]++
			Distribution.Hit[int(math.Floor(note.Time))]++
		case mapFormat.NoteTypeSlide:
			if note.Status == mapFormat.SlideStart {
				Distribution.Note[int(math.Floor(note.Time))]++
				Distribution.Hit[int(math.Floor(note.Time))]++
			} else {
				if !note.Hidden {
					Distribution.Note[int(math.Floor(note.Time))]++
				}
			}
		}
	}
	for _, item := range Distribution.Note {
		MaxScreenNPS = math.Max(MaxScreenNPS, float64(item))
	}
	return MaxScreenNPS, Distribution
}

// StandardInfoGetter 获得标准谱面信息,除了Irregular项
func StandardInfoGetter(Map mapFormat.Chart) (StandardInfo Models.MapInfoStandard) {
	StandardInfo.MapInfoBasic, StandardInfo.TotalHitNote, StandardInfo.TotalHPS, _ = basicInfoGetter(Map)
	StandardInfo.NoteCount = counter(Map)
	StandardInfo.MaxScreenNPS, StandardInfo.Distribution = distribution(Map, Map[len(Map)-1].Time)
	return StandardInfo
}

// ExtendInfoGetter 获得扩展谱面信息
/*func ExtendInfoGetter(ParsedMap mapFormat.ParsedChart) (ExtendInfo Models.MapInfoExtend) {
	var leftCount ,RightCount int
	var MaxSpeed,FingerMaxHPSLeft,FingerMaxHPSRight, FlickNoteInterval,NoteFlickInterval Models.Float64Heap
	for i,note := range ParsedMap{
		if note.Hand == mapFormat.LeftHand{
			if leftCount != 0 && !(note.Type == mapFormat.NoteTypeSlide && note.Status != mapFormat.SlideStart){
				FingerMaxHPSLeft.Push(utils.Reciprocal(note.GetIntervalFront()))
			}
			leftCount++
		}else{
			if RightCount != 0 && !(note.Type == mapFormat.NoteTypeSlide && note.Status != mapFormat.SlideStart){
				FingerMaxHPSRight.Push(utils.Reciprocal(note.GetIntervalFront()))
			}
			RightCount++
		}
		if i != 0{
			MaxSpeed.Push(note.GetGapFront() / note.GetIntervalFront())
		}
		if note.Type == mapFormat.NoteTypeSingle && note.Flick{
			FlickNoteInterval.Push(utils.Reciprocal(note.GetIntervalBack()))
			NoteFlickInterval.Push(utils.Reciprocal(note.GetIntervalFront()))
		}
	}
	totalCount := leftCount+RightCount
	ExtendInfo.LeftPercent = float64(leftCount) / float64(totalCount)
	ExtendInfo.FingerMaxHPS = math.Max(FingerMaxHPSLeft.GetTopRankAverage(),FingerMaxHPSRight.GetTopRankAverage())
	ExtendInfo.MaxSpeed = MaxSpeed.GetTopRankAverage()
	ExtendInfo.FlickNoteInterval = FlickNoteInterval.GetTopRankAverage()
	ExtendInfo.NoteFlickInterval = NoteFlickInterval.GetTopRankAverage()
	return ExtendInfo
}*/

// MapInfoGetter 获得全部的谱面信息
func MapInfoGetter(Map mapFormat.Chart, diff chartFormat.DiffType) (MapInfo Models.MapInfo) {
	//var ParsedChart mapFormat.ParsedChart
	MapInfoStandard := StandardInfoGetter(Map)
	MapDifficultyStandard := StandardDifficultyGetter(Map, diff)

	/*ParsedChart, MapInfoStandard.MapInfoBasic.IrregularInfo = ParseMap(Map)

	var MapInfoExtend, MapDifficultyExtend interface{}
	if MapInfoStandard.Irregular == Models.RegularTypeRegular {
		MapInfoExtend = ExtendInfoGetter(ParsedChart)
		MapDifficultyExtend = ExtendDifficultyGetter(ParsedChart, diff)
	}*/

	return Models.MapInfo{
		MapInfo: MapInfoStandard,
		//MapInfoExtend:         MapInfoExtend,
		MapDifficulty: MapDifficultyStandard,
		//MapDifficultyExtend:   MapDifficultyExtend,
	}
}
