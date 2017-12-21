package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type TimeTableJSON struct {
	From      string        `json:"from"`
	To        string        `json:"to"`
	TimeTable TimeTableList `json:"time-table"`
}

type TimeTableList struct {
	Weekday  []TimeTable `json:"weekday"`
	Saturday []TimeTable `json:"saturday"`
	Holiday  []TimeTable `json:"holiday"`
}

type TimeTable struct {
	Hour    int   `json:"hour"`
	Minutes []int `json:"minutes"`
}

/* Main routine */
func main() {
	var timetableJSON TimeTableJSON
	var timetableList TimeTableList

	if len(os.Args) == 2 {
		fmt.Println("目的地を指定して下さい")
		fmt.Println(os.Args[0] + " " + os.Args[1] + " <目的地>")
		os.Exit(1)
	} else if len(os.Args) != 3 {
		fmt.Println("出発地と目的地を指定して下さい")
		fmt.Println(os.Args[0] + " <出発地> <目的地>")
		os.Exit(1)
	}

	timetableJSON.From = os.Args[1]
	timetableJSON.To = os.Args[2]
	fmt.Println("出発地: " + timetableJSON.From)
	fmt.Println("目的地: " + timetableJSON.To)

	baseURL := "http://gps.iwatebus.or.jp/bls/pc/"
	fromPath := "jyosha.jsp"

	toPath, state := getNextURL(baseURL+fromPath, timetableJSON.From)
	if state == false {
		fmt.Println("Error: 指定された出発地が見つかりません")
		fmt.Println()
		os.Exit(1)
	}

	allTimetablePath, state := getNextURL(baseURL+toPath, timetableJSON.To)
	if state == false {
		fmt.Println("Error: 指定された目的地が見つかりません")
		fmt.Println()
		os.Exit(1)
	}

	weekdayTimetablePath, _ := getNextURL(baseURL+allTimetablePath, "平日")
	weekdayTimeTable := timetableParseToJSON(baseURL + weekdayTimetablePath)
	timetableList.Weekday = weekdayTimeTable

	saturdayTimetablePath, _ := getNextURL(baseURL+allTimetablePath, "土曜")
	saturdayTimeTable := timetableParseToJSON(baseURL + saturdayTimetablePath)
	timetableList.Saturday = saturdayTimeTable

	holidayTimetablePath, _ := getNextURL(baseURL+allTimetablePath, "休日")
	holidayTimeTable := timetableParseToJSON(baseURL + holidayTimetablePath)
	timetableList.Holiday = holidayTimeTable

	timetableJSON.TimeTable = timetableList

	outputJSON, _ := json.MarshalIndent(timetableJSON, "", "\t")
	writeFile("timetable_"+timetableJSON.From+"_"+timetableJSON.To+".json", outputJSON)
	fmt.Println("")
}

/* Get Next URL */
func getNextURL(URL string, target string) (foundURL string, state bool) {
	htmlSource := getHtmlSource(URL)

	regex := "<a href=\"([\\w\\/:%#\\$&\\?\\(\\)~\\.=\\+\\-]+)\" ?(target=\"_blank\")? ?>" + target + "<\\/a>"
	searchRegex := regexp.MustCompile(regex)
	found := searchRegex.FindStringSubmatch(htmlSource)

	if len(found) != 0 {
		foundURL = found[1]
		state = true
	} else {
		state = false
	}

	return
}

/* generate to JSON from diagram */
func timetableParseToJSON(URL string) (timetables []TimeTable) {
	htmlSource := getHtmlSource(URL)
	r := strings.NewReplacer("\r", "", "\n", "", "\t", "")
	htmlSource = r.Replace(htmlSource)

	for i := 1; i <= 24; i++ {
		regex := "<strong>" + checkDigit(i) + "<\\/strong><\\/div><\\/td>[ ]+<td class=\"dya-min-(even|odd)\">(\xc2\xa0)+([0-9\\s]+)(\xc2\xa0)+<\\/td>"
		searchRegex := regexp.MustCompile(regex)
		rawDiagrams := searchRegex.FindStringSubmatch(htmlSource)
		if len(rawDiagrams) == 0 {
			continue
		}

		strDiagrams := strings.Split(rawDiagrams[3], " ")
		var intDiagrams []int
		for _, strDiagram := range strDiagrams {
			intDiagram, err := strconv.Atoi(strDiagram)
			if err != nil {
				continue
			}
			intDiagrams = append(intDiagrams, intDiagram)
		}
		if len(intDiagrams) != 0 {
			var timetable TimeTable
			timetable.Hour = i
			timetable.Minutes = intDiagrams
			timetables = append(timetables, timetable)
		}
	}
	return
}

func getHtmlSource(URL string) string {
	res, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	utfBody := transform.NewReader(bufio.NewReader(res.Body), japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		panic(err)
	}
	htmlSource, _ := doc.Html()
	return html.UnescapeString(htmlSource)
}

func checkDigit(num int) string {
	if num < 10 {
		return "0" + strconv.Itoa(num)
	}
	return strconv.Itoa(num)
}

/* Create and write JSON */
func writeFile(path string, data []byte) {
	file, _ := os.Create(path)
	defer file.Close()
	file.Write(data)
	fmt.Println("JSON generate to " + path + " complete.")
	fmt.Println()
}
