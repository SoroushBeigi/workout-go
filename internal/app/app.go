package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SoroushBeigi/workout-go/internal/api"
)

type Application struct {
	Logger *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	workoutHandler:= api.NewWorkoutHandler()
	app := &Application{logger,workoutHandler}
	return app, nil
}

func (a *Application)HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "\nStatus is available\n")
}
