package Services

import (
	"ayachanV2/Models"
	"ayachanV2/Models/chartFormat"
)

// diffLevels 对应每个等级的记录最低Level，下一个为记录最高Level
var diffLevels = [5]int{5, 11, 16, 21, 29}

const (
	diffTypeTotalNPS int = iota
	diffTypeTotalHPS
	diffTypeMaxScreenNPS
	diffTypeMaxSpeed
	diffTypeFingerMaxHPS
	diffTypeFlickNoteInterval
	diffTypeNoteFlickInterval
)

// totalNPSStandard 对应NPS标准，减去最低Level为offset
var standards = [][4][]float64{
	{ // totalNPS
		{0.61, 0.85, 1.16, 1.43, 1.6, 1.77, 2.1},
		{1.62, 1.94, 2.23, 2.64, 2.96, 3.51},
		{3.23, 3.81, 4.27, 4.67, 5.06, 5.49},
		{3.84, 4.6, 4.81, 5.47, 6.13, 7.17, 8.12, 8.55, 9.85},
	},
	{ // totalHPS
		{0.58, 0.78, 0.98, 1.19, 1.33, 1.47, 1.72},
		{1.49, 1.69, 1.92, 2.22, 2.49, 2.99},
		{2.86, 3.13, 3.54, 3.84, 4.01, 4.32},
		{3, 3.4, 3.81, 4.42, 5.02, 5.83, 6.53, 7, 7.77},
	},
	{ // MaxScreenNPS
		{1.19, 1.55, 2.01, 2.37, 2.55, 2.88, 3.36},
		{2.63, 3.05, 3.43, 3.85, 4.26, 4.91},
		{4.77, 5.31, 5.75, 6.18, 6.74, 7.12},
		{5.82, 6.6, 6.7, 7.29, 8.07, 9.22, 10.5, 11.19, 13.01},
	},
	{ // MaxSpeed
		{1.19, 1.55, 2.01, 2.37, 2.55, 2.88, 3.36},
		{2.63, 3.05, 3.43, 3.85, 4.26, 4.91},
		{5.89, 7.07, 7.99, 9.58, 9.77, 11.63},
		{6.56, 8.29, 8.66, 9.64, 11.5, 13.65, 15.61, 15.93, 19.66},
	},
	{ //FingerMaxHPS
		{0.63, 1.01, 1.33, 1.71, 1.87, 2.06, 2.43},
		{2.02, 2.32, 2.55, 2.89, 3.16, 3.65},
		{3.95, 4.17, 4.7, 5.07, 5.32, 5.81},
		{4.87, 4.93, 4.97, 5.49, 5.94, 6.64, 7.35, 7.88, 8.69},
	},
	{ //FlickNoteInterval
		{},
		{},
		{1.31, 1.67, 2.06, 2.19, 2.4, 2.71},
		{1.77, 2.4, 2.39, 2.43, 3.17, 3.86, 4.54, 5.19, 5.84},
	},
	{ //NoteFlickInterval
		{},
		{},
		{2.08, 2.19, 2.6, 2.66, 3.1, 3.6},
		{2.62, 2.84, 3.08, 3.51, 4.41, 5.21, 5.52, 5.84, 7.17},
	},
}

var maxValue = [][3]float64{
	{2.45, 4, 6.45},
	{2, 3.35, 4.9},
}

// 要一定能拿到等级的
func getLevelCalc(diffType int, diff chartFormat.DiffType, value float64) (level float64) {
	if diffType == diffTypeFlickNoteInterval || diffType == diffTypeNoteFlickInterval {
		if diff <= 1 {
			return level
		}
	}
	baseLevel := float64(diffLevels[diff]) + 0.5
	standard := standards[diffType][diff]
	for i, item := range standard {
		if item > value {
			if i == 0 {
				return baseLevel
			} else {
				return (baseLevel + float64(i)) + (value-item)/(item-standard[i-1])
			}
		}
	}
	MaxLevel := float64(diffLevels[diff+1]) + 0.5
	levelPace := standard[len(standard)-1] - standard[len(standard)-2]
	MaxLevelPos := standard[len(standard)-1]
	return MaxLevel + (value-MaxLevelPos)/levelPace
}

func getTrueDiff(diff chartFormat.DiffType, totalNPS float64, totalHPS float64) (trueDiff chartFormat.DiffType) {
	if diff <= 2 && (totalNPS > maxValue[diffTypeTotalNPS][diff] || totalHPS > maxValue[diffTypeTotalHPS][diff]) {
		return getTrueDiff(diff+1, totalNPS, totalHPS)
	} else if diff == 4 {
		return chartFormat.Diff_Expert
	}
	return diff
}

func getLevelCompare(diffType int, diff chartFormat.DiffType, value float64, totalLevel float64) (status Models.DifficultyDescription) {
	if diffType == diffTypeFlickNoteInterval || diffType == diffTypeNoteFlickInterval {
		if diff <= 1 && value > 0 {
			return Models.DifficultyHigh
		} else {
			return Models.DifficultyNormal
		}
	}
	level := getLevelCalc(diffType, diff, value)
	if level-totalLevel > 0.5 {
		return Models.DifficultyHigh
	} else if level-totalLevel < -0.5 {
		return Models.DifficultyLow
	} else {
		return Models.DifficultyNormal
	}
}

// StandardDifficultyGetter 获得标准 谱面难度 信息,除了Irregular项
func StandardDifficultyGetter(Info Models.MapInfoStandard, diff chartFormat.DiffType) (StandardDifficulty Models.MapDifficultyStandard, newDiff chartFormat.DiffType) {
	newDiff = getTrueDiff(diff, Info.TotalNPS, Info.TotalHPS)
	StandardDifficulty = Models.MapDifficultyStandard{
		TotalNPS:     getLevelCalc(diffTypeTotalNPS, newDiff, Info.TotalNPS),
		TotalHPS:     getLevelCalc(diffTypeTotalHPS, newDiff, Info.TotalHPS),
		MaxScreenNPS: getLevelCalc(diffTypeMaxScreenNPS, newDiff, Info.MaxScreenNPS),
	}
	StandardDifficulty.Difficulty = (StandardDifficulty.TotalNPS*4.0 + StandardDifficulty.TotalHPS*2.0 + StandardDifficulty.MaxScreenNPS) / 7.0
	return StandardDifficulty, newDiff
}

// ExtendDifficultyGetter 获得标准 谱面难度 信息,除了Irregular项
func ExtendDifficultyGetter(Info Models.MapInfoExtend, diff chartFormat.DiffType, level float64) (ExtendDifficulty Models.MapDifficultyExtend) {
	return Models.MapDifficultyExtend{
		MaxSpeed:          getLevelCompare(diffTypeMaxSpeed, diff, Info.MaxSpeed, level),
		FingerMaxHPS:      getLevelCompare(diffTypeFingerMaxHPS, diff, Info.FingerMaxHPS, level),
		FlickNoteInterval: getLevelCompare(diffTypeFlickNoteInterval, diff, Info.FlickNoteInterval, level),
		NoteFlickInterval: getLevelCompare(diffTypeNoteFlickInterval, diff, Info.NoteFlickInterval, level),
	}
}
