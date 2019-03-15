package main

import (
	"fmt"
	"log"
	"os"

	pb "github.com/moficodes/shippy/consignment-service/proto/consignment"
	vesselpb "github.com/moficodes/shippy/vessel-service/proto/vessel"

	"github.com/micro/go-micro"
)

// type service struct {
// 	repo         Repository
// 	vesselClient vesselpb.VesselServiceClient
// }

// func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
// 	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselpb.Specification{
// 		MaxWeight: req.Weight,
// 		Capacity:  int32(len(req.Containers)),
// 	})
// 	log.Printf("found vessel: %s \n", vesselResponse.Vessel.Name)
// 	if err != nil {
// 		return err
// 	}

// 	req.VesselId = vesselResponse.Vessel.Id

// 	consignment, err := s.repo.Create(req)
// 	if err != nil {
// 		return err
// 	}
// 	res.Created = true
// 	res.Consignment = consignment
// 	return nil
// }

// func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
// 	consignments := s.repo.GetAll()
// 	res.Consignments = consignments
// 	return nil
// }

const (
	defaultHost = "localhost:27017"
)

func main() {
	// repo := &ConsignmentRepository{}

	// srv := micro.NewService(
	// 	micro.Name("go.micro.srv.consignment"),
	// 	micro.Version("latest"),
	// )

	// vesselClient := vesselpb.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	// srv.Init()

	// pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})
	// if err := srv.Run(); err != nil {
	// 	fmt.Println(err)
	// }
	host := os.Getenv("DB_HOST")

	if host == "" {
		host = defaultHost
	}

	session, err := CreateSession(host)
	defer session.Close()
	if err != nil {
		log.Panicf("could not connect to datastor with host %s - %v", host, err)
	}
	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	vesselClient := vesselpb.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	srv.Init()

	pb.RegisterShippingServiceHandler(srv.Server(), &service{session, vesselClient})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
