package routes

import (
	"errors"

	"github.com/WytQuant/fiber-gormapi/connectdb"
	"github.com/WytQuant/fiber-gormapi/models"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	// This is serializer
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// helper function
func CreateResponseUser(userModel models.User) User {
	return User{ID: userModel.ID, FirstName: userModel.FirstName, LastName: userModel.LastName}
}

// CRUD operation
// creating user
func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	connectdb.Database.Db.Create(&user)
	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

// get all users
func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}

	connectdb.Database.Db.Find(&users)
	responseUsers := []User{}
	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}

// get user by id
// helper function
func findUser(id int, user *models.User) error {
	connectdb.Database.Db.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("User does not exist")
	}

	return nil
}

func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an int")
	}

	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

// update user info
func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an int")
	}

	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateUser struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	var updateUser UpdateUser

	if err := c.BodyParser(&updateUser); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	user.FirstName = updateUser.FirstName
	user.LastName = updateUser.LastName

	connectdb.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

// delete user
func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an int")
	}

	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := connectdb.Database.Db.Delete(&user).Error; err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully Deleted User!!")
}
