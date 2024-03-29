package main

import (
	"context"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	pb "github.com/moficodes/shippy/user-service/proto/user"
)

type service struct {
	repo         Repository
	tokenService Authable
}

/*
	Create(context.Context, *User, *Response) error
	Get(context.Context, *User, *Response) error
	GetAll(context.Context, *Request, *Response) error
	Auth(context.Context, *User, *Token) error
	ValidateToken(context.Context, *Token, *Token) error
*/

func (srv *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := srv.repo.Get(req.Id)
	if err != nil {
		return err
	}
	res.User = user
	return nil
}

func (srv *service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := srv.repo.GetAll()
	if err != nil {
		return err
	}
	res.Users = users
	return nil
}

func (srv *service) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	log.Println("logging in with:", req.Email)
	user, err := srv.repo.GetByEmail(req.Email)
	log.Println(user)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := srv.tokenService.Encode(user)

	if err != nil {
		return err
	}

	res.Token = token
	return nil
}

func (srv *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	claims, err := srv.tokenService.Decode(req.Token)
	if err != nil {
		return err
	}
	log.Println(claims)
	if claims.User.Id == "" {
		return errors.New("invalid user")
	}
	res.Valid = true
	return nil
}

func (srv *service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req.Password = string(hashedPass)

	if err := srv.repo.Create(req); err != nil {
		return err
	}
	res.User = req
	return nil
}
