package card

import (
	"testing"
	"github.com/bmizerany/assert"
)

func TestSort(t *testing.T) {
	pool := NewPool()
	
	pool.ReGenerate()

	cards := NewCards()
	for i:=0; i<13; i++ {
		cards.AppendCard(pool.PopFront())
	}
	t.Log("before sort :")
	t.Log(cards, cards.Len())
	cards.Sort()
	t.Log("after sort :")
	t.Log(cards, cards.Len())

	t.Log("after sort big in front:")
	cards.Sort(BIG_CARD_IN_FRONT)
	t.Log(cards, cards.Len())

	t.Log("after random take way one card")
	card := cards.RandomTakeWayOne()
	t.Log(cards, cards.Len(), card)

	oneCards := NewCards()
	oneCards.AddAndSort(&Card{})
	oneCards.RandomTakeWayOne()
	t.Log("after random takeway one from only one card :")
	t.Log(oneCards, oneCards.Len())


	oneCards = NewCards()
	oneCards.AddAndSort(&Card{CardType:CardType_Big,CardNo:1})
	oneCards.TakeWay(&Card{CardType:CardType_Big,CardNo:1})
	t.Log("after takeway from only one card :")
	t.Log(oneCards, oneCards.Len())

	assert.Equal(t, true, true)

}

func TestCards_SameAs(t *testing.T) {
	cards1 := NewCards()
	cards2 := NewCards()
	cards1.AppendCard(&Card{CardType:CardType_Big, CardNo:1})
	cards2.AppendCard(&Card{CardType:CardType_Big, CardNo:1})
	if !cards1.SameAs(cards2) {
		t.Fatal("should be same as")
	}

	cards2.AppendCard(&Card{CardType:CardType_Big, CardNo:1})
	if cards1.SameAs(cards2) {
		t.Fatal("should not be same as")
	}
}

func TestCards_PopFront(t *testing.T) {
	card := &Card{CardType:CardType_Big, CardNo:4}
	cards := NewCards()
	cards.AppendCard(card)
	card1 := cards.PopFront()
	assert.Equal(t, card, card1)
	assert.Equal(t, cards.Len(), 0)
}

func TestCards_ComputeChiGroup(t *testing.T) {
	cards := NewCards(
		&Card{CardType:CardType_Big, CardNo:1},

		&Card{CardType:CardType_Big, CardNo:2},
		&Card{CardType:CardType_Big, CardNo:2},
		&Card{CardType:CardType_Small, CardNo:2},
		&Card{CardType:CardType_Small, CardNo:2},

		&Card{CardType:CardType_Big, CardNo:3},
		&Card{CardType:CardType_Big, CardNo:3},
		&Card{CardType:CardType_Small, CardNo:3},
		&Card{CardType:CardType_Small, CardNo:3},

		&Card{CardType:CardType_Big, CardNo:6},
		&Card{CardType:CardType_Big, CardNo:7},
		&Card{CardType:CardType_Big, CardNo:8},

		&Card{CardType:CardType_Big, CardNo:4},
		&Card{CardType:CardType_Big, CardNo:4},
		&Card{CardType:CardType_Small, CardNo:4},
		&Card{CardType:CardType_Big, CardNo:9},
		&Card{CardType:CardType_Big, CardNo:8},
		&Card{CardType:CardType_Big, CardNo:10},
	)

	t.Log(cards.computeChiGroup(&Card{CardType:CardType_Big, CardNo:10},))
	t.Log(cards.computeChiGroup(&Card{CardType:CardType_Big, CardNo:3},))
	t.Log(cards.computeChiGroup(&Card{CardType:CardType_Big, CardNo:2},))
	t.Log(cards.computeChiGroup(&Card{CardType:CardType_Small, CardNo:9},))

	t.Log(cards.computeChiGroup(&Card{CardType:CardType_Big, CardNo:9},))
}

