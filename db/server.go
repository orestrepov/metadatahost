package db

import (
	"github.com/orestrepov/metadatahost/model"
	"github.com/pkg/errors"
)

func (db *Database) GetServerById(id uint) (*model.Server, error) {
	var server model.Server
	return &server, errors.Wrap(db.First(&server, id).Error, "unable to get server")
}

func (db *Database) GetServersByHostId(hostId uint) ([]*model.Server, error) {
	var servers []*model.Server
	return servers, errors.Wrap(db.Find(&servers, model.Server{HostID: hostId}).Error, "unable to get servers")
}

func (db *Database) GetServerByAddress(address string) (*model.Server, error) {
	var server model.Server
	return &server, errors.Wrap(db.First(&server, "address = ?", address).Error, "unable to get server by address")
}

func (db *Database) CreateServer(server *model.Server) error {
	return errors.Wrap(db.Create(server).Error, "unable to create server")
}

func (db *Database) UpdateServer(server *model.Server) error {
	return errors.Wrap(db.Save(server).Error, "unable to update server")
}

func (db *Database) DeleteServerById(id uint) error {
	return errors.Wrap(db.Delete(&model.Server{}, id).Error, "unable to create server")
}
