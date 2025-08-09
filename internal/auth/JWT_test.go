package auth

import(
	"testing"
	"github.com/google/uuid"
	"time"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"errors"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "I want to be a hacker"
	expiresIn, err := time.ParseDuration("1m")	
	if err != nil {
		fmt.Printf("Something is wrong with time")
	}
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("There was an error making the token: %v", err)
	}
	validatedID, err := ValidateJWT(tokenString, tokenSecret)
	if validatedID != userID || err != nil {
		t.Errorf("Validation failed! Got ID: %s, Expected ID: %s, Error: %v", validatedID, userID, err)
	}

}

func TestExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "I want to be a hacker"
	expiresIn, err := time.ParseDuration("0s")	
	if err != nil {
		fmt.Printf("Something is wrong with time")
	}
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("There was an error making the token: %v", err)
	}
	validatedID, err := ValidateJWT(tokenString, tokenSecret)
	if !errors.Is(err, jwt.ErrTokenExpired){
		t.Errorf("Validation failed! Got ID: %s, Expected ID: %s, Error: %v", validatedID, userID, err)
	}
}

func TestWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "I want to be a hacker"
	expiresIn, err := time.ParseDuration("1m")	
	if err != nil {
		fmt.Printf("Something is wrong with time")
	}
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("There was an error making the token: %v", err)
	}
	validatedID, err := ValidateJWT(tokenString, "I will never be a hacker")
	if !errors.Is(err, jwt.ErrTokenSignatureInvalid){
		t.Errorf("Validation failed! Got ID: %s, Expected ID: %s, Error: %v", validatedID, userID, err)
	}
}