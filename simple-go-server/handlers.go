package main

import (
	"encoding/json"
	"fmt"
	"github.com/ajmclaug1/simple-go-web-app/simple-go-server/internal/data"
	"net/http"
	"strconv"
	"time"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")

	w.Write(js)

}

func (app *application) getCreateBooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		books, err := app.models.Books.GetAll()

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}
	if r.Method == http.MethodPost {
		var input struct {
			Title     string   `json:"title"`
			Published int      `json:"published,omitempty"`
			Pages     int      `json:"pages,omitempty"`
			Genres    []string `json:"genres,omitempty"`
			Rating    float32  `json:"rating,omitempty"`
		}

		err := app.readJSON(w, r, &input)
		
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
	  book := &data.Book{
      Title: input.Title,
      Published: input.Published,
      Pages: input.Pages,
      Genres: input.Genres,
      Rating: input.Rating,
    }
    err = app.models.Books.Insert(book)
	
    if err != nil {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      return
    }

    headers := make(http.Header)
    headers.Set("location", fmt.Sprintf("v1/books/%d", book.ID))

    err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
  
    if err != nil {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      return
    }
}
}

func (app *application) getUpdateDeleteBooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getBook(w, r)
	case http.MethodPut:
		app.updateBook(w, r)
	case http.MethodDelete:
		app.deleteBook(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	book := data.Book{
		ID:        idInt,
		CreatedAt: time.Now(),
		Title:     "Echoes in the Darkness",
		Published: 2019,
		Pages:     300,
		Genres:    []string{"Fiction", "Thriller"},
		Rating:    4.5,
		Version:   1,
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

	}
}

func (app *application) updateBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	var input struct {
		Title     *string  `json:"title"`
		Published *int     `json:"published"`
		Pages     *int     `json:"pages"`
		Genres    []string `json:"genres"`
		Rating    *float32 `json:"rating"`
	}

	book := data.Book{
		ID:        idInt,
		CreatedAt: time.Now(),
		Title:     "Echoes in the Darkness",
		Published: 2019,
		Pages:     300,
		Genres:    []string{"Fiction", "Thriller"},
		Rating:    4.5,
		Version:   1,
	}

	err = app.readJSON(w, r, &input)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if input.Title != nil {
		book.Title = *input.Title
	}

	if input.Published != nil {
		book.Published = *input.Published
	}
	if input.Pages != nil {
		book.Pages = *input.Pages
	}
	if len(input.Genres) > 0 {
		book.Genres = input.Genres
	}
	if input.Rating != nil {
		book.Rating = *input.Rating
	}
	fmt.Fprintf(w, "%v\n", book)
}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	fmt.Fprintf(w, "Delete a book with ID; %d", idInt)
}
