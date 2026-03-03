package taskstore

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

type Session struct {
	ID        int64
	Source    string
	Ref       string
	Goals     string
	ACText    string
	Verdict   string
	StartedAt time.Time
	EndedAt   sql.NullTime
}

type Task struct {
	ID          int64
	SessionID   int64
	RepoKey     string
	WorkBranch  string
	Title       string
	Details     string
	Severity    string
	Status      string
	DedupeKey   string
	Occurrences int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastSeenAt  time.Time
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create state dir: %w", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	store := &Store{db: db}
	if err := store.migrate(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate(ctx context.Context) error {
	const schema = `
CREATE TABLE IF NOT EXISTS qa_sessions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  source TEXT NOT NULL,
  ref TEXT NOT NULL,
  goals TEXT NOT NULL,
  ac_text TEXT NOT NULL DEFAULT '',
  verdict TEXT NOT NULL DEFAULT '',
  started_at TEXT NOT NULL,
  ended_at TEXT
);

CREATE TABLE IF NOT EXISTS tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  session_id INTEGER NOT NULL,
  repo_key TEXT NOT NULL,
  title TEXT NOT NULL,
  details TEXT NOT NULL,
  severity TEXT NOT NULL,
  status TEXT NOT NULL,
  dedupe_key TEXT NOT NULL,
  occurrences INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  last_seen_at TEXT NOT NULL,
  FOREIGN KEY(session_id) REFERENCES qa_sessions(id)
);

CREATE INDEX IF NOT EXISTS idx_tasks_session_id ON tasks(session_id);
CREATE INDEX IF NOT EXISTS idx_tasks_repo_status ON tasks(repo_key, status);
CREATE INDEX IF NOT EXISTS idx_tasks_dedupe ON tasks(repo_key, dedupe_key);

CREATE TABLE IF NOT EXISTS task_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  task_id INTEGER NOT NULL,
  event_type TEXT NOT NULL,
  payload TEXT NOT NULL,
  created_at TEXT NOT NULL,
  FOREIGN KEY(task_id) REFERENCES tasks(id)
);

CREATE INDEX IF NOT EXISTS idx_task_events_task_id ON task_events(task_id);
`
	if _, err := s.db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("migrate sqlite schema: %w", err)
	}
	if _, err := s.db.ExecContext(ctx, "ALTER TABLE tasks ADD COLUMN work_branch TEXT NOT NULL DEFAULT ''"); err != nil && !isDuplicateColumnError(err) {
		return fmt.Errorf("migrate task work_branch column: %w", err)
	}
	if _, err := s.db.ExecContext(ctx, "ALTER TABLE qa_sessions ADD COLUMN ac_text TEXT NOT NULL DEFAULT ''"); err != nil && !isDuplicateColumnError(err) {
		return fmt.Errorf("migrate qa session ac_text column: %w", err)
	}
	return nil
}

func (s *Store) StartSession(ctx context.Context, source, ref, goals, acText string, now time.Time) (Session, error) {
	started := now.UTC().Format(time.RFC3339Nano)
	res, err := s.db.ExecContext(ctx, `
INSERT INTO qa_sessions(source, ref, goals, ac_text, started_at)
VALUES(?, ?, ?, ?, ?)
`, strings.TrimSpace(source), strings.TrimSpace(ref), strings.TrimSpace(goals), strings.TrimSpace(acText), started)
	if err != nil {
		return Session{}, fmt.Errorf("insert qa session: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Session{}, fmt.Errorf("read qa session id: %w", err)
	}
	return Session{ID: id, Source: source, Ref: ref, Goals: goals, ACText: acText, StartedAt: now.UTC()}, nil
}

func (s *Store) FinishSession(ctx context.Context, sessionID int64, verdict string, now time.Time) error {
	ended := now.UTC().Format(time.RFC3339Nano)
	res, err := s.db.ExecContext(ctx, `
UPDATE qa_sessions
SET verdict = ?, ended_at = ?
WHERE id = ?
`, strings.TrimSpace(verdict), ended, sessionID)
	if err != nil {
		return fmt.Errorf("update qa session: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("qa session rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("qa session not found: %d", sessionID)
	}
	return nil
}

func (s *Store) Session(ctx context.Context, sessionID int64) (Session, error) {
	var out Session
	var started string
	var ended sql.NullString
	err := s.db.QueryRowContext(ctx, `
SELECT id, source, ref, goals, ac_text, verdict, started_at, ended_at
FROM qa_sessions
WHERE id = ?
`, sessionID).Scan(&out.ID, &out.Source, &out.Ref, &out.Goals, &out.ACText, &out.Verdict, &started, &ended)
	if err != nil {
		if err == sql.ErrNoRows {
			return Session{}, fmt.Errorf("qa session not found: %d", sessionID)
		}
		return Session{}, fmt.Errorf("query qa session: %w", err)
	}
	st, err := time.Parse(time.RFC3339Nano, started)
	if err != nil {
		return Session{}, fmt.Errorf("parse session started_at: %w", err)
	}
	out.StartedAt = st
	if ended.Valid {
		et, parseErr := time.Parse(time.RFC3339Nano, ended.String)
		if parseErr != nil {
			return Session{}, fmt.Errorf("parse session ended_at: %w", parseErr)
		}
		out.EndedAt = sql.NullTime{Valid: true, Time: et}
	}
	return out, nil
}
