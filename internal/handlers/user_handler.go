package handlers

import (
	"context"

	userpb "github.com/HackMateGolang/user-service/api/proto/v1"
	"github.com/HackMateGolang/user-service/internal/models"
	"github.com/HackMateGolang/user-service/internal/service"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := h.service.ReadUser(ctx, &models.ReadUserRequest{Login: req.Login})
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserResponse{Login: user.Login, Username: user.Username, Description: user.Description}, nil
}

func (h *UserHandler) PostUser(ctx context.Context, req *userpb.PostUserRequest) (*userpb.PostUserResponse, error) {
	login, err := h.service.CreateUser(ctx, &models.CreateUserRequest{Login: req.Login, Username: req.Username})
	if err != nil {
		return nil, err
	}
	return &userpb.PostUserResponse{Login: login}, nil
}

func (h *UserHandler) PutUser(ctx context.Context, req *userpb.PutUserRequest) (*userpb.PutUserResponse, error) {
	ok, err := h.service.ReplaceUser(ctx, &models.UpdateUserRequest{
		Login:       req.Login,
		Username:    req.Username,
		FirstName:   req.FirstName,
		SecondName:  req.SecondName,
		Patronymic:  req.Patronymic,
		Stack:       mapTechs(req.Stack),
		Description: req.Description,
		Contacts:    mapSocials(req.Contacts),
		ShortDesc:   req.ShortDesc,
		Avatar:      req.Avatar,
	})

	if err != nil {
		return nil, err
	}

	return &userpb.PutUserResponse{Ok: ok}, nil
}

func (h *UserHandler) PatchUser(ctx context.Context, req *userpb.PatchUserRequest) (*userpb.PatchUserResponse, error) {
	ok, err := h.service.PatchUser(ctx, &models.PatchUserRequest{
		Login:       req.Login,
		Username:    req.Username,
		FirstName:   req.FirstName,
		SecondName:  req.SecondName,
		Patronymic:  req.Patronymic,
		Stack:       mapTechs(req.Stack),
		Description: req.Description,
		Contacts:    mapSocials(req.Contacts),
		ShortDesc:   req.ShortDesc,
		Avatar:      req.Avatar,
	})

	if err != nil {
		return nil, err
	}

	return &userpb.PatchUserResponse{Ok: ok}, nil
}

func mapTechs(techs []*userpb.Tech) []models.Tech {
	out := make([]models.Tech, 0, len(techs))
	for _, t := range techs {
		if t == nil {
			continue
		}
		out = append(out, models.Tech{
			Name:  t.Name,
			Level: t.Level,
		})
	}

	return out
}

func mapSocials(socials []*userpb.Social) []models.Social {
	out := make([]models.Social, 0, len(socials))
	for _, s := range socials {
		if s == nil {
			continue
		}
		out = append(out, models.Social{
			Type:  s.Type,
			Url: s.Url,
		})
	}

	return out
}