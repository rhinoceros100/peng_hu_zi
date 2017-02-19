package card

import (
	"peng_hu_zi/util"
	"sort"
)

type SortType int
const (
	BIG_CARD_IN_FRONT		SortType = iota	//同样点数的情况下，大牌在前面
	SMALL_CARD_IN_FRONT						//同样点数的情况下，小牌在前面
)

type Cards struct {
	data 	[]*Card
	sortType	SortType
}

//创建一个Cards对象
func NewCards() *Cards{
	return &Cards{
		data :	make([]*Card, 0),
	}
}

//从指定的cardSlice创建一个Cards对象
func NewCardsFrom(cardSlice []*Card) *Cards{
	return &Cards{
		data: cardSlice,
	}
}

//获取cards的数据
func (cards *Cards) Data() []*Card {
	return cards.data
}

//获取第idx个牌
func (cards *Cards) At(idx int) *Card {
	if idx >= cards.Len() {
		return nil
	}
	return cards.data[idx]
}

//cards的长度，牌数
func (cards *Cards) Len() int {
	return len(cards.data)
}

//比较指定索引对应的两个牌的大小
func (cards *Cards) Less(i, j int) bool {
	cardI := cards.At(i)
	cardJ := cards.At(j)
	if cardI == nil || cardJ == nil{
		return false
	}

	if cardI.CardNo != cardJ.CardNo {
		return cardI.CardNo < cardJ.CardNo
	}

	if cards.sortType == BIG_CARD_IN_FRONT {
		if cardI.CardType > cardJ.CardType {
			return true
		}
		return false
	} else if cards.sortType == SMALL_CARD_IN_FRONT {
		if cardI.CardType < cardJ.CardType {
			return true
		}
		return false
	}

	return false
}

//交换索引为，j的两个数据
func (cards *Cards) Swap(i, j int) {
	if i == j {
		return
	}
	length := cards.Len()
	if i >= length || j >= length {
		return
	}
	swap := cards.At(i)
	cards.data[i] = cards.At(j)
	cards.data[j] = swap
}

//追加一张牌
func (cards *Cards) AppendCard(card *Card) {
	if card == nil {
		return
	}
	cards.data = append(cards.data, card)
}

//增加一张牌并排序
func (cards *Cards) AddAndSort(card *Card){
	if card == nil {
		return
	}
	cards.AppendCard(card)
	cards.Sort()//default sort
}

//追加一个cards对象
func (cards *Cards) AppendCards(other *Cards) {
	cards.data = append(cards.data, other.data...)
}

//取走一张指定的牌，并返回成功或者失败
func (cards *Cards) TakeWay(drop *Card) bool {
	if drop == nil {
		return true
	}
	for idx, card := range cards.data {
		if card.SameAs(drop) {
			cards.data = append(cards.data[0:idx], cards.data[idx+1:]...)
			return true
		}
	}
	return false
}

//取走第一张牌
func (cards *Cards) PopFront() *Card {
	if cards.Len() == 0 {
		return nil
	}
	card := cards.At(0)
	cards.data = cards.data[1:]
	return card
}

//取走最后一张牌
func (cards *Cards) PopTail() *Card {
	if cards.Len() == 0 {
		return nil
	}
	card := cards.At(cards.Len()-1)
	cards.data = cards.data[:cards.Len()-1]
	return card
}

//随机取走一张牌
func (cards *Cards) RandomTakeWayOne() *Card {
	length := cards.Len()
	if length == 0 {
		return nil
	}
	idx := util.RandomN(length)
	card := cards.At(idx)
	cards.data = append(cards.data[0:idx], cards.data[idx+1:]...)
	return card
}

//清空牌
func (cards *Cards) Clear() {
	cards.data = cards.data[0:0]
}

//排序
func (cards *Cards)Sort(sortType ...SortType) {
	if len(sortType) > 0 {
		cards.sortType = sortType[0]
	} else {
		cards.sortType = SMALL_CARD_IN_FRONT
	}
	sort.Sort(cards)
}

//是否是一样的牌组
func (cards *Cards) SameAs(other *Cards) bool {
	if cards == nil || other == nil {
		return false
	}

	length := other.Len()
	if cards.Len() != length {
		return false
	}

	for idx := 0; idx < length; idx++ {
		if !cards.At(idx).SameAs(other.At(idx)) {
			return false
		}
	}
	return true
}

func (cards *Cards) IsOkWithJiang() bool  {
	if !cards.isOkWithJiang() {
		cards.Sort(BIG_CARD_IN_FRONT)
		return cards.isOkWithJiang()
	}
	return true
}

func (cards *Cards) IsOkWithoutJiang() bool  {
	if !cards.isOkWithoutJiang() {
		if cards.Len() != 12 {
			cards.Sort(BIG_CARD_IN_FRONT)
		}
		return cards.isOkWithoutJiang()
	}
	return true
}

//带将的情况下是否OK
func (cards *Cards) isOkWithJiang() bool {
	switch cards.Len() {
	case 2:
		return Is2CardsOk(cards.data...)
	case 5:
		return Is5CardsOk(cards.data...)
	case 8:
		return Is8CardsOk(cards.data...)
	case 11:
		return Is11CardsOk(cards.data...)
	case 14:
		return Is14CardsOk(cards.data...)
	}
	return false
}

//不带将的情况下是否OK
func (cards *Cards) isOkWithoutJiang() bool {
	switch cards.Len() {
	case 3:
		return Is3CardsOk(cards.data...)
	case 6:
		return Is6CardsOk(cards.data...)
	case 9:
		return Is9CardsOk(cards.data...)
	case 12:
		return Is12CardsOk(cards.data...)
	case 15:
		return Is15CardsOk(cards.data...)
	}
	return false
}

