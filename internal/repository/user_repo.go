package repository

import (
	"database/sql"
	"log"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}
func (u UserRepository) Create(newUser *domain.User) error {
	query := `INSERT INTO users (user_id, name, email, password, profile_picture) VALUES ($1, $2, $3, $4, $5)`
	_, err := u.db.Exec(query, newUser.ID, newUser.Name, newUser.Email, newUser.Password, newUser.ProfilePicture)
	if err != nil {
		panic(err)
	} else {
		log.Println("Succsessfully added user")
	}
	return err
}

func (u UserRepository) Update(updetedUser *domain.User) error {
	query := `UPDATE users SET name=$2, email=$3, password=$4, profile_picture=$5 WHERE user_id=$1`
	_, err := u.db.Exec(query, updetedUser.ID, updetedUser.Name, updetedUser.Email, updetedUser.Password, updetedUser.ProfilePicture)
	if err != nil {
		panic(err)
	} else {
		log.Println("Succsessfully updated")
	}

	return err
}

func (u UserRepository) Delete(id int) error {
	query := `UPADATE users SET is_delete = true WHERE user_id = $1`
	_, err := u.db.Exec(query, id)
	if err != nil {
		panic(err)
	} else {
		log.Println("Succsessfully deleted data")
	}
	return err
}

func (u UserRepository) FindOne(id int) (*domain.User, error) {
	query := `SELECT name, email, password, profile_picture FROM users WHERE user_id=$1;`
	row := u.db.QueryRow(query, id)
	var user domain.User
	err := row.Scan(&user.Name, &user.Email, &user.Password, &user.ProfilePicture)
	if err != nil {
		panic(err)
	}
	user.ID = id
	return &user, nil
}

func (u UserRepository) FindAll() ([]domain.User, error) {
	query := `SELECT user_id, name, email, password, profile_picture FROM users;`
	rows, err := u.db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.ProfilePicture)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, err
}
