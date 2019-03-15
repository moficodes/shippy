package main

import (
	"context"
	"log"

	pb "github.com/moficodes/shippy/consignment-service/proto/consignment"
	vesselpb "github.com/moficodes/shippy/vessel-service/proto/vessel"
	"gopkg.in/mgo.v2"
)

type service struct {
	session      *mgo.Session
	vesselClient vesselpb.VesselServiceClient
}

func (s *service) GetRepo() Repository {
	return &ConsignmentRepository{s.session.Clone()}
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	vesserResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselpb.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if err != nil {
		return err
	}
	log.Printf("found vessel: %s \n", vesserResponse.Vessel.Name)
	req.VesselId = vesserResponse.Vessel.Id

	err = s.GetRepo().Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = req
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	defer s.GetRepo().Close()
	consignments, err := s.GetRepo().GetAll()
	if err != nil {
		return err
	}
	res.Consignments = consignments
	return nil
}
