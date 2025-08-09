package auth

import(
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"fmt"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	expiration := now.Add(expiresIn) 
	claims := jwt.RegisteredClaims{
		Issuer:	"chirpy",
		IssuedAt:	jwt.NewNumericDate(now),
		ExpiresAt:	jwt.NewNumericDate(expiration),
		Subject:	userID.String(),	
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJWT, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedJWT, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := new(jwt.RegisteredClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token)(any, error){
		return []byte(tokenSecret), nil})
	if err != nil {
		return uuid.Nil, err
	}
	tokenClaims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("token not RegisteredClaims")
	}
	tokenID, err := tokenClaims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	finalID, err := uuid.Parse(tokenID)
	if err != nil {
		return uuid.Nil, err
	}
	return finalID, nil
}