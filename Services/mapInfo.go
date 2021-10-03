package Services

import (
	"ayachanV2/Models"
	"ayachanV2/Models/chartFormat"
	"ayachanV2/Models/mapFormat"
)

// basicInfoGetter 获得最基础的谱面信息(除了Irregular项)、Hit总数、HPS
func basicInfoGetter(Map mapFormat.Chart) (info Models.MapInfoBasic, TotalHitCount int, TotalHPS float64) {
	return info, TotalHitCount, TotalHPS
}

// counter 谱面计数器，BPM计数、谱面Note计数
func counter(Map mapFormat.Chart) (BPM Models.BpmInfo, NoteCount Models.NoteCount) {
	return BPM, NoteCount
}

//distribution 谱面分布计数器，MaxScreenNPS
func distribution(Map mapFormat.Chart) (MaxScreenNPS float64, Distribution Models.Distribution) {
	return MaxScreenNPS, Distribution
}

// StandardInfoGetter 获得标准谱面信息,除了Irregular项
func StandardInfoGetter(Map mapFormat.Chart) (StandardInfo Models.MapInfoStandard) {
	StandardInfo.MapInfoBasic, StandardInfo.TotalHitNote, StandardInfo.TotalHPS = basicInfoGetter(Map)
	StandardInfo.BpmInfo, StandardInfo.NoteCount = counter(Map)
	StandardInfo.MaxScreenNPS, StandardInfo.Distribution = distribution(Map)
	return StandardInfo
}

// ExtendInfoGetter 获得扩展谱面信息
func ExtendInfoGetter(ParsedMap mapFormat.ParsedChart) (ExtendInfo Models.MapInfoExtend) {
	return ExtendInfo
}

// MapInfoGetter 获得全部的谱面信息
func MapInfoGetter(Map mapFormat.Chart, diff chartFormat.DiffType) (MapInfo Models.MapInfo) {
	var ParsedChart mapFormat.ParsedChart
	MapInfoStandard := StandardInfoGetter(Map)
	MapDifficultyStandard := StandardDifficultyGetter(Map, diff)

	ParsedChart, MapInfoStandard.MapInfoBasic.IrregularInfo = ParseMap(Map)

	var MapInfoExtend, MapDifficultyExtend interface{}
	if MapInfoStandard.Irregular == Models.RegularTypeRegular {
		MapInfoExtend = ExtendInfoGetter(ParsedChart)
		MapDifficultyExtend = ExtendDifficultyGetter(ParsedChart, diff)
	}

	return Models.MapInfo{
		MapInfoStandard:       MapInfoStandard,
		MapInfoExtend:         MapInfoExtend,
		MapDifficultyStandard: MapDifficultyStandard,
		MapDifficultyExtend:   MapDifficultyExtend,
	}
}
