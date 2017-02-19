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
		cards.AddAndSort(pool.PopFront())
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


func TestCards_Is5Card(t *testing.T) {
	cards := NewCards()
	cards.AddAndSort(&Card{CardType:CardType_Big, CardNo:1})
	cards.AddAndSort(&Card{CardType:CardType_Big, CardNo:2})
	cards.AddAndSort(&Card{CardType:CardType_Big, CardNo:3})
	if cards.IsOkWithJiang() {
		t.Fatal("it should not be ok")
	}

	cards.AddAndSort(&Card{CardType:CardType_Small, CardNo:1})
	cards.AddAndSort(&Card{CardType:CardType_Small, CardNo:1})
	if !cards.IsOkWithJiang() {
		t.Fatal("it should be ok", cards)
	}
}

func TestCards_IsOkWithJiang(t *testing.T) {
	hu5 := &Cards{
		data: []*Card{
			&Card{CardType:CardType_Big, CardNo:2},
			&Card{CardType:CardType_Big, CardNo:3},

			&Card{CardType:CardType_Big, CardNo:4},
			&Card{CardType:CardType_Big, CardNo:3},
			&Card{CardType:CardType_Small, CardNo:3},
		},
	}

	hu8 := &Cards{
		data: []*Card{
			&Card{CardType:CardType_Big, CardNo:1},
			&Card{CardType:CardType_Big, CardNo:2},
			&Card{CardType:CardType_Big, CardNo:3},

			&Card{CardType:CardType_Big, CardNo:5},
			&Card{CardType:CardType_Big, CardNo:6},
			&Card{CardType:CardType_Big, CardNo:7},

			&Card{CardType:CardType_Big, CardNo:4},
			&Card{CardType:CardType_Big, CardNo:4},
		},
	}

	hu11 := &Cards{
		data: []*Card{
			&Card{CardType:CardType_Small, CardNo:1},

			&Card{CardType:CardType_Big, CardNo:2},
			&Card{CardType:CardType_Big, CardNo:2},
			&Card{CardType:CardType_Small, CardNo:2},
			&Card{CardType:CardType_Small, CardNo:2},

			&Card{CardType:CardType_Big, CardNo:3},
			&Card{CardType:CardType_Big, CardNo:3},
			&Card{CardType:CardType_Small, CardNo:3},

			&Card{CardType:CardType_Big, CardNo:4},
			&Card{CardType:CardType_Big, CardNo:4},
			&Card{CardType:CardType_Small, CardNo:4},
		},
	}

	//12222333344445
	hu14 := &Cards{
		data: []*Card{
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
		},
	}


	hu5.Sort()
	hu8.Sort()
	hu11.Sort()
	hu14.Sort()

	if  !hu5.IsOkWithJiang() {
		t.Fatal("hu5 should ok", hu5)
	}

	if !hu8.IsOkWithJiang() {
		t.Fatal("hu8 should ok", hu8)
	}
	if !hu11.IsOkWithJiang(){
		t.Fatal("hu11 should ok", hu11)
	}
	if !hu14.IsOkWithJiang(){
		t.Fatal("hu14 should ok", hu14)
	}
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


func TestCards_IsOkWithoutJiang(t *testing.T) {
	hu3 := &Cards{
		data: []*Card{
			&Card{CardType:CardType_Big, CardNo:3},
			&Card{CardType:CardType_Big, CardNo:3},
			&Card{CardType:CardType_Small, CardNo:3},
		},
	}
	if  !hu3.IsOkWithoutJiang() {
		t.Fatal("hu3 should ok", hu3)
	}

	hu6 := &Cards{
		data: []*Card{
			&Card{CardType:CardType_Big, CardNo:2},
			&Card{CardType:CardType_Big, CardNo:3},
			&Card{CardType:CardType_Big, CardNo:3},
			&Card{CardType:CardType_Small, CardNo:3},
			&Card{CardType:CardType_Small, CardNo:3},
			&Card{CardType:CardType_Big, CardNo:4},
		},
	}
	hu6.Sort()
	if  !hu6.IsOkWithoutJiang() {
		t.Fatal("hu5 should ok", hu6)
	}


	hu9 := &Cards{
		data: []*Card{
			&Card{CardType:CardType_Big, CardNo:1},
			&Card{CardType:CardType_Big, CardNo:2},
			&Card{CardType:CardType_Big, CardNo:3},

			&Card{CardType:CardType_Big, CardNo:5},
			&Card{CardType:CardType_Big, CardNo:6},
			&Card{CardType:CardType_Big, CardNo:7},

			&Card{CardType:CardType_Big, CardNo:4},
			&Card{CardType:CardType_Big, CardNo:4},
			&Card{CardType:CardType_Small, CardNo:4},
		},
	}
	hu9.Sort()
	if !hu9.IsOkWithoutJiang() {
		t.Fatal("hu9 should ok", hu9)
	}


	hu12 := &Cards{
		data: []*Card{
			&Card{CardType:CardType_Small, CardNo:1},

			&Card{CardType:CardType_Big, CardNo:2},
			&Card{CardType:CardType_Big, CardNo:2},
			&Card{CardType:CardType_Small, CardNo:2},
			&Card{CardType:CardType_Small, CardNo:2},

			&Card{CardType:CardType_Big, CardNo:3},
			&Card{CardType:CardType_Small, CardNo:3},

			&Card{CardType:CardType_Big, CardNo:4},
			&Card{CardType:CardType_Big, CardNo:4},
			&Card{CardType:CardType_Small, CardNo:4},
			&Card{CardType:CardType_Small, CardNo:4},

			&Card{CardType:CardType_Big, CardNo:5},
		},
	}
	hu12.Sort()
	if !hu12.IsOkWithoutJiang(){
		t.Fatal("hu12 should ok", hu12)
	}

	hu15 := &Cards{
		data: []*Card{
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
		},
	}
	hu15.Sort()
	if !hu15.IsOkWithoutJiang(){
		t.Fatal("hu15 should ok", hu15)
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
	)

	t.Log(cards.computeChiGroup(&Card{CardType:CardType_Big, CardNo:10},))
	t.Log(cards.computeChiGroup(&Card{CardType:CardType_Big, CardNo:3},))
	t.Log(cards.computeChiGroup(&Card{CardType:CardType_Big, CardNo:2},))
	t.Log(cards.computeChiGroup(&Card{CardType:CardType_Small, CardNo:9},))
}