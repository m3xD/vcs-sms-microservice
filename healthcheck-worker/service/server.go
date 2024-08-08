package service

import (
	"fmt"
	"healthcheck-worker/model"
	"healthcheck-worker/repo"
	"healthcheck-worker/util"
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
	// how to create server without field id

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

func (s *ServerService) UpdateServersOn(statusMapping map[string]interface{}) error {
	log := util.NewLogger()
	keys := make([]string, len(statusMapping))
	i := 0
	for k := range statusMapping {
		keys[i] = k
		i++
	}
	if err := s.db.Table("servers").Where("ipv4 IN ?", keys).Updates(map[string]interface{}{"status": 1}).Error; err != nil {
		log.Error(fmt.Sprintf("Error updating servers: %v", err))
		return err
	}
	if err := s.db.Table("servers").Where("ipv4 NOT IN ?", keys).Updates(map[string]interface{}{"status": 0}).Error; err != nil {
		log.Error(fmt.Sprintf("Error updating servers: %v", err))
		return err
	}
	return nil
}
