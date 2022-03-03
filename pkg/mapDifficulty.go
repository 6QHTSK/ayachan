package pkg

// diffLevels 对应每个等级的记录最低Level，下一个为记录最高Level
var diffLevels = [5]int{5, 11, 16, 22, 30}

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
		{0.538, 0.621, 0.910, 1.265, 1.503, 1.789, 2.026},
		{1.334, 1.544, 2.024, 2.421, 2.866, 3.592},
		{2.628, 3.227, 3.981, 4.535, 4.959, 5.562, 5.984},
		{4.097, 4.464, 4.880, 5.438, 6.585, 7.978, 8.870, 9.470, 11.482},
	},
	{ // totalHPS
		{0.531, 0.579, 0.797, 1.054, 1.267, 1.501, 1.703},
		{1.134, 1.397, 1.695, 2.024, 2.431, 3.074},
		{1.753, 2.681, 3.229, 3.721, 4.079, 4.453, 4.829},
		{3.003, 3.295, 3.816, 4.366, 5.367, 6.481, 6.980, 7.771, 8.516},
	},
	{ // MaxScreenNPS
		{1.094, 1.161, 1.576, 2.138, 2.452, 2.857, 3.345},
		{2.219, 2.567, 3.111, 3.519, 4.207, 5.069},
		{3.833, 4.710, 5.364, 6.032, 6.500, 7.296, 7.906},
		{6.115, 6.353, 6.677, 7.296, 8.536, 10.219, 11.310, 12.704, 14.345},
	},
	{ // MaxSpeed
		{0.474, 0.520, 0.855, 1.286, 1.608, 1.937, 2.290},
		{1.771, 2.135, 2.687, 3.276, 4.214, 6.630},
		{3.639, 4.848, 6.494, 9.083, 10.980, 12.787, 14.834},
		{6.686, 7.080, 8.134, 9.687, 12.339, 15.671, 17.335, 19.757, 21.954},
	},
	{ //FingerMaxHPS
		{0.518, 0.596, 0.976, 1.397, 1.728, 2.207, 2.559},
		{1.417, 1.708, 2.200, 2.596, 3.159, 4.214},
		{2.755, 3.512, 4.215, 4.929, 5.464, 6.067, 6.501},
		{4.342, 4.546, 4.926, 5.387, 6.207, 7.242, 8.033, 9.270, 9.577},
	},
	{ //FlickNoteInterval
		{},
		{},
		{0.000, 0.000, 1.689, 2.222, 2.714, 3.111, 3.714},
		{0.000, 0.944, 1.778, 2.429, 3.333, 4.722, 5.483, 6.323, 6.825},
	},
	{ //NoteFlickInterval
		{},
		{},
		{0.000, 0.000, 2.063, 2.747, 3.431, 4.107, 4.971},
		{0.000, 1.689, 2.500, 3.514, 4.800, 5.906, 6.667, 8.629, 9.481},
	},
}

var maxValue = [][3]float64{
	{2.026, 3.592, 5.984}, // diffTypeTotalNPS
	{1.703, 3.074, 4.829}, // diffTypeTotalHPS
}

// 要一定能拿到等级的
func getLevelCalc(diffType int, diff int, value float64) (level float64) {
	baseLevel := float64(diffLevels[diff])
	standard := standards[diffType][diff]
	for i, threshold := range standard {
		if threshold > value { // 找到了较高等级
			if i == 0 {
				return baseLevel
			} else {
				lastThreshold := standard[i-1]
				return (baseLevel + float64(i-1)) + (value-lastThreshold)/(threshold-lastThreshold)
			}
		}
	}
	maxLevel := float64(diffLevels[diff+1])
	levelPace := (standard[len(standard)-1] - standard[len(standard)-3]) / 2
	maxLevelThreshold := standard[len(standard)-1]
	return maxLevel + (value-maxLevelThreshold)/levelPace
}

func getTrueDiff(diff int, totalNPS float64, totalHPS float64) (trueDiff int) {
	if diff <= 2 && (totalNPS > maxValue[diffTypeTotalNPS][diff] || totalHPS > maxValue[diffTypeTotalHPS][diff]) {
		return getTrueDiff(diff+1, totalNPS, totalHPS)
	} else if diff == 4 {
		return 3 // EXPERT
	}
	return diff
}

func getLevelCompare(diffType int, diff int, value float64, totalLevel float64) (status DifficultyDescription) {
	if diffType == diffTypeFlickNoteInterval || diffType == diffTypeNoteFlickInterval {
		if diff <= 1 && value > 0 {
			return DifficultyHigh
		} else {
			return DifficultyNormal
		}
	}
	level := getLevelCalc(diffType, diff, value)
	if level-totalLevel > 1 {
		return DifficultyHigh
	} else if level-totalLevel < -1 {
		return DifficultyLow
	} else {
		return DifficultyNormal
	}
}

// StandardDifficultyGetter 获得标准 谱面难度 信息,除了Irregular项
func StandardDifficultyGetter(metrics MapMetricsStandard, diff int) (StandardDifficulty MapDifficultyStandard, newDiff int) {
	newDiff = getTrueDiff(diff, metrics.TotalNPS, metrics.TotalHPS)
	StandardDifficulty = MapDifficultyStandard{
		TotalNPS:     getLevelCalc(diffTypeTotalNPS, newDiff, metrics.TotalNPS),
		TotalHPS:     getLevelCalc(diffTypeTotalHPS, newDiff, metrics.TotalHPS),
		MaxScreenNPS: getLevelCalc(diffTypeMaxScreenNPS, newDiff, metrics.MaxScreenNPS),
	}
	StandardDifficulty.Difficulty = (StandardDifficulty.TotalNPS*4.0 + StandardDifficulty.TotalHPS*2.0 + StandardDifficulty.MaxScreenNPS) / 7.0
	return StandardDifficulty, newDiff
}

// ExtendDifficultyGetter 获得标准 谱面难度 信息,除了Irregular项
func ExtendDifficultyGetter(metrics MapMetricsExtend, diff int, level float64) (ExtendDifficulty *MapDifficultyExtend) {
	return &MapDifficultyExtend{
		MaxSpeed:          getLevelCompare(diffTypeMaxSpeed, diff, metrics.MaxSpeed, level),
		FingerMaxHPS:      getLevelCompare(diffTypeFingerMaxHPS, diff, metrics.FingerMaxHPS, level),
		FlickNoteInterval: getLevelCompare(diffTypeFlickNoteInterval, diff, metrics.FlickNoteInterval, level),
		NoteFlickInterval: getLevelCompare(diffTypeNoteFlickInterval, diff, metrics.NoteFlickInterval, level),
	}
}
