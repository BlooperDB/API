package utils

import (
	"crypto/rand"
	"encoding/json"
	"net/http"

	"gopkg.in/validator.v2"
)

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func (tcf Block) Do() {
	if tcf.Finally != nil {

		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}

func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func GenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := GenerateRandomBytes(n)

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	return string(bytes)
}

func ValidateRequestBody(r *http.Request, s interface{}) *ErrorResponse {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(s)

	if err != nil {
		return &Error_invalid_request_data
	}

	v := validator.NewValidator()
	v.SetTag("validate")
	if err = v.Validate(s); err != nil {
		return &ErrorResponse{
			Code:    Error_invalid_request_data.Code,
			Message: Error_invalid_request_data.Message + ": " + err.Error(),
			Status:  Error_invalid_request_data.Status,
		}
	}

	return nil
}
