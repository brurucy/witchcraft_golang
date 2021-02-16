package main

import (
	"fmt"
	"math/rand"
	"time"
	"witchcraft/src"
)

func main() {

	var nums []int

	n := 10_000_000

	fmt.Println("Rng generation started")

	for i := 1; i < n; i++ {

		nums = append(nums, rand.Intn(20_000_000))
	}

	fmt.Println("Rng generation ended")

	tlist := src.NewTeleportList()

	start := time.Now()

	for i := 1; i < n; i++ {

		tlist.Add(i)

	}

	elapsed := time.Since(start)

	fmt.Println("TList Elapsed add: ", elapsed)

	start = time.Now()

	for i := 1; i < n; i++ {

		tlist.Find(nums[i-1])

	}

	elapsed = time.Since(start)

	fmt.Println("Tlist Elapsed find: ", elapsed)

	start = time.Now()

	mapp := make(map[int]bool)

	for i := 1; i < n; i++ {

		mapp[i] = true

	}

	elapsed = time.Since(start)

	fmt.Println("Hashmap Elapsed add: ", elapsed)

	start = time.Now()

	for i := 1; i < n; i++ {

		if mapp[nums[i-1]] {

		}

	}

	elapsed = time.Since(start)

	fmt.Println("Hashmap elapsed find: ", elapsed)

}
