package card

type PlayingCards struct {
	CardsInHand			*Cards		//手上的牌
	ChiCards            *Cards      //已经吃的牌, 3张牌都存
	PengCards			*Cards		//已经碰的牌，只存已碰牌的其中一张
	SaoCards			*Cards		//已经扫的牌，只存已扫牌的其中一张
	PaoCards			*Cards		//已经跑的牌，只存已跑牌的其中一张
	TiLongCards         *Cards      //已经提龙的拍，只存已提龙牌的其中一张
}

func NewPlayingCards() *PlayingCards {
	return  &PlayingCards{
		CardsInHand: NewCards(),
		ChiCards: NewCards(),
		PengCards: NewCards(),
		SaoCards: NewCards(),
		PaoCards: NewCards(),
		TiLongCards: NewCards(),
	}
}

func (playingCards *PlayingCards) Reset() {
	playingCards.CardsInHand.Clear()
	playingCards.PengCards.Clear()
	playingCards.SaoCards.Clear()
	playingCards.PaoCards.Clear()
	playingCards.TiLongCards.Clear()
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

func (playingCards *PlayingCards) DropTail() *Card {
	return playingCards.CardsInHand.PopTail()
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

	playingCards.ChiCards.AppendCards(whatGroup)
	playingCards.ChiCards.Sort()

	return true
}

//碰牌
func (playingCards *PlayingCards) Peng(whatCard *Card) bool {
	if !playingCards.CanPeng(whatCard) {
		return false
	}

	playingCards.CardsInHand.TakeWay(whatCard)
	playingCards.CardsInHand.TakeWay(whatCard)
	playingCards.PengCards.AddAndSort(whatCard)
	return true
}

//跑牌
func (playingCards *PlayingCards) Pao(whatCard *Card) bool {
	if !playingCards.CanPao(whatCard) {
		return false
	}

	playingCards.CardsInHand.TakeWay(whatCard)
	playingCards.CardsInHand.TakeWay(whatCard)
	playingCards.CardsInHand.TakeWay(whatCard)

	playingCards.PaoCards.AddAndSort(whatCard)
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
	can := playingCards.SaoCards.hasCard(whatCard)
	if can {
		return can
	}
	return playingCards.PengCards.hasCard(whatCard)
}

//检查是否能提龙
func (playingCards *PlayingCards) CanTiLong(whatCard *Card) bool {
	return playingCards.CardsInHand.canTiLong(whatCard)
}

func (playingCards *PlayingCards) IsHu() bool {
	return playingCards.CardsInHand.IsOkWithJiang()
}