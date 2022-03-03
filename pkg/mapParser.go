package pkg

import (
	"fmt"
	"github.com/6QHTSK/ayachan/internal/pkg/logrus"
	"math"
)

func getCurrentLane(time float64, lastTick *ParsedNote) (lane float64) {
	if lastTick.NoteAfter == nil {
		return lastTick.Lane
	}
	return (time-lastTick.Time)/(lastTick.NoteAfter.Time-lastTick.Time)*(lastTick.NoteAfter.Lane-lastTick.Lane) + lastTick.Lane
}

func isBelongToSlide(note *ParsedNote, pos int) bool {
	return note.Type == NoteTypeSlide && note.Pos == pos
}

func checkNoteWhileSlide(note ParsedNote, lastTickNote *ParsedNote, hand HandType) (bool, IrregularInfo) {
	if hand == LeftHand { // (ii - 2)当前为左手绿条，右手其他键，其他键应满足lane > 绿条lane
		if note.Lane <= getCurrentLane(note.Time, lastTickNote) {
			return false, IrregularInfo{
				Irregular:     RegularTypeIrregular,
				IrregularInfo: fmt.Sprintf("%.2fs 处右手跨过左手击打的绿条/与绿条重叠", note.Time),
			}
		}
	} else if hand == RightHand { // (ii - 3)当前为右手绿条，左手其他键，其他键应满足lane < 绿条lane
		if note.Lane >= getCurrentLane(note.Time, lastTickNote) {
			return false, IrregularInfo{
				Irregular:     RegularTypeIrregular,
				IrregularInfo: fmt.Sprintf("%.2fs 处左手跨过右手击打的绿条/与绿条重叠", note.Time),
			}
		}
	} else { // (ii - 4) 未知的击打手
		return false, IrregularInfo{
			Irregular:     RegularTypeIrregular,
			IrregularInfo: fmt.Sprintf("%.2fs 处出现了击打手标记错误的情况，击打手标记%d", note.Time, hand),
		}
	}
	return true, IrregularInfo{}
}

func labelNoteWhileSlide(note *ParsedNote, AnotherHandAvailable *bool, AnotherHandPos *int, AnotherHand HandType) (bool, IrregularInfo) {
	if *AnotherHandAvailable { // 如果另一只手未击打绿条，直接指定即可
		if note.Type == NoteTypeSlide && note.Status == SlideStart {
			*AnotherHandAvailable = false
			*AnotherHandPos = note.Pos
		}
	} else {
		// 检查Note是否归属于绿条
		if isBelongToSlide(note, *AnotherHandPos) { // 如果归属，就正常进行
			if note.Status == SlideEnd {
				*AnotherHandAvailable = true
			}
		} else {
			return false, IrregularInfo{
				Irregular:     RegularTypeIrregular,
				IrregularInfo: fmt.Sprintf("%.2fs 处有双绿条锁手多押", note.Time),
			}
		}
	}
	note.Hand = AnotherHand
	return true, IrregularInfo{}
}

func checkIfLabeled(note ParsedNote) bool {
	return note.Hand != LeftHand && note.Hand != RightHand
}

