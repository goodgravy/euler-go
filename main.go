package main

import (
	"fmt"
)

type ChanWithHead struct {
	head       uint
	headOk     bool
	shouldTake bool
	ch         chan uint
}

func (c *ChanWithHead) take() {
	if c.shouldTake {
		c.head, c.headOk = <-c.ch
	}
}

func generateNums(nums chan uint) {
	var current uint = 1
	for {
		nums <- current
		current++
	}
}

func selectMultiples(factor uint, nums chan uint, output chan uint) {
	for num := range nums {
		if num%factor == 0 {
			output <- num
		}
	}
}

func Dedupe(output, leftChan, rightChan chan uint) {
	left := &ChanWithHead{head: 0, headOk: false, shouldTake: true, ch: leftChan}
	right := &ChanWithHead{head: 0, headOk: false, shouldTake: true, ch: rightChan}

	for {
		left.take()
		right.take()

		// both sides expended
		if !left.headOk && !right.headOk {
			close(output)
			return
		}

		// one side or the other expended
		if !left.headOk {
			output <- right.head
			right.shouldTake = true
			continue
		} else if !right.headOk {
			output <- left.head
			left.shouldTake = true
			continue
		}

		// got two values
		if left.head < right.head {
			output <- left.head
			left.shouldTake = true
			right.shouldTake = false
		} else if left.head > right.head {
			output <- right.head
			left.shouldTake = false
			right.shouldTake = true
		} else {
			output <- left.head
			left.shouldTake = true
			right.shouldTake = true
		}
	}
}

func CapAt(limit uint, output, input chan uint) {
	for val := range input {
		if val < limit {
			output <- val
		} else {
			close(output)
			return
		}
	}
}

func main() {
	numsA := make(chan uint)
	numsB := make(chan uint)
	go generateNums(numsA)
	go generateNums(numsB)

	threes := make(chan uint)
	fives := make(chan uint)
	go selectMultiples(3, numsA, threes)
	go selectMultiples(5, numsB, fives)

	deduped := make(chan uint)
	go Dedupe(deduped, threes, fives)

	capped := make(chan uint)
	go CapAt(1000, capped, deduped)

	result := uint(0)
	for number := range capped {
		result += number
	}

	fmt.Println(result)
}
