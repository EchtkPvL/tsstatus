package main

import (
	"sort"

	"github.com/multiplay/go-ts3"
	"github.com/pkg/errors"
)

type channel struct {
	ID      int            `json:"id"`
	Name    string         `json:"name"`
	Clients []onlineClient `json:"clients"`
}

type onlineClient struct {
	Away        bool   `json:"away"`
	AwayMessage string `json:"away_message"`
	Nickname    string `json:"nickname"`
}

type serverInfo struct {
	ClientsOnline      int    `json:"clients_online"`
	QueryClientsOnline int    `json:"query_clients_online"`
	MaxClients         int    `json:"max_clients"`
	ChannelsOnline     int    `json:"total_channels"`
	Name               string `json:"name"`
	Port               int    `json:"port"`
	Status             string `json:"status"`
	Uptime             int    `json:"uptime"`
	Version            string `json:"version"`
	Platform           string `json:"platform"`
}

func serverInfoFromServer(s *ts3.Server) serverInfo {
	return serverInfo{
		ClientsOnline:      s.ClientsOnline,
		QueryClientsOnline: s.QueryClientsOnline,
		MaxClients:         s.MaxClients,
		ChannelsOnline:     s.ChannelsOnline,
		Name:               s.Name,
		Port:               s.Port,
		Status:             s.Status,
		Uptime:             s.Uptime,
		Version:            s.Version,
		Platform:           s.Platform,
	}
}

type serverStats struct {
	Server   serverInfo `json:"server"`
	Channels []channel  `json:"channels"`
}

func getServerStats() (*serverStats, error) {
	var s = &serverStats{}

	client, err := ts3.NewClient(cfg.ServerAddress)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create client")
	}
	defer client.Close()

	if err = client.Login(cfg.QueryUser, cfg.QueryPass); err != nil {
		return nil, errors.Wrap(err, "Unable to login")
	}
	defer client.Logout()

	if err = client.Use(cfg.ServerID); err != nil {
		return nil, errors.Wrap(err, "Unable to select server")
	}

	var server *ts3.Server
	server, err = client.Server.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to query server")
	}
	s.Server = serverInfoFromServer(server)

	chans, err := client.Server.ChannelList()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to query channels")
	}

	sort.Slice(chans, func(i, j int) bool { return chans[i].ChannelOrder < chans[j].ChannelOrder })

	clients, err := client.Server.ClientList()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to query clients")
	}

	for _, c := range chans {
		ch := channel{
			ID:   c.ID,
			Name: c.ChannelName,
		}

		if c.NeededSubscribePower == 0 {
			for _, cli := range clients {
				if cli.Type == 1 {
					continue
				}

				if cli.ChannelID == c.ID {
					ch.Clients = append(ch.Clients, onlineClient{
						Away:        cli.Away,
						AwayMessage: cli.AwayMessage,
						Nickname:    cli.Nickname,
					})
				}
			}
		}

		if ch.Clients != nil {
			s.Channels = append(s.Channels, ch)
		}
	}

	return s, nil
}
