package auth

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func GetIdPasswordFromRequest(r *http.Request) (id int32, password string, err error) {
	full_header := r.Header.Get("Authorization")
	auth_array := strings.Split(full_header, " ")
	if len(auth_array) < 2 {
		return 0, "", errors.New("out of bounds access error")
	}
	decoded_output, err := base64.StdEncoding.DecodeString(auth_array[1])
	if err != nil {
		return 0, "", errors.New("failed to base64 decode")
	}
	output_string := strings.Split(string(decoded_output), ":")
	if len(auth_array) != 2 {
		return 0, "", errors.New("incorrect base64 formatted string")
	}
	username, err := strconv.Atoi(output_string[0])
	if err != nil {
		return 0, "", err
	}
	password = output_string[1]
	return int32(username), password, nil
}
