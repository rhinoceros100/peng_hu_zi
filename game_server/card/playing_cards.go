package card

import (
	"fmt"
)

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
	playingCards.AlreadyChiCards.Clear()
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
	return playingCards.CardsInHand.canTiLong(whatCard) || playingCards.AlreadySaoCards.hasCard(whatCard)
}

func (playingCards *PlayingCards) IsTianHu() bool {
	if playingCards.CardsInHand.Len() == 14 {
		return playingCards.Is7Dui(playingCards.CardsInHand.data...)
	} else if playingCards.CardsInHand.Len() == 15 {
		return playingCards.IsHu()
	}
	return false
}

func (playingCards *PlayingCards) IsHu() bool {
	return playingCards.IsCardsOk(playingCards.CardsInHand.data...)
}

func (playingCards *PlayingCards) IsHuThisCard(whatCard *Card) bool {
	playingCards.CardsInHand.AddAndSort(whatCard)
	ok := playingCards.IsHu()
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
			cards.TakeWay(card)
		}
	}
}

func (playingCards *PlayingCards) GetPaoAndTiLongNum() int{
	return playingCards.AlreadyPaoCards.Len() + playingCards.AlreadyTiLongCards.Len()
}

func (playingCards *PlayingCards) GetTiLongNum() int {
	return playingCards.AlreadyTiLongCards.Len()
}

func (playingCards *PlayingCards) String() string{
	return fmt.Sprintf(
		"InHand={%v}, Chi={%v}, Peng={%v}, Sao={%v}, Pao={%v}, TiLong={%v}",
		playingCards.CardsInHand,
		playingCards.AlreadyChiCards,
		playingCards.AlreadyPengCards,
		playingCards.AlreadySaoCards,
		playingCards.AlreadyPaoCards,
		playingCards.AlreadyTiLongCards,
	)
}


func (playingCards *PlayingCards) IsCardsOk(cards ...*Card) bool {
	length := len(cards)
	if length == 2 {
		//log.Debug("IsCardsOk length==2:", cards)
		return Is2CardsOk(cards...)
	}

	if length == 3 {
		//log.Debug("IsCardsOk length==3:", cards)
		return Is3CardsOk(cards...)
	}

	for i:=0; i<length-2; i++ {
		for j:=i+1; j<length-1; j++{
			for k:=j+1; k<length; k++ {
				//log.Debug("IsCardsOk Is3CardsOk :[", i, j, k, "]" , cards[i], cards[j], cards[k])
				if !Is3CardsOk(cards[i], cards[j], cards[k]) {
					continue
				}

				otherForCheckHu := make([]*Card, 0)
				for l:=0; l<length; l++ {
					if l == i || l == j || l == k {
						continue
					}
					otherForCheckHu = append(otherForCheckHu, cards[l])
				}
				//log.Debug("IsCardsOk otherForCheckHu :", otherForCheckHu)
				if playingCards.IsCardsOk(otherForCheckHu...) {
					return true
				}
			}
		}
	}
	return false
}

func (playingCards *PlayingCards) Is7Dui(cards ...*Card) bool {
	length := len(cards)
	if length == 2 {
		return Is2CardsOk(cards[0], cards[1])
	}

	for i:=0; i<length-1; i++ {
		for j:=i+1; j<length; j++{
			if !Is2CardsOk(cards[i], cards[j]) {
				continue
			}

			otherForCheckHu := make([]*Card, 0)
			for l:=0; l<length; l++ {
				if l == i || l == j {
					continue
				}
				otherForCheckHu = append(otherForCheckHu, cards[l])
			}
			//log.Debug("IsCardsOk otherForCheckHu :", playingCards.otherForCheckHu)
			if playingCards.Is7Dui(otherForCheckHu...) {
				return true
			}
		}
	}
	return false
}