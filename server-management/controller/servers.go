package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	db "server-management/db/sqlc"
	"server-management/service"
	"server-management/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createServerRequest struct {
	Name   string `json:"name" binding:"required"`
	Status int32  `json:"status" binding:"required"`
	Ipv4   string `json:"ipv4" binding:"required"`
}

func (server *Server) createServer(ctx *gin.Context) {
	log := util.NewLogger()
	var request createServerRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Error("error when bind json: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateServerParams{
		Name:   request.Name,
		Status: request.Status,
		Ipv4:   request.Ipv4,
	}

	s, e := server.Store.CreateServer(ctx, arg)

	if e != nil {
		log.Error("error when create server: " + e.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(e))
		return
	}
	log.Info(fmt.Sprintf("create server successfully: %v", s))
	ctx.JSON(http.StatusOK, s)
}

type getServerRequest struct {
	PageID        int32  `json:"page_id" binding:"required,min=1"`
	PageSize      int32  `json:"page_size" binding:"required,min=1"`
	IsID          bool   `json:"isID"`
	ID            int32  `json:"id"`
	IsName        bool   `json:"isName"`
	Name          string `json:"name"`
	IsStatus      bool   `json:"isStatus"`
	Status        int32  `json:"status"`
	IsIpv4        bool   `json:"isIpv4"`
	Ipv4          string `json:"ipv4"`
	IDAsc         bool   `json:"isIDAsc"`
	IDDesc        bool   `json:"isIDDesc"`
	NameAsc       bool   `json:"isNameAsc"`
	NameDesc      bool   `json:"isNameDesc"`
	StatusAsc     bool   `json:"isStatusAsc"`
	StatusDesc    bool   `json:"isStatusDesc"`
	Ipv4Asc       bool   `json:"isIpv4Asc"`
	Ipv4Desc      bool   `json:"isIpv4Desc"`
	CreatedAtAsc  bool   `json:"isCreatedAtAsc"`
	CreatedAtDesc bool   `json:"isCreatedAtDesc"`
	UpdatedAtAsc  bool   `json:"isUpdatedAtAsc"`
	UpdatedAtDesc bool   `json:"isUpdatedAtDesc"`
}

