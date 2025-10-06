package pg

import "context"

func (s *Storage) Validate(ctx context.Context, currentUserID, noteID int) error {
	query := `
		SELECT user_id FROM notes WHERE id = $1
	`

	var userID int
	err := s.db.QueryRowContext(ctx, query, noteID).Scan(&userID)
	if err != nil {
		return err
	}
	if currentUserID != userID {
		return ErrForbidden
	}
	return nil
}
