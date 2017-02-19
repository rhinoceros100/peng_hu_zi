package card

import (
	"testing"
	"github.com/bmizerany/assert"
)

func TestCard(t *testing.T) {
	card1 := &Card{
		CardType: CardType_Small,
		CardNo: 1,
	}
	card2 := &Card{
		CardType: CardType_Big,
		CardNo: 1,
	}
	card3 := &Card{
		CardType: CardType_Small,
		CardNo: 1,
	}

	card4 := &Card{
		CardType: CardType_Big,
		CardNo: 2,
	}
	assert.Equal(t, card1.SameAs(card3), true)
	assert.Equal(t, card2.SameCardNoAs(card3), true)
	assert.Equal(t, card2.SameCardTypeAs(card4), true)
	assert.Equal(t, card2.PrevAt(card4), true)
}