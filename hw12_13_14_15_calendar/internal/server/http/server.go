package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	ginzap "github.com/akath19/gin-zap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/baserepo"
	"go.uber.org/zap"
)

const (
	varUUIDName       = "uuid"
	statusOk          = http.StatusOK
	statusServerError = http.StatusInternalServerError
)

var calenderApp *app.Calendar

type Server struct {
	server *http.Server
}

func NewServer(calendar *app.Calendar, listenAddress string) *Server {
	calenderApp = calendar

	server := &http.Server{
		Addr:    listenAddress,
		Handler: CreateHandler(),
	}
	return &Server{server: server}
}

func (srv Server) Start() error {
	err := srv.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (srv Server) Stop() error {
	if err := srv.server.Shutdown(context.Background()); err != nil {
		return err
	}

	return nil
}

func CreateHandler() http.Handler {
	router := gin.Default()

	zapL := zap.L()
	router.Use(ginzap.Logger(3*time.Second, zapL))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello world!",
		})
	})
	router.POST("/add", Add)
	router.PUT("/update/:id", Update)
	router.DELETE("/delete/:id", Delete)
	router.GET("/get/:id", GetByID)
	router.GET("/getlist", GetList)

	return router
}

func Add(context *gin.Context) {
	jsonData, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		log("read data error", err)
		context.Writer.WriteHeader(statusServerError)
		return
	}
	event := baserepo.Event{}
	err = json.Unmarshal(jsonData, &event)

	if err != nil {
		log("json error", err)
		context.Writer.WriteHeader(statusServerError)
	} else {
		id, err := calenderApp.Add(event)

		if err != nil {
			log("adding event error", err)
			context.Writer.WriteHeader(statusServerError)
		} else {
			context.JSON(statusOk, gin.H{varUUIDName: id.String()})
		}
	}
}

func Update(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		log("uuid error", err)
		context.Writer.WriteHeader(statusServerError)
		return
	}

	jsonData, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		log("read data error", err)
		context.Writer.WriteHeader(statusServerError)
		return
	}
	event := baserepo.Event{}
	err = json.Unmarshal(jsonData, &event)

	if err != nil {
		log("json error", err)
		context.Writer.WriteHeader(statusServerError)
	} else {
		err = calenderApp.Update(event, id)

		if err != nil {
			log("updating event error", err)
			context.Writer.WriteHeader(statusServerError)
		} else {
			context.Writer.WriteHeader(statusOk)
		}
	}
}

func Delete(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		log("uuid error", err)
		context.Writer.WriteHeader(statusServerError)
	} else {
		err = calenderApp.Delete(id)

		if err != nil {
			log("deleting event error", err)
			context.Writer.WriteHeader(statusServerError)
		} else {
			context.Writer.WriteHeader(statusOk)
		}
	}
}

func GetByID(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		log("uuid error", err)
		context.Writer.WriteHeader(statusServerError)
		return
	}

	event, err := calenderApp.GetEventByID(id)
	if err != nil {
		log("event not found", err)
		context.Writer.WriteHeader(statusServerError)
	} else {
		eventJSON, err := json.Marshal(event)
		if err != nil {
			log("event convering to json error", err)
			context.Writer.WriteHeader(statusServerError)
		} else {
			context.JSON(statusOk, eventJSON)
		}
	}
}

func GetList(context *gin.Context) {
	events, err := calenderApp.GetList()
	if err != nil {
		log("events not found", err)
		context.Writer.WriteHeader(statusServerError)
		return
	}

	eventJSON, err := json.Marshal(events)
	if err != nil {
		log("events convering to json error", err)
		context.Writer.WriteHeader(statusServerError)
	} else {
		context.JSON(statusOk, eventJSON)
	}
}

func log(message string, err error) {
	zap.L().Error(message, zap.Error(err))
}
