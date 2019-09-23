package main

import (
	"context"
	"fmt"
	"sync"
)

// Mutex is not required - just one thread uses this globals (by design)
var zeroOnes []uint64
var totalRngNumbers uint64
var guessStat []uint64
var luckyGuess, badGuess, winIfStay, winIfChange, faultIfStay, faultIfChange uint64

const printout uint64 = 1000000
const treasurePlacement uint8 = 0

func tester(ctx context.Context, wg *sync.WaitGroup, input <-chan uint8, rngTest bool) {
	defer func() {
		wg.Done()
	}()

	reset()

	for {
		select {
		case random, ok := <-input:
			if !ok {
				return
			}
			totalRngNumbers++
			if rngTest {
				statCheck(random)
			} else {
				trial(random)
			}
		case <-ctx.Done():
			return
		}
	}
}

func statCheck(randomBit uint8) {
	zeroOnes[randomBit]++
	if totalRngNumbers%printout == 0 {
		zPercent := float64(zeroOnes[0]) / float64(totalRngNumbers) * 100
		oPercent := float64(zeroOnes[1]) / float64(totalRngNumbers) * 100
		fmt.Printf("Total RNG: %d, Zeros: %d (%.3f%%), Ones: %d (%.3f%%)\n",
			totalRngNumbers,
			zeroOnes[0],
			zPercent,
			zeroOnes[1],
			oPercent)
	}
}

func trial(initialGuess uint8) {

	guessStat[initialGuess]++

	if initialGuess == treasurePlacement {
		luckyGuess++
		faultIfChange++ // change => loss
		winIfStay++     // stay => win
	} else {
		badGuess++
		winIfChange++ // change => win
		faultIfStay++ // stay => loss
	}

	if totalRngNumbers%printout == 0 {
		num := float64(totalRngNumbers)

		fmt.Printf("Total RNG: %d Lucky: %d (%.3f%%) Win if stay: %d (%.3f%%) Win if change: %d (%.3f%%) Guess on: %.3f%% %.3f%% %.3f%%\n",
			totalRngNumbers,
			luckyGuess,
			float64(luckyGuess)/num*100,
			winIfStay,
			float64(winIfStay)/num*100,
			winIfChange,
			float64(winIfChange)/num*100,
			float64(guessStat[0])/num*100,
			float64(guessStat[1])/num*100,
			float64(guessStat[2])/num*100)
	}
}

func reset() {
	zeroOnes = make([]uint64, 2, 2)
	guessStat = make([]uint64, 3, 3)
	totalRngNumbers = 0
	luckyGuess = 0
	badGuess = 0
	winIfStay = 0
	winIfChange = 0
	faultIfStay = 0
	faultIfChange = 0
}
