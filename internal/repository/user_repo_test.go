package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	user := &domain.User{
		ID:             1,
		Name:           "John Doe",
		Email:          "john.doe@example.com",
		Password:       "password123",
		ProfilePicture: "profile.jpg",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.ID, user.Name, user.Email, user.Password, user.ProfilePicture).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Implementasikan unit test lainnya sesuai dengan metode yang ada pada UserRepository

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	user := &domain.User{
		ID:             1,
		Name:           "John Doe",
		Email:          "john.doe@example.com",
		Password:       "password123",
		ProfilePicture: "profile.jpg",
	}

	mock.ExpectExec("UPDATE users").
		WithArgs(user.ID, user.Name, user.Email, user.Password, user.ProfilePicture).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	id := 1

	mock.ExpectExec("DELETE FROM users").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(id)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	id := 1

	rows := sqlmock.NewRows([]string{"name", "email", "password", "profile_picture", "is_deleted"}).
		AddRow("John Doe", "john.doe@example.com", "password123", "profile.jpg", false)

	mock.ExpectQuery("SELECT name, email, password, profile_picture, is_deleted FROM users").
		WithArgs(id).
		WillReturnRows(rows)

	user, err := repo.FindOne(id)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"user_id", "name", "email", "password", "profile_picture", "is_deleted"}).
		AddRow(1, "John Doe", "john.doe@example.com", "password123", "profile.jpg", false)

	mock.ExpectQuery("SELECT user_id, name, email, password, profile_picture, is_deleted FROM users").
		WillReturnRows(rows)

	users, err := repo.FindAll()
	assert.NoError(t, err)
	assert.NotEmpty(t, users)
	assert.Equal(t, 1, len(users))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	username := "John Doe"

	rows := sqlmock.NewRows([]string{"name", "email", "password", "profile_picture", "is_deleted"}).
		AddRow("John Doe", "john.doe@example.com", "password123", "profile.jpg", false)

	mock.ExpectQuery("SELECT name, email, password, profile_picture, is_deleted FROM users").
		WithArgs(username).
		WillReturnRows(rows)

	user, err := repo.FindByUsername(username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}