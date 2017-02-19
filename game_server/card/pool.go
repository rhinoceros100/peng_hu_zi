package card

import (
	"peng_hu_zi/util"
)

type Pool struct {
	cards *Cards
}

func NewPool() *Pool {
	pool := &Pool{
		cards:	NewCards(),
	}
	return pool
}

func (pool *Pool) generate() {
	for cardNo := 1 ; cardNo <= 10; cardNo ++ {
		for num := 0; num < 4; num ++ {
			pool.cards.AppendCard(&Card{
				CardType_Big,
				cardNo,
			})

			pool.cards.AppendCard(&Card{
				CardType_Small,
				cardNo,
			})
		}
	}
}

func (pool *Pool) ReGenerate() {
	pool.cards.Clear()
	pool.generate()
	pool.shuffle()
}

//洗牌，打乱牌
func (pool *Pool) shuffle() {
	length := pool.cards.Len()
	for cnt := 0; cnt<length; cnt++ {
		i := util.RandomN(length)
		j := util.RandomN(length)
		pool.cards.Swap(i, j)
		//log.Debug("poll shuffle swap[", i, "=>", j, "]")
	}
}

func (pool *Pool) PopFront() *Card {
	return pool.cards.PopFront()
}

func (pool *Pool) PopTail() *Card{
	return pool.cards.PopTail()
}

func (pool *Pool) At(idx int) *Card {
	return pool.cards.At(idx)
}

func (pool *Pool) GetCardNum() int {
	return pool.cards.Len()
}