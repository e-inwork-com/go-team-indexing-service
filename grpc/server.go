// Copyright 2023, e-inwork.com. All rights reserved.

package grpc

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/e-inwork-com/go-team-indexing-service/internal/data"
	"github.com/e-inwork-com/go-team-indexing-service/internal/grpc/teams"
	"github.com/e-inwork-com/go-team-indexing-service/internal/jsonlog"
	"github.com/google/uuid"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

var (
	BuildTime string
	Version   string
)

type Config struct {
	Env string

	Db struct {
		Dsn         string
		MaxOpenConn int
		MaxIdleConn int
		MaxIdleTime string
	}

	GRPCPort string
	SolrURL  string
	SolrTeam string
}

type Application struct {
	Config  Config
	Logger  *jsonlog.Logger
	Models  data.Models
	Indexes data.Indexes
}

type TeamServer struct {
	teams.UnimplementedTeamServiceServer
	Indexes data.Indexes
	Models  data.Models
}

func OpenDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Db.MaxOpenConn)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConn)

	duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (t *TeamServer) WriteTeam(ctx context.Context, req *teams.TeamRequest) (*teams.TeamResponse, error) {
	input := req.GetTeamEntry()

	id, err := uuid.Parse(input.Id)
	if err != nil {
		return nil, err
	}

	team, err := t.Models.Teams.Get(id)
	if err != nil {
		res := &teams.TeamResponse{Result: "Failed"}
		return res, err
	}

	if team.IsIndexed {
		return &teams.TeamResponse{Result: "Indexed"}, nil
	}

	resp, err := t.Indexes.Teams.Update(team)
	if err != nil || resp.StatusCode != http.StatusOK {
		return &teams.TeamResponse{Result: "Failed"}, err
	}

	if !team.IsDeleted {
		err = t.Models.Teams.IsIndexedTrue(team)
		if err != nil {
			return &teams.TeamResponse{Result: "Failed"}, err
		}
	} else {
		err = t.Models.Teams.Delete(team)
		if err != nil {
			return &teams.TeamResponse{Result: "Failed"}, err
		}
	}

	return &teams.TeamResponse{Result: "Indexed"}, nil
}

func (app *Application) GRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", app.Config.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	teams.RegisterTeamServiceServer(s, &TeamServer{
		Indexes: app.Indexes,
		Models:  app.Models})

	log.Printf("gRPC Server started on port: %v", app.Config.GRPCPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
