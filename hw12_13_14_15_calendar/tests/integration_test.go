package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cucumber/godog"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"github.com/widrik/hw/hw12_13_14_15_calendar/api/spec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CalendarClient spec.CalendarServiceClient

type CalendarTest struct {
	grpcClient       CalendarClient
	addEventResponse *spec.AddResponse
	getListResponse  *spec.GetListResponse
	responseBody     []byte
	responseStatus   codes.Code
}

func InitFeatureContext(ctx *godog.ScenarioContext) {
	calendar := &CalendarTest{}

	ctx.Step(`^I send addEvent request$`, calendar.iSendAddEventRequest)
	ctx.Step(`^I send addEvent request on date "([^"]*)"$`, calendar.iSendAddEventRequestForDate)
	ctx.Step(`^I send addEvent request on today$`, calendar.iSendAddEventRequestForToday)
	ctx.Step(`^I send updateEvent request$`, calendar.iSendUpdateEventRequest)
	ctx.Step(`^I send updateEvent request of "([^"]*)"$`, calendar.iSendUpdateEventRequestOf)
	ctx.Step(`^I send deleteEvent request$`, calendar.iSendDeleteEventRequest)
	ctx.Step(`^I send deleteEvent request for "([^"]*)"$`, calendar.iSendDeleteEventRequestOf)
	ctx.Step(`^I send getList request$`, calendar.iCallGetList)
	ctx.Step(`^I want to see event ID in response$`, calendar.wantToSeeInResponseEventID)
	ctx.Step(`^I want to see events response$`, calendar.wantToSeeEventsInResponseEventID)
	ctx.Step(`^I want receive event notification$`, calendar.iShouldReceiveEventNotification)
	ctx.Step(`^I have event on date "([^"]*)"$`, calendar.iHaveEventOnDate)
	ctx.Step(`^Response has error$`, calendar.iGetErrorResponse)
	ctx.Step(`^Response has NO errors$`, calendar.iGetSuccessResponse)

	ctx.BeforeScenario(calendar.initClient)
}

// Base methods
func (c *CalendarTest) initClient(*godog.Scenario) {
	endpoint := os.Getenv("CALENDAR_APP_ENDPOINT")
	if endpoint == "" {
		endpoint = ":8090"
	}

	connection, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to calendar service: %s", err.Error())
	}

	c.grpcClient = spec.NewCalendarServiceClient(connection)
}

// Send requests methods
func (calendar *CalendarTest) iSendAddEventRequest() error {
	var err error

	addRequest := &spec.AddRequest{
		Event: &spec.Event{
			Uuid:        uuid.New().String(),
			Title:       "Test title of event",
			Description: "Test description of event",
			Start: &timestamp.Timestamp{
				Seconds: 100,
				Nanos:   0,
			},
			Finish: &timestamp.Timestamp{
				Seconds: 1000,
				Nanos:   0,
			},
			UserId: 1,
			NotifyTime: &timestamp.Timestamp{
				Seconds: 100,
				Nanos:   0,
			},
		},
	}

	calendar.addEventResponse, err = calendar.grpcClient.Add(context.Background(), addRequest)
	calendar.responseStatus = status.Code(err)

	return err
}

func (calendar *CalendarTest) iSendAddEventRequestForDate(date string) error {
	var err error
	grpcTimestamp, err := ConvertDateToGrpcTimestamp(date)
	if err != nil {
		return err
	}

	addRequest := &spec.AddRequest{
		Event: &spec.Event{
			Uuid:        uuid.New().String(),
			Title:       "Test title of event",
			Description: "Test description of event",
			Start:       grpcTimestamp,
			Finish: &timestamp.Timestamp{
				Seconds: 1000,
				Nanos:   0,
			},
			UserId: 1,
			NotifyTime: &timestamp.Timestamp{
				Seconds: 100,
				Nanos:   0,
			},
		},
	}

	calendar.addEventResponse, err = calendar.grpcClient.Add(context.Background(), addRequest)
	calendar.responseStatus = status.Code(err)

	return err
}

func (calendar *CalendarTest) iSendAddEventRequestForToday() error {
	var err error
	grpcTimestamp, err := ConvertDateToGrpcTimestamp(time.Now().String())
	if err != nil {
		return err
	}

	addRequest := &spec.AddRequest{
		Event: &spec.Event{
			Uuid:        uuid.New().String(),
			Title:       "Test title of event",
			Description: "Test description of event",
			Start:       grpcTimestamp,
			Finish: &timestamp.Timestamp{
				Seconds: 1000,
				Nanos:   0,
			},
			UserId: 1,
			NotifyTime: &timestamp.Timestamp{
				Seconds: 100,
				Nanos:   0,
			},
		},
	}

	calendar.addEventResponse, err = calendar.grpcClient.Add(context.Background(), addRequest)
	calendar.responseStatus = status.Code(err)

	return err
}

