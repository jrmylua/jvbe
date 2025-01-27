package user

import (
	"database/sql"
	"errors"
	"slices"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/matoous/go-nanoid/v2"
)

type User struct {
	Id         string    `db:"id"`
	FullName   string    `db:"full_name"`
	ExternalId string    `db:"external_id"`
	IsAdmin    bool      `db:"is_admin"`
	CreatedAt  time.Time `db:"created_at"`
}

func (u *User) ToSessionUser() SessionUser {
	return SessionUser{
		Id: u.Id,
	}
}

type ExternalUser struct {
	Id       string `json:"sub"`
	FullName string `json:"name"`
	Picture  string `json:"picture"`
}

type SessionUser struct {
	Id          string
	Permissions []string
}

func (u SessionUser) IsAuthenticated() bool {
	return u.Id != ""
}

func (u SessionUser) hasPermission(p string) bool {
	return slices.Contains[[]string](u.Permissions, p)
}

func (u SessionUser) CanCreateEvent() bool {
	return u.hasPermission("create:event")
}

func (u SessionUser) CanDeleteEvent() bool {
	return u.hasPermission("delete:event")
}

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

var (
	ErrNoUser = errors.New("no user found")
)

func (s *Store) GetByExternalId(externalId string) (User, error) {
	stmt := `
        SELECT id, is_admin FROM user
        WHERE external_id = ?
    `
	args := []any{externalId}

	var user User
	err := s.db.Get(&user, stmt, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNoUser
	} else if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Store) GetById(id string) (User, error) {
	stmt := `
        SELECT id, full_name, external_id, is_admin, created_at FROM user
        WHERE id = ?
    `
	args := []any{id}

	var user User
	err := s.db.Get(&user, stmt, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNoUser
	} else if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Store) CreateFromExternal(externalUser ExternalUser) (User, error) {
	newId, err := gonanoid.New()
	if err != nil {
		return User{}, err
	}
	isAdmin := false

	stmt := `
        INSERT INTO user (id, full_name, external_id, is_admin, created_at)
        VALUES (?, ?, ?, ?, ?)
    `
	args := []any{
		newId,
		externalUser.FullName,
		externalUser.Id,
		isAdmin,
		time.Now().UTC(),
	}

	_, err = s.db.Exec(stmt, args...)
	if err != nil {
		return User{}, err
	}

	newUser, err := s.GetById(newId)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}
