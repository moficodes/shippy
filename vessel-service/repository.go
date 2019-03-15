package main

import (
	pb "github.com/moficodes/shippy/vessel-service/proto/vessel"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	dbName           = "shippy"
	vesselCollection = "vessels"
)

// Repository interface
type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
	Create(vessel *pb.Vessel) error
	Close()
}

// VesselRepository implements Repository interface
type VesselRepository struct {
	session *mgo.Session
}

func (repo *VesselRepository) collection() *mgo.Collection {
	return repo.session.DB(dbName).C(vesselCollection)
}

// FindAvailable queries the db for vessel
// that is greater than or equal to the spec passed to the function
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	var vessel *pb.Vessel

	err := repo.collection().Find(bson.M{
		"capacity":  bson.M{"&gte": spec.Capacity},
		"maxweight": bson.M{"&gte": spec.MaxWeight},
	}).One(&vessel)
	if err != nil {
		return nil, err
	}
	return vessel, nil
}

// Create a new record of vessel in db
func (repo *VesselRepository) Create(vessel *pb.Vessel) error {
	return repo.collection().Insert(vessel)
}

// Close the session
func (repo *VesselRepository) Close() {
	repo.session.Close()
}
