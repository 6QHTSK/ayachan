package Services

import (
	"ayachanV2/Models"
	"ayachanV2/Models/chartFormat"
	"ayachanV2/Models/mapFormat"
)

// StandardDifficultyGetter 获得标准 谱面难度 信息,除了Irregular项
func StandardDifficultyGetter(Map mapFormat.Chart, diff chartFormat.DiffType) (StandardDifficulty Models.MapDifficultyStandard) {
	return StandardDifficulty
}

// ExtendDifficultyGetter 获得标准 谱面难度 信息,除了Irregular项
func ExtendDifficultyGetter(Map mapFormat.ParsedChart, diff chartFormat.DiffType) (ExtendDifficulty Models.MapDifficultyExtend) {
	return ExtendDifficulty
}
