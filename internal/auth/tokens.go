package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AccessClaims is the JWT payload for access tokens (sub, email, exp, iat).
type AccessClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// TokenIssuer creates signed JWT access tokens.
type TokenIssuer struct {
	secret []byte
	ttl    time.Duration
}

// NewTokenIssuer constructs an issuer using JWT_SECRET and access TTL.
func NewTokenIssuer(secret string, ttl time.Duration) *TokenIssuer {
	return &TokenIssuer{secret: []byte(secret), ttl: ttl}
}

// IssueAccess builds an HS256 JWT with sub=user id, email, exp, iat.
func (t *TokenIssuer) IssueAccess(userID, email string) (string, error) {
	now := time.Now().Unix()
	exp := time.Now().Add(t.ttl).Unix()
	claims := AccessClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(time.Unix(now, 0)),
			ExpiresAt: jwt.NewNumericDate(time.Unix(exp, 0)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.secret)
}

// ParseAccess validates an access token and returns user id and email.
func ParseAccess(secret []byte, tokenStr string) (userID, email string, err error) {
	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return "", "", err
	}
	if !token.Valid {
		return "", "", fmt.Errorf("invalid token")
	}
	if claims.Subject == "" || claims.Email == "" {
		return "", "", fmt.Errorf("missing claims")
	}
	return claims.Subject, claims.Email, nil
}

// NewRefreshOpaque generates an opaque refresh token and its SHA-256 hex hash for storage.
func NewRefreshOpaque() (raw string, hash string, err error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", "", err
	}
	raw = hex.EncodeToString(b[:])
	hash = HashRefreshToken(raw)
	return raw, hash, nil
}

// HashRefreshToken returns the SHA-256 hex digest of the raw refresh token string.
func HashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
