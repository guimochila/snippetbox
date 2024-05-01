package mocks

import "snippetbox.guimochila.com/internal/models"

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "engineer@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "go@example.com" && password == "password" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentails
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
