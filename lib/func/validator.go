package lib

import (
	"errors"
	constLib "go-gateway/lib/const"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)


func ValidateParams(c *gin.Context, params interface{}) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	vtor, err := GetValidator(c)

	if err != nil {
		return err
	}

	trans, err := GetTranslation(c)

	if err != nil {
		return err
	}

	err = vtor.Struct(params)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		sliceErrs := []string{}
		for _, e := range errs {
			sliceErrs = append(sliceErrs, e.Translate(trans))
		}
		return errors.New(strings.Join(sliceErrs, ",\n"))
	}

	return nil
}

func GetValidator(c *gin.Context) (*validator.Validate, error) {
	val, ok := c.Get(constLib.ValidatorKey)

	if !ok {
		return nil, errors.New("not config validator")
	}

	validator, ok := val.(*validator.Validate)

	if !ok {
		return nil, errors.New("get validator failed")
	}

	return validator, nil
}

func GetTranslation(c *gin.Context) (ut.Translator, error) {
	trans, ok := c.Get(constLib.TranslatorKey)
	if !ok {
		return nil, errors.New("not config translator")
	}

	translator, ok := trans.(ut.Translator)

	if !ok {
		return nil, errors.New("")
	}
	return translator, nil
}