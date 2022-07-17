package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "xxxxxx"
	DB_NAME     = "go_movies_api"
)

// DB set up
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

type Movie struct {
	MovieID   string `json:"movieid"`
	MovieName string `json:"moviename"`
}

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Movie `json:"data"`
	Message string  `json:"message"`
}

type OneDataJsonResponse struct {
	Type    string `json:"type"`
	Data    Movie  `json:"data"`
	Message string `json:"message"`
}

// Main function
func main() {
	// Init the mux router
	router := mux.NewRouter()

	// Route handles & endpoints

	// Get all movies
	router.HandleFunc("/movies", getMovies).Methods("GET")
	// Get movie by id
	router.HandleFunc("/movies/{movieid}", getMovie).Methods("GET")
	// Create a movie
	router.HandleFunc("/movies", createMovie).Methods("POST")
	// Update a movie by id
	router.HandleFunc("/movies/{movieid}", updateMovie).Methods("PUT")
	// Delete a movie by id
	router.HandleFunc("/movies/{movieid}", deleteMovie).Methods("DELETE")
	// Delete all movies
	router.HandleFunc("/movies", deleteMovies).Methods("DELETE")

	// serve the app
	fmt.Println("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// Function for handling errors
func checkErr(err error) {
	fmt.Println(err)

	if err != nil {
		panic(err)
	}
}

// Get all movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := setupDB()

	// Get all movies from movies table that don't have movieID = "1"
	rows, err := db.Query("SELECT * FROM movies")

	// check for errors
	checkErr(err)

	// var response []JsonResponse
	var movies []Movie
	defer rows.Close()

	// Foreach movie
	for rows.Next() {
		var id int
		var movieID string
		var movieName string

		err := rows.Scan(&id, &movieID, &movieName)

		checkErr(err)

		movies = append(movies, Movie{MovieID: movieID, MovieName: movieName})
	}

	var response = JsonResponse{Type: "success", Data: movies}
	json.NewEncoder(w).Encode(response)
}

// Get a movie
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	movieId := params["movieid"]

	var response = OneDataJsonResponse{}
	var movie Movie

	if movieId == "" {
		response = OneDataJsonResponse{Type: "error", Message: "You are missing movieID parameter."}
	} else {
		db := setupDB()
		result, err := db.Query("SELECT * FROM movies where movieID = $1", movieId)
		// check errors
		checkErr(err)

		defer result.Close()

		for result.Next() {
			var id int
			var movieID string
			var movieName string

			err := result.Scan(&id, &movieID, &movieName)

			checkErr(err)

			movie = Movie{MovieID: movieID, MovieName: movieName}
		}
		response = OneDataJsonResponse{Type: "success", Message: "The movie has been inserted successfully!", Data: movie}

	}
	json.NewEncoder(w).Encode(response)
}

// Create a movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	fmt.Println(movie.MovieID)

	movieID := movie.MovieID
	movieName := movie.MovieName

	var response = OneDataJsonResponse{}

	if movieID == "" || movieName == "" {
		response = OneDataJsonResponse{Type: "error", Message: "You are missing movieid or moviename parameter"}
	} else {
		db := setupDB()

		fmt.Println("Inserting new movie with ID: " + movieID + " and name: " + movieName)

		var lastInsertId int

		err := db.QueryRow("INSERT INTO movies(movieID, movieName) VALUES($1, $2) returning id;", movieID, movieName).Scan(&lastInsertId)
		checkErr(err)

		movie = Movie{MovieID: movieID, MovieName: movieName}

		response = OneDataJsonResponse{Type: "success", Message: "The movie has been inserted successfully!", Data: movie}
	}
	json.NewEncoder(w).Encode(response)
}

// Update a movie
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	movieId := params["movieid"]

	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	fmt.Println(movie.MovieID)

	movieID := movie.MovieID
	movieName := movie.MovieName

	var response = OneDataJsonResponse{}

	if movieID == "" || movieName == "" {
		response = OneDataJsonResponse{Type: "error", Message: "You are missing movieid or moviename parameter"}
	} else {
		db := setupDB()

		_, err := db.Exec("UPDATE movies SET movieID = $1, movieName = $2 WHERE movieID = $3", movieID, movieName, movieId)
		checkErr(err)

		movie = Movie{MovieID: movieID, MovieName: movieName}

		response = OneDataJsonResponse{Type: "success", Message: "The movie has been updated successfully!", Data: movie}
	}
	json.NewEncoder(w).Encode(response)
}

// Delete a movie
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	movieId := params["movieid"]

	var response = JsonResponse{}
	if movieId == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID parameter."}
	} else {
		db := setupDB()
		_, err := db.Exec("DELETE FROM movies where movieID = $1", movieId)
		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The movie has been deleted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

// Delete all movies
func deleteMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := setupDB()

	_, err := db.Exec("DELETE FROM movies")

	// check errors
	checkErr(err)

	var response = JsonResponse{Type: "success", Message: "All movies have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}
