package card

//判断是否两个相同的牌: 包括AA和aa
func IsAA(card1, card2 *Card) bool {
	return card1.SameAs(card2)
}

//判断是否两个大小不同，点数一样的牌，类似于：一壹
func IsAa(card1, card2 *Card) bool {
	if card1.SameCardTypeAs(card2) {
		return false
	}
	return card1.SameCardNoAs(card2)
}

//判断是否三个相同的牌：包括AAA和aaa
func IsAAA(card1, card2, card3 *Card) bool {
	return IsAA(card1, card2) && IsAA(card2, card3)
}

//判断是否两个一样的牌且另一个与这两个大小不同，但是点数一样的牌
//类似于：一壹壹、壹壹一、一一壹、壹一一
func IsAAa(card1, card2, card3 *Card) bool {
	if IsAA(card1, card2) && IsAa(card2, card3) {
		return true
	}
	if IsAa(card1, card2) && IsAA(card2, card3) {
		return true
	}
	return false
}

//检查三张牌是不是顺子牌
//包括ABC和特殊的二七十，贰柒拾
func IsABC(card1, card2, card3 *Card) bool {
	if card1.PrevAt(card2) && card2.PrevAt(card3) {
		return true
	}

	if card1.SameCardTypeAs(card2) && card2.SameCardTypeAs(card3) {
		if card1.CardNo == 2 && card2.CardNo == 7 && card3.CardNo == 10 {
			return true
		}
	}
	return false
}

//判断四个牌是否都一样，包括AAAA和aaaa
func IsAAAA(card1, card2, card3, card4 *Card) bool {
	return IsAA(card1, card2) && IsAAA(card2, card3, card4)
}

//判断4个牌是否类似于： 一一壹壹、壹壹一一
func IsAAaa(card1, card2, card3, card4 *Card) bool {
	if !card1.SameAs(card2) {
		return false
	}
	if !card3.SameAs(card4) {
		return false
	}
	if !card1.SameCardNoAs(card3) {
		return false
	}
	return true
}

//是否6张牌是3连对
func IsAABBCC(card1, card2, card3, card4, card5, card6 *Card) bool {
	return IsABC(card1, card3, card5) &&
		IsAA(card1, card2) && IsAA(card3, card4) && IsAA(card5, card6)
}

//6连对
func IsAABBCCDDEEFF(card1, card2, card3, card4, card5, card6,
card7, card8, card9, card10, card11, card12 *Card) bool {
	return IsAA(card1, card2) && IsAA(card3, card4) && IsAA(card5, card6) &&
	IsAA(card7, card8) && IsAA(card9, card10) && IsAA(card11, card12) &&
	IsABC(card1, card3, card5) && IsABC(card5, card7, card9) && IsABC(card7, card9, card11)
}

//3张牌是否OK, ABC/AAA格式为OK
func Is3CardsOk(cards ...*Card) bool {
	if len(cards) != 3 {
		return false
	}

	if IsAAA(cards[0], cards[1], cards[2]) {
		return true
	}

	if IsABC(cards[0], cards[1], cards[2]) {
		return true
	}
	if IsAAa(cards[0], cards[1], cards[2]) {
		return true
	}
	return false
}

//胡2张牌
func Is2CardsOk(cards ...*Card) bool {
	if len(cards) != 2 {
		return false
	}
	return IsAA(cards[0], cards[1])
}
