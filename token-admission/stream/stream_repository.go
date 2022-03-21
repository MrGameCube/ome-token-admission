package stream

import (
	"database/sql"
	"time"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS streams(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	app_name TEXT NOT NULL,
	owner_name TEXT,
	owner_id TEXT,
	public INTEGER,
	creation_date TEXT NOT NULL,
	);`
	_, err := r.db.Exec(query)
	return err
}
func (r *SQLiteRepository) Create(stream StreamEntity) (*StreamEntity, error) {
	query := `INSERT INTO streams (name, app_name, owner_name, owner_id, public, creation_date) VALUES(?,?,?,?,?,?)`
	res, err := r.db.Exec(query, stream.StreamName, stream.ApplicationName, stream.OwnerName, stream.OwnerID, stream.Public, stream.CreationDate.Format(time.RFC3339))
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
func (r *SQLiteRepository) All() ([]StreamEntity, error) {
	res, err := r.db.Query("SELECT * FROM streams")
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var all []StreamEntity
	for res.Next() {
		var stream StreamEntity
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
func (r *SQLiteRepository) FindByName(streamName string, appName string) (*StreamEntity, error) {
	res, err := r.db.Query("SELECT * FROM streams WHERE name=? and app_name=?", streamName, appName)
	if !res.Next() || err != nil {
		return nil, err
	}
	defer res.Close()

	var streamInfo *StreamEntity
	var dbISODate string
	err = res.Scan(&streamInfo.ID, &streamInfo.StreamName, &streamInfo.ApplicationName, &streamInfo.OwnerName, &streamInfo.OwnerID, &streamInfo.Public, &dbISODate)
	if err != nil {
		return nil, err
	}

	streamInfo.CreationDate, err = time.Parse(time.RFC3339, dbISODate)
	if err != nil {
		return nil, err
	}

	return streamInfo, nil

}
func (r *SQLiteRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM streams WHERE id=?", id)
	return err
}