func (server *Server) getServer(ctx *gin.Context) {
	log := util.NewLogger()
	var request getServerRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Error("error when bind json: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetServerParams{
		Limit:         request.PageSize,
		Offset:        (request.PageID - 1) * request.PageSize,
		IsID:          request.IsID,
		ID:            request.ID,
		IsName:        request.IsName,
		Name:          request.Name,
		IsStatus:      request.IsStatus,
		Status:        request.Status,
		IsIpv4:        request.IsIpv4,
		Ipv4:          request.Ipv4,
		IDAsc:         request.IDAsc,
		IDDesc:        request.IDDesc,
		NameAsc:       request.NameAsc,
		NameDesc:      request.NameDesc,
		StatusAsc:     request.StatusAsc,
		StatusDesc:    request.StatusAsc,
		Ipv4Asc:       request.Ipv4Asc,
		Ipv4Desc:      request.Ipv4Desc,
		CreatedAtAsc:  request.CreatedAtAsc,
		CreatedAtDesc: request.CreatedAtDesc,
		UpdatedAtAsc:  request.UpdatedAtAsc,
		UpdatedAtDesc: request.UpdatedAtDesc,
	}

	if arg.IsID && arg.ID == 0 {
		log.Error("error when get server: invalid ID")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid ID"})
		return
	} else if arg.IsIpv4 && arg.Ipv4 == "" {
		log.Error("error when get server: invalid ipv4")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid ipv4"})
		return
	} else if arg.IsName && arg.Name == "" {
		log.Error("error when get server: invalid name")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid name"})
		return
	}

	s, err := server.Store.GetServer(ctx, arg)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("error when get server: " + err.Error())
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.Error("error when get server: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Info(fmt.Sprintf("get server successfully: %v", s))
	ctx.JSON(http.StatusOK, s)
}

type updateServerRequest struct {
	ID        int64  `json:"id" binding:"required"`
	SetName   bool   `json:"setName"`
	Name      string `json:"name"`
	SetStatus bool   `json:"setStatus"`
	Status    int32  `json:"status"`
	SetIpv4   bool   `json:"setIpv4"`
	Ipv4      string `json:"ipv4"`
}

func (server *Server) updateServer(ctx *gin.Context) {
	var request updateServerRequest
	log := util.NewLogger()
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Error("error when bind json: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdateServerParams{
		ID:        request.ID,
		SetName:   request.SetName,
		Name:      request.Name,
		SetStatus: request.SetStatus,
		Status:    request.Status,
		SetIpv4:   request.SetIpv4,
		Ipv4:      request.Ipv4,
	}

	if args.SetName && args.Name == "" {
		log.Error("error when update server: invalid name")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid name"})
		return
	} else if args.SetIpv4 && args.Ipv4 == "" {
		log.Error("error when update server: invalid ipv4")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid ipv4"})
		return
	}

	s, err := server.Store.UpdateServer(ctx, args)

	if err != nil {
		log.Error("error when update server: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Info(fmt.Sprintf("update server successfully: %v", s))
	ctx.JSON(http.StatusOK, s)

}

type deleteServerRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteServer(ctx *gin.Context) {
	var request deleteServerRequest
	log := util.NewLogger()

	if err := ctx.ShouldBindUri(&request); err != nil {
		log.Error("error when bind uri: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.Store.DeleteServer(ctx, request.ID)

	if err != nil {
		log.Error("error when delete server: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Info(fmt.Sprintf("delete server successfully: %v", request.ID))
	ctx.JSON(http.StatusOK, "deleted")
}

type importExcelRequest struct {
	Path string `json:"path"`
}

func (server *Server) importExcel(ctx *gin.Context) {
	var request importExcelRequest
	log := util.NewLogger()
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Error("error when bind json: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		ctx.Abort()
		return
	}

	i, err := service.ImportServer(request.Path)
	if err != nil {
		log.Error("error when import server: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	servers, err := server.Store.GetAllServers(ctx)

	markedName := make(map[string]bool)
	markedID := make(map[int64]bool)

	for index := 0; index < len(servers); index++ {
		markedID[servers[index].ID] = true
		markedName[servers[index].Name] = true
	}

	if err != nil {
		log.Error("error when get all server: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var success []string
	var fail []string
	s := 0
	f := 0
	for index, data := range i {
		if index == 0 {
			continue
		}
		id := data[0]
		name := data[1]
		actualID, _ := strconv.Atoi(id)
		st, _ := strconv.Atoi(data[3])
		if markedID[int64(actualID)] || markedName[name] {
			fail = append(fail, name)
			f++
		} else {
			success = append(success, name)
			s++
			_, err := server.Store.CreateServer(ctx, db.CreateServerParams{
				Name:   name,
				Status: int32(st),
				Ipv4:   data[2],
			})
			if err != nil {
				log.Error("error when create server: " + err.Error())
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
	}
	log.Info(fmt.Sprintf("import server successfully: %v", s))
	ctx.JSON(http.StatusOK, gin.H{
		"message": "import successfully",
		"success": gin.H{
			"count":        s,
			"success_name": success,
		},
		"fail": gin.H{
			"count":     f,
			"fail_name": fail,
		},
	})
}

func (server *Server) exportExcel(ctx *gin.Context) {
	log := util.NewLogger()
	var request getServerRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Error("error when bind json: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetServerParams{
		Limit:         request.PageSize,
		Offset:        (request.PageID - 1) * request.PageSize,
		IsID:          request.IsID,
		ID:            request.ID,
		IsName:        request.IsName,
		Name:          request.Name,
		IsStatus:      request.IsStatus,
		Status:        request.Status,
		IsIpv4:        request.IsIpv4,
		Ipv4:          request.Ipv4,
		IDAsc:         request.IDAsc,
		IDDesc:        request.IDDesc,
		NameAsc:       request.NameAsc,
		NameDesc:      request.NameDesc,
		StatusAsc:     request.StatusAsc,
		StatusDesc:    request.StatusAsc,
		Ipv4Asc:       request.Ipv4Asc,
		Ipv4Desc:      request.Ipv4Desc,
		CreatedAtAsc:  request.CreatedAtAsc,
		CreatedAtDesc: request.CreatedAtDesc,
		UpdatedAtAsc:  request.UpdatedAtAsc,
		UpdatedAtDesc: request.UpdatedAtDesc,
	}

	if arg.IsID && arg.ID == 0 {
		log.Error("error when get server: invalid ID")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid ID"})
		return
	} else if arg.IsIpv4 && arg.Ipv4 == "" {
		log.Error("error when get server: invalid ipv4")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid ipv4"})
		return
	} else if arg.IsName && arg.Name == "" {
		log.Error("error when get server: invalid name")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid name"})
		return
	}

	s, err := server.Store.GetServer(ctx, arg)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("error when get server: " + err.Error())
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.Error("error when get server: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	path, _ := service.ExportServer(s)
	log.Info(fmt.Sprintf("export server successfully: %v", path))
	ctx.JSON(http.StatusOK, path)
}
