package repository

import (
	"a21hc3NpZ25tZW50/model"
	"database/sql"
	"fmt"
)

type SessionsRepository interface {
	AddSessions(session model.Session) error
	DeleteSession(token string) error
	UpdateSessions(session model.Session) error
	SessionAvailName(name string) error
	SessionAvailToken(token string) (model.Session, error)

	FetchByID(id int) (*model.Session, error)
}

type sessionsRepoImpl struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *sessionsRepoImpl {
	return &sessionsRepoImpl{db}
}

func (s *sessionsRepoImpl) AddSessions(session model.Session) error {
	_, err := s.db.Exec(
		"INSERT INTO sessions (token, username, expiry) VALUES ($1, $2, $3)",
		session.Token,
		session.Username,
		session.Expiry,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionsRepoImpl) DeleteSession(token string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE token = $1", token)
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionsRepoImpl) UpdateSessions(session model.Session) error {
	_, err := s.db.Exec(
		"UPDATE sessions SET token = $1, username = $2, expiry = $3 WHERE id = $4",
		session.Token,
		session.Username,
		session.Expiry,
		session.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionsRepoImpl) SessionAvailName(name string) error {
	row := s.db.QueryRow("SELECT id FROM sessions WHERE username = $1", name)

	var id int
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return fmt.Errorf("session for user %s already exists", name)
}

func (s *sessionsRepoImpl) SessionAvailToken(token string) (model.Session, error) {
	row := s.db.QueryRow("SELECT id, token, username, expiry FROM sessions WHERE token = $1", token)

	var session model.Session
	err := row.Scan(&session.ID, &session.Token, &session.Username, &session.Expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Session{}, nil 
		}
		return model.Session{}, err
	}

	return session, nil
}

func (u *sessionsRepoImpl) FetchByID(id int) (*model.Session, error) {
	row := u.db.QueryRow("SELECT id, token, username, expiry FROM sessions WHERE id = $1", id)

	var session model.Session
	err := row.Scan(&session.ID, &session.Token, &session.Username, &session.Expiry)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
