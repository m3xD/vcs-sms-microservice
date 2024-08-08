package service

import (
	"fmt"
	"healthcheck-server/model"
	"healthcheck-server/repo"
	"healthcheck-server/util"
)

type ServerService struct {
	db repo.Database
}

func NewServerService(db repo.Database) *ServerService {
	return &ServerService{db: db}
}

func (s *ServerService) GetAllServers() []model.Server {
	logger := util.NewLogger()
	var servers []model.Server
	err := s.db.Find(&servers).Error
	if err != nil {
		logger.Error("error when get all server")
		return []model.Server{}
	}
	return servers
}

func (s *ServerService) CreateServer(server *model.Server) error {
	logger := util.NewLogger()
	err := s.db.Create(server).Error
	if err != nil {
		logger.Error("error when create server")
	}
	return err
}

func (s *ServerService) UpdateServer(server *model.Server) error {
	logger := util.NewLogger()
	err := s.db.Save(server).Error
	if err != nil {
		logger.Error("error when update server")
	}
	return err
}

func (s *ServerService) GetServerByID(id string) interface{} {
	logger := util.NewLogger()
	server := &model.Server{}
	err := s.db.First(server, id).Error
	if err != nil {
		logger.Error(fmt.Sprintf("error when get server by id: %v", id))
		return nil
	}
	return server
}

func (s *ServerService) GetServerByIP(ip string) interface{} {
	logger := util.NewLogger()
	server := &model.Server{}
	err := s.db.Where("ipv4 = ?", ip).First(server).Error
	if err != nil {
		logger.Error(fmt.Sprintf("error when get server by ip: %v", ip))
		return nil
	}
	return server
}
