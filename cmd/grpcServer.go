package main

import (
	"context"
	"log"
	"net"
	"tritchgo/internal/store"

	"tritchgo/proto/stream_stats"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type StatsServer struct {
	stream_stats.UnimplementedStreamStatsServiceServer
	db    *pgxpool.Pool
	store *store.Storage
}

func NewStatsServer(ctx context.Context, db *pgxpool.Pool) *StatsServer {
	statsStore := store.NewStorage(db)
	return &StatsServer{db: db, store: &statsStore}
}

func (s *StatsServer) GetUserStats(ctx context.Context, req *stream_stats.UserStatsRequest) (*stream_stats.UserStatsResponse, error) {
	stats, err := s.store.Stats.GetUserStatsById(ctx, req.UserId)
	if err != nil {
		log.Printf("Err fetch user stats  %v", err)
		return nil, err
	}

	var protoStats []*stream_stats.StreamStats
	for _, stat := range stats {
		protoStats = append(protoStats, &stream_stats.StreamStats{
			Id:             stat.ID.String(),
			StreamId:       stat.StreamID,
			UserId:         stat.UserID,
			GameId:         stat.GameID,
			Date:           stat.Date.Format("2006-01-02"),
			Airtime:        int32(stat.Airtime),
			PeakViewers:    int32(stat.PeakViewers),
			AverageViewers: int32(stat.AverageViewers),
			HoursWatched:   int32(stat.HoursWatched),
		})
	}

	return &stream_stats.UserStatsResponse{Stats: protoStats}, nil
}

func (s *StatsServer) GetStreamStats(ctx context.Context, req *stream_stats.StreamStatsRequest) (*stream_stats.StreamStatsResponse, error) {
	stats, err := s.store.Stats.GetStreamStatsById(ctx, req.StreamId)
	if err != nil {
		return nil, err
	}
	var protoStats []*stream_stats.StreamStats
	for _, stat := range stats {
		protoStats = append(protoStats, &stream_stats.StreamStats{
			Id:             stat.ID.String(),
			StreamId:       stat.StreamID,
			UserId:         stat.UserID,
			GameId:         stat.GameID,
			Date:           stat.Date.Format("2006-01-02"),
			Airtime:        int32(stat.Airtime),
			PeakViewers:    int32(stat.PeakViewers),
			AverageViewers: int32(stat.AverageViewers),
			HoursWatched:   int32(stat.HoursWatched),
		})

	}
	return &stream_stats.StreamStatsResponse{Stats: protoStats}, nil
}

func StartGRPCServer(db *pgxpool.Pool) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("err listen")
	}

	grpcServer := grpc.NewServer()
	stream_stats.RegisterStreamStatsServiceServer(grpcServer, NewStatsServer(context.Background(), db))
	reflection.Register(grpcServer)

	log.Println("gRPC is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
