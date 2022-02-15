package Services

import (
	"fmt"
	"github.com/6QHTSK/ayachan/Log"
	"github.com/6QHTSK/ayachan/Models"
	"github.com/6QHTSK/ayachan/Models/MapFormat"
	"math"
)

func getCurrentLane(time float64, lastTick *MapFormat.ParsedNote) (lane float64) {
	if lastTick.NoteAfter == nil {
		return lastTick.Lane
	}
	return (time-lastTick.Time)/(lastTick.NoteAfter.Time-lastTick.Time)*(lastTick.NoteAfter.Lane-lastTick.Lane) + lastTick.Lane
}

func isBelongToSlide(note *MapFormat.ParsedNote, pos int) bool {
	return note.Type == MapFormat.NoteTypeSlide && note.Pos == pos
}

func checkNoteWhileSlide(note MapFormat.ParsedNote, lastTickNote *MapFormat.ParsedNote, hand MapFormat.HandType) (bool, Models.IrregularInfo) {
	if hand == MapFormat.LeftHand { // (ii - 2)当前为左手绿条，右手其他键，其他键应满足lane > 绿条lane
		if note.Lane <= getCurrentLane(note.Time, lastTickNote) {
			return false, Models.IrregularInfo{
				Irregular:     Models.RegularTypeIrregular,
				IrregularInfo: fmt.Sprintf("%.2fs 处右手跨过左手击打的绿条/与绿条重叠", note.Time),
			}
		}
	} else if hand == MapFormat.RightHand { // (ii - 3)当前为右手绿条，左手其他键，其他键应满足lane < 绿条lane
		if note.Lane >= getCurrentLane(note.Time, lastTickNote) {
			return false, Models.IrregularInfo{
				Irregular:     Models.RegularTypeIrregular,
				IrregularInfo: fmt.Sprintf("%.2fs 处左手跨过右手击打的绿条/与绿条重叠", note.Time),
			}
		}
	} else { // (ii - 4) 未知的击打手
		return false, Models.IrregularInfo{
			Irregular:     Models.RegularTypeIrregular,
			IrregularInfo: fmt.Sprintf("%.2fs 处出现了击打手标记错误的情况，击打手标记%d", note.Time, hand),
		}
	}
	return true, Models.IrregularInfo{}
}

func labelNoteWhileSlide(note *MapFormat.ParsedNote, AnotherHandAvailable *bool, AnotherHandPos *int, AnotherHand MapFormat.HandType) (bool, Models.IrregularInfo) {
	if *AnotherHandAvailable { // 如果另一只手未击打绿条，直接指定即可
		if note.Type == MapFormat.NoteTypeSlide && note.Status == MapFormat.SlideStart {
			*AnotherHandAvailable = false
			*AnotherHandPos = note.Pos
		}
	} else {
		// 检查Note是否归属于绿条
		if isBelongToSlide(note, *AnotherHandPos) { // 如果归属，就正常进行
			if note.Status == MapFormat.SlideEnd {
				*AnotherHandAvailable = true
			}
		} else {
			return false, Models.IrregularInfo{
				Irregular:     Models.RegularTypeIrregular,
				IrregularInfo: fmt.Sprintf("%.2fs 处有双绿条锁手多押", note.Time),
			}
		}
	}
	note.Hand = AnotherHand
	return true, Models.IrregularInfo{}
}

func checkIfLabeled(note MapFormat.ParsedNote) bool {
	return note.Hand != MapFormat.LeftHand && note.Hand != MapFormat.RightHand
}

