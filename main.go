package main

import (
	"database/sql"
	"strconv"

	"fmt"

	_ "github.com/lib/pq" //Go postgres driver for Go's database/sql package
	"gopkg.in/yaml.v3"    //The yaml package enables Go programs to comfortably encode and decode YAML values

	"encoding/json"

	"net/http"

	"github.com/gorilla/mux" //Package gorilla/mux implements a request router and dispatcher for matching incoming requests to their respective handler
)

const (
	host     = "host.docker.internal" //we will connect to postgresql from Docker
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	r := newRouter()

	http.ListenAndServe(":8181", r)
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	staticFileDirectory := http.Dir("./static/")
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/static/").Handler(staticFileHandler).Methods("GET")
	r.HandleFunc("/student", getStudentHandler).Methods("GET")
	r.HandleFunc("/api/v1/StudentCreate", studentCreate).Methods("POST")

	r.HandleFunc("/mystudent", studentRouter)

	return r
}

func studentCreate(w http.ResponseWriter, r *http.Request) {

	NewStudent := Student{}
	NewStudent.ID, _ = strconv.Atoi(r.FormValue("id"))

	NewStudent.Name = r.FormValue("name")

	fmt.Println(NewStudent)

	output, _ := json.Marshal(NewStudent)

	fmt.Println(string(output))

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	database, err := sql.Open("postgres", psqlconn)

	if err != nil {
		fmt.Println("Something went wrong!")
	}

	defer database.Close()
	sql := "INSERT INTO school.student (id,name) values (" + strconv.Itoa(NewStudent.ID) + ", '" + NewStudent.Name + "' );"
	q, err := database.Exec(sql)
	if err != nil {
		fmt.Println(err)
		getError(w, err)
	}
	fmt.Println(q)

	http.Redirect(w, r, "/student", http.StatusSeeOther)

}

func studentRouter(w http.ResponseWriter, r *http.Request) {
	//A sample data for testing purpose
	myStudent := Student{}
	myStudent.ID = 2005
	myStudent.Name = "King Kong"

	output, _ := yaml.Marshal(&myStudent)
	fmt.Fprintln(w, string(output))

}

type Student struct {
	ID   int    `json:"id"` //key is id
	Name string `json:"name"`
}

func getError(w http.ResponseWriter, err error) {
	fmt.Println("Error", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))

}

func getStudentHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Pragma", "no-cache")
	var id int
	var name string

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		getError(w, err)
		return
	}
	CheckError(err)
	fmt.Println("Connected!")
	rows, err := db.Query("select id, name from school.student ")
	if err != nil {
		getError(w, err)
		return
	}
	defer rows.Close()
	students := []Student{}
	for rows.Next() {
		rows.Scan(&id, &name)
		student := Student{id, name}
		students = append(students, student)

	}
	fmt.Printf("Length  %d Students : %+v", len(students), students)

	defer db.Close()

	err = db.Ping()
	CheckError(err)

	studentListBytes, err := json.Marshal(students)

	if err != nil {
		getError(w, err)
		return
	}

	w.Write(studentListBytes)
}

func CheckError(err error) {
	if err != nil {
		panic(err)

	}
}
