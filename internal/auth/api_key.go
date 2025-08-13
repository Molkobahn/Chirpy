package auth

import(
	"net/http"
	"strings"
	"fmt"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("Header not found")
	}
	var tokenString string
	for _, header := range authHeader {
		tokenString, _ = strings.CutPrefix(header, "ApiKey ")
	}
	return tokenString, nil
}