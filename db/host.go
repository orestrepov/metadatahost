package db

import (
	"github.com/orestrepov/metadatahost/model"
	"github.com/pkg/errors"
)

func (db *Database) GetHostById(id uint) (*model.Host, error) {
	var host model.Host
	return &host, errors.Wrap(db.First(&host, id).Error, "unable to get host")
}

func (db *Database) GetHostByName(name string) (*model.Host, error) {
	var host model.Host
	return &host, errors.Wrap(db.First(&host, "name = ?", name).Error, "unable to get host by name")
}

func (db *Database) GetHosts() ([]*model.Host, error) {
	var hosts []*model.Host
	return hosts, errors.Wrap(db.Find(&hosts).Error, "unable to get hosts")
}

func (db *Database) CreateHost(host *model.Host) error {
	return errors.Wrap(db.Create(host).Error, "unable to create host")
}

func (db *Database) UpdateHost(host *model.Host) error {
	return errors.Wrap(db.Save(host).Error, "unable to update host")
}

func (db *Database) DeleteHostById(id uint) error {
	return errors.Wrap(db.Delete(&model.Host{}, id).Error, "unable to create host")
}