// ParseMap 拆谱拆谱、输入的谱排序好先
func ParseMap(Map Chart) (parsedMap ParsedChart, irregularInfo IrregularInfo) {
	defer func() {
		err := recover()
		if err != nil {
			logrus.Log.Errorf("分析谱面出现异常！异常情况：%s\n", err)
			parsedMap = nil
			irregularInfo = IrregularInfo{
				Irregular:     RegularTypeUnknown,
				IrregularInfo: fmt.Sprint(err),
			}
		}
	}()

	parsedMap = Map.InitParseChart()

	mapLen := parsedMap.Len()
	// 第一部分，检查是否存在多压（不考虑绿条锁手）
	for i := 2; i < mapLen; i++ {
		if parsedMap[i].Beat == parsedMap[i-1].Beat && parsedMap[i-1].Beat == parsedMap[i-2].Beat {
			return parsedMap, IrregularInfo{
				Irregular:     RegularTypeIrregular,
				IrregularInfo: fmt.Sprintf("%.2fs 处出现了三押", parsedMap[i].Time),
			}
		} else if parsedMap[i].Beat == parsedMap[i-1].Beat {
			if parsedMap[i].Lane == parsedMap[i-1].Lane {
				return parsedMap, IrregularInfo{
					Irregular:     RegularTypeIrregular,
					IrregularInfo: fmt.Sprintf("%.2fs 处出现了Note重叠", parsedMap[i].Time),
				}
			} else {
				// 由于排序，右边的Note在左边的后面,这里标记双压带来的锁手
				parsedMap[i].Hand = RightHand
				parsedMap[i-1].Hand = LeftHand
			}
		}
	}

	// 第二部分，识别谱面中的绿条
	for i, note := range parsedMap {
		// 检出绿条头键后
		if note.Type == NoteTypeSlide && note.Status == SlideStart {
			// 该绿条还未被标记
			if note.Hand == UnknownHand {
				//1.串联整个绿条
				hand := UnknownHand
				lastTick := i
				var lastNoteBeat float64
				for j := i + 1; j < parsedMap.Len(); j++ {
					if parsedMap[j].Type == NoteTypeSlide && parsedMap[j].Pos == note.Pos {
						if parsedMap[lastTick].Beat == parsedMap[j].Beat {
							return parsedMap, IrregularInfo{
								Irregular:     RegularTypeIrregular,
								IrregularInfo: fmt.Sprintf("%.2fs 处绿条内有多个键重叠", parsedMap[j].Time),
							}
						}
						parsedMap[lastTick].NoteAfter = &parsedMap[j]
						parsedMap[j].NotePrevious = &parsedMap[lastTick]
						lastTick = j
						if parsedMap[j].Status == SlideEnd {
							lastNoteBeat = parsedMap[j].Beat // 记录下最后一个Note的信息，便于搜索退出
							break
						}
					}
				}
				// 2. 遍历持续范围内的所有键
				lastTickNote := &parsedMap[i] // lastTickNote 记录上一个绿条节点，便于推算了绿条当前位置
				AnotherHand := UnknownHand    // AnotherHand 记录另一只手是什么
				AnotherHandAvailable := true  // AnotherHandAvailable 记录另一只手是否被绿条占用
				AnotherHandPos := 0           // 另一只手被绿条占用时的绿条编号Pos
				for j := i - 1; j < parsedMap.Len(); j++ {
					// (1)无前一键或前一键不在绿条范围内，跳过
					if j < 0 || (j == i-1 && parsedMap[j].Beat != parsedMap[i].Beat) {
						continue
					}
					// (2)超出最后一个Note的范围，退出
					if parsedMap[j].Beat > lastNoteBeat {
						break
					}
					// (3)搜到属于本绿条的内容，记录上一个Note
					if isBelongToSlide(&parsedMap[j], note.Pos) {
						lastTickNote = &parsedMap[j]
						continue
					}
					// 下面确定为绿条持续时间内，其它的note
					// (4)如果上一个note不属于绿条内，则认为出现了锁手双押
					if j >= 1 && parsedMap[j].Beat == parsedMap[j-1].Beat && !isBelongToSlide(&parsedMap[j-1], note.Pos) {
						return parsedMap, IrregularInfo{
							Irregular:     RegularTypeIrregular,
							IrregularInfo: fmt.Sprintf("%.2fs 处有单绿条锁手多押", parsedMap[j].Time),
						}
					}
					// （5）第一次出现其他Note，记录下来
					if hand == UnknownHand {
						currentLane := getCurrentLane(parsedMap[j].Time, lastTickNote)
						if parsedMap[j].Lane < currentLane {
							hand = RightHand
							AnotherHand = LeftHand
						} else if parsedMap[j].Lane > currentLane {
							hand = LeftHand
							AnotherHand = RightHand
						} else {
							return parsedMap, IrregularInfo{
								Irregular:     RegularTypeIrregular,
								IrregularInfo: fmt.Sprintf("%.2fs 处有键与绿条重叠", parsedMap[j].Time),
							}
						}
					} else { // (6）出现了其他note，检查是否合理
						status, IrregularInfo := checkNoteWhileSlide(parsedMap[j], lastTickNote, hand)
						if !status {
							return parsedMap, IrregularInfo
						}
					}
					// (7)给其他键指定击打的手
					status, IrregularInfo := labelNoteWhileSlide(&parsedMap[j], &AnotherHandAvailable, &AnotherHandPos, AnotherHand)
					if !status {
						return parsedMap, IrregularInfo
					}
				}
				// 4. 给开头的绿条刚才的手分配，确定
				currentTick := &parsedMap[i]
				for currentTick != nil {
					currentTick.Hand = hand
					currentTick = currentTick.NoteAfter
				}
			} else { // 该绿条被标记，
				hand := note.Hand
				lastTick := i
				var lastNoteBeat float64
				//1.串联整个绿条
				for j := i + 1; j < parsedMap.Len(); j++ {
					if parsedMap[j].Type == NoteTypeSlide && parsedMap[j].Pos == note.Pos {
						if parsedMap[lastTick].Beat == parsedMap[j].Beat {
							return parsedMap, IrregularInfo{
								Irregular:     RegularTypeIrregular,
								IrregularInfo: fmt.Sprintf("%.2fs 处绿条内有多个键重叠", parsedMap[j].Time),
							}
						}
						parsedMap[lastTick].NoteAfter = &parsedMap[j]
						parsedMap[j].NotePrevious = &parsedMap[lastTick]
						lastTick = j
						parsedMap[j].Hand = hand
						if parsedMap[j].Status == SlideEnd {
							lastNoteBeat = parsedMap[j].Beat // 记录下最后一个Note的信息，便于搜索退出
							break
						}
					}
				}
				// 2. 遍历持续范围内的所有键
				var lastTickNote *ParsedNote // lastTickNote 记录上一个绿条节点，便于推算了绿条当前位置
				var AnotherHand HandType     // AnotherHand 记录另一只手是什么
				if hand == LeftHand {
					AnotherHand = RightHand
				} else if hand == RightHand {
					AnotherHand = LeftHand
				}
				AnotherHandAvailable := true // AnotherHandAvailable 记录另一只手是否被绿条占用
				AnotherHandPos := 0          // 另一只手被绿条占用时的绿条编号Pos
				for j := i; j < parsedMap.Len(); j++ {
					// (2)超出最后一个Note的范围，退出
					if parsedMap[j].Beat > lastNoteBeat {
						break
					}
					// (3)搜到属于本绿条的内容，记录上一个Note
					if isBelongToSlide(&parsedMap[j], note.Pos) {
						lastTickNote = &parsedMap[j]
						continue
					}
					// 下面确定为绿条持续时间内，其它的note
					status, IrregularInfo := checkNoteWhileSlide(parsedMap[j], lastTickNote, hand)
					if !status {
						return parsedMap, IrregularInfo
					}
					status, IrregularInfo = labelNoteWhileSlide(&parsedMap[j], &AnotherHandAvailable, &AnotherHandPos, AnotherHand)
					if !status {
						return parsedMap, IrregularInfo
					}
				}
			}
		}
	}

	// 第三部分，对未标记的，交互的识别和标记，注意此处的绿条、双压、绿条范围内的音符均已被标记。
	for i := range parsedMap {
		// 跳过前两个音符
		if i < 2 {
			continue
		}
		// 检查是否被标记
		if checkIfLabeled(parsedMap[i]) &&
			checkIfLabeled(parsedMap[i-1]) &&
			checkIfLabeled(parsedMap[i-2]) {

			interval1 := parsedMap[i].Time - parsedMap[i-1].Time
			interval2 := parsedMap[i-1].Time - parsedMap[i-2].Time
			// 检查是否间隔相差是否不大(10ms)且时长较短(200ms)
			if interval1 < 0.2 && interval2 < 0.2 &&
				math.Abs(interval2-interval1) < 0.01 {

				// 识别结构 left-right-left 小三角
				if parsedMap[i-1].Lane > parsedMap[i].Lane &&
					parsedMap[i-1].Lane > parsedMap[i-2].Lane {

					parsedMap[i].Hand = TryLeftHand
					parsedMap[i-1].Hand = TryRightHand
					parsedMap[i-2].Hand = TryLeftHand
				}

				// 识别结构 right-left-right 小三角
				if parsedMap[i-1].Lane > parsedMap[i].Lane &&
					parsedMap[i-1].Lane > parsedMap[i-2].Lane {

					parsedMap[i].Hand = TryRightHand
					parsedMap[i-1].Hand = TryLeftHand
					parsedMap[i-2].Hand = TryRightHand
				}
			}
		}
	}
	// 第四部分，全盘采用第三部分的拆分建议
	for i := range parsedMap {
		if parsedMap[i].Hand == TryLeftHand {
			parsedMap[i].Hand = LeftHand
		} else if parsedMap[i].Hand == TryRightHand {
			parsedMap[i].Hand = RightHand
		}
	}

	// 第五部分，对其他部分采用贪心算法拆谱。可以证明，未标记的Note均为单个独立的出现，直接根据左手、右手的最终note就行了。连接左右手的Note
	var lastLeftHandNote, lastRightHandNote *ParsedNote
	for i := range parsedMap {
		// 如果还没有指定哪个手
		if parsedMap[i].Hand == UnknownHand {
			//Case 1 那只手近用哪只手
			var GapLeft, GapRight, Interval1, Interval2 float64

			if lastLeftHandNote == nil {
				Interval1 = 100.0
			} else {
				Interval1 = math.Abs(parsedMap[i].Time - lastLeftHandNote.Time)
			}

			if lastRightHandNote == nil {
				Interval2 = 100.0
			} else {
				Interval2 = math.Abs(parsedMap[i].Time - lastRightHandNote.Time)
			}

			if Interval1 > 0.5 {
				GapLeft = math.Abs(parsedMap[i].Lane - 1.5)
			} else {
				GapLeft = math.Abs(parsedMap[i].Lane - lastLeftHandNote.Lane)
			}

			if Interval2 > 0.5 {
				GapRight = math.Abs(parsedMap[i].Lane - 4.5)
			} else {
				GapRight = math.Abs(parsedMap[i].Lane - lastRightHandNote.Lane)
			}

			if GapLeft < GapRight {
				parsedMap[i].Hand = LeftHand
			} else if GapLeft > GapRight {
				parsedMap[i].Hand = RightHand
			} else {
				// Case 2 在左边用左手，右边用右手
				if parsedMap[i].Lane < 2.9 {
					parsedMap[i].Hand = LeftHand
				} else if parsedMap[i].Lane > 3.1 {
					parsedMap[i].Hand = RightHand
				} else {
					// Case 3 那边距上个键的间隔长，用哪只手

					if Interval1 > Interval2 {
						parsedMap[i].Hand = LeftHand
					} else {
						parsedMap[i].Hand = RightHand
					}
				}
			}
		}
		if parsedMap[i].Hand == LeftHand {
			if lastLeftHandNote != nil {
				lastLeftHandNote.NoteAfter = &parsedMap[i]
			}
			parsedMap[i].NotePrevious = lastLeftHandNote
			lastLeftHandNote = &parsedMap[i]
		} else if parsedMap[i].Hand == RightHand {
			if lastRightHandNote != nil {
				lastRightHandNote.NoteAfter = &parsedMap[i]
			}
			parsedMap[i].NotePrevious = lastRightHandNote
			lastRightHandNote = &parsedMap[i]
		}
	}
	return parsedMap, IrregularInfo{Irregular: RegularTypeRegular}
}
