package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

type OptionsType struct {
	CaseInsensitive bool
	WholeWord       bool
	WholeLine       bool
	FilenameOnly    bool
	FilesWoMatches  bool
}

func (a *App) SelectFolder() string {
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Folder",
	})
	if err != nil {
		log.Println("Error selecting a folder")
	}
	return selection
}

func (a *App) SelectFile() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select File",
	})
	if err != nil {
		log.Println("Error selecting a file")
	}
	return selection
}

func (a *App) Search(path string, pattern string, options OptionsType) string {
	if path == "" {
		_, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:         "ERROR",
			Message:       "No path was entered",
			Buttons:       []string{"OK"},
			DefaultButton: "OK",
		})
		if err != nil {
			log.Println(err)
			return ""
		}
		return ""
	}
	if pattern == "" {
		_, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:         "ERROR",
			Message:       "No pattern was entered",
			Buttons:       []string{"OK"},
			DefaultButton: "OK",
		})
		if err != nil {
			log.Println(err)
			return ""
		}
		return ""
	}

	if options.CaseInsensitive {
		pattern = "(?i)" + pattern
	}
	if options.WholeWord {
		pattern = `\b` + pattern + `\b`
	}
	if options.WholeLine {
		pattern = "^" + pattern + "$"
	}

	results, err := walkDir(path, pattern, options)
	if err != nil {
		log.Fatalln("ERROR walking the directory!")
	}
	if len(results) > 0 {
		return results
	} else {
		return "Not results were found"
	}
}

func walkDir(dirToWalk string, pattern string, options OptionsType) (string, error) {
	var matches []string
	err := filepath.Walk(dirToWalk, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			matchesFromFile, err2 := checkFileForPattern(path, pattern, options)
			if err2 != nil {
				log.Printf("Failed opening file: %s", err2)
			} else {
				matches = append(matches, matchesFromFile...)
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	results := strings.Join(matches, "")
	return results, nil
}

func isBinary(fileToRead string) (bool, error) {
	data := make([]byte, 256)
	file, err := os.Open(fileToRead)
	if err != nil {
		return false, nil
	}
	defer file.Close()
	count, err := file.Read(data)
	if err != nil {
		return false, err
	}
	for i := 0; i < count; i++ {
		if data[i] == 0 {
			return true, nil
		}
	}
	return false, nil
}

func checkFileForPattern(fileToRead string, pattern string, options OptionsType) ([]string, error) {
	var matches []string
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fileToRead)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Size() == 0 {
		return nil, nil
	}

	fileIsBinary, err := isBinary(fileToRead)
	if err != nil {
		return nil, err
	}
	if fileIsBinary {
		log.Printf("%s is binary\n", fileToRead)
		return nil, nil
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtLines []string
	for scanner.Scan() {
		txtLines = append(txtLines, scanner.Text())
	}
	if len(txtLines) == 0 {
		log.Printf("%s has no new line control characters.\n", fileToRead)
		return nil, nil
	}
	fileToRead = strings.ReplaceAll(fileToRead, `\`, `\\`)
	numNonMatches := 0
	numLinesMatches := 0
	for lineNum, line := range txtLines {
		if r.MatchString(line) {
			numLinesMatches++
			if !options.FilesWoMatches {
				if options.FilenameOnly {
					if numLinesMatches == 1 {
						match := fmt.Sprintf("\n%s\n", fileToRead)
						matches = append(matches, match)
						break
					}
				}
				var printableLine string
				var sb strings.Builder
				for _, r := range line {
					if int(r) >= 32 && int(r) != 127 {
						if r == '\\' || r == '"' {
							sb.WriteRune('\\')
						}
						sb.WriteRune(r)
					}
				}
				printableLine = sb.String()
				if numLinesMatches == 1 {
					match := fmt.Sprintf("\n%s\n\t%d:  %s\n", fileToRead, lineNum+1, printableLine)
					matches = append(matches, match)
				} else {
					match := fmt.Sprintf("\t%d:  %s\n", lineNum+1, printableLine)
					matches = append(matches, match)
				}
			}
		} else {
			numNonMatches++
		}
	}
	if options.FilesWoMatches && numNonMatches == len(txtLines) {
		match := fmt.Sprintf("\n%s\n", fileToRead)
		matches = append(matches, match)
	}
	return matches, nil
}
