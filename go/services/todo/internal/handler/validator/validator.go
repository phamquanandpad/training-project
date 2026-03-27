package validator

import (
	"bytes"
	_ "embed"
	"time"

	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/locales/ja"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/errors"
)

type ValidUserAttributes struct {
	UserId int64 `validate:"gt=0"`
}

//go:embed embed/ja.json
var jaMessages []byte

// TODO Please reconsider whether to disable 'cyclop'.
//
//nolint:cyclop
func InitValidator() (*validator.Validate, error) {
	uni := ut.New(en_US.New(), ja.New())
	translator, _ := uni.GetTranslator("ja")

	if err := uni.ImportByReader(ut.FormatJSON, bytes.NewReader(jaMessages)); err != nil {
		return nil, errors.NewInternalError("failed to init validator", err)
	}
	if err := uni.VerifyTranslations(); err != nil {
		return nil, errors.NewInternalError("failed to init validator", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := ja_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		return nil, errors.NewInternalError("failed to init validator", err)
	}

	if err := validate.RegisterValidation("datestring", func(fl validator.FieldLevel) bool {
		dateLayout := "2006-01-02"
		_, err := time.Parse(dateLayout, fl.Field().String())
		return err == nil
	}); err != nil {
		return nil, errors.NewInternalError("failed to init validator", err)
	}

	return validate, nil
}
