package main

import (
	"log"
	"os"
	"testing"

	"github.com/cucumber/godog"
)

func TestMain(m *testing.M) {
	log.Printf("Main app test...")

	status := godog.TestSuite{
		Name:                "icalendar_test",
		ScenarioInitializer: InitFeatureContext,
		Options: &godog.Options{
			Format:    "pretty",
			Paths:     []string{"features"},
			Randomize: 0,
		},
	}.Run()

	if runStatus := m.Run(); status > 0 {
		os.Exit(runStatus)
		return
	}

	os.Exit(0)
}
