package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	id int
}

func (s *Server) HandleGetId(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct {
		Result string
		Id     int
	}{"OK", s.id})
}

func (s *Server) HandleSetId(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		json.NewEncoder(w).Encode(struct {
			Result string
		}{"error, can't parse int"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.id = id

	json.NewEncoder(w).Encode(struct {
		Result string
		Id     int
	}{"OK", id})
}

type Client struct{}

func (c Client) echoRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error requesting url: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(fmt.Errorf("error closing body: %w", err))
		}
	}()

	respBody, _ := io.ReadAll(resp.Body)

	return string(respBody), nil
}
