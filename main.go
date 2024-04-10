package main

import (
	"FreeTime/model"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	TotalWeeks = 22
	StartDay   = "2024年02月26日" //学期开始时间
	shortForm  = "2006年01月02日" //用于格式化
)

// ExtractHourFromTime 用于解析08:20-10:00这种时间
func ExtractHourFromTime(timeStr string) (string, string) {
	// 分割时间段字符串，获取开始时间和结束时间
	parts := strings.Split(timeStr, "-")
	if len(parts) != 2 {
		return "", "" // 如果分割结果不是两部分，返回空字符串
	}

	startTime := parts[0]
	endTime := parts[1]

	// 提取开始时间和结束时间的小时部分（去掉前导零）
	startHour := strings.TrimLeft(strings.Split(startTime, ":")[0], "0")
	endHour := strings.TrimLeft(strings.Split(endTime, ":")[0], "0")

	return startHour, endHour
}

// CheckExamTime 格式化考试时间 ，返回的参数分别为：在第几周考试，星期几，考试开始小时，考试结束小时
func CheckExamTime(exam string) *model.Exam {
	str := `2024年04月18日 08:20-10:00 旗山西3-303`

	dateRe := regexp.MustCompile(`\d{4}年\d{2}月\d{2}日`)
	timeRe := regexp.MustCompile(`\d{2}:\d{2}-\d{2}:\d{2}`)

	date := dateRe.FindString(str)
	examTime := timeRe.FindString(str)
	startHour, endHour := ExtractHourFromTime(examTime)
	startHours, _ := strconv.Atoi(startHour)
	endHours, _ := strconv.Atoi(endHour)

	// 解析日期
	t1, _ := time.Parse(shortForm, StartDay)
	t2, _ := time.Parse(shortForm, date)

	// 计算日期差
	week := int(t2.Sub(t1).Hours() / 24)
	weekday := int(t2.Weekday())
	week /= 7
	exams := &model.Exam{
		Week:      week + 1,
		Weekday:   weekday,
		StartHour: startHours,
		EndHour:   endHours,
	}

	return exams
}

// GetFreeTimeInOneWeek 根据weekId获取单独一周的空闲时间的交集
func GetFreeTimeInOneWeek(weekId int, table []*model.Table) []string {
	return table[weekId].FindFreeTime()
}

// ExgCourse 格式化输入课表
func ExgCourse(time string) *model.Course {
	res := exg(time)
	course := &model.Course{
		StartWeek: res[0],
		EndWeek:   res[1],
		Day:       res[2],
		StartTime: res[3],
		EndTime:   res[4],
	}
	return course
}

// ExgAdjust 格式化输入调课信息
func ExgAdjust(time string) (oldCourse, newCourse *model.Course) {
	res := exg(time)
	oldCourse = &model.Course{
		StartWeek: res[0],
		EndWeek:   res[1],
		Day:       res[2],
		StartTime: res[3],
		EndTime:   res[4],
	}
	newCourse = &model.Course{
		StartWeek: res[5],
		EndWeek:   res[6],
		Day:       res[7],
		StartTime: res[8],
		EndTime:   res[9],
	}
	return oldCourse, newCourse
}

// 把时间信息转成[]int数组
func exg(time string) []int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(time, -1)
	res := make([]int, 0)
	// 去除前导零并转换为整数
	for _, match := range matches {
		match = strings.TrimLeft(match, "0")
		num, err := strconv.Atoi(match)
		if err != nil {
			fmt.Printf("Error converting string to int: %v\n", err)
			continue
		}
		res = append(res, num)
	}
	return res
}

// 格式化输出
func readData(week int, data []string) {

	groupedData := make(map[int][]int)

	for _, item := range data {
		parts := strings.Split(item, "-")
		if len(parts) == 2 {
			day, _ := strconv.Atoi(parts[0])
			hour, _ := strconv.Atoi(parts[1])

			groupedData[day] = append(groupedData[day], hour)
		}
	}

	for day, hours := range groupedData {
		start := hours[0]
		for k, i := range hours {
			if k > 0 {
				//看小时是否连续，如果不连续，说明这段时间是忙碌的，可以把之前空闲的时间提取出来
				if hours[k]-hours[k-1] != 1 {
					fmt.Printf("第%d周：星期%d,%d-%d\n", week+1, day, start, hours[k-1])
					start = i
				}
			}
			//到24时截止
			if k == len(hours)-1 {
				fmt.Printf("第%d周：星期%d,%d-%d\n", week+1, day, start, hours[k])
			}
		}
	}
}

func main() {
	tables := make([]*model.Table, 0)

	for i := 0; i < TotalWeeks; i++ {
		table := model.NewTable()
		tables = append(tables, table)
	}

	//这个是课表的时间
	s := make([]string, 0)
	s = append(s, "02-03 星期4:5-6节")
	s = append(s, "10-11 星期6:1-8节")
	s = append(s, "01-16 星期1:5-6节")

	//m把节转成小时
	m := map[int]int{
		1:  8,
		2:  9,
		3:  10,
		4:  11,
		5:  14,
		6:  15,
		7:  16,
		8:  17,
		9:  19,
		10: 20,
		11: 21,
	}
	//格式化输入并添加忙碌时间
	for _, busyTime := range s {
		course := ExgCourse(busyTime)
		for i := course.StartWeek; i <= course.EndWeek; i++ {
			tables[i-1].AddBusyTime(m[course.StartTime], m[course.EndTime], course.Day)
		}

	}
	//添加考试时间
	exams := CheckExamTime("")
	//fmt.Println(exams.Week, exams.StartHour, exams.EndHour, exams.Weekday)
	tables[exams.Week].AddBusyTime(exams.StartHour, exams.EndHour, exams.Weekday)
	//获取结果并格式化
	for i := 0; i < TotalWeeks; i++ {
		readData(i, tables[i].FindFreeTime())
	}
}
