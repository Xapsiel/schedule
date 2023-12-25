package API

import (
	"os"
	"schedule/schedule_api/db"
	"schedule/schedule_api/scrapper"
)

func Update(URL string, old_link string) string {
	schedule, teacher_schedule, update_link := scrapper.Parse(URL, old_link)
	if update_link == old_link {
		return update_link
	}
	file, err := os.Create("schedule_api/link.txt")
	if err != nil {
		panic(err)
	}
	file.Write([]byte(update_link))

	db.Make_db(schedule, teacher_schedule)
	return update_link
}

func Get_info_about(group string, year int, month int, day int) string {
	return db.Info_about(group, year, month, day)
}
func NextStudentPair(group string) string {
	return db.NextStudentPair(group)
}
func NextTeacherPair(name string) string {
	return db.NextTeacherPair(name)
}

func Teacher(name string, year int, month int, day int) string {
	return db.About_teacher(name, year, month, day)
}
