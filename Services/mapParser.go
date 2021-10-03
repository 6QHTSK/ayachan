package Services

import (
	"ayachanV2/Models"
	"ayachanV2/Models/mapFormat"
)

// ParseMap 拆谱拆谱
func ParseMap(Map mapFormat.Chart) (ParsedMap mapFormat.ParsedChart, IrregularInfo Models.IrregularInfo) {
	ParsedMap = Map.InitParseChart()

	//TODO 拆谱部分

	return ParsedMap, IrregularInfo
}
