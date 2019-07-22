package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Value struct {
	Timestamp time.Time
	Value     float64
}

func averageOfChan(in chan Value) float64 {
	var sum float64
	var count int
	for v := range in {
		sum += v.Value
		count++
	}
	return sum / float64(count)
}

func averageOfSlice(in []Value) float64 {
	var sum float64
	var count int
	for _, v := range in {
		sum += v.Value
		count++
	}
	return sum / float64(count)
}

func main() {
	// Create a large array of random numbers
	input := make([]Value, 1e7)
	for i := 0; i < 1e7; i++ {
		input[i] = Value{time.Unix(int64(i), 0), rand.Float64()}
	}

	for i := 0; i < 20; i++ {
		func() {
			st := time.Now()
			in := make(chan Value, 1000)
			go func() {
				defer close(in)
				for _, v := range input {
					in <- v
				}
			}()
			averageOfChan(in)
			fmt.Println("Channel version took", time.Since(st))
		}()
	}

	for i := 0; i < 20; i++ {
		func() {
			st := time.Now()
			averageOfSlice(input)
			fmt.Println("Slice version took", time.Since(st))
		}()
	}
}
