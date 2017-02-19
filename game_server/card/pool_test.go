package card

import (
	"testing"
//	"time"
	"github.com/bmizerany/assert"
)

func TestPool(t *testing.T) {
	pool := NewPool()

	pool.ReGenerate()

	t.Log(pool.cards)

	//t.Log(pool.GetCardNum())
	beforeGet := NewCards()
	beforeGet.AppendCards(pool.cards)
	newCards := NewCards()
	for{
		card := pool.PopFront()
		if card == nil {
			break
		}
		t.Log(card, pool.GetCardNum())
		newCards.AppendCard(card)
		//time.Sleep(time.Second)
	}
	if pool.GetCardNum() != 0 {
		t.Fatal("card num of pool should be 0")
	}

	assert.Equal(t, newCards.SameAs(beforeGet), true)
}
