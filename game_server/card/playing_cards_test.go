package card

import (
	"testing"
	"github.com/bmizerany/assert"
)

func TestPlayingCards_Chi(t *testing.T) {
	playingCards := NewPlayingCards()
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:1})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:3})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:5})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:7})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:9})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:2})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:4})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:8})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})

	group := NewCards()
	group.AppendCard(&Card{CardType:CardType_Big, CardNo:4})
	group.AppendCard(&Card{CardType:CardType_Big, CardNo:5})
	group.AppendCard(&Card{CardType:CardType_Big, CardNo:6})
	chi := playingCards.Chi(&Card{CardType:CardType_Big, CardNo:5}, group)

	assert.Equal(t, chi, true)
	assert.Equal(t, playingCards.cardsInHand[CardType_Big].Len(), 7)
	assert.Equal(t, playingCards.cardsAlreadyChi[CardType_Big].SameAs(group), true)

	groupTong := NewCards()
	groupTong.AppendCard(&Card{CardType:CardType_Small, CardNo:6})
	groupTong.AppendCard(&Card{CardType:CardType_Small, CardNo:7})
	groupTong.AppendCard(&Card{CardType:CardType_Small, CardNo:8})
	chiTongErr := playingCards.Chi(&Card{CardType:CardType_Small, CardNo:8}, groupTong)

	assert.Equal(t, chiTongErr, false)
	assert.Equal(t, playingCards.cardsInHand[CardType_Small].Len(), 2)

	chiTongOk := playingCards.Chi(&Card{CardType:CardType_Small, CardNo:7}, groupTong)
	assert.Equal(t, chiTongOk, true)
	assert.Equal(t, playingCards.cardsInHand[CardType_Small].Len(), 0)
	assert.Equal(t, playingCards.cardsAlreadyChi[CardType_Small].SameAs(groupTong), true)

//	t.Log(playingCards.ToString())
}


func TestPlayingCards_Peng(t *testing.T) {
	playingCards := NewPlayingCards()
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:1})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:1})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:1})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})

	pengWan := playingCards.Peng(&Card{CardType:CardType_Big, CardNo:1})
	assert.Equal(t, pengWan, true)
	assert.Equal(t, playingCards.cardsInHand[CardType_Big].Len(), 1)
	assert.Equal(t, playingCards.cardsAlreadyPeng[CardType_Big].Len(), 3)

	pengTong := playingCards.Peng(&Card{CardType:CardType_Small, CardNo:6})
	assert.Equal(t, pengTong, true)
	assert.Equal(t, playingCards.cardsInHand[CardType_Small].Len(), 0)
	assert.Equal(t, playingCards.cardsAlreadyPeng[CardType_Small].Len(), 3)

	pengJian := playingCards.Peng(&Card{CardType:CardType_Small, CardNo:1})

	assert.Equal(t, pengJian, false)
	//t.Log(playingCards.ToString())
}


func TestPlayingCards_Gang(t *testing.T) {
	playingCards := NewPlayingCards()
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:1})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:1})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:1})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})

	gangWan := playingCards.Gang(&Card{CardType:CardType_Big, CardNo:1})
	assert.Equal(t, gangWan, true)
	assert.Equal(t, playingCards.cardsInHand[CardType_Big].Len(), 0)
	assert.Equal(t, playingCards.cardsAlreadyGang[CardType_Big].Len(), 4)

	gangTong := playingCards.Gang(&Card{CardType:CardType_Small, CardNo:6})
	assert.Equal(t, gangTong, false)
	assert.Equal(t, playingCards.cardsInHand[CardType_Small].Len(), 2)
	assert.Equal(t, playingCards.cardsAlreadyGang[CardType_Small].Len(), 0)

	gangJian := playingCards.Peng(&Card{CardType:CardType_Small, CardNo:1})

	assert.Equal(t, gangJian, false)
	//t.Log(playingCards.ToString())
}

func TestPlayingCards_Reset(t *testing.T) {
	playingCards := NewPlayingCards()
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:1})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:1})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:1})

	t.Log(playingCards.GetCardsInHandByType(CardType_Big))
	assert.Equal(t, playingCards.GetCardsInHandByType(CardType_Big).Len(), 3)

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:1})
	assert.Equal(t, playingCards.DropTail().String(), "一")
	playingCards.Reset()
	assert.Equal(t, playingCards.GetCardsInHandByType(CardType_Big).Len(), 0)
	t.Log(playingCards.GetCardsInHandByType(CardType_Big))
}