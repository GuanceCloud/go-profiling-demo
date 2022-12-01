package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

var movies = func() []Movie {
	movies, err := readMovies()
	if err != nil {
		panic(err)
	}
	return movies
}()

type Movie struct {
	Title       string  `json:"title"`
	VoteAverage float64 `json:"vote_average"`
	ReleaseDate string  `json:"release_date"`
}

type Movies []Movie

func (m Movies) Len() int {
	return len(m)
}

func (m Movies) Less(i, j int) bool {

	t1, err := time.Parse("2006-01-02", m[i].ReleaseDate)
	if err != nil {
		return false
	}
	t2, err := time.Parse("2006-01-02", m[j].ReleaseDate)
	if err != nil {
		return true
	}
	return t1.After(t2)
}

func (m Movies) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

var _ sort.Interface = Movies{}

func readMovies() ([]Movie, error) {
	f, err := os.Open("./movies5000.json.gz")
	if err != nil {
		return nil, fmt.Errorf("open movies data file fail: %w", err)
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("gzip new reader from *FILE fail: %w", err)
	}
	defer r.Close()

	var movies []Movie

	if err := json.NewDecoder(r).Decode(&movies); err != nil {
		return nil, fmt.Errorf("json unmarshal fail: %w", err)
	}

	return movies, nil
}

func main() {

	err := profiler.Start(
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,

			// The profiles below are disabled by default to keep overhead
			// low, but can be enabled as needed.
			profiler.BlockProfile,
			profiler.MutexProfile,

			profiler.GoroutineProfile,
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer profiler.Stop()

	http.HandleFunc("/movies", func(w http.ResponseWriter, req *http.Request) {
		q := req.URL.Query().Get("q")

		moviesCopy := make(Movies, len(movies))
		copy(moviesCopy, movies)

		sort.Sort(moviesCopy)

		if q != "" {
			q = strings.ToUpper(q)
			matchCount := 0
			for idx, m := range moviesCopy {
				if strings.Contains(strings.ToUpper(m.Title), q) && idx != matchCount {
					moviesCopy[matchCount] = moviesCopy[idx]
					matchCount++
				}
			}
			moviesCopy = moviesCopy[:matchCount]
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(moviesCopy); err != nil {
			log.Printf("encode into json fail: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
