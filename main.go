package main

import (
	"context"
	cryptorng "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wgWriters sync.WaitGroup
	var wgReaders sync.WaitGroup

	trialsNum := 10000000

	// Test the randomness

	fmt.Println("Statistics check")

	numChannel := make(chan uint8)
	wgReaders.Add(1)
	go tester(ctx, &wgReaders, numChannel, true)
	pseudoRandom(ctx, &wgWriters, numChannel, 1, trialsNum)

	wgWriters.Wait()
	close(numChannel)
	wgReaders.Wait()

	fmt.Println("Statistics check completed")

	// Run the trial

	fmt.Println("Three boxes trial")

	numChannel = make(chan uint8)
	wgReaders.Add(1)
	go tester(ctx, &wgReaders, numChannel, false)
	pseudoRandom(ctx, &wgWriters, numChannel, 2, trialsNum)

	wgWriters.Wait()
	close(numChannel)
	cancel() // just because I can do it!
	wgReaders.Wait()
	fmt.Println("Trial completed")
}

type cryptoSrc struct{}

func (s cryptoSrc) Seed(seed int64) {}

func (s cryptoSrc) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSrc) Uint64() (v uint64) {
	err := binary.Read(cryptorng.Reader, binary.LittleEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func pseudoRandom(ctx context.Context, wg *sync.WaitGroup, out chan<- uint8, numBits uint8, numTrials int) {
	var rngSrc cryptoSrc
	rnd := rand.New(rngSrc)
	rnd.Seed(11)

	upperLim := int(uint8(255) >> (8 - numBits))
	if upperLim == 1 {
		upperLim++
	}
	streamer := func() {
		defer func() {
			wg.Done()
		}()
		for i := 1; i < numTrials; i++ {
			random := uint8(rnd.Intn(upperLim))
			select {
			case out <- random:
			case <-ctx.Done():
				return
			}
		}
	}

	wg.Add(1)
	go streamer()
}
