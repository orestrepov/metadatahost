package app

import (
	"encoding/json"
	"github.com/orestrepov/metadatahost/model"
	"github.com/sirupsen/logrus"
)

func (ctx *Context) GetHostById(id uint) (*model.Host, error) {

	host, err := ctx.Database.GetHostById(id)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func (ctx *Context) GetHostByName(name string) (*model.Host, error) {

	host, err := ctx.Database.GetHostByName(name)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func (ctx *Context) CreateHost(host *model.Host) error {

	return ctx.Database.CreateHost(host)
}

func (ctx *Context) UpdateHost(host *model.Host) error {

	if host.ID == 0 {
		return &ValidationError{"cannot update"}
	}

	return ctx.Database.UpdateHost(host)
}

func (ctx *Context) DeleteHostById(id uint) error {

	_, err := ctx.GetHostById(id)
	if err != nil {
		return err
	}

	return ctx.Database.DeleteHostById(id)
}

// List the previously searched domains
func (ctx *Context) DomainSearchHistory() (*model.HistoryResponse, error) {

	// Build the domain host data response
	historyResponse, err := BuildDomainHistoryResponse(ctx)
	if err != nil {
		return nil, err
	}

	return historyResponse, nil
}

// Get the configuration data of a domain
func (ctx *Context) SearchDomain(name string) (*model.HostResponse, error) {

	// parse the domain name
	name = getDomain(name)

	// Get of domain data/configuration
	domainData, err := GetDomainData(ctx, name)
	if err != nil {
		return nil, err
	}

	// Process of domain data recovery from SSLLabs
	host, err := ProcessDomainData(ctx, name, domainData)
	if err != nil {
		return nil, err
	}

	// Build the domain host data response
	hostResponse, err := BuildDomainSearchResponse(ctx, host)
	if err != nil {
		return nil, err
	}

	return hostResponse, nil
}

// Get the host data from SSL Labs
func GetDomainData(ctx *Context, name string) ([]byte, error) {

	APIURL := ctx.Config.SSLLabsAPIURL + "/api/" + ctx.Config.SSLLabsAPIVersion + "/analyze?host=" + name
	data, err := HTTPClient(APIURL)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Process and save the host data getting from domain resources
func ProcessDomainData(ctx *Context, domain string, data []byte) (*model.Host, error) {

	var sslHost model.SSLHost
	var host model.Host
	var lowestGrade model.Grade

	// get and validate if the domain exist in db
	previousHost, _ := ctx.Database.GetHostByName(domain)
	if previousHost.ID != 0 {
		host = *previousHost
		logrus.Infof("The domain %s exist, id: %d", domain, previousHost.ID)
	}

	// we unmarshal our byteArray which contains our jsonFile's content into 'sshHost'
	json.Unmarshal(data, &sslHost)

	// set the host domain data
	host.Name = sslHost.Name
	host.ServersChanged = false
	// scraping the web page to get the title and logo
	webPage, errWebPage := ScrapingWebPage(domain)
	if errWebPage != nil {
		host.IsDown = true
	} else {
		host.Logo = webPage.Logo
		host.Title = webPage.Title
		host.IsDown = false
	}

	// if the domain not exist in db create it
	if previousHost.ID == 0 {
		ctx.Database.CreateHost(&host)
		logrus.Infof("The domain %s was created, id: %d", domain, host.ID)
	}

	// process and save the domain servers
	ProcessServersData(ctx, host, sslHost.Endpoints, &lowestGrade)

	if previousHost.ID != 0 {
		// validate if the servers changed
		if previousHost.SslGrade != lowestGrade.Name {
			host.ServersChanged = true
			// set the previous grade
			host.PreviousSslGrade = previousHost.SslGrade
		}
	}

	// update host data
	if len(lowestGrade.Name) > 0 {
		host.SslGrade = lowestGrade.Name
	}
	logrus.Infof("The domain %s was updated, id: %d", host.Name, host.ID)
	err := ctx.Database.UpdateHost(&host)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

// Process and save the servers data
func ProcessServersData(ctx *Context, host model.Host, endpoints []model.SSLEndpoint, lowestGrade *model.Grade) {

	// get Endpoints data
	for i := 0; i < len(endpoints); i++ {
		var serverDB model.Server
		// validate if exist the server
		previousServer, _ := ctx.Database.GetServerByAddress(endpoints[i].IPAddress)
		if previousServer.ID != 0 {
			serverDB = *previousServer
		}
		// get/set the whois data
		whoisModel, err := WhoisQuery(endpoints[i].IPAddress)
		if err == nil {
			serverDB.Country = whoisModel.Country
			serverDB.Owner = whoisModel.OrgName
		}
		// set other server data
		serverDB.Address = endpoints[i].IPAddress
		if len(endpoints[i].Grade) > 0 {
			serverDB.SslGrade = endpoints[i].Grade
		}
		serverDB.HostID = host.ID
		serverDB.Host = host
		// validate and set the lowest ssl grade between servers
		if i == 0 {
			lowestGrade.Score = model.Grades[serverDB.SslGrade]
			lowestGrade.Name = serverDB.SslGrade
		}
		if model.Grades[serverDB.SslGrade] > lowestGrade.Score {
			lowestGrade.Score = model.Grades[serverDB.SslGrade]
			lowestGrade.Name = serverDB.SslGrade
		}
		// save server in db
		if previousServer.ID == 0 {
			ctx.Database.CreateServer(&serverDB)
			logrus.Infof("The server %s was created, id: %d", serverDB.Address, serverDB.ID)
		} else {
			serverDB.ID = previousServer.ID
			ctx.Database.UpdateServer(&serverDB)
			logrus.Infof("The server %s exist and it was updated, id: %d", serverDB.Address, serverDB.ID)
		}
	}
}

// Built the domain response json object
func BuildDomainSearchResponse(ctx *Context, host *model.Host) (*model.HostResponse, error) {

	var hostResponse model.HostResponse

	// get the servers associated with the domain
	servers, err := ctx.Database.GetServersByHostId(host.ID)
	serversResponse := make([]model.ServerResponse, len(servers))
	for i := 0; i < len(servers); i++ {
		var serverResponse model.ServerResponse
		serverResponse.SslGrade = servers[i].SslGrade
		serverResponse.Owner = servers[i].Owner
		serverResponse.Country = servers[i].Country
		serverResponse.Address = servers[i].Address
		serversResponse[i] = serverResponse
	}
	// set host domain values
	hostResponse.Name = host.Name
	hostResponse.SslGrade = host.SslGrade
	hostResponse.ServersChanged = host.ServersChanged
	hostResponse.PreviousSslGrade = host.PreviousSslGrade
	hostResponse.IsDown = host.IsDown
	hostResponse.Title = host.Title
	hostResponse.Logo = host.Logo
	hostResponse.TotalServers = len(servers)
	hostResponse.LastSearch = host.UpdatedAt.String()
	hostResponse.Servers = serversResponse

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &hostResponse, nil
}

// Built the domain history json object
func BuildDomainHistoryResponse(ctx *Context) (*model.HistoryResponse, error) {

	// get the hosts history
	hosts, err := ctx.Database.GetHosts()
	hostsResponse := make([]model.HostResponse, len(hosts))
	for i := 0; i < len(hosts); i++ {
		var hostResponse model.HostResponse
		hostResponse.Name = hosts[i].Name
		hostResponse.SslGrade = hosts[i].SslGrade
		hostResponse.PreviousSslGrade = hosts[i].PreviousSslGrade
		hostResponse.Title = hosts[i].Title
		hostResponse.Logo = hosts[i].Logo
		hostResponse.IsDown = hosts[i].IsDown
		hostResponse.ServersChanged = hosts[i].ServersChanged

		servers, _ := ctx.Database.GetServersByHostId(hosts[i].ID)
		hostResponse.TotalServers = len(servers)
		hostResponse.LastSearch = hosts[i].UpdatedAt.String()

		hostsResponse[i] = hostResponse
	}

	var historyResponse model.HistoryResponse
	historyResponse.Items = hostsResponse

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &historyResponse, nil
}
