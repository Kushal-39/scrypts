package storage

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func Init(path string) error {
	var err error
	db, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	// keep connections small for the sqlite driver in this simple app
	db.SetMaxOpenConns(1)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// enable foreign keys and WAL mode for better durability/concurrency
	if _, err = tx.Exec(`PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;`); err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
  CREATE TABLE IF NOT EXISTS users (
  username TEXT PRIMARY KEY,
  password_hash TEXT NOT NULL,
  wrapped_key BLOB,
  wrapped_nonce BLOB,
  created_at INTEGER NOT NULL
);

    if err != nil {
        tx.Rollback()
        return err
    }
    return tx.Commit()
  nonce BLOB NOT NULL,
  created INTEGER NOT NULL,
  modified INTEGER NOT NULL,
  FOREIGN KEY(owner) REFERENCES users(username)
);

CREATE INDEX IF NOT EXISTS idx_notes_owner ON notes(owner);

	`)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}

type User struct {
	Username     string
	PasswordHash string
	WrappedKey   []byte
	WrappedNonce []byte
	CreatedAt    int64
}

func CreateUser(u User) error {
	if db == nil {
		return errors.New("DB not initialized")
	}
	_, err := db.Exec(`INSERT INTO users(username,password_hash,wrapped_key,wrapped_nonce,created_at) VALUES (?,?,?,?,?)`, u.Username, u.PasswordHash, u.WrappedKey, u.WrappedNonce, u.CreatedAt)
	return err
}

func GetUser(username string) (User, error) {
	var u User
	if db == nil {
		return User{}, errors.New("DB not initialized")
	}
	row := db.QueryRow(`SELECT username, password_hash, wrapped_key, wrapped_nonce, created_at FROM users WHERE username = ?`, username)
	var wk, wn []byte
	if err := row.Scan(&u.Username, &u.PasswordHash, &wk, &wn, &u.CreatedAt); err != nil {
		return User{}, err
	}
	u.WrappedKey = wk
	u.WrappedNonce = wn
	return u, nil
}

func SaveWrappedKey(username string, wrapped, nonce []byte) error {
	if db == nil {
		return errors.New("DB not initialized")
	}
	_, err := db.Exec(`UPDATE users SET wrapped_key = ?, wrapped_nonce = ? WHERE username = ?`, wrapped, nonce, username)
	return err
}

type Note struct {
	ID       string
	Owner    string
	Content  []byte
	Nonce    []byte
	Created  int64
	Modified int64
}

func SaveNote(n Note) error {
	if db == nil {
		return errors.New("db not initialized")
	}
	_, err := db.Exec(`INSERT INTO notes(id,owner,content,nonce,created,modified) VALUES(?,?,?,?,?,?)`,
		n.ID, n.Owner, n.Content, n.Nonce, n.Created, n.Modified)
	return err
}

func UpdateNote(n Note) error {
	if db == nil {
		return errors.New("db not initialized")
	}
	_, err := db.Exec(`UPDATE notes SET content = ?, nonce = ?, modified = ? WHERE id = ? AND owner = ?`,
		n.Content, n.Nonce, n.Modified, n.ID, n.Owner)
	return err
}

func DeleteNote(id, owner string) error {
	if db == nil {
		return errors.New("db not initialized")
	}
	_, err := db.Exec(`DELETE FROM notes WHERE id = ? AND owner = ?`, id, owner)
	return err
}

func GetNotesByOwner(owner string) ([]Note, error) {
	if db == nil {
		return nil, errors.New("db not initialized")
	}
	rows, err := db.Query(`SELECT id, owner, content, nonce, created, modified FROM notes WHERE owner = ?`, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Owner, &n.Content, &n.Nonce, &n.Created, &n.Modified); err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
