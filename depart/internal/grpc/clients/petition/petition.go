package petition

import (
	"context"
	"database/sql"

	pb "depart/internal/grpc/clients/petition/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type PetitionGrpcClient struct {
	grpcClient pb.PetitionClient
	conn       *grpc.ClientConn
}

func NewPetitionGrpcClient(ctx context.Context, serviceUrl string) (*PetitionGrpcClient, error) {
	conn, err := grpc.NewClient(
		serviceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	grpcClient := pb.NewPetitionClient(conn)

	return &PetitionGrpcClient{
		grpcClient: grpcClient,
		conn:       conn}, nil
}

func (s *PetitionGrpcClient) UpdatePetition(ctx context.Context, petition_id uuid.UUID,
	done_at sql.NullTime, report_id uuid.UUID, content_job string) (bool, error) {
	var DoneAt *timestamppb.Timestamp
	if done_at.Valid {
		DoneAt = timestamppb.New(done_at.Time)
	}
	req := &pb.UpdatePetitionRequest{
		PetitionId: petition_id.String(),
		DoneAt:     DoneAt,
		ReportId:   wrapperspb.String(report_id.String()),
		ContentJob: wrapperspb.String(content_job),
	}
	_, err := s.grpcClient.UpdatePetition(ctx, req)
	if err != nil {
		return false, err
	}
	return true, nil
}
