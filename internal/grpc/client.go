package grpc

import (
	"context"
	"log"

	"backend_institutions/internal/loggerpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	LoggerClient loggerpb.LoggerServiceClient
	Conn         *grpc.ClientConn
)

func ConnectLogger() error {

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	Conn = conn
	LoggerClient = loggerpb.NewLoggerServiceClient(conn)

	log.Println("Connected to Logger Service on port 50051")

	return nil
}

func SendLog(
	service string,
	method string,
	endpoint string,
	request string,
	response string,
	status int32,
) error {

	_, err := LoggerClient.SaveLog(
		context.Background(),
		&loggerpb.LogRequest{
			Service:  service,
			Method:   method,
			Endpoint: endpoint,
			Request:  request,
			Response: response,
			Status:   status,
		},
	)

	if err != nil {
		return err
	}

	return nil
}