package auth
import (
	"time"
	"os"
	"log"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
)

func getKey() []byte {
	envKeyStr := os.Getenv("HMAC_KEY")
	if envKeyStr == "" {
		log.Fatal("HMAC_KEY environment variable is not set")
	}
	secretKey, err := base64.StdEncoding.DecodeString(envKeyStr)
	if err != nil {
		log.Fatalf("Failed to decode Base64 string: %v", err)
	}
	return secretKey 
}

func CreateJWT(user string) string {
	claims := jwt.RegisteredClaims{
		Issuer:    	"go-web-app",
    	Audience: 	jwt.ClaimStrings{"billing-api"},
		IssuedAt:  	jwt.NewNumericDate(time.Now()),
		ExpiresAt: 	jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
		NotBefore: 	jwt.NewNumericDate(time.Now()),
		Subject:   	user,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	
	tokenString, err := token.SignedString(getKey())
	if err != nil {
		log.Fatal("Failed to sign token")
	}
	return tokenString
}

func VerifyJWT(cookie string) bool {
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (any, error) {
		return getKey(), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		log.Printf("ERROR: Failed to parse JWT: %v", err)
		return false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		aud := claims["Audience"]
		log.Printf("claim aud: %s", aud)

	}
	return true
}