package handlers

import (
	"egit90/urlShortner/model"
	"encoding/json"
	"fmt"
	"net/http"
)

type db interface {
	GetValue(string) (string, error)
	InsertKeyValue(model.Data) error
	GetAll() ([]model.Data, error)
}

func DbHandler(myDb db, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if value, err := myDb.GetValue(path); err == nil {
			http.Redirect(w, r, value, http.StatusFound)
		}
		fallback.ServeHTTP(w, r)
	}
}

func SaveHandler(myDb db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var data map[string]string

		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for key, val := range data {
			k := []byte(key)
			if k[0] != '/' {
				http.Error(w, "path must start with '/' ", http.StatusBadRequest)
				return
			}
			v := []byte(val)
			if err := myDb.InsertKeyValue(model.Data{Key: k, Value: v}); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		w.WriteHeader(http.StatusCreated)

	}
}

func GetAllHandler(myDB db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		data, err := myDB.GetAll()

		if err != nil {
			http.Error(w, "Error getting data", http.StatusInternalServerError)
			return
		}

		// Convert byte slices to strings
		res := []map[string]string{}

		for _, item := range data {
			keyString := string(item.Key)
			valueString := string(item.Value)
			res = append(res, map[string]string{"Key": keyString, "Value": valueString})
		}
		jsonData, err := json.Marshal(res)
		if err != nil {
			http.Error(w, "Error converting data to JSON", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)

	}
}

func DefaultMux(mydb db) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/save", SaveHandler(mydb))
	mux.HandleFunc("/getall", GetAllHandler(mydb))
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}
