package utils

import (
	"bufio"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"

	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/crypto/bcrypt"
)

var validVideoTypeRegEx regexp.Regexp = *regexp.MustCompile("^video/mp4$")

func ValidateFormField(fieldName string, field *multipart.Part, fieldErr error, isLastField bool) (string, int, error) {
	if fieldName != field.FormName() {
		return fmt.Sprintf("'%s' field not found", fieldName), http.StatusBadRequest, fmt.Errorf("Missing form field: %s", fieldName)
	} else if fieldErr != nil && fieldErr != io.EOF {
		return fmt.Sprintf("Err parsing formField: %s, err: %s", fieldName, fieldErr), http.StatusBadRequest, fmt.Errorf("Error parsing request's form field")
	} else if isLastField == false && fieldErr == io.EOF {
		return fmt.Sprintf("Incomplete form/data request, missing formfields"), http.StatusBadRequest, fmt.Errorf("Incomplete form/data request")
	}

	return "", 0, nil
}

func ValidVideoFileType(file *multipart.Part) (bool, io.Reader, *mimetype.MIME) {
	return ValidFileType(file, validVideoTypeRegEx)
}

func ValidFileType(file *multipart.Part, fileExtRegEx regexp.Regexp) (bool, io.Reader, *mimetype.MIME) {
	var buf *bufio.Reader = bufio.NewReader(file)
	var (
		sniff []byte
		err   error
	)
	sniff, err = buf.Peek(512)
	if err != nil {
		return false, nil, nil
	}

	var fileContentType = mimetype.Detect(sniff)
	if fileExtRegEx.MatchString(fileContentType.String()) {
		return true, io.MultiReader(buf, file), fileContentType
	}

	return false, nil, nil
}

func HashPassword(password string) (string, error) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPw), nil
}

func VerifyPassword(hashPw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPw), []byte(pw)) == nil
}
