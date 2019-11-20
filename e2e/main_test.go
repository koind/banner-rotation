package main

import (
	"github.com/DATA-DOG/godog"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	time.Sleep(5 * time.Second)

	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    "progress",
		Paths:     []string{"features"},
		Randomize: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
