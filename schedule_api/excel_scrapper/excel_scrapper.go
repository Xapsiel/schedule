package excel_scrapper

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func ERROR(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func read_file(PATH string) [][]string {
	f, err := excelize.OpenFile(PATH)
	ERROR(err)
	defer func() {
		err := f.Close()
		ERROR(err)
	}()
	rows, err := f.GetRows("курс 1 ПИ ")
	ERROR(err)
	return rows
}
func to_pretty(elem []string) []string {
	for i := 0; i < 20; i++ {
		elem = append(elem, "")
	}
	elem = elem[0:22]
	elem = slices.Delete(elem, 11, 14) //"сращиваем обе части таблицы(четную и нечетную)
	lenght := len(elem)
	for i := 0; i < 20-lenght; i++ {
		elem = append(elem, "") //дополняем массив, чтобы в случае чего были обнаружены "окна" в парах
	}
	return elem
}

func read_schedule() map[string]map[int][][]string {
	PATH := "schedule_api/excel_scrapper/PI.xlsx"
	rows := read_file(PATH)
	week := []string{"понедельник", "вторник", "среда", "четверг", "пятница", "суббота"}
	lessons_type := []string{"лк", "пз"} //тип пары
	//день недели
	lt := ""                                                                                                                                                                                                                       // Тип пары_2
	together := []string{"высшая математика", "основы российской государственности", "иностранный язык(немецкий,французкий)", "физическая культура", "история россии/с 13.10", "история россии/с 09.10", "история россии/с 16.10"} //предметы для которых практические вместе

	var all_info map[string]map[int][][]string //вся информация. КЛЮЧ_1 - День. Ключ_2-номер пары, значение-массив pair
	all_info = make(map[string]map[int][][]string)

	comment := ""   //комментарий преподавателя
	day := ""       //день недели
	number := 1     //номер пары
	subject_1 := "" //предмет первой подгруппы
	teacher_1 := "" //преподаватель первой подгруппы
	address_1 := "" //адрес первой подгруппы

	subject_2 := "" //предмет второй подгруппы
	teacher_2 := "" //преподаватель второй подгруппы
	address_2 := "" //адрес второй подгруппы

	var pair map[int][][]string //массив всех предметов за number пару у каждой группы/подгруппы
	pair = make(map[int][][]string)

	var new_day_index []int //при смене дня недели записывает индекс строки, где он меняется
	for index, elem := range rows {

		if index%5 != 4 { //ускоряет работы, минуя просмотр ненужных строк
			continue
		}
		if len(elem) == 0 { //если длина строки столбца равна 0, то мы ее скипаем
			continue
		}

		teach := rows[index+1]
		addr := rows[index+2]

		teach = to_pretty(teach)
		addr = to_pretty(addr)
		elem = to_pretty(elem)

		if slices.Contains(week, strings.ToLower(elem[0])) {
			new_day_index = append(new_day_index, index) //если первый элемент строки-день неделю, то его индекс записываем
		}
		if len(new_day_index) != 0 { //если длинна не равна 0, то значит, что как минимум понедельник уже был
			d := strings.ToLower(rows[new_day_index[len(new_day_index)-1]][0])
			if day != d { //сравниваем нынешний день с последним записанным

				if day != "" { //если день не пуст, то значит предыдущий день закончился
					all_info[strings.ToLower(day)] = pair //записываем данные
					pair = make(map[int][][]string)       //создаем новый массив для этого дня

				}
				number = 1                                                          //начинаем с 1 пары
				day = strings.ToLower(rows[new_day_index[len(new_day_index)-1]][0]) //начинаем новый день
				if day == "суббота" {
					return all_info
				}
			} else if strings.Contains("1234567890", elem[1]) {

				number++ //если в элементе есть числа-элемент содержит номер пары -> увеличиваем на единицу
			}

			for i := 2; i < len(elem); i++ { //начинаем сразу со второго индекса, чтобы начать с "типа лекции"
				nlt := strings.ToLower(elem[i]) //нынешний элемент

				if slices.Contains(lessons_type, nlt) { //если он показывает тип, то выясняем какой именно
					if strings.Contains(" лк ", nlt) {
						lt = "лк"
						continue
					} else if strings.Contains(" пз ", nlt) {
						lt = "пз"
						continue
					}
				}
				if lt == "лк" || (lt == "пз" && slices.Contains(together, nlt)) { //если тип лекции или ПЗ,но те,которые проходят всей группой,то
					subject_1 = elem[i] //записываем их для обеих подгрупп
					teacher_1 = teach[i]
					address_1 = addr[i]

					subject_2 = subject_1
					teacher_2 = teacher_1
					address_2 = address_1

				} else if lt == "пз" { //если обычное пз, то взависимости от того, кому принадлежат
					teach = append(teach, "") //чтобы избежать ошибки при попытке доступа к элементу
					addr = append(addr, "")   //добавляем пустую строку
					subject_1 = elem[i]
					teacher_1 = teach[i]
					address_1 = addr[i]

					subject_2 = elem[i+1]
					teacher_2 = teach[i+1]
					address_2 = addr[i+1]

				} else if lt == "" && (i == 2 || i == 5 || i == 8 || i == 11 || i == 14 || i == 17) { // если данная ячейка-это ячейка, в которой хранится тип
					//в таком случае, если она пуста, то пар на эту пару нет
					subject_1 = "" //предмет первой подгруппы
					teacher_1 = "" //преподаватель первой подгруппы
					address_1 = "" //адрес первой подгруппы

					subject_2 = "" //предмет второй подгруппы
					teacher_2 = "" //преподаватель второй подгруппы
					address_2 = "" //адрес второй подгруппы
				} else if lt == "" {
					continue
				}
				pair[number] = append(pair[number], []string{subject_1, teacher_1, address_1, lt, comment})
				pair[number] = append(pair[number], []string{subject_2, teacher_2, address_2, lt, comment})
				lt = ""
			} //Информация за number пару

		}
	}
	return all_info

}

type Teacher_info struct {
	Lessons      map[int][][]string
	Teacher_name string
	WeekType     string
	Day          string
	Date_day     int
	Date_month   int
	Date_year    int
}
type Student_info struct {
	Group      string
	WeekType   string
	Day        string
	Lessons    map[int][]string
	Date_day   int
	Date_month int
	Date_year  int
}

func Code(ti Teacher_info) string {
	date_day := strconv.Itoa(ti.Date_day)
	date_month := strconv.Itoa(ti.Date_month)
	date_year := strconv.Itoa(ti.Date_year)
	name := ti.Teacher_name
	return date_day + "-" + date_month + "-" + date_year + "-" + name
}

func Update() ([]Student_info, []Teacher_info) { //делаем более красивый вид ДБ
	day_number := map[string]int{
		"понедельник": 0,
		"вторник":     1,
		"среда":       2,
		"четверг":     3,
		"пятница":     4,
	}
	timing := map[int]string{
		1: "8:00",
		2: "9:50",
		3: "11:30",
		4: "13:20",
		5: "15:00",
		6: "16:40",
		7: "18:20",
		8: "20:00",
	}
	var teachers map[string]Teacher_info
	teachers = make(map[string]Teacher_info)
	//можно было сразу, но мне уже страшно переделывать код
	all_info := read_schedule()
	var res []Student_info
	when_monday := int(time.Now().Weekday()) - 1
	current_date := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-when_monday, 12, 0, 0, 0, time.UTC)
	if !is_nechet(current_date) {
		current_date = current_date.Add(time.Hour * (-24) * 7)
	}

	groups := []string{"231-1", "231-2", "232-1", "232-2", "233-1", "233-2"}
	for times := 0; times < 3; times++ {
		if !is_nechet(current_date) {
			current_date = current_date.Add(time.Hour * (-24) * 7)
		}
		for i := 0; i < len(groups); i++ {
			for day := range all_info {
				num_of_day := day_number[day]

				var chet_stud_schedule, nechet_stud_schedule Student_info
				elem := groups[i]

				chet_stud_schedule.Date_day = current_date.Add(time.Hour * time.Duration(24*7+24*num_of_day)).Day()
				chet_stud_schedule.Date_month = int(current_date.Add(time.Hour * time.Duration(24*7+24*num_of_day)).Month())
				chet_stud_schedule.Date_year = current_date.Add(time.Hour * time.Duration(24*7+24*num_of_day)).Year()
				chet_stud_schedule.Day = day
				chet_stud_schedule.Group = elem
				chet_stud_schedule.Lessons = make(map[int][]string)
				chet_stud_schedule.WeekType = "четная"

				nechet_stud_schedule.Date_day = current_date.Add(time.Hour * time.Duration(24*num_of_day)).Day()
				nechet_stud_schedule.Date_month = int(current_date.Add(time.Hour * time.Duration(24*num_of_day)).Month())
				nechet_stud_schedule.Date_year = current_date.Add(time.Hour * time.Duration(24*num_of_day)).Year()
				nechet_stud_schedule.Group = elem
				nechet_stud_schedule.Day = day
				nechet_stud_schedule.Lessons = make(map[int][]string)
				nechet_stud_schedule.WeekType = "нечетная"

				for number, subject := range all_info[day] {
					if subject[i+6][1] == "" && subject[i+6][3] == "лк" {
						subject[i+6][1] = subject[i][1]
					}
					chet_stud_schedule.Lessons[number] = subject[i+6]
					chet_stud_schedule.Lessons[number] = append(chet_stud_schedule.Lessons[number], timing[number])
					nechet_stud_schedule.Lessons[number] = subject[i]
					nechet_stud_schedule.Lessons[number] = append(nechet_stud_schedule.Lessons[number], timing[number])

					if subject[i+6][1] != "" {
						chet_teacher := Teacher_info{Teacher_name: strings.Join(strings.Split(subject[i+6][1], " "), ""),
							WeekType:   chet_stud_schedule.WeekType,
							Day:        chet_stud_schedule.Day,
							Date_day:   chet_stud_schedule.Date_day,
							Date_month: chet_stud_schedule.Date_month,
							Date_year:  chet_stud_schedule.Date_year,
						}
						cod := Code(chet_teacher)

						if len(teachers[cod].Lessons) == 0 {
							chet_teacher.Lessons = make(map[int][][]string)
							teachers[cod] = chet_teacher
						}
						teachers[cod].Lessons[number] = append(teachers[cod].Lessons[number], []string{subject[i+6][0], elem, subject[i+6][2], subject[i+6][3], subject[i+6][4], timing[number]})
					}
					if subject[i][1] != "" {
						nechet_teacher := Teacher_info{Teacher_name: strings.Join(strings.Split(subject[i][1], " "), ""),
							WeekType:   nechet_stud_schedule.WeekType,
							Day:        nechet_stud_schedule.Day,
							Date_day:   nechet_stud_schedule.Date_day,
							Date_month: nechet_stud_schedule.Date_month,
							Date_year:  nechet_stud_schedule.Date_year,
						}
						cod := Code(nechet_teacher)
						if len(teachers[cod].Lessons) == 0 {
							nechet_teacher.Lessons = make(map[int][][]string)
							teachers[cod] = nechet_teacher
						}
						teachers[cod].Lessons[number] = append(teachers[cod].Lessons[number], []string{subject[i][0], elem, subject[i][2], subject[i][3], subject[i][4], timing[number]})
					}
				}
				res = append(res, chet_stud_schedule) //добавляем инфу
				res = append(res, nechet_stud_schedule)

			}

		}
		current_date = current_date.Add(time.Hour * 24 * 14)

	}
	var teacher_result []Teacher_info
	for _, key := range teachers {
		teacher_result = append(teacher_result, key)
	}

	student, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	teacher, err := json.Marshal(teacher_result)
	if err != nil {
		panic(err)
	}
	file, err := os.Create("schedule_api/db/students.json")
	if err != nil {
		panic(err)
	}
	file.Write(student)

	file, err = os.Create("schedule_api/db/teachers.json")
	if err != nil {
		panic(err)
	}
	file.Write(teacher)
	return res, teacher_result
}

