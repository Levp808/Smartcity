package orchestrator

import (
	"context"

	"petition/internal/grpc/clients/orchestrator/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrchestratorGrpcClient struct {
	grpcClient pb.OrchestratorServiceClient
	conn       *grpc.ClientConn
}

func NewOrchestratorGrpcClient(ctx context.Context, serviceUrl string) (*OrchestratorGrpcClient, error) {
	conn, err := grpc.NewClient(
		serviceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	grpcClient := pb.NewOrchestratorServiceClient(conn)

	return &OrchestratorGrpcClient{
		grpcClient: grpcClient,
		conn:       conn}, nil
}

func (s *OrchestratorGrpcClient) GetDepartment(ctx context.Context, petition_id string,
	location string, description string) (string, error) {
	req := &pb.ModerationRequest{Id: petition_id, Location: location, Description: description}
	depart, err := s.grpcClient.Moderation(ctx, req)
	if err != nil {
		return "", err
	}
	return depart.Depart, nil
}
