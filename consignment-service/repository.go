package main

import (
	pb "github.com/moficodes/shippy/consignment-service/proto/consignment"
	"gopkg.in/mgo.v2"
)

const (
	dbName                = "shippy"
	consignmentCollection = "consignments"
)

// Repository interface
// ConsignmentRepository implements it
type Repository interface {
	Create(*pb.Consignment) error
	GetAll() ([]*pb.Consignment, error)
	Close()
}

// ConsignmentRepository holds the session
type ConsignmentRepository struct {
	session *mgo.Session
}

// unexported method collection
// gives back a collection from the session
// with given dbname and collectionname
func (repo *ConsignmentRepository) collection() *mgo.Collection {
	return repo.session.DB(dbName).C(consignmentCollection)
}

// Create a new consignment entry on the db
func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) error {
	return repo.collection().Insert(consignment)
}

// GetAll consignments from the db
func (repo *ConsignmentRepository) GetAll() ([]*pb.Consignment, error) {
	var consignments []*pb.Consignment

	err := repo.collection().Find(nil).All(&consignments)
	return consignments, err
}

// Close the session connection
func (repo *ConsignmentRepository) Close() {
	repo.session.Close()
}
