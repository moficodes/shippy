package main

import (
	"context"
	"fmt"
	"log"

	"github.com/micro/go-micro"
	pb "github.com/moficodes/shippy/consignment-service/proto/consignment"
	vesselpb "github.com/moficodes/shippy/vessel-service/proto/vessel"
)

type Repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type ConsignmentRepository struct {
	consignments []*pb.Consignment
}

func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *ConsignmentRepository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo         Repository
	vesselClient vesselpb.VesselServiceClient
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselpb.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	log.Printf("found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}

	req.VesselId = vesselResponse.Vessel.Id

	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &ConsignmentRepository{}

	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	vesselClient := vesselpb.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	srv.Init()

	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
