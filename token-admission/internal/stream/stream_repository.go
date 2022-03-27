package stream

import (
	"database/sql"
	"errors"
	"github.com/MrGameCube/ome-token-admission/token-admission/ta-models"
	"time"
)

var (
	ErrNotFound = errors.New("stream not found")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) (*SQLiteRepository, error) {
	repo := SQLiteRepository{db: db}
	err := repo.Migrate()
	return &repo, err
}

func (r *SQLiteRepository) Migrate() error {
	query := `CREATE TABLE IF NOT EXISTS streams(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	app_name TEXT NOT NULL,
	owner_name TEXT,
	owner_id TEXT,
	public INTEGER,
	creation_date TEXT NOT NULL
	);`
	_, err := r.db.Exec(query)
	return err
}
func (r *SQLiteRepository) Create(streamParams *ta_models.StreamParameters) (*ta_models.StreamEntity, error) {
	stream := ta_models.StreamEntity{
		Title:           streamParams.Title,
		StreamName:      streamParams.StreamName,
		ApplicationName: streamParams.ApplicationName,
		CreationDate:    time.Now(),
		OwnerName:       streamParams.OwnerName,
		OwnerID:         streamParams.OwnerID,
		Public:          streamParams.Public,
	}
	query := `INSERT INTO streams (name, app_name, owner_name, owner_id, public, creation_date) VALUES(?,?,?,?,?,?)`
	res, err := r.db.Exec(query, streamParams.StreamName, streamParams.ApplicationName, streamParams.OwnerName, streamParams.OwnerID, streamParams.Public, stream.CreationDate.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	stream.ID = id
	return &stream, nil
}
func (r *SQLiteRepository) All() ([]ta_models.StreamEntity, error) {
	res, err := r.db.Query("SELECT * FROM streams")
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var all []ta_models.StreamEntity
	for res.Next() {
		var stream ta_models.StreamEntity
		var dbISODate string
		if err := res.Scan(&stream.ID,
			&stream.StreamName,
			&stream.ApplicationName,
			&stream.OwnerName,
			&stream.OwnerID,
			&stream.Public,
			&dbISODate); err != nil {
			return nil, err
		}

		stream.CreationDate, err = time.Parse(time.RFC3339, dbISODate)
		if err != nil {
			return nil, err
		}

		all = append(all, stream)
	}
	return all, nil
}
func (r *SQLiteRepository) FindByName(streamName string, appName string) (*ta_models.StreamEntity, error) {
	res, err := r.db.Query("SELECT * FROM streams WHERE name=? and app_name=?", streamName, appName)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if !res.Next() {
		return nil, ErrNotFound
	}
	var streamInfo ta_models.StreamEntity
	var dbISODate string
	err = res.Scan(&streamInfo.ID, &streamInfo.StreamName, &streamInfo.ApplicationName, &streamInfo.OwnerName, &streamInfo.OwnerID, &streamInfo.Public, &dbISODate)
	if err != nil {
		return nil, err
	}

	streamInfo.CreationDate, err = time.Parse(time.RFC3339, dbISODate)
	if err != nil {
		return nil, err
	}

	return &streamInfo, nil

}
func (r *SQLiteRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM streams WHERE id=?", id)
	return err
}
