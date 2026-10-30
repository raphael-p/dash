package database

import (
	"math/rand/v2"
	"time"

	"github.com/raphael-p/dash/pkg/generate"
	"github.com/raphael-p/dash/pkg/logger"
)

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

	createdTime0 := time.Now().Add(-5 * time.Hour)
	updatedTime0 := createdTime0.Add(time.Second * 124)
	timeSpent0 := 3600
	_, err := DB.Exec(insertTask, "datajack override", "bypass corporate firewalls and extract encrypted payloads", createdTime0, updatedTime0, timeSpent0)
	if err != nil {
		return err
	}

	createdTime1 := time.Now().Add(-50 * time.Hour)
	updatedTime1 := createdTime1.Add(-time.Second * 1590)
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

	createdTime3 := time.Now().Add(-200 * time.Hour)
	updatedTime3 := createdTime3.Add(time.Second * 3542)
	timeSpent3 := 304
	_, err = DB.Exec(insertTask, "chrome smuggler", "transport illegal cyberware through police scanners", createdTime3, updatedTime3, timeSpent3)
	if err != nil {
		return err
	}

	createdTime4 := time.Now().Add(-12 * time.Minute)
	updatedTime4 := createdTime4.Add(time.Second * 132)
	timeSpent4 := 186
	_, err = DB.Exec(insertTask, "ghostline breach", "penetrate abandoned server farms for lost archives", createdTime4, updatedTime4, timeSpent4)
	if err != nil {
		return err
	}

	createdTime5 := time.Now().Add(-24 * 70 * time.Hour)
	updatedTime5 := createdTime5.Add(time.Second * 28443)
	timeSpent5 := 60
	_, err = DB.Exec(insertTask, "black ICE hunter", "neutralize hostile AI defenses in corporate systems", createdTime5, updatedTime5, timeSpent5)
	if err != nil {
		return err
	}

	createdTime6 := time.Now().Add(-24 * 200 * time.Hour)
	updatedTime6 := createdTime6.Add(time.Second * 31431)
	_, err = DB.Exec(insertTask, "memory forge", "fabricate false identity packets for underground networks", createdTime6, updatedTime6, 0)
	if err != nil {
		return err
	}

	createdTime7 := time.Now().Add(-24 * 20 * time.Hour)
	updatedTime7 := createdTime7.Add(time.Second * 4175)
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

	createdTime9 := time.Now().Add(-24 * 6 * time.Hour)
	updatedTime9 := createdTime9.Add(time.Second * 2142)
	_, err = DB.Exec(insertTaskNoTimeSpent, "circuit prophet", "decode predictive algorithms from seized tech cult servers", createdTime9, updatedTime9)
	if err != nil {
		return err
	}

	rng := rand.New(rand.NewChaCha8([32]byte{1}))
	for range randomEntryCount {
		_, err = DB.Exec(insertTask, randomTaskGenerator(rng, false)...)
		if err != nil {
			return err
		}
	}

	insertCompleted := `
	INSERT INTO tasks (name, description, created_at, updated_at, time_spent_seconds, completed_at) 
	VALUES (?, ?, ?, ?, ?, ?)
	`
	for range randomEntryCount {
		_, err = DB.Exec(insertCompleted, randomTaskGenerator(rng, true)...)
		if err != nil {
			return err
		}
	}

	insertPaused := `
	INSERT INTO tasks (name, description, created_at, updated_at, time_spent_seconds, paused_at) 
	VALUES (?, ?, ?, ?, ?, ?)
	`
	for range randomEntryCount {
		_, err = DB.Exec(insertPaused, randomTaskGenerator(rng, true)...)
		if err != nil {
			return err
		}
	}

	return nil
}

func randomTaskGenerator(rng *rand.Rand, isCompleted bool) []any {
	createdTime := generate.TimeAfter(rng, time.Unix(0, 0))
	updatedTime := generate.TimeAfter(rng, createdTime)
	args := []any{
		generate.Words(rng, 2, 3),
		generate.Words(rng, 4, 10),
		createdTime,
		updatedTime,
		rng.Uint64N(2000),
	}
	if isCompleted {
		args = append(args, updatedTime)
	}
	return args
}
