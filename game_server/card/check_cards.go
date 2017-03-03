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

//6张牌是否OK, 3连对/Is3CardsOk * 2/A + BBBB + C
//TODO
func Is6CardsOk(cards ...*Card) bool {
	if len(cards) != 6 {
		return false
	}
	//左右两边都是Is3CardOk
	if Is3CardsOk(cards[0:3]...) && Is3CardsOk(cards[3:6]...) {
		return true
	}

	//3连对
	if IsAABBCC(cards[0], cards[1], cards[2], cards[3], cards[4], cards[5]) {
		return true
	}

	//中间4个点数相同，但是和左右各一个组成顺子 类似于：一 二 二 贰 贰 三
	if IsAAaa(cards[1], cards[2], cards[3], cards[4]) &&
		IsABC(cards[0], cards[1], cards[5]){
		return true
	}

	//类似于：叁 四 四 肆 肆 伍
	if IsAAaa(cards[1], cards[2], cards[3], cards[4]) &&
		IsABC(cards[0], cards[4], cards[5]){
		return true
	}
	//一二二三三四  :013 245
	//一二贰三叁肆 :013 245
	//一贰二叁三肆 :025 134
	//壹贰二叁三四 ：013 245
	return false
}

//9张牌是否OK
func Is9CardsOk(cards ...*Card) bool {
	if len(cards) != 9 {
		return false
	}

	// 3 + 6
	if Is3CardsOk(cards[0:3]...) &&
		Is6CardsOk(cards[3:9]...) {
		return true
	}

	//6 + 3
	if Is6CardsOk(cards[0:6]...) &&
		Is3CardsOk(cards[6:9]...){
		return true
	}
	return false
}

//12张牌是否OK
func Is12CardsOk(cards ...*Card) bool {
	if len(cards) != 12 {
		return false
	}

	//3 + 9
	if Is3CardsOk(cards[0:3]...) &&
		Is9CardsOk(cards[3:12]...) {
		return true
	}

	//6 + 6
	if Is6CardsOk(cards[0:6]...) &&
		Is6CardsOk(cards[6:12]...) {
		return true
	}

	//9 + 3
	if Is9CardsOk(cards[0:9]...) &&
		Is3CardsOk(cards[9:12]...){
		return true
	}

	return false
}

//TODO 两个提龙，七对
//TODO 天胡和普通胡牌分开
func Is15CardsOk(cards ...*Card) bool {
	if len(cards) != 15 {
		return false
	}

	// 3 + 12
	if Is3CardsOk(cards[0:3]...) &&
		Is12CardsOk(cards[3:15]...){
		return true
	}

	// 6 + 9
	if Is6CardsOk(cards[0:6]...) &&
		Is9CardsOk(cards[6:15]...){
		return true
	}

	// 9 + 6
	if Is9CardsOk(cards[0:9]...) &&
		Is6CardsOk(cards[9:15]...){
		return true
	}

	//12 + 3
	if Is12CardsOk(cards[0:12]...) &&
		Is3CardsOk(cards[12:15]...) {
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

//胡5张牌
func Is5CardsOk(cards ...*Card) bool {
	if len(cards) != 5 {
		return false
	}

	//AA+BCD, 2 + 3
	if Is2CardsOk(cards[0:2]...) && Is3CardsOk(cards[2:5]...){
		return true
	}

	//ABC + DD, 3 + 2
	if Is3CardsOk(cards[0:3]...) && Is2CardsOk(cards[3:5]...) {
		return true
	}

	//A + Bbb + C 1
	if Is2CardsOk(cards[2], cards[3]) &&
		Is3CardsOk(cards[0], cards[1], cards[4]){
		return true
	}

	return false
}

//胡8张牌
//TODO
func Is8CardsOk(cards ...*Card) bool {
	if len(cards) != 8 {
		return false
	}
	//2 + 6
	if Is2CardsOk(cards[0:2]...) &&
		Is6CardsOk(cards[2:8]...) {
		return true
	}

	//6 + 2
	if Is6CardsOk(cards[0:6]...)&&
		Is2CardsOk(cards[6:8]...) {
		return true
	}

	//3 + 2 + 3
	if Is3CardsOk(cards[0:3]...) && Is2CardsOk(cards[3:5]...) &&
		Is3CardsOk(cards[5:8]...) {
		return true
	}
	return false
}

//胡11张牌
//TODO
func Is11CardsOk(cards ...*Card) bool {
	if len(cards) != 11 {
		return false
	}

	//最左边的两个为眼， 2 + 9
	if Is2CardsOk(cards[0:2]...) &&
		Is9CardsOk(cards[2:11]...) {
		return true
	}

	//最右边的两个为眼， 9 + 2
	if Is9CardsOk(cards[0:9]...) &&
		Is2CardsOk(cards[9:11]...) {
		return true
	}

	//中间左边两个为眼， 3 + 2 + 6
	if Is3CardsOk(cards[0:3]...) && Is2CardsOk(cards[3:5]...) &&
		Is6CardsOk(cards[5:11]...) {
		return true
	}

	//中间右边两个为眼， 6 + 2 + 3
	if Is6CardsOk(cards[0:6]...) && Is2CardsOk(cards[6:8]...) &&
		Is3CardsOk(cards[8:11]...){
		return true
	}
	return false
}

//胡14张牌
//TODO 小七对，双龙
func Is14CardsOk(cards ...*Card) bool {
	if len(cards)!= 14 {
		return false
	}

	// 2 + 12
	if Is2CardsOk(cards[0:2]...) && Is12CardsOk(cards[2:14]...) {
		return true
	}

	// 3 + 2 + 9
	if Is3CardsOk(cards[0:3]...) && Is2CardsOk(cards[3:5]...) &&
		Is9CardsOk(cards[5:14]...) {
		return true
	}

	// 6 + 2 +6
	if Is6CardsOk(cards[0:6]...) && Is2CardsOk(cards[6:8]...) &&
		Is6CardsOk(cards[8:14]...) {
		return true
	}

	// 9 + 2 + 3
	if Is9CardsOk(cards[0:9]...) && Is2CardsOk(cards[9:11]...) &&
		Is3CardsOk(cards[11:14]...) {
		return true
	}

	// 12 + 2
	if Is12CardsOk(cards[0:12]...) && Is2CardsOk(cards[12:14]...) {
		return true
	}
	return false
}