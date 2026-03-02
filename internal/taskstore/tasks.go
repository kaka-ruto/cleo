package taskstore

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func (s *Store) UpsertOpenTask(ctx context.Context, in Task, now time.Time) (Task, bool, error) {
	found, err := s.openTaskByDedupe(ctx, in.RepoKey, in.DedupeKey)
	if err != nil {
		return Task{}, false, err
	}
	if found != nil {
		updated := now.UTC().Format(time.RFC3339Nano)
		res, execErr := s.db.ExecContext(ctx, `
UPDATE tasks
SET occurrences = occurrences + 1,
    details = ?,
    severity = ?,
    updated_at = ?,
    last_seen_at = ?
WHERE id = ?
`, strings.TrimSpace(in.Details), strings.TrimSpace(in.Severity), updated, updated, found.ID)
		if execErr != nil {
			return Task{}, false, fmt.Errorf("update duplicate task: %w", execErr)
		}
		if _, affErr := res.RowsAffected(); affErr != nil {
			return Task{}, false, fmt.Errorf("duplicate task rows affected: %w", affErr)
		}
		out, getErr := s.Task(ctx, found.ID)
		if getErr != nil {
			return Task{}, false, getErr
		}
		return out, false, nil
	}
	created := now.UTC().Format(time.RFC3339Nano)
	res, err := s.db.ExecContext(ctx, `
INSERT INTO tasks(
  session_id,
  repo_key,
  work_branch,
  title,
  details,
  severity,
  status,
  dedupe_key,
  occurrences,
  created_at,
  updated_at,
  last_seen_at
)
VALUES(?, ?, '', ?, ?, ?, 'open', ?, 1, ?, ?, ?)
`, in.SessionID, strings.TrimSpace(in.RepoKey), strings.TrimSpace(in.Title), strings.TrimSpace(in.Details), strings.TrimSpace(in.Severity), strings.TrimSpace(in.DedupeKey), created, created, created)
	if err != nil {
		return Task{}, false, fmt.Errorf("insert task: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Task{}, false, fmt.Errorf("read task id: %w", err)
	}
	out, err := s.Task(ctx, id)
	if err != nil {
		return Task{}, false, err
	}
	return out, true, nil
}

func (s *Store) openTaskByDedupe(ctx context.Context, repoKey, dedupeKey string) (*Task, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, session_id, repo_key, work_branch, title, details, severity, status, dedupe_key, occurrences, created_at, updated_at, last_seen_at
FROM tasks
WHERE repo_key = ? AND dedupe_key = ? AND status IN ('open', 'in_progress')
ORDER BY id ASC
LIMIT 1
`, strings.TrimSpace(repoKey), strings.TrimSpace(dedupeKey))
	task, err := scanTask(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query task by dedupe: %w", err)
	}
	return &task, nil
}

func (s *Store) Task(ctx context.Context, id int64) (Task, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, session_id, repo_key, work_branch, title, details, severity, status, dedupe_key, occurrences, created_at, updated_at, last_seen_at
FROM tasks
WHERE id = ?
`, id)
	task, err := scanTask(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, fmt.Errorf("task not found: %d", id)
		}
		return Task{}, fmt.Errorf("query task: %w", err)
	}
	return task, nil
}

func (s *Store) TasksBySession(ctx context.Context, sessionID int64) ([]Task, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, session_id, repo_key, work_branch, title, details, severity, status, dedupe_key, occurrences, created_at, updated_at, last_seen_at
FROM tasks
WHERE session_id = ?
ORDER BY id ASC
`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("query tasks by session: %w", err)
	}
	defer rows.Close()
	var out []Task
	for rows.Next() {
		task, scanErr := scanTask(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan task row: %w", scanErr)
		}
		out = append(out, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks by session: %w", err)
	}
	return out, nil
}

func (s *Store) ListTasks(ctx context.Context, status string) ([]Task, error) {
	status = strings.TrimSpace(status)
	query := `
SELECT id, session_id, repo_key, work_branch, title, details, severity, status, dedupe_key, occurrences, created_at, updated_at, last_seen_at
FROM tasks`
	args := []any{}
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}
	query += " ORDER BY id ASC"
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()
	var out []Task
	for rows.Next() {
		task, scanErr := scanTask(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan listed task: %w", scanErr)
		}
		out = append(out, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listed tasks: %w", err)
	}
	return out, nil
}

func (s *Store) UpdateTaskStatus(ctx context.Context, id int64, status string, now time.Time) error {
	res, err := s.db.ExecContext(ctx, `
UPDATE tasks
SET status = ?, updated_at = ?
WHERE id = ?
`, strings.TrimSpace(status), now.UTC().Format(time.RFC3339Nano), id)
	if err != nil {
		return fmt.Errorf("update task status: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("task status rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found: %d", id)
	}
	return nil
}

func (s *Store) SetTaskWorkBranch(ctx context.Context, id int64, branch string, now time.Time) error {
	res, err := s.db.ExecContext(ctx, `
UPDATE tasks
SET work_branch = ?, updated_at = ?
WHERE id = ?
`, strings.TrimSpace(branch), now.UTC().Format(time.RFC3339Nano), id)
	if err != nil {
		return fmt.Errorf("set task work branch: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("task work branch rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found: %d", id)
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTask(row scanner) (Task, error) {
	var out Task
	var created string
	var updated string
	var seen string
	err := row.Scan(&out.ID, &out.SessionID, &out.RepoKey, &out.WorkBranch, &out.Title, &out.Details, &out.Severity, &out.Status, &out.DedupeKey, &out.Occurrences, &created, &updated, &seen)
	if err != nil {
		return Task{}, err
	}
	createdAt, err := time.Parse(time.RFC3339Nano, created)
	if err != nil {
		return Task{}, fmt.Errorf("parse task created_at: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, updated)
	if err != nil {
		return Task{}, fmt.Errorf("parse task updated_at: %w", err)
	}
	seenAt, err := time.Parse(time.RFC3339Nano, seen)
	if err != nil {
		return Task{}, fmt.Errorf("parse task last_seen_at: %w", err)
	}
	out.CreatedAt = createdAt
	out.UpdatedAt = updatedAt
	out.LastSeenAt = seenAt
	return out, nil
}
