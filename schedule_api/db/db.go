package db

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"schedule/schedule_api/excel_scrapper"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var st []excel_scrapper.Student_info
var tc []excel_scrapper.Teacher_info

func main() {

}

func Make_db(data []excel_scrapper.Student_info, teachers_data []excel_scrapper.Teacher_info) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://projectlaim2023:jTpSqRamIKn3UTT2@cluster0.lxtqivz.mongodb.net/").SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	client.Database("CFU").Collection("schedule").Drop(context.TODO()) //ощищаем нашу коллекцию с расписанием
	collection := client.Database("CFU").Collection("schedule")
	for _, e := range data {
		collection.InsertOne(context.Background(), e) //записываем в нее новое
	}
	client.Database("CFU").Collection("teachers").Drop(context.TODO()) //ощищаем нашу коллекцию с расписанием

	collection = client.Database("CFU").Collection("teachers")
	for _, e := range teachers_data {
		collection.InsertOne(context.Background(), e) //записываем в нее новое
	}
}

func Info_about(group string, year int, month int, day int) string {
	file, err := ioutil.ReadFile("schedule_api/db/students.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &st)
	if err != nil {
		panic(err)
	}

	for _, elem := range st {
		if elem.Group == group && elem.Date_day == day && elem.Date_month == month && elem.Date_year == year {
			json_data, err := json.Marshal(elem)
			if err != nil {
				var a []byte

				return string(a)
			}

			return string(json_data)
		}
	}
	var a []byte
	return string(a)
}

func NextStudentPair(group string) string {
	file, err := ioutil.ReadFile("schedule_api/db/students.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &st)
	if err != nil {
		panic(err)
	}

	timing := []int{8 * 60, 9*60 + 50, 11*60 + 30, 13*60 + 20, 15 * 60, 16*60 + 40, 18*60 + 20, 20 * 60}

	current_day := time.Now()
	minute := current_day.Minute() + current_day.Hour()*60
	timing = append(timing, minute)
	sort.Ints(timing)
	pair_number := -1
	for index := 0; index < len(timing); index++ {
		if timing[index] == minute {
			pair_number = index + 1
		}
	}
	if pair_number == 9 {
		pair_number = 2
	}
	if pair_number != -1 {
		for _, elem := range st {
			if elem.Group == group && elem.Date_day == current_day.Day() && elem.Date_month == int(current_day.Month()) && elem.Date_year == current_day.Year() {
				address := elem.Lessons[pair_number]
				json_data, err := json.Marshal(address)
				if err != nil {
					var a []byte
					return string(a)
				}

				return string(json_data)
			}
		}
	}

	var a []byte
	return string(a)

}

func About_teacher(name string, year int, month int, day int) string {
	file, err := ioutil.ReadFile("schedule_api/db/teachers.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &tc)
	if err != nil {
		panic(err)
	}
	for _, elem := range tc {
		if elem.Teacher_name == name && elem.Date_day == day && elem.Date_month == month && elem.Date_year == year {
			json_data, err := json.Marshal(elem)
			if err != nil {
				var a []byte

				return string(a)
			}

			return string(json_data)
		}
	}

	var a []byte

	return string(a)

}

func NextTeacherPair(name string) string {
	file, err := ioutil.ReadFile("schedule_api/db/teachers.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &tc)
	if err != nil {
		panic(err)
	}
	timing := []int{8 * 60, 9*60 + 50, 11*60 + 30, 13*60 + 20, 15 * 60, 16*60 + 40, 18*60 + 20, 20 * 60}

	current_day := time.Now()
	minute := current_day.Minute() + current_day.Hour()*60
	timing = append(timing, minute)
	sort.Ints(timing)
	pair_number := -1
	for index := 0; index < len(timing); index++ {
		if timing[index] == minute {
			pair_number = index + 1
		}
	}
	if pair_number == 9 {
		pair_number = 2
	}
	if pair_number != -1 {
		if pair_number != -1 {
			for _, elem := range tc {
				if elem.Teacher_name == name && elem.Date_day == current_day.Day() && elem.Date_month == int(current_day.Month()) && elem.Date_year == current_day.Year() {
					address := elem.Lessons[pair_number]
					json_data, err := json.Marshal(address)
					if err != nil {
						var a []byte
						return string(a)
					}

					return string(json_data)
				}
			}
		}
	}

	var a []byte
	return string(a)

}
