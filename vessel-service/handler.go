package main

import (
	"context"

	pb "github.com/moficodes/shippy/vessel-service/proto/vessel"
	"gopkg.in/mgo.v2"
)

// service implements handler from proto
type service struct {
	session *mgo.Session
}

// GetRepo returns a vessel repository with a clone of the session
func (s *service) GetRepo() Repository {
	return &VesselRepository{s.session.Clone()}
}

// FindAvailable rpc method from proto.
func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	defer s.GetRepo().Close()
	vessel, err := s.GetRepo().FindAvailable(req)
	if err != nil {
		return err
	}

	res.Vessel = vessel
	return nil
}

func (s *service) Create(ctx context.Context, req *pb.Vessel, res *pb.Response) error {
	defer s.GetRepo().Close()
	if err := s.GetRepo().Create(req); err != nil {
		return err
	}
	res.Vessel = req
	res.Created = true
	return nil
}
