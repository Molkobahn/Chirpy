package auth

import(
	"fmt"
	"strings"
	"net/http"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("Header not found")
	}
	var tokenString string	
	for _, header := range authHeader {
		tokenString, _ = strings.CutPrefix(header, "Bearer ")
	}
	return tokenString, nil
}