package actions

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/raphael-p/datashard/database"
)

var days = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
}

// regex for deadline expressed relative to today (e.g "1d 2w" or "1d")
var reRelativeDeadline = regexp.MustCompile(`^(\b(\d+)d)?\s?(\b(\d+)w)?$`)

func interpretDatePrompt(datePrompt string) (time.Time, error) {
	var parsedDate time.Time
	var err error

	// try to interpret in yyyy-mm-dd format
	if datePrompt != "" {
		parsedDate, err = time.Parse(time.DateOnly, datePrompt)
		if err == nil {
			return parsedDate, nil
		}
	}

	datePrompt = strings.ToLower(datePrompt)
	currentDate := time.Now()

	// check for "today" or "tomorrow"
	if datePrompt == "today" {
		return currentDate, nil
	} else if datePrompt == "tomorrow" {
		dueDate := currentDate.AddDate(0, 0, 1)
		return dueDate, nil
	}

	// check for weekday
	dueWeekday, ok := days[datePrompt]
	if ok {
		currentWeekday := currentDate.Weekday()
		diff := (int(dueWeekday) - int(currentWeekday) + 7) % 7
		if diff == 0 {
			diff = 7
		}

		dueDate := currentDate.AddDate(0, 0, diff)
		return dueDate, nil
	}

	// check for days or weeks from today
	var diff uint64
	diffMatch := reRelativeDeadline.FindStringSubmatch(datePrompt)
	if len(diffMatch) == 5 {
		daysFromToday, _ := strconv.ParseUint(diffMatch[2], 10, 64)
		diff += daysFromToday

		weeksFromToday, _ := strconv.ParseUint(diffMatch[4], 10, 64)
		diff += (weeksFromToday * 7)
	}
	if diff > 0 {
		dueDate := currentDate.AddDate(0, 0, int(diff))
		return dueDate, nil
	}

	err = fmt.Errorf("could not determine a date from '%s'", datePrompt)
	return parsedDate, err
}

// AddNote inserts a new note and its metadata into the database
func AddNote(title, content string, importance int, dueDatePrompt string) {
	tx, err := database.DB.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	// Insert into notes table
	insertNote := `INSERT INTO notes (title, content, created_at, updated_at) VALUES (?, ?, ?, ?)`
	createdAt := time.Now()
	updatedAt := createdAt

	res, err := database.ExecWithLazyInit(tx, insertNote, title, content, createdAt, updatedAt)
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to insert into notes: %v", err)
	}

	noteID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to retrieve last insert ID: %v", err)
	}

	// Insert into meta table
	var dueDate any
	if dueDatePrompt != "" {
		dueDateParsed, err := interpretDatePrompt(dueDatePrompt)
		if err != nil {
			tx.Rollback()
			log.Fatalf("unvalid due date format, use YYYY-MM-DD")
		}
		dueDate = dueDateParsed.Format(time.DateOnly)
	}

	insertMeta := `INSERT INTO meta (note_id, importance, due_date) VALUES (?, ?, ?)`
	_, err = tx.Exec(insertMeta, noteID, importance, dueDate)
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to insert into meta: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}

	fmt.Printf("note '%s' added successfully with ID %d\n", title, noteID)
}
