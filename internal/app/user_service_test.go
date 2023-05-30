package app_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	user := &domain.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	// Mengharapkan eksekusi query INSERT yang berhasil
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))

	err = usecase.InsertUser(user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertUser_ValidationError(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	user := &domain.User{
		Name:     "",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	err = usecase.InsertUser(user)
	assert.EqualError(t, err, "name is required")
}

func TestInsertUser_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	user := &domain.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	// Mengharapkan error ketika eksekusi query INSERT gagal
	mock.ExpectExec("INSERT INTO users").WillReturnError(errors.New("database error"))

	err = usecase.InsertUser(user)
	assert.EqualError(t, err, "database error")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	user := &domain.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	// Mengharapkan eksekusi query UPDATE yang berhasil
	mock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))

	err = usecase.UpdateUser(user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_ValidationError(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	user := &domain.User{
		ID:       1,
		Name:     "",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	err = usecase.UpdateUser(user)
	assert.EqualError(t, err, "name is required")
}

func TestUpdateUser_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	user := &domain.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	// Mengharapkan error ketika eksekusi query UPDATE gagal
	mock.ExpectExec("UPDATE users").WillReturnError(errors.New("database error"))

	err = usecase.UpdateUser(user)
	assert.EqualError(t, err, "database error")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindOne_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	userID := 1

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(userID, "John Doe", "john.doe@example.com")

	// Mengharapkan eksekusi query SELECT yang mengembalikan data pengguna
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").WithArgs(userID).WillReturnRows(rows)

	user, err := usecase.FindOne(userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john.doe@example.com", user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindOne_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	userID := 1

	// Mengharapkan eksekusi query SELECT yang tidak mengembalikan data pengguna
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").WithArgs(userID).WillReturnRows(sqlmock.NewRows([]string{}))

	user, err := usecase.FindOne(userID)
	assert.EqualError(t, err, "user not found")
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "John Doe", "john.doe@example.com").
		AddRow(2, "Jane Smith", "jane.smith@example.com")

	// Mengharapkan eksekusi query SELECT yang mengembalikan beberapa pengguna
	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

	users, err := usecase.FindAll()
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAll_NoUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	// Mengharapkan eksekusi query SELECT yang tidak mengembalikan pengguna
	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(sqlmock.NewRows([]string{}))

	users, err := usecase.FindAll()
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 0)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	userID := 1

	// Mengharapkan eksekusi query DELETE yang berhasil
	mock.ExpectExec("DELETE FROM users WHERE id = ?").WithArgs(userID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = usecase.Delete(userID)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	userID := 1

	// Mengharapkan eksekusi query DELETE yang tidak mempengaruhi baris data
	mock.ExpectExec("DELETE FROM users WHERE id = ?").WithArgs(userID).WillReturnResult(sqlmock.NewResult(0, 0))

	err = usecase.Delete(userID)
	assert.EqualError(t, err, "user not found")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByUsername_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	username := "john.doe"

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "John Doe", "john.doe@example.com")

	// Mengharapkan eksekusi query SELECT yang mengembalikan data pengguna berdasarkan username
	mock.ExpectQuery("SELECT (.+) FROM users WHERE username = ?").WithArgs(username).WillReturnRows(rows)

	user, err := usecase.FindByUsername(username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john.doe@example.com", user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByUsername_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Membuat objek repository mock
	repo := repository.NewUserRepository(db)

	// Membuat objek usecase dengan repository mock
	usecase := app.NewUserUsecase(repo)

	username := "john.doe"

	// Mengharapkan eksekusi query SELECT yang tidak mengembalikan data pengguna berdasarkan username
	mock.ExpectQuery("SELECT (.+) FROM users WHERE username = ?").WithArgs(username).WillReturnRows(sqlmock.NewRows([]string{}))

	user, err := usecase.FindByUsername(username)
	assert.EqualError(t, err, "user not found")
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}