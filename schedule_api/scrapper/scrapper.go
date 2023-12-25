package scrapper

import (
	"schedule/schedule_api/excel_scrapper"
	"strings"

	"github.com/gocolly/colly"
)

func getDowloandLink(URL string) string {
	scrapper := colly.NewCollector()
	var link string
	scrapper.OnHTML(".card-body", func(elem *colly.HTMLElement) {

		if strings.Contains(elem.Text, "Программная") {
			elem.ForEach("a", func(_ int, f_elem *colly.HTMLElement) {
				link = f_elem.Attr("href")
			})
		}
	})
	scrapper.Visit(URL)
	return link
}

func dowloand(Link string, old_link string) string {
	URL := getDowloandLink(Link)
	scrapper := colly.NewCollector()
	var link string
	flag := true
	scrapper.OnHTML("#downloadFile", func(elem *colly.HTMLElement) {
		link = elem.Attr("href")
		if old_link == link {
			flag = false
			return
		}
		dowl_scrapper := colly.NewCollector()
		dowl_scrapper.OnResponse(func(r *colly.Response) {
			if strings.Contains(link, "xlsx") {
				r.Save("schedule_api/excel_scrapper/PI.xlsx")
			}
		})
		dowl_scrapper.Visit(link)
	})

	scrapper.Visit(URL)
	if !flag {
		return old_link
	}
	return link
}

func Parse(Link string, old_link string) ([]excel_scrapper.Student_info, []excel_scrapper.Teacher_info, string) {
	update_link := dowloand(Link, old_link)
	if update_link == old_link {
		var t_1 []excel_scrapper.Student_info
		var t_2 []excel_scrapper.Teacher_info
		return t_1, t_2, old_link
	}
	schedule, teachers := excel_scrapper.Update()
	return schedule, teachers, update_link
}