//检查是否能吃
func (cards *Cards) canChi(whatCard *Card, whatGroup *Cards) bool {
	groups := cards.computeChiGroup(whatCard)
	for _, group := range groups {
		if group.SameAs(whatGroup) {
			return true
		}
	}
	return false
}

//检查是否能碰
func (cards *Cards) canPeng(whatCard *Card) bool  {
	return cards.calcSameCardNum(whatCard) >= 2
}

//检查是否能杠
func (cards *Cards) canGang(whatCard *Card) bool {
	return cards.calcSameCardNum(whatCard) >= 3
}

//计算与指定牌一样的牌的数量
func (cards *Cards) calcSameCardNum(whatCard *Card) int {
	num := 0
	for _, card := range cards.data {
		if card.SameAs(whatCard) {
			num++
		}
	}
	return num
}

/*	计算指定的牌可以吃牌的组合
*	假设要吃的牌为C，则需要检查是否存在如下组合：
*	ABC、BCD、BCCD、BCCCD、BCCCD、CDE
*	如果存在AB,则添加组合ABC
*	如果存在BD/BCD/BCCD/BCCCD, 则添加组合BCD
*	如果存在DE,则添加组合CDE
*/
func (cards *Cards) computeChiGroup(card *Card) []*Cards {
	return nil
	/*
	length := cards.Len()
	if length < 2 {
		return nil
	}
	cardsSlice := make([]*Cards, 0)

	//检查AB/BD/DE组合
	for idx := 0; idx < length-1; idx++ {
		//if AB组合，加上card后相当于ABC
		if IsABC(cards.At(idx), cards.At(idx+1), card) {
			tmp := NewCards()
			tmp.AppendCard(cards.At(idx))
			tmp.AppendCard(cards.At(idx+1))
			tmp.AppendCard(card)
			cardsSlice = append(cardsSlice, tmp)
		}

		//if BD组合，加上card后相当于BCD
		if IsABC(cards.At(idx), card, cards.At(idx+1))  {
			tmp := NewCards()
			tmp.AppendCard(cards.At(idx))
			tmp.AppendCard(card)
			tmp.AppendCard(cards.At(idx+1))
			cardsSlice = append(cardsSlice, tmp)
		}

		//if DE组合，加上card后相当于CDE
		if IsABC(card, cards.At(idx), cards.At(idx+1)) {
			tmp := NewCards()
			tmp.AppendCard(card)
			tmp.AppendCard(cards.At(idx))
			tmp.AppendCard(cards.At(idx+1))
			cardsSlice = append(cardsSlice, tmp)
		}

		//if BCD 组合，加上card后相当于BCCD
		if IsABBC(cards.At(idx), card, cards.At(idx+1), cards.At(idx+2)) {
			tmp := NewCards()
			tmp.AppendCard(cards.At(idx))
			tmp.AppendCard(card)
			tmp.AppendCard(cards.At(idx+2))
			cardsSlice = append(cardsSlice, tmp)
		}

		//if BCCD 组合，加上card后相当于BCCCD
		if IsABBBC(cards.At(idx), card, cards.At(idx+1), cards.At(idx+2), cards.At(idx+3)) {
			tmp := NewCards()
			tmp.AppendCard(cards.At(idx))
			tmp.AppendCard(card)
			tmp.AppendCard(cards.At(idx+3))
			cardsSlice = append(cardsSlice, tmp)
		}

		//if BCCCD 组合，加上card后相当于BCCCCD
		if IsABBBBC(cards.At(idx), card, cards.At(idx+1), cards.At(idx+2), cards.At(idx+3), cards.At(idx+4)) {
			tmp := NewCards()
			tmp.AppendCard(cards.At(idx))
			tmp.AppendCard(card)
			tmp.AppendCard(cards.At(idx+4))
			cardsSlice = append(cardsSlice, tmp)
		}
	}
	return cardsSlice
	*/
}


func (cards *Cards) String() string {
	str := ""
	for _, card := range cards.data{
		str += card.String() + ","
	}
	return str
}

//把牌分成左右2份：[:idx], [idx+1:]
func (cards *Cards) Split(idx int) (left, right *Cards){
	length := cards.Len()
	if length <= idx {
		return cards, nil
	}
	left = &Cards{
		data:	cards.data[0 : length-idx],
	}
	right = &Cards{
		data:	cards.data[length-idx:],
	}
	return left, right
}

//是否所有的牌都是同一个牌
func (cards *Cards) IsAllCardSame() bool {
	length := cards.Len()
	for idx := 1; idx < length; idx++ {
		if !cards.At(0).SameAs(cards.At(idx)) {
			return false
		}
	}
	return true
}

//是否所有的牌都是同一个类型
func (cards *Cards) IsAllCardSameCardType() bool {
	length := cards.Len()
	for idx := 1; idx < length; idx++ {
		if !cards.At(0).SameCardTypeAs(cards.At(idx)) {
			return false
		}
	}
	return true
}

//是否所有的牌都是同样的点数
func (cards *Cards) IsAllCardSameCardNo() bool {
	length := cards.Len()
	for idx := 1; idx < length; idx++ {
		if !cards.At(0).SameCardNoAs(cards.At(idx)) {
			return false
		}
	}
	return true
}

//获取不同牌的类型的数量
func (cards *Cards) CalcDiffCardCnt() int {
	has := make(map[int64]bool)
	for _, card := range cards.data {
		has[card.MakeKey()] = true
	}
	return len(has)
}

func (cards *Cards) CalcCardCntAsSameCardType(cardType int) int {
	cnt := 0
	tmp := &Card{CardType:cardType}
	for _, card := range cards.data {
		if card.SameCardTypeAs(tmp) {
			cnt++
		}
	}
	return cnt
}