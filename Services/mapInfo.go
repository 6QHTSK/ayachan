package Services

import (
	"container/heap"
	"github.com/6QHTSK/ayachan/Models"
	"github.com/6QHTSK/ayachan/Models/ChartFormat"
	"github.com/6QHTSK/ayachan/Models/MapFormat"
	"github.com/6QHTSK/ayachan/utils"
	"math"
)

// basicMetricsGetter 获得最基础的谱面信息(除了Irregular项)、Hit总数、HPS
func basicMetricsGetter(Map MapFormat.Chart) (info Models.MapMetricsBasic, TotalHitCount int, TotalHPS float64, BPMInfo Models.BpmInfo, err error) {
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
		case MapFormat.NoteTypeSingle:
			if noteFlag {
				firstNoteTime = note.Time
				noteFlag = false
			}
			info.TotalNote++
			TotalHitCount++
			if note.Flick && note.Direction != 0 {
				info.SPRhythm = true
			}
		case MapFormat.NoteTypeSlide:
			if noteFlag {
				firstNoteTime = note.Time
				noteFlag = false
			}
			if !note.Hidden {
				info.TotalNote++
				if note.Status == MapFormat.SlideStart {
					TotalHitCount++
				}
			} else {
				info.SPRhythm = true
			}
		case MapFormat.NoteTypeBpm:
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
func counter(Map MapFormat.Chart) (NoteCount Models.NoteCount) {
	for _, note := range Map {
		switch note.Type {
		case MapFormat.NoteTypeSingle:
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
		case MapFormat.NoteTypeSlide:
			switch note.Status {
			case MapFormat.SlideStart:
				NoteCount.SlideStart++
			case MapFormat.SlideTick:
				if note.Hidden {
					NoteCount.SlideHidden++
				} else {
					NoteCount.SlideTick++
				}
			case MapFormat.SlideEnd:
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
func distribution(Map MapFormat.Chart, totalTime float64) (MaxScreenNPS float64, Distribution Models.Distribution) {
	Distribution.Note = make([]int, int(math.Ceil(totalTime+0.01)))
	Distribution.Hit = make([]int, int(math.Ceil(totalTime+0.01)))
	var MaxScreenNPSHeap Models.Float64Heap
	heap.Init(&MaxScreenNPSHeap)
	for _, note := range Map {
		switch note.Type {
		case MapFormat.NoteTypeSingle:
			Distribution.Note[int(math.Floor(note.Time))]++
			Distribution.Hit[int(math.Floor(note.Time))]++
		case MapFormat.NoteTypeSlide:
			if note.Status == MapFormat.SlideStart {
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
func StandardInfoGetter(Map MapFormat.Chart) (StandardInfo Models.MapMetricsStandard) {
	StandardInfo.MapMetricsBasic, StandardInfo.TotalHitNote, StandardInfo.TotalHPS, StandardInfo.BpmInfo, _ = basicMetricsGetter(Map)
	StandardInfo.NoteCount = counter(Map)
	StandardInfo.MaxScreenNPS, StandardInfo.Distribution = distribution(Map, Map[len(Map)-1].Time)
	return StandardInfo
}

// ExtendMetricsGetter 获得扩展谱面信息
func ExtendMetricsGetter(ParsedMap MapFormat.ParsedChart) (ExtendInfo *Models.MapMetricsExtend) {
	var leftCount, RightCount int
	var MaxSpeed, FingerMaxHPSLeft, FingerMaxHPSRight, FlickNoteInterval, NoteFlickInterval Models.Float64Heap
	heap.Init(&MaxSpeed)
	heap.Init(&FingerMaxHPSLeft)
	heap.Init(&FingerMaxHPSRight)
	heap.Init(&FlickNoteInterval)
	heap.Init(&NoteFlickInterval)
	for i, note := range ParsedMap {
		if note.Hand == MapFormat.LeftHand {
			if leftCount != 0 && !(note.Type == MapFormat.NoteTypeSlide && note.Status != MapFormat.SlideStart) {
				heap.Push(&FingerMaxHPSLeft, utils.Reciprocal(note.GetIntervalFront()))
			}
			leftCount++
		} else {
			if RightCount != 0 && !(note.Type == MapFormat.NoteTypeSlide && note.Status != MapFormat.SlideStart) {
				heap.Push(&FingerMaxHPSRight, utils.Reciprocal(note.GetIntervalFront()))
			}
			RightCount++
		}
		if i != 0 {
			heap.Push(&MaxSpeed, math.Abs(note.GetGapFront())/note.GetIntervalFront())
		}
		if note.Type == MapFormat.NoteTypeSingle && note.Flick {
			heap.Push(&FlickNoteInterval, utils.Reciprocal(note.GetIntervalBack()))
			heap.Push(&NoteFlickInterval, utils.Reciprocal(note.GetIntervalFront()))
		}
	}
	totalCount := leftCount + RightCount
	return &Models.MapMetricsExtend{
		LeftPercent:       float64(leftCount) / float64(totalCount),
		MaxSpeed:          math.Max(FingerMaxHPSLeft.GetTopRankAverage(), FingerMaxHPSRight.GetTopRankAverage()),
		FingerMaxHPS:      MaxSpeed.GetTopRankAverage(),
		FlickNoteInterval: FlickNoteInterval.GetTopRankAverage(),
		NoteFlickInterval: NoteFlickInterval.GetTopRankAverage(),
	}
}

// MapInfoGetter 获得全部的谱面信息
func MapInfoGetter(Map MapFormat.Chart, diff ChartFormat.DiffType) (MapInfo Models.MapInfo) {
	MapInfoStandard := StandardInfoGetter(Map)
	MapDifficultyStandard, diff := StandardDifficultyGetter(MapInfoStandard, diff)

	var ParsedMap MapFormat.ParsedChart
	ParsedMap, MapInfoStandard.MapMetricsBasic.IrregularInfo = ParseMap(Map)

	var MapInfoExtend *Models.MapMetricsExtend
	var MapDifficultyExtend *Models.MapDifficultyExtend
	if MapInfoStandard.Irregular == Models.RegularTypeRegular {
		MapInfoExtend = ExtendMetricsGetter(ParsedMap)
		MapDifficultyExtend = ExtendDifficultyGetter(*MapInfoExtend, diff, MapDifficultyStandard.Difficulty)
	}

	return Models.MapInfo{
		MapMetrics:          &MapInfoStandard,
		MapMetricsExtend:    MapInfoExtend,
		MapDifficulty:       &MapDifficultyStandard,
		MapDifficultyExtend: MapDifficultyExtend,
	}
}