func is_nechet(date time.Time) bool {
	nechet_weeks := "11.09-17.09; 25.09-01.10; 09.10-15.10; 23.10-29.10; 06.11-12.11; 20.11-26.11; 04.12-10.12; 18.12-24.12"
	weeks := strings.Split(nechet_weeks, "; ")
	for _, value := range weeks {
		period := strings.Split(value, "-")
		start := strings.Split(period[0], ".")
		s_day, _ := strconv.Atoi(start[0])
		s_month, _ := strconv.Atoi(start[1])
		end := strings.Split(period[1], ".")
		e_day, _ := strconv.Atoi(end[0])
		e_month, _ := strconv.Atoi(end[1])
		start_date := time.Date(date.Year(), time.Month(s_month), s_day, 0, 0, 0, 0, time.UTC)
		end_date := time.Date(date.Year(), time.Month(e_month), e_day, 0, 0, 0, 0, time.UTC)
		if e_month < s_month {
			end_date = time.Date(end_date.Year()+1, time.Month(e_month), e_day, 24, 0, 0, 0, time.UTC)
		}
		if date.Before(end_date) && start_date.Before(date) {
			return true
		}

	}
	return false
}
func is_chet(date time.Time) bool {
	chet_weeks := "18.09-24.09; 02.10-08.10; 16.10-22.10; 30.10-05.11; 13.11-19.11; 27.11-03.12; 11.12-17.12; 25.12-31.12"
	weeks := strings.Split(chet_weeks, "; ")
	for _, value := range weeks {
		period := strings.Split(value, "-")
		start := strings.Split(period[0], ".")
		s_day, _ := strconv.Atoi(start[0])
		s_month, _ := strconv.Atoi(start[1])
		end := strings.Split(period[1], ".")
		e_day, _ := strconv.Atoi(end[0])
		e_month, _ := strconv.Atoi(end[1])
		start_date := time.Date(date.Year(), time.Month(s_month), s_day, 0, 0, 0, 0, time.UTC)
		end_date := time.Date(date.Day(), time.Month(e_month), e_day, 0, 0, 0, 0, time.UTC)
		if e_month < s_month {
			end_date = time.Date(end_date.Year()+1, time.Month(e_month), e_day, 0, 0, 0, 0, time.UTC)
		}
		if date.Before(end_date) && start_date.Before(date) {
			return true
		}
	}
	return false
}
