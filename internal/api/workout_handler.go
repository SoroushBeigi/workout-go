package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/SoroushBeigi/workout-go/internal/store"
	"github.com/SoroushBeigi/workout-go/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {

	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIdParam %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}
	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutById %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": workout})

}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: decoding createWorkout %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: createWorkout %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"})
		return
	}
	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: updateWorkout %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutById <- updateWorkout %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	var updateWorkoutRequest struct {
		Title           *string          `json:"title"`
		Description     *string          `json:"description"`
		DurationMinutes *int             `json:"duration_minutes"`
		CaloriesBurned  *int             `json:"calories_burned"`
		Exercises       []store.Exercise `json:"exercises"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		wh.logger.Printf("ERROR: updateWorkout - decoder %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}

	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}

	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}

	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}

	if updateWorkoutRequest.Exercises != nil {
		existingWorkout.Exercises = updateWorkoutRequest.Exercises
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		wh.logger.Printf("ERROR: updateWorkout %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to update workout"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})

}

func (wh *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: deleteWorkout %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutID)
	if err == sql.ErrNoRows {
		wh.logger.Printf("ERROR: deleteWorkout sql %v", err)
		utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}
	if err != nil {
		wh.logger.Printf("ERROR: deleteWorkout %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to delete workout"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"message": "workout deleted successfully"})

}
