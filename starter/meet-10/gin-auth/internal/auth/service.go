package auth

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrInvalidToken       = errors.New("token tidak valid atau kedaluwarsa")
)

type Service struct {
	repo   *Repository
	secret []byte
}

func NewService(repo *Repository, secret []byte) *Service {
	return &Service{repo: repo, secret: secret}
}

func (s *Service) Register(ctx context.Context, name, email, password string) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user := User{Name: strings.TrimSpace(name), Email: strings.ToLower(email), PasswordHash: string(hash)}
	if err := s.repo.Create(ctx, &user); err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, strings.ToLower(email))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}
	now := time.Now()
	// ponytail: JWT berlaku satu jam; tambah revocation store saat butuh logout server-side.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "gin-auth",
		Subject:   strconv.FormatUint(uint64(user.ID), 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
	})
	return token.SignedString(s.secret)
}

func (s *Service) Authenticate(ctx context.Context, rawToken string) (User, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
		return s.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired(), jwt.WithIssuer("gin-auth"))
	if err != nil {
		return User{}, ErrInvalidToken
	}
	id, err := strconv.ParseUint(claims.Subject, 10, strconv.IntSize)
	if err != nil || id == 0 {
		return User{}, ErrInvalidToken
	}
	user, err := s.repo.GetByID(ctx, uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, ErrInvalidToken
	}
	return user, err
}
