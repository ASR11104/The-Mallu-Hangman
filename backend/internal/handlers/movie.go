package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ASR11104/the-mallu-hangman/internal/config"
	"github.com/ASR11104/the-mallu-hangman/internal/session"
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
	Results      []Details `json:"results"`
	Page         int       `json:"page"`
	TotalPages   int       `json:"total_pages"`
	TotalResults int       `json:"total_results"`
}

func Movies(w http.ResponseWriter, r *http.Request, cfg config.Config, sessionManager *session.Manager) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.URL.Query().Get("session_id")
	difficulty := r.URL.Query().Get("difficulty")
	fmt.Println("Difficulty:", difficulty)
	language := r.URL.Query().Get("language")
	page := 1
	filters := make(map[string]string)
	switch difficulty {
	case "easy":
		filters = map[string]string{
			"with_original_language": language,
			"vote_average.gte":       "5",
			"vote_count.gte":         "5",
			"sort_by":                "vote_average.desc",
		}
		totalPAges := getTotalPages(filters, cfg.TheMovieDBToken)
		page = utils.RandomNumber(1, totalPAges)
		filters["page"] = fmt.Sprint(page)
	case "medium":
		filters = map[string]string{
			"with_original_language": language,
			"vote_average.gte":       "1",
			"vote_count.gte":         "1",
			"sort_by":                "vote_average.desc",
		}
		totalPAges := getTotalPages(filters, cfg.TheMovieDBToken)
		page = utils.RandomNumber(1, totalPAges)
		filters["page"] = fmt.Sprint(page)
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
	randomMovie := chooseRandomMovie(response.Results, sessionID, sessionManager)
	fmt.Printf("Selected Movie: %+v", randomMovie)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(randomMovie)
}

func getTotalPages(filters map[string]string, token string) (totalPages int) {
	release_date_filter := "primary_release_date.gte=2000-01-01"
	result_filter := ""
	for key, value := range filters {
		if key == "page" {
			continue
		}
		result_filter += "&" + key + "=" + value
	}
	result_filter += "&" + release_date_filter
	url := "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&" + result_filter
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

	var responseData Response
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	fmt.Println("Total Pages:", responseData.TotalPages)
	fmt.Println("Total Results:", responseData.TotalResults)
	return responseData.TotalPages
}

func getMovies(filters map[string]string, token string) (responseData Response) {
	release_date_filter := "primary_release_date.gte=2000-01-01"
	result_filter := ""
	for key, value := range filters {
		if key == "page" {
			continue
		}
		result_filter += "&" + key + "=" + value
	}
	result_filter += "&" + release_date_filter
	url := "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&" + result_filter
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
	// fmt.Println(responseData)
	return responseData
}

func chooseRandomMovie(movies []Details, sessionID string, sessionManager *session.Manager) Details {
	// If no session ID provided, return a random movie without tracking
	if sessionID == "" {
		randomIndex := utils.RandomNumber(0, len(movies)-1)
		return movies[randomIndex]
	}

	// Ensure session exists
	sess := sessionManager.GetSession(sessionID)
	if sess == nil {
		sess = sessionManager.CreateSession(sessionID)
	}

	// Filter out already used movies
	availableMovies := make([]Details, 0, len(movies))
	for _, movie := range movies {
		if !sessionManager.IsMovieUsed(sessionID, movie.ID) {
			availableMovies = append(availableMovies, movie)
		}
	}

	// If all movies in this page are used, try to get another page
	maxAttempts := 3
	for len(availableMovies) == 0 && maxAttempts > 0 {
		maxAttempts--
		// This is a simplified approach - in production you might want to fetch more pages
		break
	}

	if len(availableMovies) == 0 {
		// All movies in current response are used, return a random one anyway
		randomIndex := utils.RandomNumber(0, len(movies)-1)
		return movies[randomIndex]
	}

	randomIndex := utils.RandomNumber(0, len(availableMovies)-1)
	selectedMovie := availableMovies[randomIndex]

	// Mark the movie as used in this session
	sessionManager.MarkMovieAsUsed(sessionID, selectedMovie.ID)

	return selectedMovie
}