func (calendar *CalendarTest) iHaveEventOnDate(date string) error {
	var err error
	grpcTimestamp, err := ConvertDateToGrpcTimestamp(time.Now().String())
	if err != nil {
		return err
	}

	addRequest := &spec.AddRequest{
		Event: &spec.Event{
			Uuid:        uuid.New().String(),
			Title:       "Test title of event",
			Description: "Test description of event",
			Start:       grpcTimestamp,
			Finish: &timestamp.Timestamp{
				Seconds: 1000,
				Nanos:   0,
			},
			UserId: 1,
			NotifyTime: &timestamp.Timestamp{
				Seconds: 100,
				Nanos:   0,
			},
		},
	}

	calendar.addEventResponse, err = calendar.grpcClient.Add(context.Background(), addRequest)
	calendar.responseStatus = status.Code(err)

	return err
}

func (calendar *CalendarTest) iSendUpdateEventRequest() error {
	updateRequest := &spec.UpdateRequest{
		Uuid: calendar.addEventResponse.Uuid,
		Event: &spec.Event{
			Title:       "Test update of title of event",
			Description: "Test description of event",
			Start: &timestamp.Timestamp{
				Seconds: 100,
				Nanos:   0,
			},
			Finish: &timestamp.Timestamp{
				Seconds: 1000,
				Nanos:   0,
			},
			UserId: 1,
			NotifyTime: &timestamp.Timestamp{
				Seconds: 100,
				Nanos:   0,
			},
		},
	}

	_, err := calendar.grpcClient.Update(context.Background(), updateRequest)
	calendar.responseStatus = status.Code(err)

	return nil
}

func (calendar *CalendarTest) iSendUpdateEventRequestOf(id string) error {
	updateRequest := &spec.UpdateRequest{
		Uuid: id,
		Event: &spec.Event{
			Title:       "Test update of title of event",
			Description: "Test description of event",
			Start: &timestamp.Timestamp{
				Seconds: 100,
				Nanos:   0,
			},
			Finish: &timestamp.Timestamp{
				Seconds: 1000,
				Nanos:   0,
			},
			UserId: 1,
			NotifyTime: &timestamp.Timestamp{
				Seconds: 100,
				Nanos:   0,
			},
		},
	}

	_, err := calendar.grpcClient.Update(context.Background(), updateRequest)
	calendar.responseStatus = status.Code(err)

	return nil
}

func (calendar *CalendarTest) iSendDeleteEventRequest() error {
	deleteRequest := &spec.DeleteRequest{
		Uuid: calendar.addEventResponse.Uuid,
	}

	_, err := calendar.grpcClient.Delete(context.Background(), deleteRequest)
	calendar.responseStatus = status.Code(err)

	return nil
}

func (calendar *CalendarTest) iSendDeleteEventRequestOf(id string) error {
	deleteRequest := &spec.DeleteRequest{
		Uuid: id,
	}

	_, err := calendar.grpcClient.Delete(context.Background(), deleteRequest)
	calendar.responseStatus = status.Code(err)

	return nil
}

func (calendar *CalendarTest) iCallGetList() error {
	var err error

	getListRequest := &spec.GetListRequest{}

	calendar.getListResponse, err = calendar.grpcClient.GetList(context.Background(), getListRequest)
	calendar.responseStatus = status.Code(err)

	return err
}

// Expectations
func (c *CalendarTest) iShouldReceiveEventNotification() error {
	// todo
}

// Check responses
func (calendar *CalendarTest) wantToSeeInResponseEventID() error {
	if calendar.addEventResponse == nil {
		return errors.New("expected response, got empty response")
	}

	if calendar.addEventResponse.Uuid == "" || calendar.addEventResponse.Uuid == uuid.Nil.String() {
		return errors.New("expected ID in response, got empty ID")
	}

	return nil
}

func (calendar *CalendarTest) wantToSeeEventsInResponseEventID() error {
	if calendar.addEventResponse == nil {
		return errors.New("expected response, got empty response")
	}

	if calendar.addEventResponse.Uuid == "" || calendar.addEventResponse.Uuid == uuid.Nil.String() {
		return errors.New("expected ID in response, got empty ID")
	}

	return nil
}

func (calendar *CalendarTest) iGetSuccessResponse() error {
	if calendar.responseStatus != codes.OK {
		return errors.New("expected 200, got " + string(calendar.responseStatus))
	}
	return nil
}

func (calendar *CalendarTest) iGetErrorResponse() error {
	if calendar.responseStatus == codes.OK {
		return errors.New("expected error, got " + string(calendar.responseStatus))
	}
	return nil
}