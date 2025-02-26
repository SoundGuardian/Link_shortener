package postgre

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/config"
	"url-shortener/internal/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.DB) (*Storage, error) {
	const op = "storage.postgre.NewStorage"
	DbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Name, cfg.Password)
	fmt.Println("DEBUG", DbInfo, "DEBUG")
	db, err := sql.Open("postgres", DbInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS url(
        id int NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        alias TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL UNIQUE);
    CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
    `)
	// stmt, err := db.Query(`
	// CREATE TABLE IF NOT EXISTS url(
	//     id INTEGER PRIMARY KEY,
	//     alias TEXT NOT NULL UNIQUE,
	//     url TEXT NOT NULL);
	// CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	// `)
	// fmt.Print(stmt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// _, err = stmt.Exec()
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgre.SaveURL"
	var id int64

	// stmt, err := s.db.Prepare("INSERT INTO url(url,alias) values($1,$2) RETURNING id")
	// if err != nil {
	// 	return 0, fmt.Errorf("%s: prepare statment: %w", op, err)
	// }

	res := s.db.QueryRow("INSERT INTO url(url,alias) values($1,$2) RETURNING id", urlToSave, alias).Scan(&id)
	// res, err := stmt.Exec(urlToSave, alias)

	if res != nil {
		return 0, fmt.Errorf("%s: exicute statment: %s", res, op)
	}

	// id, err := res.LastInsertId()
	// if err != nil {
	// 	return 0, fmt.Errorf("%s: %w", op, err)
	// }
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgre.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = $1")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) (string, error) {
	const op = "storage/postgre/DeleteURL"

	stmt, err := s.db.Prepare(`DELETE FROM url WHERE alias = $1`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	// var resURL string
	_, err = stmt.Query(alias)
	// err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "", nil
}
