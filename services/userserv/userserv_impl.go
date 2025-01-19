package userserv

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ARF-DEV/image-processing-api/model"
	"github.com/ARF-DEV/image-processing-api/repos/userrepo"
	"github.com/ARF-DEV/image-processing-api/utils/httputils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type UserServImpl struct {
	userRepo userrepo.UserRepo
}

func New(userRepo userrepo.UserRepo) UserServ {
	return &UserServImpl{userRepo: userRepo}
}

func (s *UserServImpl) Login(ctx context.Context, user model.User) (model.AutheticationResponse, error) {
	userSrc, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.AutheticationResponse{}, httputils.ErrUnauthorized
		}
		return model.AutheticationResponse{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userSrc.Password), []byte(user.Password)); err != nil {
		log.Println("error when comparing password: ", err)
		return model.AutheticationResponse{}, httputils.ErrUnauthorized
	}
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(viper.GetString("SECRET_KEY")))
	if err != nil {
		return model.AutheticationResponse{}, err
	}
	return model.AutheticationResponse{
		AccessToken: tokenStr,
	}, nil
}

func (s *UserServImpl) Register(ctx context.Context, user model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	if err != nil {
		return err
	}
	userSrc, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if userSrc.Email == user.Email {
		return fmt.Errorf("email already registered")
	}
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return err
	}
	return nil
}
