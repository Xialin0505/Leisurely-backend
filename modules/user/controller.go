package user

import (
	"leisurely/database/models"
	"leisurely/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func UserRouter(app fiber.Router) {
	app.Post("/initUser", InitUser)
	app.Post("/updateUser", UpdateUser)
	app.Get("/users/uid/:uid", GetUserByUID)
	app.Get("/users/username/:username", GetUserByName)
	app.Get("/userLogin/:email/:password", Login)
	app.Delete("/deleteAccount/:uid", DeleteAccount)
}

func Login(c *fiber.Ctx) error {
	email := c.Params("email")
	password := c.Params("password")

	user, err := GetUserLogin(email, password)

	if err != nil && strings.Contains(err.Error(), "record not found") {
		return c.Status(200).JSON(fiber.Map{"success": false, "message": "Login fail", "user": nil})
	} else if err != nil {
		return c.Status(409).JSON(fiber.Map{"success": false, "message": "Other error occurs", "user": nil})
	}

	plan, err := GetPlanByUID(user.UID)
	
	if err != nil {
		return c.Status(201).JSON(fiber.Map{"success": true, "message": "Login successfully, cannot get plan", "user": user})
	}
	
	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Login successfully", "user": user, "plan": plan, "PlanNum": len(plan)})
}

func InitUser(c *fiber.Ctx) error {
	// Take params from client
	var form InitUserDTO

	if err := c.BodyParser(&form); err != nil {
		log.Error().Stack().Err(err).Msg("UserController: Failed to parse client request")
		return err
	}

	// If the user already exists, just return it
	user, err := GetUserByEmail(form.Email)
	if user.UID != 0 {
		return c.Status(409).JSON(fiber.Map{"message:": "A user with this email already exist"})
	}

	if res := utils.ValidateDTO(form); len(res) != 0 {
		return c.Status(400).JSON(res)
	}

	// if form.Username != "" && IsUsernameTaken(form.Username) {
	// 	return c.Status(409).JSON(fiber.Map{"message": "User name exists, choose a new one"})
	// }

	currentUID, err := GetCurrentUserID()
	if err != nil {
		log.Error().Stack().Err(err).Msg("UserController: Database errors cannot get largest UID")
	}

	newUser := models.User{
		UID:      currentUID,
		Name:     "",
		UserName: strings.ToLower(form.Username),
		Email:    form.Email,
		Password: form.Password,
		//Phone:    &form.PhoneNumber,
		PhotoUrl: form.PhotoUrl,
		Gender:   models.Gender(form.Gender),
		//Flags:    0,
		//Bio:      form.Bio,
	}

	if time, err := time.Parse(YYYYMMDD, form.Birthday); err == nil {
		log.Error().Stack().Err(err).Msg(time.GoString())
		newUser.Birthday = time
	}

	if err := CreateUserProfile(&newUser); err != nil {
		return fiber.ErrForbidden
	}

	return c.Status(201).JSON(newUser)

}

func UpdateUser(c *fiber.Ctx) error {
	var form UpdateUserDTO

	if err := c.BodyParser(&form); err != nil {
		log.Error().Stack().Err(err).Msg("UserController: Failed to parse client request")
		return err
	}

	if res := utils.ValidateDTO(form); len(res) != 0 {
		return c.Status(400).JSON(res)
	}

	// User got from context is read only, so need to fetch from DB again
	user, err := GetUserProfileByID(form.Uid)
	if err != nil {
		return err
	}

	// if form.Username != "" && form.Username != user.UserName && IsUsernameTaken(form.Username) {
	// 	return c.Status(409).JSON(fiber.Map{"message": "Username exists, choose a new one"})
	// }

	if form.Username != "" {
		user.UserName = form.Username
	}

	if form.Email != "" {
		user.Email = form.Email
	}

	if form.Name != "" {
		user.Name = form.Name
	}

	if form.Gender != 0 {
		user.Gender = models.Gender(form.Gender)
	}

	if form.PhotoUrl != "" {
		user.PhotoUrl = form.PhotoUrl
	}

	if form.Birthday != "" {
		if time, err := time.Parse(YYYYMMDD, form.Birthday); err == nil {
			user.Birthday = time
		}
	}

	if form.Password != "" {
		user.Password = form.Password
	}

	if err := UpdateUserProfile(user); err != nil {
		// Username duplicated
		if err == gorm.ErrInvalidValue {
			return c.JSON(fiber.Map{"success": false, "message": "Username not available", "data": user})
		}
		return err
	}

	return c.JSON(fiber.Map{"success": true, "message": "Updated successfully", "data": user})
}

func GetUserByUID(c *fiber.Ctx) error {
	uid, err := strconv.Atoi(c.Params("uid"))
	user, err := GetUserProfileByID(uid)
	
	if err != nil {
		return fiber.ErrNotFound
	}

	plans, err := GetPlanByUID(uid)
	
	if err != nil {
		return c.Status(200).JSON(fiber.Map{"user": user})
	}

	return c.JSON(fiber.Map{"user": user, "plan": plans, "PlanNum": len(plans)})
}

func GetUserByName(c *fiber.Ctx) error {
	username := c.Params("username")

	user, err := GetUserProfileByName(username)
	if err != nil {
		return fiber.ErrNotFound
	}
	
	plans, err := GetPlanByUID(user.UID)
	
	if err != nil {
		return c.Status(200).JSON(fiber.Map{"user": user})
	}

	return c.JSON(fiber.Map{"user": user, "plan": plans, "PlanNum": len(plans)})
}

func DeleteAccount(c *fiber.Ctx) error {
	uid, err := strconv.Atoi(c.Params("uid"))

	user, _ := GetUserProfileByID(uid)
	err = DeleteUserProfile(user)

	if err != nil {
		c.Status(400).JSON(fiber.Map{
			"message": "Failed to delete user: " + err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": "true",
	})
}
