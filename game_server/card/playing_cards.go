package card

type PlayingCards struct {
	CardsInHand			*Cards		//手上的牌
	AlreadyChiCards            *Cards      //已经吃的牌, 3张牌都存
	AlreadyPengCards			*Cards		//已经碰的牌，只存已碰牌的其中一张
	AlreadySaoCards			*Cards		//已经扫的牌，只存已扫牌的其中一张
	AlreadyPaoCards			*Cards		//已经跑的牌，只存已跑牌的其中一张
	AlreadyTiLongCards         *Cards      //已经提龙的拍，只存已提龙牌的其中一张
}

func NewPlayingCards() *PlayingCards {
	return  &PlayingCards{
		CardsInHand: NewCards(),
		AlreadyChiCards: NewCards(),
		AlreadyPengCards: NewCards(),
		AlreadySaoCards: NewCards(),
		AlreadyPaoCards: NewCards(),
		AlreadyTiLongCards: NewCards(),
	}
}

func (playingCards *PlayingCards) Reset() {
	playingCards.CardsInHand.Clear()
	playingCards.AlreadyPengCards.Clear()
	playingCards.AlreadySaoCards.Clear()
	playingCards.AlreadyPaoCards.Clear()
	playingCards.AlreadyTiLongCards.Clear()
}

func (playingCards *PlayingCards) AddCards(cards *Cards) {
	playingCards.CardsInHand.AppendCards(cards)
	playingCards.CardsInHand.Sort()
}


func (playingCards *PlayingCards) DropCards(cards *Cards) {
	for _, card := range cards.Data() {
		playingCards.DropCard(card)
	}
}

//增加一张牌
func (playingCards *PlayingCards) AddCard(card *Card) {
	playingCards.CardsInHand.AddAndSort(card)
}

//丢弃一张牌
func (playingCards *PlayingCards) DropCard(card *Card) bool {
	return playingCards.CardsInHand.TakeWay(card)
}

func (playingCards *PlayingCards) Tail() *Card {
	return playingCards.CardsInHand.Tail()
}

//吃牌，要吃whatCard，以及吃哪个组合whatGroup
func (playingCards *PlayingCards) Chi(whatCard *Card, whatGroup *Cards) bool {
	if !playingCards.CanChi(whatCard, whatGroup) {
		return false
	}

	for _, card := range whatGroup.Data() {//移动除了whatCard以外的card到cardsAlreadyChi
		if card.SameAs(whatCard) {
			continue
		}
		playingCards.CardsInHand.TakeWay(card)
	}

	playingCards.AlreadyChiCards.AppendCards(whatGroup)
	return true
}

//碰牌
func (playingCards *PlayingCards) Peng(whatCard *Card) bool {
	if !playingCards.CanPeng(whatCard) {
		return false
	}

	playingCards.CardsInHand.TakeWay(whatCard)
	playingCards.CardsInHand.TakeWay(whatCard)
	playingCards.AlreadyPengCards.AddAndSort(whatCard)
	return true
}

//跑牌
func (playingCards *PlayingCards) Pao(whatCard *Card) bool {
	if playingCards.AlreadySaoCards.hasCard(whatCard) {
		playingCards.AlreadySaoCards.TakeWay(whatCard)
		playingCards.AlreadyPaoCards.AddAndSort(whatCard)
		return true
	} else if playingCards.AlreadyPengCards.hasCard(whatCard) {
		playingCards.AlreadyPengCards.TakeWay(whatCard)
		playingCards.AlreadyPaoCards.AddAndSort(whatCard)
		return true
	}
	return false
}

//提龙
func (playingCards *PlayingCards) TiLong(whatCard *Card) bool {
	if !playingCards.CanTiLong(whatCard) {
		return false
	}
	for i:=0; i<4; i++ {
		playingCards.CardsInHand.TakeWay(whatCard)
	}
	playingCards.AlreadyTiLongCards.AddAndSort(whatCard)
	return true
}

//扫
func (playingCards *PlayingCards) Sao(whatCard *Card) bool {
	if !playingCards.CanSao(whatCard) {
		return false
	}
	for i:=0; i<3; i++ {
		playingCards.CardsInHand.TakeWay(whatCard)
	}
	playingCards.AlreadySaoCards.AddAndSort(whatCard)
	return true
}

/*	计算指定的牌可以吃牌的组合
*/
func (playingCards *PlayingCards) ComputeChiGroup(card *Card) []*Cards {
	return playingCards.CardsInHand.computeChiGroup(card)
}

//检查是否能吃
func (playingCards *PlayingCards) CanChi(whatCard *Card, whatGroup *Cards) bool {
	return playingCards.CardsInHand.canChi(whatCard, whatGroup)
}

//检查是否能碰
func (playingCards *PlayingCards) CanPeng(whatCard *Card) bool  {
	return playingCards.CardsInHand.canPeng(whatCard)
}

//检查是否能扫
func (playingCards *PlayingCards) CanSao(whatCard *Card) bool {
	return playingCards.CardsInHand.canSao(whatCard)
}

//检查是否能跑
func (playingCards *PlayingCards) CanPao(whatCard *Card) bool {
	if playingCards.AlreadySaoCards.hasCard(whatCard) {
		return true
	} else if playingCards.AlreadyPengCards.hasCard(whatCard) {
		return true
	}
	return false
}

//检查是否能提龙
func (playingCards *PlayingCards) CanTiLong(whatCard *Card) bool {
	return playingCards.CardsInHand.canTiLong(whatCard)
}

func (playingCards *PlayingCards) IsHu() bool {
	paoAndTLCnt := playingCards.GetPaoAndTiLongNum()
	if paoAndTLCnt >= 2 {
		ok := playingCards.CardsInHand.IsOkWithJiang()
		if ok {
			return true
		}
	}
	return false
}

func (playingCards *PlayingCards) IsHuThisCard(whatCard *Card) bool {
	playingCards.CardsInHand.AddAndSort(whatCard)
	ok := playingCards.CardsInHand.IsOkWithoutJiang()
	playingCards.CardsInHand.TakeWay(whatCard)
	return ok
}

//计算提龙
func (playingCards *PlayingCards) ComputeTiLong() {
	cards := playingCards.CardsInHand
	start := 0
	for ; start < cards.Len()-3; {
		if cards.At(start).SameAs(cards.At(start + 3)) {
			playingCards.AlreadyTiLongCards.AppendCard(cards.At(start))
			start += 3
		} else {
			start++
		}
	}

	for _, card := range playingCards.AlreadyTiLongCards.data {
		for i := 0; i < 4; i ++ {
			playingCards.CardsInHand.TakeWay(card)
		}
	}
}

//计算扫牌
func (playingCards *PlayingCards) ComputeSao() {
	cards := playingCards.CardsInHand
	start := 0
	for ; start < cards.Len()-2; {
		if cards.At(start).SameAs(cards.At(start + 2)) {
			playingCards.AlreadySaoCards.AppendCard(cards.At(start))
			start += 3
		} else {
			start++
		}
	}

	for _, card := range playingCards.AlreadySaoCards.data {
		for i := 0; i < 3; i ++ {
			playingCards.CardsInHand.TakeWay(card)
		}
	}
}

func (playingCards *PlayingCards) GetPaoAndTiLongNum() int{
	return playingCards.AlreadyPaoCards.Len() + playingCards.AlreadyTiLongCards.Len()
}