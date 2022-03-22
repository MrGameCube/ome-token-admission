package token

import (
	"database/sql"
	"errors"
	"github.com/MrGameCube/ome-token-admission/token-admission/ta-models"
	"time"
)

var (
	ErrNotFound = errors.New("No entry found")
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
	// TODO: Delete expired tokens with trigger
	query := `
	CREATE TABLE IF NOT EXISTS tokens(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	token TEXT NOT NULL UNIQUE,
	app_name TEXT NOT NULL,
	stream_name TEXT NOT NULL,
	direction TEXT CHECK(direction IN ('outgoing', 'incoming')) NOT NULL ,
	expires_at TEXT NOT NULL
	);`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Create(token ta_models.TokenEntity) (*ta_models.TokenEntity, error) {
	query := `INSERT INTO tokens (token, app_name, stream_name, direction, expires_at) VALUES(?,?,?,?,?)`
	res, err := r.db.Exec(query, token.Token, token.Application, token.Stream, token.Direction, token.ExpiresAt.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	token.ID = id
	return &token, nil
}

func (r *SQLiteRepository) All() ([]ta_models.TokenEntity, error) {
	res, err := r.db.Query("SELECT * FROM tokens WHERE expired_at < datetime()")
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var all []ta_models.TokenEntity
	for res.Next() {
		var token ta_models.TokenEntity
		var dbDate string
		if err := res.Scan(&token.ID,
			&token.Token,
			&token.Application,
			&token.Stream,
			&token.Direction,
			&dbDate); err != nil {
			return nil, err
		}

		token.ExpiresAt, err = time.Parse(time.RFC3339, dbDate)
		if err != nil {
			return nil, err
		}

		all = append(all, token)
	}
	return all, nil
}
func (r *SQLiteRepository) FindByToken(token string) (*ta_models.TokenEntity, error) {
	res, err := r.db.Query("SELECT * FROM tokens WHERE token=? and datetime() < expires_at", token)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	if !res.Next() {
		return nil, ErrNotFound
	}

	var tokenEnt ta_models.TokenEntity
	var dbISODate string
	err = res.Scan(&tokenEnt.ID, &tokenEnt.Token, &tokenEnt.Application, &tokenEnt.Stream, &tokenEnt.Direction, &dbISODate)
	if err != nil {
		return nil, err
	}

	tokenEnt.ExpiresAt, err = time.Parse(time.RFC3339, dbISODate)
	if err != nil {
		return nil, err
	}

	return &tokenEnt, nil

}
func (r *SQLiteRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM tokens WHERE id=?", id)
	return err
}