// ParseMap 拆谱拆谱、输入的谱排序好先
func ParseMap(Map MapFormat.Chart) (ParsedMap MapFormat.ParsedChart, IrregularInfo Models.IrregularInfo) {
	defer func() {
		err := recover()
		if err != nil {
			Log.Log.Errorf("分析谱面出现异常！异常情况：%s\n", err)
			ParsedMap = nil
			IrregularInfo = Models.IrregularInfo{
				Irregular:     Models.RegularTypeUnknown,
				IrregularInfo: fmt.Sprint(err),
			}
		}
	}()

	ParsedMap = Map.InitParseChart()

	mapLen := ParsedMap.Len()
	// 第一部分，检查是否存在多压（不考虑绿条锁手）
	for i := 2; i < mapLen; i++ {
		if ParsedMap[i].Beat == ParsedMap[i-1].Beat && ParsedMap[i-1].Beat == ParsedMap[i-2].Beat {
			return ParsedMap, Models.IrregularInfo{
				Irregular:     Models.RegularTypeIrregular,
				IrregularInfo: fmt.Sprintf("%.2fs 处出现了三押", ParsedMap[i].Time),
			}
		} else if ParsedMap[i].Beat == ParsedMap[i-1].Beat {
			if ParsedMap[i].Lane == ParsedMap[i-1].Lane {
				return ParsedMap, Models.IrregularInfo{
					Irregular:     Models.RegularTypeIrregular,
					IrregularInfo: fmt.Sprintf("%.2fs 处出现了Note重叠", ParsedMap[i].Time),
				}
			} else {
				// 由于排序，右边的Note在左边的后面,这里标记双压带来的锁手
				ParsedMap[i].Hand = MapFormat.RightHand
				ParsedMap[i-1].Hand = MapFormat.LeftHand
			}
		}
	}

	// 第二部分，识别谱面中的绿条
	for i, note := range ParsedMap {
		// 检出绿条头键后
		if note.Type == MapFormat.NoteTypeSlide && note.Status == MapFormat.SlideStart {
			// 该绿条还未被标记
			if note.Hand == MapFormat.UnknownHand {
				//1.串联整个绿条
				hand := MapFormat.UnknownHand
				lastTick := i
				var lastNoteBeat float64
				for j := i + 1; j < ParsedMap.Len(); j++ {
					if ParsedMap[j].Type == MapFormat.NoteTypeSlide && ParsedMap[j].Pos == note.Pos {
						if ParsedMap[lastTick].Beat == ParsedMap[j].Beat {
							return ParsedMap, Models.IrregularInfo{
								Irregular:     Models.RegularTypeIrregular,
								IrregularInfo: fmt.Sprintf("%.2fs 处绿条内有多个键重叠", ParsedMap[j].Time),
							}
						}
						ParsedMap[lastTick].NoteAfter = &ParsedMap[j]
						ParsedMap[j].NotePrevious = &ParsedMap[lastTick]
						lastTick = j
						if ParsedMap[j].Status == MapFormat.SlideEnd {
							lastNoteBeat = ParsedMap[j].Beat // 记录下最后一个Note的信息，便于搜索退出
							break
						}
					}
				}
				// 2. 遍历持续范围内的所有键
				lastTickNote := &ParsedMap[i]        // lastTickNote 记录上一个绿条节点，便于推算了绿条当前位置
				AnotherHand := MapFormat.UnknownHand // AnotherHand 记录另一只手是什么
				AnotherHandAvailable := true         // AnotherHandAvailable 记录另一只手是否被绿条占用
				AnotherHandPos := 0                  // 另一只手被绿条占用时的绿条编号Pos
				for j := i - 1; j < ParsedMap.Len(); j++ {
					// (1)无前一键或前一键不在绿条范围内，跳过
					if j < 0 || (j == i-1 && ParsedMap[j].Beat != ParsedMap[i].Beat) {
						continue
					}
					// (2)超出最后一个Note的范围，退出
					if ParsedMap[j].Beat > lastNoteBeat {
						break
					}
					// (3)搜到属于本绿条的内容，记录上一个Note
					if isBelongToSlide(&ParsedMap[j], note.Pos) {
						lastTickNote = &ParsedMap[j]
						continue
					}
					// 下面确定为绿条持续时间内，其它的note
					// (4)如果上一个note不属于绿条内，则认为出现了锁手双押
					if j >= 1 && ParsedMap[j].Beat == ParsedMap[j-1].Beat && !isBelongToSlide(&ParsedMap[j-1], note.Pos) {
						return ParsedMap, Models.IrregularInfo{
							Irregular:     Models.RegularTypeIrregular,
							IrregularInfo: fmt.Sprintf("%.2fs 处有单绿条锁手多押", ParsedMap[j].Time),
						}
					}
					// （5）第一次出现其他Note，记录下来
					if hand == MapFormat.UnknownHand {
						currentLane := getCurrentLane(ParsedMap[j].Time, lastTickNote)
						if ParsedMap[j].Lane < currentLane {
							hand = MapFormat.RightHand
							AnotherHand = MapFormat.LeftHand
						} else if ParsedMap[j].Lane > currentLane {
							hand = MapFormat.LeftHand
							AnotherHand = MapFormat.RightHand
						} else {
							return ParsedMap, Models.IrregularInfo{
								Irregular:     Models.RegularTypeIrregular,
								IrregularInfo: fmt.Sprintf("%.2fs 处有键与绿条重叠", ParsedMap[j].Time),
							}
						}
					} else { // (6）出现了其他note，检查是否合理
						status, IrregularInfo := checkNoteWhileSlide(ParsedMap[j], lastTickNote, hand)
						if !status {
							return ParsedMap, IrregularInfo
						}
					}
					// (7)给其他键指定击打的手
					status, IrregularInfo := labelNoteWhileSlide(&ParsedMap[j], &AnotherHandAvailable, &AnotherHandPos, AnotherHand)
					if !status {
						return ParsedMap, IrregularInfo
					}
				}
				// 4. 给开头的绿条刚才的手分配，确定
				currentTick := &ParsedMap[i]
				for currentTick != nil {
					currentTick.Hand = hand
					currentTick = currentTick.NoteAfter
				}
			} else { // 该绿条被标记，
				hand := note.Hand
				lastTick := i
				var lastNoteBeat float64
				//1.串联整个绿条
				for j := i + 1; j < ParsedMap.Len(); j++ {
					if ParsedMap[j].Type == MapFormat.NoteTypeSlide && ParsedMap[j].Pos == note.Pos {
						if ParsedMap[lastTick].Beat == ParsedMap[j].Beat {
							return ParsedMap, Models.IrregularInfo{
								Irregular:     Models.RegularTypeIrregular,
								IrregularInfo: fmt.Sprintf("%.2fs 处绿条内有多个键重叠", ParsedMap[j].Time),
							}
						}
						ParsedMap[lastTick].NoteAfter = &ParsedMap[j]
						ParsedMap[j].NotePrevious = &ParsedMap[lastTick]
						lastTick = j
						ParsedMap[j].Hand = hand
						if ParsedMap[j].Status == MapFormat.SlideEnd {
							lastNoteBeat = ParsedMap[j].Beat // 记录下最后一个Note的信息，便于搜索退出
							break
						}
					}
				}
				// 2. 遍历持续范围内的所有键
				var lastTickNote *MapFormat.ParsedNote // lastTickNote 记录上一个绿条节点，便于推算了绿条当前位置
				var AnotherHand MapFormat.HandType     // AnotherHand 记录另一只手是什么
				if hand == MapFormat.LeftHand {
					AnotherHand = MapFormat.RightHand
				} else if hand == MapFormat.RightHand {
					AnotherHand = MapFormat.LeftHand
				}
				AnotherHandAvailable := true // AnotherHandAvailable 记录另一只手是否被绿条占用
				AnotherHandPos := 0          // 另一只手被绿条占用时的绿条编号Pos
				for j := i; j < ParsedMap.Len(); j++ {
					// (2)超出最后一个Note的范围，退出
					if ParsedMap[j].Beat > lastNoteBeat {
						break
					}
					// (3)搜到属于本绿条的内容，记录上一个Note
					if isBelongToSlide(&ParsedMap[j], note.Pos) {
						lastTickNote = &ParsedMap[j]
						continue
					}
					// 下面确定为绿条持续时间内，其它的note
					status, IrregularInfo := checkNoteWhileSlide(ParsedMap[j], lastTickNote, hand)
					if !status {
						return ParsedMap, IrregularInfo
					}
					status, IrregularInfo = labelNoteWhileSlide(&ParsedMap[j], &AnotherHandAvailable, &AnotherHandPos, AnotherHand)
					if !status {
						return ParsedMap, IrregularInfo
					}
				}
			}
		}
	}

	// 第三部分，对未标记的，交互的识别和标记，注意此处的绿条、双压、绿条范围内的音符均已被标记。
	for i := range ParsedMap {
		// 跳过前两个音符
		if i < 2 {
			continue
		}
		// 检查是否被标记
		if checkIfLabeled(ParsedMap[i]) &&
			checkIfLabeled(ParsedMap[i-1]) &&
			checkIfLabeled(ParsedMap[i-2]) {

			interval1 := ParsedMap[i].Time - ParsedMap[i-1].Time
			interval2 := ParsedMap[i-1].Time - ParsedMap[i-2].Time
			// 检查是否间隔相差是否不大(10ms)且时长较短(200ms)
			if interval1 < 0.2 && interval2 < 0.2 &&
				math.Abs(interval2-interval1) < 0.01 {

				// 识别结构 left-right-left 小三角
				if ParsedMap[i-1].Lane > ParsedMap[i].Lane &&
					ParsedMap[i-1].Lane > ParsedMap[i-2].Lane {

					ParsedMap[i].Hand = MapFormat.TryLeftHand
					ParsedMap[i-1].Hand = MapFormat.TryRightHand
					ParsedMap[i-2].Hand = MapFormat.TryLeftHand
				}

				// 识别结构 right-left-right 小三角
				if ParsedMap[i-1].Lane > ParsedMap[i].Lane &&
					ParsedMap[i-1].Lane > ParsedMap[i-2].Lane {

					ParsedMap[i].Hand = MapFormat.TryRightHand
					ParsedMap[i-1].Hand = MapFormat.TryLeftHand
					ParsedMap[i-2].Hand = MapFormat.TryRightHand
				}
			}
		}
	}
	// 第四部分，全盘采用第三部分的拆分建议
	for i := range ParsedMap {
		if ParsedMap[i].Hand == MapFormat.TryLeftHand {
			ParsedMap[i].Hand = MapFormat.LeftHand
		} else if ParsedMap[i].Hand == MapFormat.TryRightHand {
			ParsedMap[i].Hand = MapFormat.RightHand
		}
	}

	// 第五部分，对其他部分采用贪心算法拆谱。可以证明，未标记的Note均为单个独立的出现，直接根据左手、右手的最终note就行了。连接左右手的Note
	var lastLeftHandNote, lastRightHandNote *MapFormat.ParsedNote
	for i := range ParsedMap {
		// 如果还没有指定哪个手
		if ParsedMap[i].Hand == MapFormat.UnknownHand {
			//Case 1 那只手近用哪只手
			var GapLeft, GapRight, Interval1, Interval2 float64

			if lastLeftHandNote == nil {
				Interval1 = 100.0
			} else {
				Interval1 = math.Abs(ParsedMap[i].Time - lastLeftHandNote.Time)
			}

			if lastRightHandNote == nil {
				Interval2 = 100.0
			} else {
				Interval2 = math.Abs(ParsedMap[i].Time - lastRightHandNote.Time)
			}

			if Interval1 > 0.5 {
				GapLeft = math.Abs(ParsedMap[i].Lane - 1.5)
			} else {
				GapLeft = math.Abs(ParsedMap[i].Lane - lastLeftHandNote.Lane)
			}

			if Interval2 > 0.5 {
				GapRight = math.Abs(ParsedMap[i].Lane - 4.5)
			} else {
				GapRight = math.Abs(ParsedMap[i].Lane - lastRightHandNote.Lane)
			}

			if GapLeft < GapRight {
				ParsedMap[i].Hand = MapFormat.LeftHand
			} else if GapLeft > GapRight {
				ParsedMap[i].Hand = MapFormat.RightHand
			} else {
				// Case 2 在左边用左手，右边用右手
				if ParsedMap[i].Lane < 2.9 {
					ParsedMap[i].Hand = MapFormat.LeftHand
				} else if ParsedMap[i].Lane > 3.1 {
					ParsedMap[i].Hand = MapFormat.RightHand
				} else {
					// Case 3 那边距上个键的间隔长，用哪只手

					if Interval1 > Interval2 {
						ParsedMap[i].Hand = MapFormat.LeftHand
					} else {
						ParsedMap[i].Hand = MapFormat.RightHand
					}
				}
			}
		}
		if ParsedMap[i].Hand == MapFormat.LeftHand {
			if lastLeftHandNote != nil {
				lastLeftHandNote.NoteAfter = &ParsedMap[i]
			}
			ParsedMap[i].NotePrevious = lastLeftHandNote
			lastLeftHandNote = &ParsedMap[i]
		} else if ParsedMap[i].Hand == MapFormat.RightHand {
			if lastRightHandNote != nil {
				lastRightHandNote.NoteAfter = &ParsedMap[i]
			}
			ParsedMap[i].NotePrevious = lastRightHandNote
			lastRightHandNote = &ParsedMap[i]
		}
	}
	return ParsedMap, Models.IrregularInfo{Irregular: Models.RegularTypeRegular}
}
