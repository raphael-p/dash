package utils

import (
	"time"

	"github.com/raphael-p/datashard/database"
	"github.com/raphael-p/datashard/logger"
)

func FillDatabaseWithSampleData() {
	database.Wipe()
	err := fillTasksWithSampleData()
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func fillTasksWithSampleData() error {
	insertTask := `
	INSERT INTO tasks (name, description, created_at, updated_at) 
	VALUES (?, ?, ?, ?)
	`

	createdTime0 := time.Date(2015, time.September, 21, 19, 53, 20, 0, time.UTC)
	updatedTime0 := createdTime0.Add(time.Second * 124590)
	_, err := database.DB.Exec(insertTask, "datajack override", "bypass corporate firewalls and extract encrypted payloads", createdTime0, updatedTime0)
	if err != nil {
		return err
	}

	createdTime1 := time.Date(2015, time.September, 21, 19, 53, 20, 0, time.UTC)
	updatedTime1 := createdTime1.Add(time.Second * 124590)
	_, err = database.DB.Exec(insertTask, "neon trace", "track digital footprints across the city's undernet", createdTime1, updatedTime1)
	if err != nil {
		return err
	}

	createdTime2 := time.Date(2015, time.September, 21, 19, 53, 20, 0, time.UTC)
	updatedTime2 := createdTime2.Add(time.Second * 124590)
	_, err = database.DB.Exec(insertTask, "synthwave uplink", "Establish a covert connection to a rogue satellite", createdTime2, updatedTime2)
	if err != nil {
		return err
	}

	createdTime3 := time.Date(2015, time.September, 21, 19, 53, 20, 0, time.UTC)
	updatedTime3 := createdTime3.Add(time.Second * 124590)
	_, err = database.DB.Exec(insertTask, "chrome smuggler", "transport illegal cyberware through police scanners", createdTime3, updatedTime3)
	if err != nil {
		return err
	}

	createdTime4 := time.Date(2015, time.September, 21, 19, 53, 20, 0, time.UTC)
	updatedTime4 := createdTime4.Add(time.Second * 124590)
	_, err = database.DB.Exec(insertTask, "ghostline breach", "penetrate abandoned server farms for lost archives", createdTime4, updatedTime4)
	if err != nil {
		return err
	}

	createdTime5 := time.Date(2015, time.September, 21, 19, 53, 20, 0, time.UTC)
	updatedTime5 := createdTime5.Add(time.Second * 124590)
	_, err = database.DB.Exec(insertTask, "black ICE hunter", "neutralize hostile AI defenses in corporate systems", createdTime5, updatedTime5)
	if err != nil {
		return err
	}

	createdTime6 := time.Date(2015, time.September, 21, 19, 53, 20, 0, time.UTC)
	updatedTime6 := createdTime6.Add(time.Second * 124590)
	_, err = database.DB.Exec(insertTask, "memory forge", "fabricate false identity packets for underground networks", createdTime6, updatedTime6)
	if err != nil {
		return err
	}

	createdTime7 := time.Date(2015, time.September, 21, 19, 53, 20, 0, time.UTC)
	updatedTime7 := createdTime7.Add(time.Second * 124590)
	_, err = database.DB.Exec(insertTask, "holojack hijack", "seize control of public hologram billboards for rebel propaganda", createdTime7, updatedTime7)
	if err != nil {
		return err
	}

	createdTime8 := time.Date(2015, time.September, 21, 19, 53, 20, 0, time.UTC)
	updatedTime8 := createdTime8.Add(time.Second * 124590)
	_, err = database.DB.Exec(insertTask, "byte runner", "deliver volatile data drives through surveillance zones", createdTime8, updatedTime8)
	if err != nil {
		return err
	}

	createdTime9 := time.Date(2015, time.September, 21, 19, 53, 20, 0, time.UTC)
	updatedTime9 := createdTime9.Add(time.Second * 124590)
	_, err = database.DB.Exec(insertTask, "circuit prophet", "decode predictive algorithms from seized tech cult servers", createdTime9, updatedTime9)
	if err != nil {
		return err
	}

	return nil
}
