package service

import (
	"log"
	"milestone3/be/internal/dto"
	"milestone3/be/internal/entity"
	"milestone3/be/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(user *entity.Users) error
	GetByEmail(email string) (user entity.Users, err error)
	GetById(id int) (user entity.Users, err error)
}

type UserServ struct {
	userRepo UserRepository
}

func NewUserService(ur UserRepository) *UserServ {
	return &UserServ{userRepo: ur}
}

func (us *UserServ) CreateUser(req dto.UserRequest) (res dto.UserResponse, err error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error encrypt password")
		return dto.UserResponse{}, err
	}

	req.Password = string(passHash)

	user := entity.Users{
		Name: req.Name,
		Email: req.Email,
		Password: req.Password,
	}

	if err := us.userRepo.Create(&user); err != nil {
		return dto.UserResponse{}, err
	}

	//get user id to show in the response
	userInfo, err := us.GetUserById(user.Id)
	if err != nil {
		log.Println("failed get user by id")
		return dto.UserResponse{}, err
	}

	return userInfo, nil
}

func (us *UserServ) GetUserById(id int) (res dto.UserResponse, err error) {
	user, err := us.userRepo.GetById(id)
	if err != nil {
		log.Println("failed get user by id")
		return dto.UserResponse{}, err
	}

	userInfo := dto.UserResponse{
		Id: user.Id,
		Name: user.Name,
		Email: user.Email,
		Role: user.Role,
	}

	return userInfo, nil
}

func (us *UserServ) GetUserByEmail(email, password string) (accessToken string, err error) {
	user, err := us.userRepo.GetByEmail(email) 
	if err != nil {
		log.Println("failed get user by email", err.Error())
		return "", err
	}		

	// WIP auth/validation
	//validation??
	//mailjet validation kalo ada 

	//compare hash pass and input pass
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println("failed to compare password hash")
		return "", err
	}

	token, err := utils.GenerateJwtToken(email, user.Role, user.Id)
	if err != nil {
		log.Println("failed to generate jwt token")
		return "", err
	}

	return token, nil
}