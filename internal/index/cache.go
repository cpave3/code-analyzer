package index

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // pure-Go SQLite driver

	"github.com/cameron/code-evaluator/internal/model"
)

// Cache is the SQLite repository index. Parse results (symbols, metrics,
// fingerprints, imports, references) are stored per (adapter, adapter
// version, blob id), so unchanged files are never reparsed across snapshots,
// commits or runs. It lives outside the evaluated repository.
type Cache struct {
	db *sql.DB
	// Hits and Misses count lookups for reporting.
	Hits, Misses int
}

// OpenCache opens (or creates) the cache in dir. An empty dir gives an
// in-memory cache that still deduplicates within the run.
func OpenCache(dir string) (*Cache, error) {
	dsn := "file::memory:?cache=shared"
	if dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
		dsn = "file:" + filepath.Join(dir, "index.sqlite") + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS parsed_file (
		adapter TEXT NOT NULL, version TEXT NOT NULL, blob TEXT NOT NULL, data BLOB NOT NULL,
		PRIMARY KEY (adapter, version, blob))`); err != nil {
		db.Close()
		return nil, fmt.Errorf("init cache: %w", err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS symbol (
		adapter TEXT NOT NULL, version TEXT NOT NULL, blob TEXT NOT NULL,
		name TEXT, kind TEXT, container TEXT, start_line INT, end_line INT,
		exact_fp TEXT, struct_fp TEXT, cyclomatic INT)`); err != nil {
		db.Close()
		return nil, err
	}
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS symbol_struct ON symbol(struct_fp)`)
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS symbol_blob ON symbol(adapter, version, blob)`)
	return &Cache{db: db}, nil
}

// Close closes the database.
func (c *Cache) Close() error { return c.db.Close() }

// Get looks up a parse result.
func (c *Cache) Get(adapter, version, blob string) (*model.ParsedFile, bool) {
	var data []byte
	err := c.db.QueryRow(`SELECT data FROM parsed_file WHERE adapter=? AND version=? AND blob=?`, adapter, version, blob).Scan(&data)
	if err != nil {
		c.Misses++
		return nil, false
	}
	var pf model.ParsedFile
	if json.Unmarshal(data, &pf) != nil {
		c.Misses++
		return nil, false
	}
	c.Hits++
	return &pf, true
}

// PutEntry is one parse result to store.
type PutEntry struct {
	Adapter, Version, Blob string
	Parsed                 *model.ParsedFile
}

// PutAll stores parse results in one transaction, including a queryable
// symbol table.
func (c *Cache) PutAll(entries []PutEntry) error {
	if len(entries) == 0 {
		return nil
	}
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	pfStmt, err := tx.Prepare(`INSERT OR REPLACE INTO parsed_file(adapter,version,blob,data) VALUES(?,?,?,?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	symDel, _ := tx.Prepare(`DELETE FROM symbol WHERE adapter=? AND version=? AND blob=?`)
	symStmt, err := tx.Prepare(`INSERT INTO symbol(adapter,version,blob,name,kind,container,start_line,end_line,exact_fp,struct_fp,cyclomatic) VALUES(?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, e := range entries {
		data, err := json.Marshal(e.Parsed)
		if err != nil {
			tx.Rollback()
			return err
		}
		if _, err := pfStmt.Exec(e.Adapter, e.Version, e.Blob, data); err != nil {
			tx.Rollback()
			return err
		}
		_, _ = symDel.Exec(e.Adapter, e.Version, e.Blob)
		for _, s := range e.Parsed.Symbols {
			cc := 0
			if s.Metrics != nil {
				cc = s.Metrics.Cyclomatic
			}
			if _, err := symStmt.Exec(e.Adapter, e.Version, e.Blob, s.Name, s.Kind, s.Container, s.StartLine, s.EndLine, s.ExactFP, s.StructFP, cc); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit()
}
