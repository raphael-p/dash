package database

import (
	"math/rand/v2"
	"strings"
	"time"

	"github.com/raphael-p/datashard/pkg/logger"
)

const letterSet = "abcd3fgh1jklmn0pqrstuvwxyz"
const minWordLength = 2
const maxWordLength = 10

func FillWithSampleData(randomEntryCount int) {
	Wipe()
	err := fillTasksWithSampleData(randomEntryCount)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func fillTasksWithSampleData(randomEntryCount int) error {
	insertTask := `
	INSERT INTO tasks (name, description, created_at, updated_at, time_spent_seconds) 
	VALUES (?, ?, ?, ?, ?)
	`

	createdTime0 := time.Date(2015, time.September, 21, 19, 53, 20, 0, time.UTC)
	updatedTime0 := createdTime0.Add(time.Second * 124590)
	timeSpent0 := 3600
	_, err := DB.Exec(insertTask, "datajack override", "bypass corporate firewalls and extract encrypted payloads", createdTime0, updatedTime0, timeSpent0)
	if err != nil {
		return err
	}

	createdTime1 := time.Date(2016, time.April, 6, 9, 56, 2, 0, time.UTC)
	updatedTime1 := createdTime1.Add(-time.Second * 124590)
	timeSpent1 := 829
	_, err = DB.Exec(insertTask, "neon trace", "track digital footprints across the city's undernet", createdTime1, updatedTime1, timeSpent1)
	if err != nil {
		return err
	}

	createdTime2 := time.Date(2055, time.January, 1, 1, 1, 1, 0, time.UTC)
	updatedTime2 := createdTime2.Add(time.Second * 1)
	timeSpent2 := 451
	_, err = DB.Exec(insertTask, "synthwave uplink", "Establish a covert connection to a rogue satellite", createdTime2, updatedTime2, timeSpent2)
	if err != nil {
		return err
	}

	createdTime3 := time.Date(2015, time.September, 21, 19, 53, 24, 0, time.UTC)
	updatedTime3 := createdTime3.Add(time.Second * 43542)
	timeSpent3 := 304
	_, err = DB.Exec(insertTask, "chrome smuggler", "transport illegal cyberware through police scanners", createdTime3, updatedTime3, timeSpent3)
	if err != nil {
		return err
	}

	createdTime4 := time.Date(2015, time.September, 21, 19, 54, 26, 0, time.UTC)
	updatedTime4 := createdTime4.Add(time.Second * 132)
	timeSpent4 := 186
	_, err = DB.Exec(insertTask, "ghostline breach", "penetrate abandoned server farms for lost archives", createdTime4, updatedTime4, timeSpent4)
	if err != nil {
		return err
	}

	createdTime5 := time.Date(2015, time.September, 21, 19, 55, 29, 0, time.UTC)
	updatedTime5 := createdTime5.Add(time.Second * 84143)
	timeSpent5 := 60
	_, err = DB.Exec(insertTask, "black ICE hunter", "neutralize hostile AI defenses in corporate systems", createdTime5, updatedTime5, timeSpent5)
	if err != nil {
		return err
	}

	createdTime6 := time.Date(2015, time.September, 21, 20, 06, 15, 0, time.UTC)
	updatedTime6 := createdTime6.Add(time.Second * 31431)
	_, err = DB.Exec(insertTask, "memory forge", "fabricate false identity packets for underground networks", createdTime6, updatedTime6, 0)
	if err != nil {
		return err
	}

	createdTime7 := time.Date(2015, time.September, 21, 23, 43, 47, 0, time.UTC)
	updatedTime7 := createdTime7.Add(time.Second * 134175)
	_, err = DB.Exec(insertTask, "holojack hijack", "seize control of public hologram billboards for rebel propaganda", createdTime7, updatedTime7, 0)
	if err != nil {
		return err
	}

	insertTaskNoTimeSpent := `
	INSERT INTO tasks (name, description, created_at, updated_at) 
	VALUES (?, ?, ?, ?)
	`
	createdTime8 := time.Date(2015, time.September, 22, 3, 25, 30, 0, time.UTC)
	updatedTime8 := createdTime8.Add(time.Second * 6243)
	_, err = DB.Exec(insertTaskNoTimeSpent, "byte runner", "deliver volatile data drives through surveillance zones", createdTime8, updatedTime8)
	if err != nil {
		return err
	}

	createdTime9 := time.Date(2015, time.September, 22, 13, 11, 11, 0, time.UTC)
	updatedTime9 := createdTime9.Add(time.Second * 42142)
	_, err = DB.Exec(insertTaskNoTimeSpent, "circuit prophet", "decode predictive algorithms from seized tech cult servers", createdTime9, updatedTime9)
	if err != nil {
		return err
	}

	rng := rand.New(rand.NewChaCha8([32]byte{1}))
	for range randomEntryCount {
		_, err = DB.Exec(insertTask, randomTaskGenerator(rng)...)
		if err != nil {
			return err
		}
	}

	return nil
}

func randomTaskGenerator(rng *rand.Rand) []any {
	createdTime := time.Unix(rng.Int64N(time.Now().Unix()), 0)
	return []any{
		generateWords(rng, 2, 3),
		generateWords(rng, 4, 10),
		createdTime,
		time.Unix(rng.Int64N(time.Now().Unix()-createdTime.Unix())+createdTime.Unix(), 0),
		rng.Uint64N(2000),
	}
}

func generateWords(rng *rand.Rand, minWords, maxWords uint) string {
	wordsCount := int(rng.UintN(maxWords-minWords+1) + minWords)

	words := make([]string, wordsCount)
	for i := 0; i < wordsCount; i++ {
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
