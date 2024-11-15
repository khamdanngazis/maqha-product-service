package repository

import (
	"context"
	"errors"
	"maqhaa/library/logging"
	"maqhaa/library/middleware"
	"maqhaa/product_service/external/model"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type UserRepository interface {
	GetUser(ctx context.Context, token string) (*model.UserData, error)
}

type userRepository struct {
	connetionURl string
}

func NewUserRepository(connetionURl string) UserRepository {
	return &userRepository{connetionURl: connetionURl}
}

func (r *userRepository) GetUser(ctx context.Context, token string) (*model.UserData, error) {
	requestID, _ := ctx.Value(middleware.RequestIDKey).(string)
	req := &model.GetUserRequest{
		Token: token, // Replace with a valid product ID for your test data
	}
	conn, err := grpc.Dial(r.connetionURl, grpc.WithInsecure())
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error GetUser  %s", err.Error())
		return nil, err
	}
	defer conn.Close()

	client := model.NewUserClient(conn)

	resp, err := client.GetUser(context.Background(), req)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error GetUser  %s", err.Error())
		return nil, err
	}
	if resp.Code != 0 {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error GetUser  %v", resp)
		return nil, errors.New(resp.Message)
	}

	if resp.Data == nil {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error GetUser  %s", errors.New("Data User Nill"))
		return nil, errors.New("Data User Nill")
	}
	return resp.Data, nil
}
