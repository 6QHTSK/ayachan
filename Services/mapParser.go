package Services

import (
	"ayachanV2/Models"
	"ayachanV2/Models/mapFormat"
	"fmt"
)

func getCurrentLane(time float64, lastTick *mapFormat.ParsedNote) (lane float64) {
	if lastTick.NoteAfter == nil {
		return lastTick.Lane
	}
	return (time-lastTick.Time)/(lastTick.NoteAfter.Time-lastTick.Time)*(lastTick.NoteAfter.Lane-lastTick.Lane) + lastTick.Lane
}

// ParseMap 拆谱拆谱
func ParseMap(Map mapFormat.Chart) (ParsedMap mapFormat.ParsedChart, IrregularInfo Models.IrregularInfo) {
	ParsedMap = Map.InitParseChart()

	// 第一部分，识别谱面中的绿条
	for i, note := range ParsedMap {
		if note.Type == mapFormat.NoteTypeSlide && note.Status == mapFormat.SlideStart {
			// 检出绿条头键后
			//1.串联整个绿条
			hand := mapFormat.UnknownHand
			lastTick := i
			var lastNoteBeat float64
			for j := i + 1; j < ParsedMap.Len(); j++ {
				if ParsedMap[j].Type == mapFormat.NoteTypeSlide && ParsedMap[j].Pos == note.Pos {

					if ParsedMap[j].Beat == ParsedMap[lastTick].Beat {
						return ParsedMap, Models.IrregularInfo{
							Irregular:     Models.RegularTypeIrregular,
							IrregularInfo: fmt.Sprintf("%.2fs 处出现了绿条内同时出现的键", ParsedMap[j].Time),
						}
					}

					ParsedMap[lastTick].NoteAfter = &ParsedMap[j]
					ParsedMap[j].NotePrevious = &ParsedMap[lastTick]
					lastTick = j
					if ParsedMap[j].Status == mapFormat.SlideEnd {
						lastNoteBeat = ParsedMap[j].Beat
						break
					}
				}
			}
			// 2. 向前寻找键，与绿条共存的另一组键我们先指定手，后面再连接起来
			for j := i - 1; j >= 0; j-- {
				if j != i-1 {
					return ParsedMap, Models.IrregularInfo{
						Irregular:     Models.RegularTypeIrregular,
						IrregularInfo: fmt.Sprintf("%.2fs 处出现了多压", note.Time),
					}
				}
				if ParsedMap[j].Beat < note.Beat {
					break
				}
				if ParsedMap[j].Lane < note.Lane {
					hand = mapFormat.RightHand
				} else {
					hand = mapFormat.LeftHand
				}
			}
			// 3. 向后寻找键
			lastTickNote := &ParsedMap[i]
			for j := i + 1; j <= ParsedMap.Len(); j++ {
				if ParsedMap[j].Beat > lastNoteBeat {
					break
				}
				if ParsedMap[j].Type == mapFormat.NoteTypeSlide && ParsedMap[j].Pos == note.Pos {
					lastTickNote = &ParsedMap[j]
					continue
				}
				if ParsedMap[j].Beat == ParsedMap[j-1].Beat {
					return ParsedMap, Models.IrregularInfo{
						Irregular:     Models.RegularTypeIrregular,
						IrregularInfo: fmt.Sprintf("%.2fs 处击打绿条同时要处理同时出现的音符", ParsedMap[j].Time),
					}
				}
				if hand == mapFormat.UnknownHand {
					if ParsedMap[j].Lane < getCurrentLane(ParsedMap[j].Time, lastTickNote) {
						hand = mapFormat.RightHand
					} else {
						hand = mapFormat.LeftHand
					}
				} else if hand == mapFormat.LeftHand {
					if ParsedMap[j].Lane < getCurrentLane(ParsedMap[j].Time, lastTickNote) {
						return ParsedMap, Models.IrregularInfo{
							Irregular:     Models.RegularTypeIrregular,
							IrregularInfo: fmt.Sprintf("%.2fs 处右手跨过左手的绿条", ParsedMap[j].Time),
						}
					}
				} else {
					if ParsedMap[j].Lane > getCurrentLane(ParsedMap[j].Time, lastTickNote) {
						return ParsedMap, Models.IrregularInfo{
							Irregular:     Models.RegularTypeIrregular,
							IrregularInfo: fmt.Sprintf("%.2fs 处左手跨过右手的绿条", ParsedMap[j].Time),
						}
					}
				}
			}
			// 4. 给开头的绿条刚才的手分配，确定
			currentTick := &ParsedMap[i]
			for currentTick != nil {
				currentTick.Hand = hand
				currentTick = currentTick.NoteAfter
			}
		}
	}

	// 第二部分，识别谱面中的双压

	//TODO 拆谱部分

	return ParsedMap, IrregularInfo
}
