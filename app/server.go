package app

import "github.com/orestrepov/metadatahost/model"

func (ctx *Context) GetServerById(id uint) (*model.Server, error) {

	server, err := ctx.Database.GetServerById(id)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (ctx *Context) GetServersByHostId(hostId uint) ([]*model.Server, error) {

	return ctx.Database.GetServersByHostId(hostId)
}

func (ctx *Context) CreateServer(server *model.Server) (int, error) {

	return ctx.Database.CreateServer(server)
}

func (ctx *Context) UpdateServer(server *model.Server) (int, error) {

	if server.ID == 0 {
		return 0, &ValidationError{"cannot update"}
	}

	return ctx.Database.UpdateServer(server)
}

func (ctx *Context) DeleteServerById(id uint) (int, error) {

	_, err := ctx.GetServerById(id)
	if err != nil {
		return 0, err
	}

	return ctx.Database.DeleteServerById(id)
}
