package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func WalletAddress(fl validator.FieldLevel) bool {
	walletAddress := fl.Field().String()
	ok, _ := regexp.MatchString(`^T\w{33,}`, walletAddress)
	return ok
}
