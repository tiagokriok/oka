package repositories

import "github.com/tiagokriok/oka/internal/storages"

type User struct {
	ID       string `json:"id,omitempty" param:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type UserRepository struct {
	db *storages.MysqlDB
}

func NewUserRepository(db *storages.MysqlDB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (ur *UserRepository) CreateUser(user *User) error {
	_, err := ur.db.Exec(
		"INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)",
		user.ID,
		user.Name,
		user.Email,
		user.Password,
	)

	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetUserById(id string) (*User, error) {
	var user User

	err := ur.db.Get(&user, "SELECT u.id, u.name, u.email FROM users u WHERE id=?", id)

	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (ur *UserRepository) UpdateUser(id string, user *User) (*User, error) {
	name := nullString(user.Name)
	email := nullString(user.Email)
	password := nullString(user.Password)

	err := ur.db.Get(user, `
    UPDATE users
    SET name = COALESCE(?, name),
    email = COALESCE(?, email),
    password = COALESCE(?, password)
    WHERE id = ?
    RETURNING id, name, email;
  `, name, email, password, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) DeleteUser(id string) error {
	_, err := ur.db.Exec("DELETE FROM users where id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
