package streamsdk

import (
	"context"
	"fmt"
	"time"

	"github.com/GetStream/getstream-go"
)

type StreamClient struct {
	client *getstream.Stream
}

func NewStreamClient(apiKey, apiSecret string) (*StreamClient, error) {
	client, err := getstream.NewClient(apiKey, apiSecret)
	if err != nil {
		return nil, err
	}
	return &StreamClient{client}, nil
}

type CreateUserParams struct {
	Role  string
	Name  string
	Email string
	Id    int64
}

func (s *StreamClient) CreateUser(ctx context.Context, params CreateUserParams) error {
	customId := fmt.Sprintf("%s_%d", params.Email, params.Id)
	_, err := s.client.UpdateUsers(ctx, &getstream.UpdateUsersRequest{
		Users: map[string]getstream.UserRequest{
			customId: {
				ID:   customId,
				Role: getstream.PtrTo(params.Role),
				Name: getstream.PtrTo(params.Name),
			},
		},
	})
	return err
}

func (s *StreamClient) CreateToken(customId string) (string, error) {
	return s.client.CreateToken(customId, getstream.WithExpiration(time.Hour*72))
}
