package http

import (
	"context"
	"github.com/google/uuid"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/baserepo"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
	"github.com/akath19/gin-zap"
	"github.com/gin-gonic/gin"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/app"
)

const varUuidName = "uuid"
const statusOk = http.StatusOK
const statusServerError = http.StatusInternalServerError

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
	router.Use(ginzap.Logger(3 * time.Second, zapL))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello world!",
		})
	})
	router.POST("/add", Add)
	router.PUT("/update/:id", Update)
	router.DELETE("/delete/:id", Delete)
	router.GET("/get/:id", GetById)
	router.GET("/getlist", GetList)


	return router
}

func Add(context *gin.Context) {
	jsonData, err := ioutil.ReadAll(context.Request.Body)

	if err != nil {
		log("read data error", err)
		context.Writer.WriteHeader(statusServerError)
	} else {
		event := baserepo.Event{}
		err = json.Unmarshal(jsonData, event)

		if err != nil {
			log("json error", err)
			context.Writer.WriteHeader(statusServerError)
		} else {
			id, err := calenderApp.Add(event)

			if err != nil {
				log("adding event error", err)
				context.Writer.WriteHeader(statusServerError)
			} else {
				context.JSON(statusOk, gin.H{varUuidName: id.String()})
			}
		}
	}
}

func Update(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		log("uuid error", err)
		context.Writer.WriteHeader(statusServerError)
	} else {
		jsonData, err := ioutil.ReadAll(context.Request.Body)

		if err != nil {
			log("read data error", err)
			context.Writer.WriteHeader(statusServerError)
		} else {
			event := baserepo.Event{}
			err = json.Unmarshal(jsonData, event)

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

func GetById(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		log("uuid error", err)
		context.Writer.WriteHeader(statusServerError)
	} else {
		event, err := calenderApp.GetEventByID(id)

		if err != nil {
			log("event not found", err)
			context.Writer.WriteHeader(statusServerError)
		} else {
			eventJson, err := json.Marshal(event)
			if err != nil {
				log("event convering to json error", err)
				context.Writer.WriteHeader(statusServerError)
			} else {
				context.JSON(statusOk, eventJson)
			}
		}
	}
}

func GetList(context *gin.Context) {
	events, err := calenderApp.GetList()

	if err != nil {
		log("events not found", err)
		context.Writer.WriteHeader(statusServerError)
	} else {
		eventJson, err := json.Marshal(events)
		if err != nil {
			log("events convering to json error", err)
			context.Writer.WriteHeader(statusServerError)
		} else {
			context.JSON(statusOk, eventJson)
		}
	}
}

func log(message string, err error) {
	zap.L().Error(message, zap.Error(err))
}