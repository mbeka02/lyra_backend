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
	client, err := getstream.NewClient(apiKey, apiSecret, getstream.WithTimeout(10000*time.Millisecond))
	if err != nil {
		return nil, err
	}
	return &StreamClient{client}, nil
}

type CreateStreamUserParams struct {
	Name   string
	Email  string
	UserID int64
}

func (s *StreamClient) CreateUser(ctx context.Context, params CreateStreamUserParams) error {
	userIdString := fmt.Sprintf("%d", params.UserID)
	_, err := s.client.UpdateUsers(ctx, &getstream.UpdateUsersRequest{
		Users: map[string]getstream.UserRequest{
			userIdString: {
				ID:   userIdString,
				Name: getstream.PtrTo(params.Name),
			},
		},
	})
	return err
}

func (s *StreamClient) CreateToken(userIdString string) (string, error) {
	return s.client.CreateToken(userIdString, getstream.WithExpiration(time.Hour*72))
}
