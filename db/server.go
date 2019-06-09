package db

import (
	"github.com/orestrepov/metadatahost/model"
	"github.com/pkg/errors"
	"time"
)

func (db *Database) GetServerById(id uint) (*model.Server, error) {
	var server model.Server
	row := db.QueryRow(`SELECT id, address, ssl_grade, country, owner, host_id, updated_at 
								FROM servers 
								WHERE id = $1`, id)
	if err := row.Scan(&server.ID,
		&server.Address,
		&server.SslGrade,
		&server.Country,
		&server.Owner,
		&server.HostID,
		&server.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "unable to get server by id")
	}
	return &server, nil
}

func (db *Database) GetServersByHostId(hostId uint) ([]*model.Server, error) {
	var servers []*model.Server
	rows, err := db.Query(`SELECT id, address, ssl_grade, country, owner, host_id, updated_at 
									FROM servers 
									WHERE host_id = $1`, hostId)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get servers by host_id")
	}
	defer rows.Close()
	for rows.Next() {
		var server model.Server
		if err := rows.Scan(&server.ID,
			&server.Address,
			&server.SslGrade,
			&server.Country,
			&server.Owner,
			&server.HostID,
			&server.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "unable to get servers by host_id")
		}
		servers = append(servers, &server)
	}
	return servers, nil
}

func (db *Database) GetServerByAddress(address string) (*model.Server, error) {
	var server model.Server
	row := db.QueryRow(`SELECT id, address, ssl_grade, country, owner, host_id, updated_at 
								FROM servers 
								WHERE address = $1`, address)
	if err := row.Scan(&server.ID,
		&server.Address,
		&server.SslGrade,
		&server.Country,
		&server.Owner,
		&server.HostID,
		&server.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "unable to get server by address")
	}
	return &server, nil
}

func (db *Database) CreateServer(server *model.Server) (int, error) {
	stmt, err := db.Prepare(
		`INSERT INTO 
    			servers (address, ssl_grade, country, owner, host_id, created_at, updated_at) 
    			VALUES ($1, $2, $3, $4, $5, $6, $7)
    			RETURNING id`)
	if err != nil {
		return 0, errors.Wrap(err, "unable to prepare create server")
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(server.Address, server.SslGrade, server.Country,
		server.Owner, server.HostID, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "unable to run insert server")
	}
	return id, nil
}

func (db *Database) UpdateServer(server *model.Server) (int, error) {
	stmt, err := db.Prepare(
		`UPDATE servers SET 
                   (address, ssl_grade, country, owner, host_id, updated_at) 
                    = ($1, $2, $3, $4, $5, $6) 
                    WHERE id = $7
    				RETURNING id`)
	if err != nil {
		return 0, errors.Wrap(err, "unable to prepare update server")
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(
		server.Address, server.SslGrade, server.Country, server.Owner,
		server.HostID, time.Now(), server.ID).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "unable to run update server")
	}
	return id, nil
}

func (db *Database) DeleteServerById(id uint) (int, error) {
	// TODO: implement logical delete
	stmt, err := db.Prepare(
		`DELETE FROM servers WHERE id = $1 RETURNING id`)
	if err != nil {
		return 0, errors.Wrap(err, "unable to prepare delete server")
	}
	defer stmt.Close()
	var deletedId int
	err = stmt.QueryRow(id).Scan(&deletedId)
	if err != nil {
		return 0, errors.Wrap(err, "unable to run delete server")
	}
	return deletedId, nil
}
