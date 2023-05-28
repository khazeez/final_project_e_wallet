package repository

import (
	"database/sql"
	"log"


	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(newUser *domain.User) error
	Update(updatedUser *domain.User) error
	Delete(id int) error
	FindOne(id int) (*domain.User, error)
	FindAll() ([]domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(newUser *domain.User) error {
	query := `INSERT INTO users (name, email, password, profile_picture) VALUES ($1, $2, $3, $4)`
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Simpan pengguna ke dalam database
	_, err = u.db.Exec(query, newUser.Name, newUser.Email, hashedPassword, newUser.ProfilePicture)
	if err != nil {
		return err
	}

	log.Println("Successfully added user")
	return nil
}

func (u *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := u.db.Exec(query, id)
	if err != nil {
		log.Println("Failed to delete user:", err)
		return err
	}

	log.Println("Successfully deleted user")
	return nil
}

func (u *userRepository) FindOne(id int) (*domain.User, error) {
	query := `SELECT user_id, name, email, password, profile_picture, is_deleted FROM users WHERE user_id=$1`
	row := u.db.QueryRow(query, id)
	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.ProfilePicture, &user.IsDeleted)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) FindAll() ([]domain.User, error) {
	query := `SELECT user_id, name, email, password, profile_picture, is_deleted FROM users`
	rows, err := u.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.ProfilePicture, &user.IsDeleted)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
func (u *userRepository) Update(updatedUser *domain.User) error {
	query := `UPDATE users SET name=$1, email=$2, profile_picture=$3 WHERE user_id=$4`

	_, err := u.db.Exec(query, updatedUser.Name, updatedUser.Email, updatedUser.ProfilePicture, updatedUser.ID)
	if err != nil {
		log.Println("Failed to update user:", err)
		return err
	}

	log.Println("Successfully updated user")
	return nil
}
