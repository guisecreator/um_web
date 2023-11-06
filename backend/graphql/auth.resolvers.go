package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"
	"errors"
	"fmt"
	"github.com/guisecreator/um_web/db/dbmodels"
	"github.com/guisecreator/um_web/graphql/model"
	"github.com/guisecreator/um_web/internal/authpayload"
	"github.com/guisecreator/um_web/internal/sessions"
	"github.com/guisecreator/um_web/internal/token"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// Signup is the resolver for the signup field.
func (r *mutationResolver) Signup(ctx context.Context, email string, password string) (*model.AuthPayload, error) {
	panic(fmt.Errorf("not implemented: Signup - signup"))
}

// Login is the resolver for the login field.
//func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.AuthPayload, error) {
//	const incorrectMessage = "incorrect login and/or password"
//
//	dbUser := dbmodels.User{}
//
//	log.Printf("Попытка входа с email: %s и паролем: %s", email, password)
//
//	errScan := r.Db.
//		NewSelect().
//		Model(&dbUser).
//		Where("email = ?", email).
//		Scan(ctx)
//	if errScan != nil {
//		log.Println(errScan)
//	}
//
//	create_at := time.Now().Format(time.RFC3339)
//	update_at := time.Now().Format(time.RFC3339)
//
//	if errScan != nil {
//		log.Println(errScan.Error())
//		if dbUser == (dbmodels.User{}) {
//			return nil, errors.New("err scan: " + incorrectMessage)
//		}
//		return nil, errScan
//	}
//
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dbUser.Password), bcrypt.DefaultCost)
//	if err != nil {
//		log.Println(err)
//		return nil, errors.New("password error")
//	}
//
//	errHash := bcrypt.CompareHashAndPassword(
//		hashedPassword,
//		[]byte(password),
//	)
//	if errHash != nil {
//		log.Println(errHash)
//		return nil, errors.New("compare hashed password error")
//	}
//
//	generateToken, errToken := token.GenerateToken(dbUser.ID)
//	if errToken != nil {
//		return nil, errToken
//	}
//
//	cookieValue := generateToken
//	fmt.Printf("Cookie value: %s\n", cookieValue)
//	r.Sessions.AddSession(
//		cookieValue, &sessions.Session{
//			Id:    dbUser.ID,
//			Login: model.User{Email: dbUser.Email},
//			Roles: model.Roles(dbUser.Roles),
//		},
//	)
//
//	authInfo := &model.AuthInfo{Token: cookieValue}
//	expiration := time.Now().Add(365 * 24 * time.Hour)
//
//	cookie := http.Cookie{
//		Name:     "auth_cookie",
//		Value:    cookieValue,
//		Path:     "/",
//		SameSite: http.SameSiteLaxMode,
//		Domain:   "localhost",
//		HttpOnly: true,
//		Secure:   true,
//		Expires:  expiration,
//	}
//	s := authpayload.ForContext(ctx)
//
//	http.SetCookie(s.ResponseWriter, &cookie)
//	fmt.Printf("Set cookie %v=%s\n", cookie.Name, cookieValue)
//
//	id := dbUser.ID
//	return &model.AuthPayload{
//		User: &model.User{
//			ID:       id,
//			Email:    dbUser.Email,
//			CreateAt: create_at,
//			UpdateAt: update_at,
//			Role:     model.Roles(dbUser.Roles),
//		},
//		Info: authInfo,
//	}, nil
//}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.AuthPayload, error) {
	const incorrectMessage = "incorrect login and/or password"

	dbUser := dbmodels.User{}

	errScan := r.Db.
		NewSelect().
		Model(&dbUser).
		Where("email = ?", email).
		Scan(ctx)
	if errScan != nil {
		log.Println(errScan)
	}

	create_at := time.Now().Format(time.RFC3339)
	update_at := time.Now().Format(time.RFC3339)

	if errScan != nil {
		log.Println(errScan.Error())
		if dbUser == (dbmodels.User{}) {
			return nil, errors.New("err scan: " + incorrectMessage)
		}
		return nil, errScan
	}

	hashedPassword, err := bcrypt.
		GenerateFromPassword([]byte(
			dbUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return nil, errors.New("hashed password error")
	}

	errHash := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(password))
	if errHash != nil {
		log.Println(errHash)
		return nil, errors.New("compare hashed password error")
	}

	generateToken, errToken := token.GenerateToken(dbUser.ID)
	if errToken != nil {
		return nil, errToken
	}

	cookieValue := generateToken
	r.Sessions.AddSession(
		cookieValue, &sessions.Session{
			Id:    dbUser.ID,
			Login: model.User{Email: dbUser.Email},
			Roles: model.Roles(dbUser.Roles),
		},
	)

	authInfo := &model.AuthInfo{Token: cookieValue}
	expiration := time.Now().Add(365 * 24 * time.Hour)

	cookie := http.Cookie{
		Name:     "auth_cookie",
		Value:    cookieValue,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
		Secure:   true,
		Expires:  expiration,
	}
	s := authpayload.ForContext(ctx)

	http.SetCookie(s.ResponseWriter, &cookie)
	fmt.Printf("Set cookie %v=%s\n", cookie.Name, cookieValue)

	id := dbUser.ID
	return &model.AuthPayload{
		User: &model.User{
			ID:       id,
			Email:    dbUser.Email,
			CreateAt: create_at,
			UpdateAt: update_at,
			Role:     model.Roles(dbUser.Roles),
		},
		Info: authInfo,
	}, nil
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context, email string) (*string, error) {
	user := model.User{}

	db := r.Db.NewSelect().Model(&user).Scan(ctx)
	if db != nil {
		return nil, db
	}

	tokenFromContext, err := authpayload.
		GetSessionTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	r.Sessions.RemoveSession(tokenFromContext)

	user = model.User{
		Email: email,
	}

	return &tokenFromContext, nil
}
