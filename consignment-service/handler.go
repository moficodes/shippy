package main

import (
	"context"
	"log"

	pb "github.com/moficodes/shippy/consignment-service/proto/consignment"
	vesselpb "github.com/moficodes/shippy/vessel-service/proto/vessel"
	"gopkg.in/mgo.v2"
)

// service holds the vesselClient and the session
// implememnts handler interface from consignment protobuf
type service struct {
	session      *mgo.Session
	vesselClient vesselpb.VesselServiceClient
}

// GetRepo returns a ConsignmentRepository which implements Repostitory
func (s *service) GetRepo() Repository {
	return &ConsignmentRepository{s.session.Clone()}
}

// CreateConsignment takes a Consignment and Response and input.
// Uses the vesselclient from service to get vessel information
// then insert the consignment to the db
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	s.GetRepo().Close()
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselpb.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if err != nil {
		return err
	}
	log.Printf("found vessel: %s \n", vesselResponse.Vessel.Name)
	req.VesselId = vesselResponse.Vessel.Id

	err = s.GetRepo().Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = req
	return nil
}

// GetConsignment queries the db for consignments
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	defer s.GetRepo().Close()
	consignments, err := s.GetRepo().GetAll()
	if err != nil {
		return err
	}
	res.Consignments = consignments
	return nil
}
