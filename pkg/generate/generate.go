package generate

import (
	"math/rand/v2"
	"strings"
	"time"
)

const letterSet = "abcd3fgh1jklmn0pqrstuvwxyz"
const minWordLength = 2
const maxWordLength = 10

func TimeAfter(rng *rand.Rand, t time.Time) time.Time {
	now := time.Now()

	var from time.Time
	switch rng.IntN(6) {
	case 0:
		from = now.Add(-5 * time.Minute)
	case 1:
		from = now.Add(-7 * 24 * time.Hour)
	case 2:
		from = now.AddDate(0, -1, 0)
	case 3:
		from = now.AddDate(0, -9, 0)
	case 4:
		from = now.AddDate(-3, 0, 0)
	case 5:
		from = t
	}

	if from.Before(t) {
		from = t
	}

	start, end := from.Unix(), now.Unix()
	if start >= end {
		return now
	}

	return time.Unix(start+rng.Int64N(end-start), 0)
}

func Words(rng *rand.Rand, minWords, maxWords uint) string {
	wordsCount := int(rng.UintN(maxWords-minWords+1) + minWords)

	words := make([]string, wordsCount)
	for i := range wordsCount {
		words[i] = generateWord(rng)
	}
	return strings.Join(words, " ")
}

func generateWord(rng *rand.Rand) string {
	wordLength := int(rng.UintN(maxWordLength-minWordLength+1) + minWordLength)

	word := make([]byte, wordLength)
	for i := range wordLength {
		word[i] = letterSet[rng.UintN(uint(len(letterSet)))]
	}
	return string(word)
}
