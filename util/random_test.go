package util

import (
	"testing"
)

func TestRandom(t *testing.T) {
	for i := 0; i<10000; i++{
		if Random(0,0) != 0 {
			t.Fatal("should be 0")
		}
		
		if Random(0,1) != 0 {
			t.Fatal("should be 0")
		}
		min := Random(0, 1000)
		max := Random(1000, 2000)
		r := Random(min, max)
		//t.Log(min, max, r)
		if min<0 || min >1000 {
			t.Fatal("min err")
		}
		
		if max<1000 || min >2000 {
			t.Fatal("max err")
		}
		
		if r<min || r>max {
			t.Fatal("r err")
		}
	}
}


func TestRandomN(t *testing.T) {
	for i := 0; i<100000000; i++{
		if RandomN(0) != 0 {
			t.Fatal("should be 0")
		}
		r := RandomN(1)
		if r>=1 {
			t.Fatal("r err", r)
		}
		min := RandomN(1000)
		max := RandomN(2000)
		if min >=1000 {
			t.Fatal("min err", min)
		}
		
		if min >=2000 {
			t.Fatal("max err", max)
		}
	}
}
