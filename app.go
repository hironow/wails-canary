package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// RandomImage get result
type RandomImage struct {
	Message string
	Status  string
}

// AllBreeds get result
type AllBreeds struct {
	Message map[string]map[string][]string
	Status  string
}

// ImagesByBreed get result
type ImagesByBreed struct {
	Message []string
	Status  string
}

func (a *App) GetRandomImageUrl() string {
	res, err := http.Get("https://dog.ceo/api/breeds/image/random")
	if err != nil {
		log.Fatal(err)
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data RandomImage
	json.Unmarshal(resData, &data)
	return data.Message
}

func (a *App) GetBreedList() []string {
	var breeds []string

	res, err := http.Get("https://dog.ceo/api/breeds/list/all")
	if err != nil {
		log.Fatal(err)
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data AllBreeds
	json.Unmarshal(resData, &data)

	for k := range data.Message {
		breeds = append(breeds, k)
	}
	sort.Strings(breeds)
	return breeds
}

func (a *App) GetImageUrlsByBreed(breed string) []string {
	url := fmt.Sprintf("%s%s%s%s", "https://dog.ceo/api/", "breed/", breed, "/images")
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data ImagesByBreed
	json.Unmarshal(resData, &data)
	return data.Message
}
