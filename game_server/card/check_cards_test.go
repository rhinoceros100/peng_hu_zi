package card

import (
	"testing"
	"time"
	"github.com/bmizerany/assert"
)

func TestIsCardsOk(t *testing.T) {
	start := time.Now()
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

	hu1 := IsCardsOk(hu15.data...)
	hu2 := IsCardsOk(hu15.data...)
	hu3 := IsCardsOk(hu15.data...)
	hu4 := IsCardsOk(hu15.data...)
	assert.Equal(t, hu1, true)
	assert.Equal(t, hu2, true)
	assert.Equal(t, hu3, true)
	assert.Equal(t, hu4, true)
	t.Log(time.Now().Sub(start))
}
