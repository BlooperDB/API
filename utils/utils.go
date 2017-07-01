package utils

import (
	"crypto/rand"
	"encoding/json"
	"net/http"

	"encoding/base64"
	"reflect"

	"errors"

	"bytes"
	"compress/zlib"
	"io"

	"crypto/sha256"
	"fmt"

	"gopkg.in/validator.v2"
)

var v *validator.Validator

func Initialize() {
	v = validator.NewValidator()
	v.SetTag("validate")
	v.SetValidationFunc("blueprint_string", validBlueprintString)
}

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

	if err = v.Validate(s); err != nil {
		return &ErrorResponse{
			Code:    Error_invalid_request_data.Code,
			Message: Error_invalid_request_data.Message + ": " + err.Error(),
			Status:  Error_invalid_request_data.Status,
		}
	}

	return nil
}

func validBlueprintString(v interface{}, _ string) error {
	s := reflect.ValueOf(v).String()

	decoded, err := base64.StdEncoding.DecodeString(s[1:])
	if err != nil {
		return errors.New("Not valid blueprint string")
	}

	b := bytes.NewReader(decoded)
	r, err := zlib.NewReader(b)

	if err != nil {
		return errors.New("Not valid blueprint string")
	}

	var out bytes.Buffer
	io.Copy(&out, r)
	r.Close()

	var js map[string]interface{}
	err = json.Unmarshal([]byte(out.String()), &js)

	if err != nil {
		return errors.New("Not valid blueprint string")
	}

	return nil
}

func SHA265(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
