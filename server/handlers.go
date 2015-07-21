package main

import (
	"encoding/json"
	"github.com/abates/disgo"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type appHandler func(http.ResponseWriter, *http.Request) (interface{}, error)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i, err := fn(w, r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		if err == disgo.ErrNotFound {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode("Resource not found")
		} else {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(err.Error())
		}
	} else {
		w.WriteHeader(http.StatusOK)
		if e := json.NewEncoder(w).Encode(i); e != nil {
			log.Printf("Failed to encode output: %v", e)
		}
	}
}

func decode(r *http.Request, limit int64, v interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, limit))
	if err != nil {
		return err
	}

	if err := r.Body.Close(); err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func Search(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var c disgo.SearchCriteria
	if err := decode(r, 10240, &c); err != nil {
		return nil, err
	}
	return disgo.SearchByHash(c.Hash, c.Distance)
}

func Add(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var i disgo.ImageInfo
	if err := decode(r, 10240, &i); err != nil {
		return nil, err
	}

	if err := disgo.AddLocation(i.Location, i.Hash); err != nil {
		return nil, err
	}

	return "Added " + i.Location + " to the database", nil
}

func Show(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	hash, err := strconv.ParseUint(vars["hash"], 10, 64)
	if err != nil {
		return nil, err
	}
	return disgo.Find(disgo.PHash(hash))
}
