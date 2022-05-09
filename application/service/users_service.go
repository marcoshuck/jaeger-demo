package service

import (
	"context"
	"errors"
	"github.com/marcoshuck/jaeger-demo/api"
	"go.opentelemetry.io/otel/trace"
)

type UsersV1 interface {
	api.UsersServiceServer
}

type usersV1 struct {
	api.UnimplementedUsersServiceServer
	tracer trace.Tracer
}

func (u *usersV1) Get(ctx context.Context, id *api.UserID) (*api.User, error) {
	ctx, span := u.tracer.Start(ctx, "get")
	defer span.End()

	if id.GetId() != "test" {
		return nil, errors.New("user not found")
	}
	return &api.User{
		Id:          "test",
		Name:        "Marcos Huck",
		Location:    "Paraná, ER",
		Url:         "huck.com.ar",
		Description: "This is my user",
		Verified:    false,
	}, nil
}

func NewUsersV1(tracer trace.Tracer) UsersV1 {
	return &usersV1{
		tracer: tracer,
	}
}
