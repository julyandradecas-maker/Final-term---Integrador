package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"chat-system/backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Pin      string `json:"pin"`
	jwt.RegisteredClaims
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Pin      string `json:"pin"`
	IsActive bool   `json:"is_active"`
}

func VerifyJWT(tokenString string, cfg *config.Config) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritmo inesperado: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("no se pudieron extraer los claims del token")
	}

	return claims, nil
}

func VerifyWithFastAPI(tokenString string, cfg *config.Config) (*UserResponse, error) {
	url := fmt.Sprintf("%s/users/me", cfg.APIURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error llamando a FastAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("token rechazado")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FastAPI respondió status %d", resp.StatusCode)
	}

	var user UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("error leyendo respuesta de FastAPI: %w", err)
	}

	return &user, nil
}
