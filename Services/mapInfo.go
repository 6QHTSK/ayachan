package Services

import (
	"ayachanV2/Models"
	"ayachanV2/Models/chartFormat"
	"ayachanV2/Models/mapFormat"
	"ayachanV2/utils"
	"container/heap"
	"math"
)

// basicInfoGetter 获得最基础的谱面信息(除了Irregular项)、Hit总数、HPS
func basicInfoGetter(Map mapFormat.Chart) (info Models.MapInfoBasic, TotalHitCount int, TotalHPS float64, BPMInfo Models.BpmInfo, err error) {
	var BPMList map[float64]float64
	var firstNoteTime float64
	BPMList = make(map[float64]float64)
	noteFlag := true // 检查前置区间内是否无note
	BPMInfo.BPMLow = math.MaxFloat64
	BPMInfo.BPMHigh = -1.0
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
				BPMInfo.BPMLow = math.Min(BPMInfo.BPMLow, currentBPM)
				BPMInfo.BPMHigh = math.Max(BPMInfo.BPMHigh, currentBPM)
				BPMList[currentBPM] += note.Time - currentBPMStartTime
				if BPMList[currentBPM] > MainBPMTime {
					BPMInfo.MainBPM = currentBPM
					MainBPMTime = BPMList[currentBPM]
				}
			}
			currentBPM = note.BPM
			currentBPMStartTime = note.Time
		}
	}
	// Append最后一个BPM数据
	BPMInfo.BPMLow = math.Min(BPMInfo.BPMLow, currentBPM)
	BPMInfo.BPMHigh = math.Max(BPMInfo.BPMHigh, currentBPM)
	BPMList[currentBPM] += Map[len(Map)-1].Time - currentBPMStartTime
	if BPMList[currentBPM] > MainBPMTime {
		BPMInfo.MainBPM = currentBPM
		MainBPMTime = BPMList[currentBPM]
	}
	info.TotalTime = math.Max(20.0, Map[len(Map)-1].Time-firstNoteTime)
	info.TotalNPS = float64(info.TotalNote) / info.TotalTime
	TotalHPS = float64(TotalHitCount) / info.TotalTime
	return info, TotalHitCount, TotalHPS, BPMInfo, nil
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
	var MaxScreenNPSHeap Models.Float64Heap
	heap.Init(&MaxScreenNPSHeap)
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
		//MaxScreenNPS = math.Max(MaxScreenNPS, float64(item))
		heap.Push(&MaxScreenNPSHeap, float64(item))
	}
	return MaxScreenNPSHeap.GetTopRankAverage(), Distribution
}

// StandardInfoGetter 获得标准谱面信息,除了Irregular项
func StandardInfoGetter(Map mapFormat.Chart) (StandardInfo Models.MapInfoStandard) {
	StandardInfo.MapInfoBasic, StandardInfo.TotalHitNote, StandardInfo.TotalHPS, StandardInfo.BpmInfo, _ = basicInfoGetter(Map)
	StandardInfo.NoteCount = counter(Map)
	StandardInfo.MaxScreenNPS, StandardInfo.Distribution = distribution(Map, Map[len(Map)-1].Time)
	return StandardInfo
}

// ExtendInfoGetter 获得扩展谱面信息
func ExtendInfoGetter(ParsedMap mapFormat.ParsedChart) (ExtendInfo Models.MapInfoExtend) {
	var leftCount, RightCount int
	var MaxSpeed, FingerMaxHPSLeft, FingerMaxHPSRight, FlickNoteInterval, NoteFlickInterval Models.Float64Heap
	heap.Init(&MaxSpeed)
	heap.Init(&FingerMaxHPSLeft)
	heap.Init(&FingerMaxHPSRight)
	heap.Init(&FlickNoteInterval)
	heap.Init(&NoteFlickInterval)
	for i, note := range ParsedMap {
		if note.Hand == mapFormat.LeftHand {
			if leftCount != 0 && !(note.Type == mapFormat.NoteTypeSlide && note.Status != mapFormat.SlideStart) {
				heap.Push(&FingerMaxHPSLeft, utils.Reciprocal(note.GetIntervalFront()))
			}
			leftCount++
		} else {
			if RightCount != 0 && !(note.Type == mapFormat.NoteTypeSlide && note.Status != mapFormat.SlideStart) {
				heap.Push(&FingerMaxHPSRight, utils.Reciprocal(note.GetIntervalFront()))
			}
			RightCount++
		}
		if i != 0 {
			heap.Push(&MaxSpeed, math.Abs(note.GetGapFront())/note.GetIntervalFront())
		}
		if note.Type == mapFormat.NoteTypeSingle && note.Flick {
			heap.Push(&FlickNoteInterval, utils.Reciprocal(note.GetIntervalBack()))
			heap.Push(&NoteFlickInterval, utils.Reciprocal(note.GetIntervalFront()))
		}
	}
	totalCount := leftCount + RightCount
	ExtendInfo.LeftPercent = float64(leftCount) / float64(totalCount)
	ExtendInfo.FingerMaxHPS = math.Max(FingerMaxHPSLeft.GetTopRankAverage(), FingerMaxHPSRight.GetTopRankAverage())
	ExtendInfo.MaxSpeed = MaxSpeed.GetTopRankAverage()
	ExtendInfo.FlickNoteInterval = FlickNoteInterval.GetTopRankAverage()
	ExtendInfo.NoteFlickInterval = NoteFlickInterval.GetTopRankAverage()
	return ExtendInfo
}

// MapInfoGetter 获得全部的谱面信息
func MapInfoGetter(Map mapFormat.Chart, diff chartFormat.DiffType) (MapInfo Models.MapInfo) {
	//var ParsedChart mapFormat.ParsedChart
	MapInfoStandard := StandardInfoGetter(Map)
	MapDifficultyStandard, diff := StandardDifficultyGetter(MapInfoStandard, diff)

	var ParsedMap mapFormat.ParsedChart
	ParsedMap, MapInfoStandard.MapInfoBasic.IrregularInfo = ParseMap(Map)
	//_, MapInfoStandard.MapInfoBasic.IrregularInfo = ParseMap(Map)

	var MapInfoExtend, MapDifficultyExtend interface{}
	if MapInfoStandard.Irregular == Models.RegularTypeRegular {
		MapInfoExtend = ExtendInfoGetter(ParsedMap)
		MapDifficultyExtend = ExtendDifficultyGetter(MapInfoExtend.(Models.MapInfoExtend), diff, MapDifficultyStandard.Difficulty)
	}

	return Models.MapInfo{
		MapInfo:             MapInfoStandard,
		MapInfoExtend:       MapInfoExtend,
		MapDifficulty:       MapDifficultyStandard,
		MapDifficultyExtend: MapDifficultyExtend,
	}
}
