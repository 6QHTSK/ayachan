package pkg

import (
	"container/heap"
	"math"
)

// basicMetricsGetter 获得最基础的谱面信息(除了Irregular项)、Hit总数、HPS
func basicMetricsGetter(Map Chart) (info MapMetricsBasic, TotalHitCount int, TotalHPS float64, BPMInfo BpmInfo, err error) {
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
		case NoteTypeSingle:
			if noteFlag {
				firstNoteTime = note.Time
				noteFlag = false
			}
			info.TotalNote++
			TotalHitCount++
			if note.Flick && note.Direction != 0 {
				info.SPRhythm = true
			}
		case NoteTypeSlide:
			if noteFlag {
				firstNoteTime = note.Time
				noteFlag = false
			}
			if !note.Hidden {
				info.TotalNote++
				if note.Status == SlideStart {
					TotalHitCount++
				}
			} else {
				info.SPRhythm = true
			}
		case NoteTypeBpm:
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
func counter(Map Chart) (NoteCount NoteCount) {
	for _, note := range Map {
		switch note.Type {
		case NoteTypeSingle:
			if note.Flick {
				if note.Direction > 0 {
					NoteCount.DirectionRight++
				} else if note.Direction < 0 {
					NoteCount.DirectionLeft++
				} else {
					NoteCount.Flick++
				}
			} else {
				NoteCount.Single++
			}
		case NoteTypeSlide:
			switch note.Status {
			case SlideStart:
				NoteCount.SlideStart++
			case SlideTick:
				if note.Hidden {
					NoteCount.SlideHidden++
				} else {
					NoteCount.SlideTick++
				}
			case SlideEnd:
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
func distribution(Map Chart, totalTime float64) (MaxScreenNPS float64, Distribution Distribution) {
	Distribution.Note = make([]int, int(math.Ceil(totalTime+0.01)))
	Distribution.Hit = make([]int, int(math.Ceil(totalTime+0.01)))
	var MaxScreenNPSHeap Float64Heap
	heap.Init(&MaxScreenNPSHeap)
	for _, note := range Map {
		switch note.Type {
		case NoteTypeSingle:
			Distribution.Note[int(math.Floor(note.Time))]++
			Distribution.Hit[int(math.Floor(note.Time))]++
		case NoteTypeSlide:
			if note.Status == SlideStart {
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

// reciprocal 去除0的倒数为inf的问题
func reciprocal(num float64) (r float64) {
	if num == 0.0 {
		return r
	}
	return 1.0 / num
}

// extendMetricsGetter 获得扩展谱面信息
func extendMetricsGetter(ParsedMap ParsedChart) *MapMetricsExtend {
	var leftCount, RightCount int
	var MaxSpeed, FingerMaxHPSLeft, FingerMaxHPSRight, FlickNoteInterval, NoteFlickInterval Float64Heap
	heap.Init(&MaxSpeed)
	heap.Init(&FingerMaxHPSLeft)
	heap.Init(&FingerMaxHPSRight)
	heap.Init(&FlickNoteInterval)
	heap.Init(&NoteFlickInterval)
	for i, note := range ParsedMap {
		if note.Hand == LeftHand {
			if leftCount != 0 && !(note.Type == NoteTypeSlide && note.Status != SlideStart) {
				heap.Push(&FingerMaxHPSLeft, reciprocal(note.GetIntervalFront()))
			}
			leftCount++
		} else {
			if RightCount != 0 && !(note.Type == NoteTypeSlide && note.Status != SlideStart) {
				heap.Push(&FingerMaxHPSRight, reciprocal(note.GetIntervalFront()))
			}
			RightCount++
		}
		if i != 0 {
			heap.Push(&MaxSpeed, math.Abs(note.GetGapFront())/note.GetIntervalFront())
		}
		if note.Type == NoteTypeSingle && note.Flick {
			heap.Push(&FlickNoteInterval, reciprocal(note.GetIntervalBack()))
			heap.Push(&NoteFlickInterval, reciprocal(note.GetIntervalFront()))
		}
	}
	totalCount := leftCount + RightCount
	return &MapMetricsExtend{
		LeftPercent:       float64(leftCount) / float64(totalCount),
		MaxSpeed:          math.Max(FingerMaxHPSLeft.GetTopRankAverage(), FingerMaxHPSRight.GetTopRankAverage()),
		FingerMaxHPS:      MaxSpeed.GetTopRankAverage(),
		FlickNoteInterval: FlickNoteInterval.GetTopRankAverage(),
		NoteFlickInterval: NoteFlickInterval.GetTopRankAverage(),
	}
}

// StandardInfoGetter 获得标准谱面信息,除了Irregular项
func StandardInfoGetter(Map Chart) (StandardInfo MapMetricsStandard) {
	StandardInfo.MapMetricsBasic, StandardInfo.TotalHitNote, StandardInfo.TotalHPS, StandardInfo.BpmInfo, _ = basicMetricsGetter(Map)
	StandardInfo.NoteCount = counter(Map)
	StandardInfo.MaxScreenNPS, StandardInfo.Distribution = distribution(Map, Map[len(Map)-1].Time)
	return StandardInfo
}

// MapInfoGetter 获得全部的谱面信息
func MapInfoGetter(inputMap InputMap, diff int) (MapInfo, error) {
	suc, err := inputMap.MapCheck()
	if !suc {
		return MapInfo{}, err
	}
	Map := inputMap.Decode()
	MapInfoStandard := StandardInfoGetter(Map)
	MapDifficultyStandard, diff := StandardDifficultyGetter(MapInfoStandard, diff)

	var ParsedMap ParsedChart
	ParsedMap, MapInfoStandard.MapMetricsBasic.IrregularInfo = ParseMap(Map)

	var MapInfoExtend *MapMetricsExtend
	var MapDifficultyExtend *MapDifficultyExtend
	if MapInfoStandard.Irregular == RegularTypeRegular {
		MapInfoExtend = extendMetricsGetter(ParsedMap)
		MapDifficultyExtend = ExtendDifficultyGetter(*MapInfoExtend, diff, MapDifficultyStandard.Difficulty)
	}

	return MapInfo{
		MapMetrics:          &MapInfoStandard,
		MapMetricsExtend:    MapInfoExtend,
		MapDifficulty:       &MapDifficultyStandard,
		MapDifficultyExtend: MapDifficultyExtend,
	}, nil
}
