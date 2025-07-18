package postgre

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"project_1/internal/storage"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storageLink string) (*Storage, error) {
	const op = "storage.postgre.New"

	db, err := sql.Open("postgres", storageLink)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS url (
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
			`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_alias ON url (alias);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgre.SaveURL"
	var id int64

	stmt, err := s.db.Prepare("INSERT INTO url (url, alias) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRow(urlToSave, alias).Scan(&id)
	if err != nil {
		if pqErr, _ := err.(*pq.Error); errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	slog.Info("URL saved successfully", slog.String("id", fmt.Sprintf("%d", id)))

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgre.GetURL"
	var resultURL string

	row := s.db.QueryRow("SELECT url FROM url WHERE alias = $1", alias)
	err := row.Scan(&resultURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	slog.Info("URL retrieved successfully", slog.String("alias", alias), slog.String("url", resultURL))

	return resultURL, nil
}

func (s *Storage) AliasExists(alias string) (bool, error) {
	const op = "storage.postgre.AliasExists"
	var resultURL bool

	row := s.db.QueryRow("SELECT EXISTS (SELECT 1 FROM url WHERE username = $1)", alias)
	err := row.Scan(&resultURL)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resultURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgre.DeleteURL"
	var deletedURL string
	var deletedId string

	deletedRow := s.db.QueryRow("SELECT url FROM url WHERE alias = $1", alias)
	err := deletedRow.Scan(&deletedURL)
	if err != nil {
		return fmt.Errorf("%s: %w; alias: %s ", op, err, alias)
	}

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = $1 RETURNING id")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRow(alias).Scan(&deletedId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	slog.Info("URL was deleted succsessfully", slog.String("id", deletedId), slog.String("URL", deletedURL), slog.String("alias", alias))

	return nil
}
