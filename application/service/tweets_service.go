package service

import (
	"context"
	"fmt"
	"github.com/marcoshuck/jaeger-demo/api"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TweetsV1 interface {
	api.TweetsServiceServer
}

type tweetsV1 struct {
	api.UnimplementedTweetsServiceServer
	users api.UsersServiceClient
}

func (t *tweetsV1) Get(ctx context.Context, id *api.TweetID) (*api.Tweet, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(attribute.Key("tweet_id").String(id.GetId()))

	span.AddEvent("Getting user")
	user, err := t.users.Get(ctx, &api.UserID{Id: "test"})
	if err != nil {
		err = fmt.Errorf("tweets.get: %w", err)
		span.RecordError(err)
		return nil, err
	}
	span.AddEvent("Got user")

	return &api.Tweet{
		CreatedAt: timestamppb.Now(),
		Id:        "test",
		Text:      "This is a tweet!",
		User:      user,
	}, nil
}

func NewTweetsV1(users api.UsersServiceClient) TweetsV1 {
	return &tweetsV1{
		users: users,
	}
}
