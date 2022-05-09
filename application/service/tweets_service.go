package service

import (
	"context"
	"fmt"
	"github.com/marcoshuck/jaeger-demo/api"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TweetsV1 interface {
	api.TweetsServiceServer
}

type tweetsV1 struct {
	api.UnimplementedTweetsServiceServer
	users  api.UsersServiceClient
	tracer trace.Tracer
}

func (t *tweetsV1) Get(ctx context.Context, id *api.TweetID) (*api.Tweet, error) {
	ctx, span := t.tracer.Start(ctx, "get")
	defer span.End()

	user, err := t.users.Get(ctx, &api.UserID{Id: "test"})
	if err != nil {
		return nil, fmt.Errorf("tweets.get: %w", err)
	}
	return &api.Tweet{
		CreatedAt: timestamppb.Now(),
		Id:        "test",
		Text:      "This is a tweet!",
		User:      user,
	}, nil
}

func NewTweetsV1(users api.UsersServiceClient, tracer trace.Tracer) TweetsV1 {
	return &tweetsV1{
		users:  users,
		tracer: tracer,
	}
}
