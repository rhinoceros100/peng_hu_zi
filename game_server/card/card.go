package card

type Card struct {
	CardType int 	//牌类型
	CardNo   int 	//牌编号
}

//card是否在other的前一位牌
func (card *Card) PrevAt(other *Card) bool{
	if !card.SameCardTypeAs(other) {
		return false
	}
	return card.CardNo + 1 == other.CardNo
}

//是否同一类型的牌
func (card *Card) SameCardTypeAs(other *Card) bool {
	if other == nil || card == nil {
		return false
	}
	return other.CardType == card.CardType
}

func (card *Card) SameCardNoAs(other *Card) bool {
	if other == nil || card == nil {
		return false
	}
	return other.CardNo == card.CardNo
}

func (card *Card) SameAs(other *Card) bool {
	if other == nil || card == nil {
		return false
	}
	if other.CardType != card.CardType {
		return false
	}
	if other.CardNo != card.CardNo {
		return false
	}
	return true
}

func (card *Card) CopyFrom(other *Card) {
	if other == nil || card == nil {
		return
	}
	card.CardType = other.CardType
	card.CardNo = other.CardNo
}

func (card *Card) MakeKey() int64 {
	var ret int64
	ret = int64(card.CardType ) | int64(card.CardNo << 32)
	return ret
}

func (card *Card) String() string {
	if card == nil {
		return "nil"
	}
	cardNameMap := cardNameMap()
	noNameMap, ok1 := cardNameMap[card.CardType]
	if !ok1 {
		return "unknow card type"
	}

	name, ok2 := noNameMap[card.CardNo]
	if !ok2 {
		return "unknow card no"
	}
	return name
}

func cardNameMap() map[int]map[int]string {
	return map[int]map[int]string{
		CardType_Big: {
			1: 		"壹",
			2:  	"贰",
			3:   	"叁",
			4:  	"肆",
			5:		"伍",
			6:		"陸",
			7:		"柒",
			8:		"捌",
			9:		"玖",
			10:		"拾",
		},
		CardType_Small: {
			1: 		"一",
			2:  	"二",
			3:   	"三",
			4:  	"四",
			5:		"五",
			6:		"六",
			7:		"七",
			8:		"八",
			9:		"九",
			10:		"十",
		},
	}
}
