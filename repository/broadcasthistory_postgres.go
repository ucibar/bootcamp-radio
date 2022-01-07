package repository

import (
	"database/sql"
	"errors"
	"github.com/uCibar/bootcamp-radio/entity"
)

type BroadcastHistoryPostgresRepository struct {
	db *sql.DB
}

func NewBroadcastHistoryPostgresRepository(db *sql.DB) *BroadcastHistoryPostgresRepository {
	return &BroadcastHistoryPostgresRepository{db: db}
}

func (repository *BroadcastHistoryPostgresRepository) Create(bh *entity.BroadcastHistory) error {
	var insertedId int64

	err := repository.db.QueryRow("INSERT INTO broadcast_history(user_id, title, max_viewers, begin_at, end_at) VALUES($1, $2, $3, $4, $5) RETURNING id;",
		bh.UserID, bh.Title, bh.MaxViewers, bh.BeginAt, bh.EndAt).Scan(&insertedId)
	if err != nil {
		return err
	}

	bh.ID = insertedId

	return err
}

func (repository *BroadcastHistoryPostgresRepository) GetByID(id int64) (*entity.BroadcastHistory, error) {
	var bh entity.BroadcastHistory

	err := repository.db.QueryRow("SELECT * FROM broadcast_history WHERE id = $1", id).
		Scan(&bh.ID, &bh.UserID, &bh.Title, &bh.MaxViewers, &bh.BeginAt, &bh.EndAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrBroadcastNotFound
	} else if err != nil {
		return nil, err
	}

	return &bh, nil
}

func (repository *BroadcastHistoryPostgresRepository) AllByUserID(userID int64) ([]*entity.BroadcastHistory, error) {
	bhs := make([]*entity.BroadcastHistory, 0)

	rows, err := repository.db.Query("SELECT * FROM broadcast_history WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var bh entity.BroadcastHistory
		err = rows.Scan(&bh.ID, &bh.UserID, &bh.Title, &bh.MaxViewers, &bh.BeginAt, &bh.EndAt)
		if err != nil {
			return nil, err
		}

		bhs = append(bhs, &bh)
	}

	return bhs, nil
}
