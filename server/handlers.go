package main

import (
	"encoding/json"
	"github.com/abates/disgo"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type appHandler func(http.ResponseWriter, *http.Request) (interface{}, error)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i, err := fn(w, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(i); err != nil {
		panic(err)
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

type searchCriteria struct {
	Hash     disgo.PHash `json:"hash"`
	Distance uint        `json:"distance"`
}

func Search(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var c searchCriteria
	if err := decode(r, 10240, &c); err != nil {
		return nil, err
	}
	return disgo.SearchByHash(c.Hash, c.Distance)
}

type imageInfo struct {
	Hash     disgo.PHash `json:"hash"`
	Location string      `json:"location"`
}

func Add(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var i imageInfo
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
