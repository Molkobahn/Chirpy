package auth

import(
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	header := make(map[string][]string)
	header["Authorization"] = append(header["Authorization"], "Bearer Hello World")
	tokenString, err := GetBearerToken(header)
	if tokenString != "Hello World" || err != nil {
		t.Errorf("token doesn't match: %s Error: %v", tokenString, err)
	}

	
}