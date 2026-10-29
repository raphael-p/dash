package sqlitebackup

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// Backup creates a daily snapshot in a "backups" directory beside dbPath.
// Call it at CLI startup, after opening the database.
func Backup(db *sql.DB, dbPath string) error {
	dbPath, err := filepath.Abs(dbPath)
	if err != nil {
		return fmt.Errorf("resolve database path: %w", err)
	}

	dir := filepath.Join(filepath.Dir(dbPath), "backups")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create backup directory: %w", err)
	}

	prefix := filepath.Base(dbPath) + "."
	name := prefix + time.Now().Format("2006-01-02") + ".bak"

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("list backups: %w", err)
	}

	var backups []string
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		n := entry.Name()
		if !strings.HasPrefix(n, prefix) || !strings.HasSuffix(n, ".bak") {
			continue
		}
		date := strings.TrimSuffix(strings.TrimPrefix(n, prefix), ".bak")
		if _, err := time.Parse("2006-01-02", date); err != nil {
			continue
		}
		if n == name {
			return nil
		}
		backups = append(backups, n)
	}

	f, err := os.CreateTemp(dir, ".sqlitebackup-*")
	if err != nil {
		return fmt.Errorf("create temporary backup: %w", err)
	}
	tmp := f.Name()
	defer os.Remove(tmp)

	if err := f.Close(); err != nil {
		return fmt.Errorf("close temporary backup: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := db.ExecContext(ctx, "VACUUM main INTO ?", tmp); err != nil {
		return fmt.Errorf("create backup: %w", err)
	}
	if err := os.Rename(tmp, filepath.Join(dir, name)); err != nil {
		return fmt.Errorf("publish backup: %w", err)
	}

	backups = append(backups, name)
	slices.Sort(backups) // Date-based names sort chronologically.
	for len(backups) > 10 {
		if err := os.Remove(filepath.Join(dir, backups[0])); err != nil {
			return fmt.Errorf("remove old backup: %w", err)
		}
		backups = backups[1:]
	}
	return nil
}
