package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"schedule/schedule_api/API"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/schedule/update", updateHendler)

	router.HandleFunc("/schedule/student/next/{group_number}", nextStudentHendler)
	router.HandleFunc("/schedule/student/pairInfo/{year}/{month}/{day}/{group_number}", groupHendler)

	router.HandleFunc("/schedule/teacher/pairInfo/{year}/{month}/{day}/{name}", teacherHandler)
	router.HandleFunc("/schedule/teacher/next/{name}", nextTeacherHandler)
	http.ListenAndServe(":80", router)
}

func nextStudentHendler(rw http.ResponseWriter, req *http.Request) {
	if time.Now().Hour() == 0 {
		updateHendler(rw, req)
	}
	value := mux.Vars(req)
	group := value["group_number"]
	json_data := API.NextStudentPair(group)
	rw.Header().Set("Content-Type", "application/json") //устанавливаем какой будет контент страницы
	rw.Write([]byte(json_data))
}
func teacherHandler(rw http.ResponseWriter, req *http.Request) {
	if time.Now().Hour() == 0 {
		updateHendler(rw, req)
	}
	value := mux.Vars(req)
	year, _ := strconv.Atoi(value["year"])
	month, _ := strconv.Atoi(value["month"])
	day, _ := strconv.Atoi(value["day"])
	name := value["name"]
	json_data := API.Teacher(name, year, month, day)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(json_data))
}
func nextTeacherHandler(rw http.ResponseWriter, req *http.Request) {
	if time.Now().Hour() == 0 {
		updateHendler(rw, req)
	}
	value := mux.Vars(req)
	name := value["name"]

	json_data := API.NextTeacherPair(name)
	rw.Header().Set("Content-Type", "application/json") //устанавливаем какой будет контент страницы
	rw.Write([]byte(json_data))
}

func groupHendler(rw http.ResponseWriter, req *http.Request) {
	if time.Now().Hour() == 0 {
		updateHendler(rw, req)
	}
	value := mux.Vars(req)
	group := value["group_number"]
	year, _ := strconv.Atoi(value["year"])
	month, _ := strconv.Atoi(value["month"])
	day, _ := strconv.Atoi(value["day"])
	json_data := API.Get_info_about(group, year, month, day) //получаем наш документ формата json
	rw.Header().Set("Content-Type", "application/json")      //устанавливаем какой будет контент страницы
	rw.Write([]byte(json_data))                              //переводим в массив байтов

}
func updateHendler(_ http.ResponseWriter, _ *http.Request) {
	file, err := ioutil.ReadFile("schedule_api/link.txt")
	if err != nil {
		os.Create("schedule_api/link.txt")
	}

	old_link := string(file)
	old_link = API.Update("https://cfuv.ru/raspisanie-fakultativov-fiziko-tekhnicheskogo-instituta", old_link)

}
