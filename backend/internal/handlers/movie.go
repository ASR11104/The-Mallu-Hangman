package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ASR11104/the-mallu-hangman/internal/config"
	"github.com/ASR11104/the-mallu-hangman/internal/utils"
)

type Details struct {
	ID               int64   `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	Overview         string  `json:"overview"`
	Title            string  `json:"title"`
	ReleaseDate      string  `json:"release_date"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int64   `json:"vote_count"`
	Popularity       float64 `json:"popularity"`
}

type Response struct {
	Results    []Details `json:"results"`
	Page       int       `json:"page"`
	TotalPages int       `json:"total_pages"`
}

func Movies(w http.ResponseWriter, r *http.Request, cfg config.Config) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	difficulty := r.URL.Query().Get("difficulty")
	fmt.Println("Difficulty:", difficulty)
	language := r.URL.Query().Get("language")
	page := 1
	filters := make(map[string]string)
	switch difficulty {
	case "easy":
		page = utils.RandomNumber(1, utils.Easy)
		filters = map[string]string{
			"with_original_language": language,
			"page":                   fmt.Sprint(page),
			"vote_average.gte":       "7.5",
			"vote_count.gte":         "500",
			"sort_by":                "vote_average.desc",
		}
	case "medium":
		page = utils.RandomNumber(utils.Easy, utils.Medium)
		filters = map[string]string{
			"with_original_language": language,
			"page":                   fmt.Sprint(page),
		}
	case "hard":
		page = utils.RandomNumber(utils.Medium, utils.Hard)
		filters = map[string]string{
			"with_original_language": language,
			"page":                   fmt.Sprint(page),
		}
	default:
		http.Error(w, "Invalid difficulty level", http.StatusBadRequest)
		return
	}

	response := getMovies(filters, cfg.TheMovieDBToken)
	randomMovie := chooseRandomMovie(response.Results)
	fmt.Printf("Selected Movie: %+v", randomMovie)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(randomMovie)
}

func getMovies(filters map[string]string, token string) (responseData Response) {
	release_date_filter := "primary_release_date.gte=2000-01-01"
	for key, value := range filters {
		fmt.Println("Adding filter:", key, "=", value)
		release_date_filter += "&" + key + "=" + value
	}
	fmt.Println("Filters:", filters)
	url := "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&" + release_date_filter
	fmt.Println("Request URL:", url)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	// fmt.Println(string(body))
	fmt.Println(responseData)
	return responseData
}

func chooseRandomMovie(movies []Details) Details {
	randomIndex := utils.RandomNumber(0, len(movies)-1)
	return movies[randomIndex]
}
