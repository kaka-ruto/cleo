package taskstore

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (s *Store) AddTaskEvent(ctx context.Context, taskID int64, eventType string, payload string, now time.Time) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO task_events(task_id, event_type, payload, created_at)
VALUES(?, ?, ?, ?)
`, taskID, strings.TrimSpace(eventType), strings.TrimSpace(payload), now.UTC().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("insert task event: %w", err)
	}
	return nil
}
