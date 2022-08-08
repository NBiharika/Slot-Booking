package manager

import (
	"Slot_booking/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user entity.User) (entity.User, error)
	Find(userID uint64) (entity.User, error)
	FindUsingEmail(user entity.User) (entity.User, error)
	FindAll() ([]entity.User, error)
	UpdateToUser(email string) error
	UpdateToAdmin(email string) error
	UpdateToBlockUser(email string) error
}

type UserDB struct {
	connection *gorm.DB
}

func UserRepo() UserRepository {
	return &UserDB{
		connection: dbClient,
	}
}

func (db *UserDB) Create(user entity.User) (entity.User, error) {
	err := db.connection.Create(&user).Error
	return user, err
}

func (db *UserDB) Find(userID uint64) (entity.User, error) {
	var user entity.User
	user.ID = userID
	err := db.connection.First(&user).Error

	return user, err
}

func (db *UserDB) FindUsingEmail(user entity.User) (entity.User, error) {
	err := db.connection.Where("email = ?", user.Email).First(&user).Error
	return user, err
}

func (db *UserDB) FindAll() ([]entity.User, error) {
	var users []entity.User
	err := db.connection.Find(&users).Error
	return users, err
}

func (db *UserDB) UpdateToUser(email string) error {
	err := db.connection.Model(&entity.User{}).Where("email=?", email).Update("role", "user").Error
	return err
}

func (db *UserDB) UpdateToAdmin(email string) error {
	err := db.connection.Model(&entity.User{}).Where("email=?", email).Update("role", "admin").Error
	return err
}

func (db *UserDB) UpdateToBlockUser(email string) error {
	err := db.connection.Model(&entity.User{}).Where("email=?", email).Update("role", "blocked_user").Error
	return err
}
