package utils

import (
    "errors"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID uint64 `json:"user_id"`
    Email  string `json:"email"`
    Role   int8   `json:"role"`
    jwt.RegisteredClaims
}

var jwtSecret = []byte("your-secret-key-change-in-production")

func GenerateToken(userID uint64, email string, role int8) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "gin-auth-demo",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}