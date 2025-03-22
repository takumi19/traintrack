package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"traintrack/internal/database"
)

func (a *Api) handleListExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := a.db.ListExercises()
	if err != nil {
    a.l.Level(ERROR).Println("Failed to list exercises:", err.Error())
		WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: err.Error()})
		return
	}
	WriteJSON(w, http.StatusOK, exercises)
}

func (a *Api) handleGetExerciseByID(w http.ResponseWriter, r *http.Request) {
	exerciseId, err := strconv.ParseInt(r.PathValue("exercise_id"), 10, 64)
  if err != nil {
    a.l.Level(ERROR).Println(err.Error())
    WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: err.Error()})
    return
  }

  exercise, err := a.db.GetExerciseByID(exerciseId)
  if err != nil {
    a.l.Level(ERROR).Println(err.Error())
    WriteJSON(w, http.StatusInternalServerError, err.Error())
    return
  }

  WriteJSON(w, http.StatusOK, exercise)
}

func (a *Api) handleGetExerciseByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("exercise_name")

  exercise, err := a.db.GetExerciseByName(name)
  if err != nil {
    a.l.Level(ERROR).Println(err.Error())
    WriteJSON(w, http.StatusInternalServerError, err.Error())
    return
  }

  WriteJSON(w, http.StatusOK, exercise)
}

// TODO: Make this update on conflict
func (a *Api) handleAddExercise(w http.ResponseWriter, r *http.Request) {
  var exercise database.ExerciseInfo
  if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
    a.l.Level(ERROR).Println("Failed to decode exercise info from request:", err.Error())
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  found, err := a.db.GetExerciseByName(exercise.Name)
  if found != nil || err != nil {
    WriteJSON(w, http.StatusConflict, &ApiError{Error: "This exercise already exists"})
    return
  }

  id, err := a.db.AddExercise(&exercise)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  WriteJSON(w, http.StatusOK, map[string]int64{
    "Id": id,
  })
}
