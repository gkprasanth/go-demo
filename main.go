package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Person struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

type Response struct {
	Message string `json:"message,omitempty"`
	Data    Person
	Error   string
}

var datalayer map[int]Person
var personID int

func init() {
	datalayer = make(map[int]Person)
	personID = 0

	log.Println("datalayer is initialized..")
}

func main() {
	log.Println("server started o port: 4444")

	http.HandleFunc("/", healthHandler)

	http.HandleFunc("/person", personhandler)
	log.Fatal(http.ListenAndServe(":4444", nil))

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Server is running"))
}

func personhandler(w http.ResponseWriter, r *http.Request) {
	var res Response

	//send the response  at end
	defer func() {
		json.NewEncoder(w).Encode(res)
	}()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)

	switch r.Method {
	case http.MethodPost:
		//read json only and convert to struct
		var p Person

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			res.Error = err.Error()
			return
		}
		id := createPerson(p)
		res.Message = fmt.Sprint("New person created wiht id: %d", id)

	case http.MethodGet:
		strId := r.URL.Query().Get("id")

		if strId == "" {
			log.Println("invalid id")
			res.Error = "invalid id"
			w.WriteHeader(400)
			return
		}

		//coonvert id tp int
		id, err := strconv.Atoi(strId)

		if err != nil {
			log.Printf("invalid id", err)
			res.Error = err.Error()
			w.WriteHeader(400)
			return
		}

		p, err := readPerson(id)

		if err != nil {
			log.Printf("read error", err)
			res.Error = err.Error()
			w.WriteHeader(400)
			return
		}
		res.Data = p

	case http.MethodPut:
		var p Person
		strId := r.URL.Query().Get("id")

		if strId == "" {
			res.Error = "invalid id"
			w.WriteHeader(400)
			log.Printf("invalid id")
		}
		id, err := strconv.Atoi(strId)

		if err != nil {
			log.Printf("error", err)
			w.WriteHeader(400)
			res.Error = err.Error()
		}

		json.NewDecoder(r.Body).Decode(&p)

		 err = updatePerson(id, p)

		if err != nil {
			log.Println("update err", err)
			res.Error = err.Error()
			w.WriteHeader(400)
		}

		res.Message = "Person updated successfully"
		res.Data = p

	case http.MethodDelete:
		strId := r.URL.Query().Get("id")

		if strId != "" {
			log.Println("invalid id")
			w.WriteHeader(400)
			res.Message = "invalid id"
			res.Error = "id error"
			return
		}

		id, err := strconv.Atoi(strId)

		if err != nil {
			w.WriteHeader(400)
			log.Println("Error", err)
			res.Error = err.Error()
			return
		}

		err = deletePerson(id)

		if err != nil {
			w.WriteHeader(400)
			log.Println("error occured:", err)
			res.Error = err.Error()
			return
		}

		res.Message = "Deleted successfully"

	}

}

//Data Layer

func createPerson(p Person) int {
	pId := personID + 1
	datalayer[pId] = p
	personID = pId
	return pId
}

func readPerson(id int) (Person,error) {
	var person Person

	for key, value := range datalayer {
		if key == id {
			person = value
			return person, nil
		}
	}

	return Person{}, fmt.Errorf("person id not found. id: %d", id)
}

func updatePerson(id int, p Person) error {
	for key, _ := range datalayer {
		if key == id {
			datalayer[id] = p
			return nil
		}
	}

	//for invalid id
	return fmt.Errorf("person id is not found to update. id: %d", id)
}

func deletePerson(id int) error {

	for key, _ := range datalayer {
		if key == id {
			delete(datalayer, id)
			return nil
		}
	}

	return fmt.Errorf("person Id not found. id: %d", id)
}
