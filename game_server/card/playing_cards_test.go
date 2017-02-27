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
	assert.Equal(t, playingCards.CardsInHand.Len(), 9)
	assert.Equal(t, playingCards.AlreadyChiCards.SameAs(group), true)

	groupTong := NewCards()
	groupTong.AppendCard(&Card{CardType:CardType_Small, CardNo:6})
	groupTong.AppendCard(&Card{CardType:CardType_Small, CardNo:7})
	groupTong.AppendCard(&Card{CardType:CardType_Small, CardNo:8})
	chiTongErr := playingCards.Chi(&Card{CardType:CardType_Small, CardNo:8}, groupTong)

	assert.Equal(t, chiTongErr, false)
	assert.Equal(t, playingCards.CardsInHand.Len(), 9)

	chiTongOk := playingCards.Chi(&Card{CardType:CardType_Small, CardNo:7}, groupTong)
	assert.Equal(t, chiTongOk, true)
	assert.Equal(t, playingCards.CardsInHand.Len(), 7)
	group.AppendCards(groupTong)
	assert.Equal(t, playingCards.AlreadyChiCards.SameAs(group), true)

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
	assert.Equal(t, playingCards.CardsInHand.Len(), 3)
	assert.Equal(t, playingCards.AlreadyPengCards.Len(), 1)

	pengTong := playingCards.Peng(&Card{CardType:CardType_Small, CardNo:6})
	assert.Equal(t, pengTong, true)
	assert.Equal(t, playingCards.CardsInHand.Len(), 1)
	assert.Equal(t, playingCards.AlreadyPengCards.Len(), 2)

	pengJian := playingCards.Peng(&Card{CardType:CardType_Small, CardNo:1})

	assert.Equal(t, pengJian, false)
	//t.Log(playingCards.ToString())
}


func TestPlayingCards_Pao(t *testing.T) {
	playingCards := NewPlayingCards()
	playingCards.AlreadySaoCards.AddAndSort(&Card{CardType:CardType_Big, CardNo:1})
	playingCards.AlreadySaoCards.AddAndSort(&Card{CardType:CardType_Big, CardNo:1})
	playingCards.AlreadySaoCards.AddAndSort(&Card{CardType:CardType_Big, CardNo:1})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})

	pao1 := playingCards.Pao(&Card{CardType:CardType_Big, CardNo:1})
	assert.Equal(t, pao1, true)
	assert.Equal(t, playingCards.CardsInHand.Len(), 2)
	assert.Equal(t, playingCards.AlreadyPaoCards.Len(), 1)

	pao6 := playingCards.Pao(&Card{CardType:CardType_Small, CardNo:6})
	assert.Equal(t, pao6, false)
	assert.Equal(t, playingCards.CardsInHand.Len(), 2)
	assert.Equal(t, playingCards.AlreadyPaoCards.Len(), 1)

	peng := playingCards.Peng(&Card{CardType:CardType_Small, CardNo:1})

	assert.Equal(t, peng, false)
	//t.Log(playingCards.ToString())
}

func TestPlayingCards_CanTiLong(t *testing.T) {
	playingCards := NewPlayingCards()

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	assert.Equal(t, playingCards.CanTiLong(&Card{CardType:CardType_Small, CardNo:6}), true)
	//t.Log(playingCards.ToString())
}

func TestPlayingCards_ComputePao(t *testing.T) {
	playingCards := NewPlayingCards()

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:7})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:7})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.ComputeSao()
	assert.Equal(t, playingCards.CardsInHand.Len(), 2)
	assert.Equal(t, playingCards.AlreadySaoCards.Len(), 2)
	assert.Equal(t, playingCards.AlreadySaoCards.At(0).SameAs(&Card{CardType:CardType_Small, CardNo:6}), true)
	assert.Equal(t, playingCards.AlreadySaoCards.At(1).SameAs(&Card{CardType:CardType_Small, CardNo:8}), true)
}

func TestPlayingCards_ComputeTiLong(t *testing.T) {
	playingCards := NewPlayingCards()

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:7})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:7})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.ComputeTiLong()
	assert.Equal(t, playingCards.CardsInHand.Len(), 2)
	assert.Equal(t, playingCards.AlreadyTiLongCards.Len(), 2)
	assert.Equal(t, playingCards.AlreadyTiLongCards.At(0).SameAs(&Card{CardType:CardType_Small, CardNo:6}), true)
	assert.Equal(t, playingCards.AlreadyTiLongCards.At(1).SameAs(&Card{CardType:CardType_Small, CardNo:8}), true)
}

func TestPlayingCards_IsHu(t *testing.T) {
	playingCards := NewPlayingCards()
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:7})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:7})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})

	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:7})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:9})

	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:4})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:4})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:4})


	playingCards.ComputeTiLong()
	playingCards.ComputeSao()

	whatCard := &Card{CardType:CardType_Small, CardNo:10}
	result := playingCards.IsHuThisCard(whatCard)
	t.Log(playingCards.CardsInHand, whatCard)
	t.Log(result)
	assert.Equal(t, result, true)
}

func TestPlayingCards_IsHu2(t *testing.T) {
	playingCards := NewPlayingCards()
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:7})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:7})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})

	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:7})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:9})

	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:4})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:4})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:4})


	playingCards.ComputeTiLong()
	playingCards.ComputeSao()

	whatCard := &Card{CardType:CardType_Small, CardNo:7}
	result := playingCards.IsHuThisCard(whatCard)
	t.Log(playingCards.CardsInHand, whatCard)
	t.Log(result)
	assert.Equal(t, result, true)
}

func TestPlayingCards_TestCard(t *testing.T) {
	playingCards := NewPlayingCards()
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:6})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:7})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:7})

	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:8})

	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:7})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:8})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:9})

	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:4})
	playingCards.AddCard(&Card{CardType:CardType_Small, CardNo:4})
	playingCards.AddCard(&Card{CardType:CardType_Big, CardNo:4})


	playingCards.ComputeTiLong()
	playingCards.ComputeSao()

	whatCard := &Card{CardType:CardType_Small, CardNo:7}
	result := playingCards.IsHuThisCard(whatCard)
	t.Log(playingCards.CardsInHand, whatCard)
	t.Log(result)
	assert.Equal(t, result, true)
	assert.Equal(t, playingCards.AlreadyTiLongCards.Len(), 1)
	assert.Equal(t, playingCards.AlreadySaoCards.Len(), 1)
}