package service

import (
	"github.com/knowhunger/ortoo/commons/log"
	"github.com/knowhunger/ortoo/server/mongodb"
)

type OrtooService struct {
	mongo *mongodb.RepositoryMongo
}

func NewOrtooService(mongoConf *mongodb.Config) (*OrtooService, error) {
	mongo, err := mongodb.New(mongoConf)
	if err != nil {
		return nil, log.OrtooError(err, "fail to connect to MongoDB")
	}
	return &OrtooService{
		mongo: mongo,
	}, nil
}

func (o *OrtooService) Initialize() error {
	if err := o.mongo.InitializeCollections(); err != nil {
		return log.OrtooError(err, "fail to initialize mongoDB")
	}
	return nil
}
