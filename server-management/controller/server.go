package api

import (
	"fmt"
	"os"
	db "server-management/db/sqlc"
	"server-management/enum"
	"server-management/middleware"
	"server-management/service"
	"server-management/token"
	"server-management/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Config util.Config
	Store  db.Store
	Token  token.Maker
	Router *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWT(os.Getenv("SECRET_KEY"))
	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker %w", err)
	}
	server := &Server{Store: store}
	router := gin.Default()
	server.Config = config
	server.Token = tokenMaker

	userCtrl := &UserController{
		store: store,
		token: tokenMaker,
	}

	scopeCtrl := &ScopeController{
		store: store,
	}

	reportService := service.NewReportService(service.NewElastic(), service.NewRedis(), store, service.NewMailService(), service.NewKafkaService(util.GetConsumer(), util.GetProducer()))
	report := NewReport(reportService)

	router.POST("/users/login", userCtrl.loginUser)

	authRoutes := router.Group("/").Use(middleware.AuthMiddleWare(server.Token))

	router.POST("/users", userCtrl.createUser)

	authRoutes.POST("/create_server", middleware.CheckScope(enum.API_SERVER_WRITE), server.createServer)
	authRoutes.GET("/get_server", middleware.CheckScope(enum.API_SERVER_READ), server.getServer)
	authRoutes.PUT("/update_server", middleware.CheckScope(enum.API_SERVER_WRITE), server.updateServer)
	authRoutes.DELETE("/delete_server/:id", middleware.CheckScope(enum.API_SERVER_WRITE), server.deleteServer)
	authRoutes.POST("/import_excel", middleware.CheckScope(enum.API_SERVER_WRITE), server.importExcel)
	authRoutes.POST("/export_excel", middleware.CheckScope(enum.API_SERVER_WRITE), server.exportExcel)
	authRoutes.POST("/send_report", middleware.CheckScope(enum.API_REPORT_READ), report.SendReport)

	authRoutes.POST("/update_role", middleware.CheckScope(enum.API_USER_WRITE), userCtrl.updateRole)

	authRoutes.POST("/get_scope", middleware.CheckScope(enum.API_USER_READ), scopeCtrl.getScope)
	authRoutes.POST("/create_scope", middleware.CheckScope(enum.API_USER_WRITE), scopeCtrl.createScope)
	authRoutes.POST("/delete_scope", middleware.CheckScope(enum.API_USER_WRITE), scopeCtrl.deleteScope)
	authRoutes.POST("/update_scope", middleware.CheckScope(enum.API_USER_WRITE), scopeCtrl.updateScope)
	server.Router = router
	if err != nil {
		return nil, fmt.Errorf("error when send every day task %w", err)
	}
	return server, nil
}

func (server *Server) Start() error {
	return server.Router.Run(":8081")
}

func errorResponse(e error) gin.H {
	return gin.H{"error": e.Error()}
}
