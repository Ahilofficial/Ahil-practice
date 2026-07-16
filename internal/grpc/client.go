package grpc

import (
	"context"
	"log"
	"os"

	"backend_institutions/internal/loggerpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	LoggerClient loggerpb.LoggerServiceClient
	Conn         *grpc.ClientConn
)

func ConnectLogger() error {
	port := os.Getenv("LOGGER_GRPC_PORT")
	if port == "" {
		port = "15051"
	}


	conn, err := grpc.NewClient(
		"localhost:"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	Conn = conn
	LoggerClient = loggerpb.NewLoggerServiceClient(conn)

	log.Printf("Connected to Logger Service on port %s\n", port)

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