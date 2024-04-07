package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// model for course - file
type Course struct {
	CourseId    string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"courseprice"`
	Author      *Author `json:"author"`
}

type Author struct {
	FullName string `json:"fullname"`
	Website  string `json:"website"`
}

// fake DB
var courses []Course

// middleware, helper - file
func (c *Course) IsEmpty() bool {
	return c.CourseName == ""
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to the home page</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all courses")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get single course")
	w.Header().Set("content-type", "application/json")

	// grab id from request
	params := mux.Vars(r)
	fmt.Println(params)

	// loop through courses to find the course with the id
	for _, course := range courses {
		if course.CourseId == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("Course not found")
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("create course")
	w.Header().Set("content-type", "application/json")

	// if body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send course details")
		return
	}

	// if {}
	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course)

	// if course.IsEmpty() {
	// 	json.NewEncoder(w).Encode("Please send course with id and name")
	// 	return
	// }

	// generate unique id, string
	// append course to courses

	course.CourseId = strconv.Itoa(rand.Intn(1000000))
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
}

func updateCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update course")
	w.Header().Set("content-type", "application/json")

	// grab id from request
	params := mux.Vars(r)

	// loop through courses to find the course with the id
	// reomve the course
	// add with my id and new course details
	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			course.CourseId = params["id"]
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}

	// if course not found
	json.NewEncoder(w).Encode("Course not found")
}

func deleteCourse(w http.ResponseWriter, r *http.Request) {

	fmt.Println("delete course")
	w.Header().Set("content-type", "application/json")

	// grab id from request
	params := mux.Vars(r)

	// loop through courses to find the course with the id
	// reomve the course
	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode("Course deleted")
			return
		}
	}

	// if course not found
	json.NewEncoder(w).Encode("Course not found")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/courses/{id}", getCourse).Methods("GET")
	r.HandleFunc("/courses", createCourse).Methods("POST")
	r.HandleFunc("/courses/{id}", updateCourse).Methods("PUT")
	r.HandleFunc("/courses/{id}", deleteCourse).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":4000", r))
}
