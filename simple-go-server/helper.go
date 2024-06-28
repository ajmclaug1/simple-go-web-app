package main

import (
	"encoding/json"
	"net/http"
  "io"
  "errors"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil

}

func (app *application) readJSON (w http.ResponseWriter, r *http.Request, dst any) error {
  maxBytes := 1_048_576
  r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

  dec := json.NewDecoder(r.Body)
  dec.DisallowUnknownFields()

  if err := dec.Decode(dst); err != nil {
    return err
  }

  err := dec.Decode(&struct{}{})

  if err != io.EOF {
    return errors.New("body must only contain a single JSON object")
  }

  return nil
}