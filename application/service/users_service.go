package service

import (
	"context"
	"errors"
	"github.com/marcoshuck/jaeger-demo/api"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UsersV1 interface {
	api.UsersServiceServer
}

type usersV1 struct {
	api.UnimplementedUsersServiceServer
}

func (u *usersV1) Get(ctx context.Context, id *api.UserID) (*api.User, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(attribute.Key("user_id").String(id.GetId()))

	span.AddEvent("Checking user is valid")
	if id.GetId() != "test" {
		err := errors.New("user not found")
		span.RecordError(err)
		return nil, err
	}

	span.AddEvent("Returning user")
	return &api.User{
		Id:          "test",
		Name:        "Marcos Huck",
		Location:    "Paraná, ER",
		Url:         "huck.com.ar",
		Description: "This is my user",
		Verified:    false,
	}, nil
}

func NewUsersV1() UsersV1 {
	return &usersV1{}
}
