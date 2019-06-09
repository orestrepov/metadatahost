package db

import (
	"github.com/orestrepov/metadatahost/model"
	"github.com/pkg/errors"
	"time"
)

func (db *Database) GetHostById(id uint) (*model.Host, error) {
	var host model.Host
	row := db.QueryRow(`SELECT id, name, servers_changed, ssl_grade, previous_ssl_grade, 
								logo, title, is_down, updated_at 
								FROM hosts 
								WHERE id = $1`, id)
	if err := row.Scan(&host.ID,
		&host.Name,
		&host.ServersChanged,
		&host.SslGrade,
		&host.PreviousSslGrade,
		&host.Logo,
		&host.Title,
		&host.IsDown,
		&host.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "unable to get host by id")
	}
	return &host, nil
}

func (db *Database) GetHostByName(name string) (*model.Host, error) {
	var host model.Host
	row := db.QueryRow(`SELECT id, name, servers_changed, ssl_grade, 
								previous_ssl_grade, logo, title, is_down, updated_at 
								FROM hosts WHERE name = $1`, name)
	if err := row.Scan(&host.ID,
		&host.Name,
		&host.ServersChanged,
		&host.SslGrade,
		&host.PreviousSslGrade,
		&host.Logo,
		&host.Title,
		&host.IsDown,
		&host.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "unable to get host by name")
	}
	return &host, nil
}

func (db *Database) GetHosts() ([]*model.Host, error) {
	var hosts []*model.Host
	rows, err := db.Query(`SELECT id, name, servers_changed, ssl_grade, previous_ssl_grade, 
									logo, title, is_down, updated_at 
									FROM hosts`)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get hosts")
	}
	defer rows.Close()
	for rows.Next() {
		var host model.Host
		if err := rows.Scan(&host.ID,
			&host.Name,
			&host.ServersChanged,
			&host.SslGrade,
			&host.PreviousSslGrade,
			&host.Logo,
			&host.Title,
			&host.IsDown,
			&host.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "unable to get hosts")
		}
		hosts = append(hosts, &host)
	}
	return hosts, nil
}

func (db *Database) CreateHost(host *model.Host) (int, error) {
	stmt, err := db.Prepare(
		`INSERT INTO 
    				hosts (name, servers_changed, ssl_grade, previous_ssl_grade, logo, title, is_down, created_at, updated_at) 
    				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    				RETURNING id`)
	if err != nil {
		return 0, errors.Wrap(err, "unable to prepare create host")
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(
		host.Name, host.ServersChanged, host.SslGrade, host.PreviousSslGrade,
		host.Logo, host.Title, host.IsDown, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "unable to run insert host")
	}
	return id, nil
}

func (db *Database) UpdateHost(host *model.Host) (int, error) {
	stmt, err := db.Prepare(
		`UPDATE hosts 
    				SET (name, servers_changed, ssl_grade, previous_ssl_grade, logo, title, is_down, updated_at) 
    				= ($1, $2, $3, $4, $5, $6, $7, $8)
    				WHERE id = $9
    				RETURNING id`)
	if err != nil {
		return 0, errors.Wrap(err, "unable to prepare update host")
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(
		host.Name, host.ServersChanged, host.SslGrade, host.PreviousSslGrade,
		host.Logo, host.Title, host.IsDown, time.Now(), host.ID).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "unable to run update host")
	}
	return id, nil
}

func (db *Database) DeleteHostById(id uint) (int, error) {
	// TODO: implement logical delete
	stmt, err := db.Prepare(
		`DELETE FROM hosts WHERE id = $1 RETURNING id`)
	if err != nil {
		return 0, errors.Wrap(err, "unable to prepare delete host")
	}
	defer stmt.Close()
	var deletedId int
	err = stmt.QueryRow(id).Scan(&deletedId)
	if err != nil {
		return 0, errors.Wrap(err, "unable to run delete host")
	}
	return deletedId, nil
}
