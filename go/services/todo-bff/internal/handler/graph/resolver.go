package graph

import (
	"regexp"

	validator_v10 "github.com/go-playground/validator/v10"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/handler/graph/generated"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/usecase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	userLogin    usecase.UserLogin
	userRegister usecase.UserRegister
	tokenRefresh usecase.TokenRefresh
	tokenVerify  usecase.TokenVerify

	todoGetter  usecase.TodoGetter
	todoLister  usecase.TodoLister
	todoCreator usecase.TodoCreator
	todoUpdater usecase.TodoUpdater
	todoDeleter usecase.TodoDeleter

	validate *validator_v10.Validate
}

func New(
	todoGetter usecase.TodoGetter,
	todoLister usecase.TodoLister,
	todoCreator usecase.TodoCreator,
	todoUpdater usecase.TodoUpdater,
	todoDeleter usecase.TodoDeleter,
) generated.Config {
	return generated.Config{
		Resolvers: &Resolver{
			todoGetter:  todoGetter,
			todoLister:  todoLister,
			todoCreator: todoCreator,
			todoUpdater: todoUpdater,
			todoDeleter: todoDeleter,
		},
	}
}

func NewAuth(
	userLogin usecase.UserLogin,
	userRegister usecase.UserRegister,
) generated.Config {
	validate := newValidate()
	return generated.Config{
		Resolvers: &Resolver{
			userLogin:    userLogin,
			userRegister: userRegister,

			validate: validate,
		},
	}
}

// nolint: errcheck
func newValidate() *validator_v10.Validate {
	validate := validator_v10.New(validator_v10.WithRequiredStructEnabled())
	validate.RegisterValidation("is-valid-password", passwordValidator)
	validate.RegisterValidation("is-valid-email", emailValidator)
	return validate
}

func passwordValidator(fl validator_v10.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	ruleCounter := 0
	regexPatterns := []string{
		"([A-Z]+)",
		"([a-z]+)",
		"([0-9]+)",
		"([!@#$%^&*]+)",
	}

	for _, rp := range regexPatterns {
		matched, err := regexp.MatchString(rp, password)
		if err != nil {
			panic(err)
		}
		if matched {
			ruleCounter += 1
		}
	}

	return ruleCounter >= 3
}

func emailValidator(fl validator_v10.FieldLevel) bool {
	email := fl.Field().String()
	regexPattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(regexPattern, email)
	if err != nil {
		panic(err)
	}
	return matched
}
