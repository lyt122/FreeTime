package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	TotalWeeks  = 20 //一共有多少周
	DayPerWeek  = 7
	HoursPerDay = 24
)

type Table struct {
	table [][]int
}

func NewTable() *Table {
	table := make([][]int, DayPerWeek+1)
	for i := range table {
		table[i] = make([]int, HoursPerDay+1)
	}
	return &Table{table}
}

func (ag *Table) AddBusyTime(startHour, endHour, day int) {
	for i := startHour; i <= endHour; i++ {
		ag.table[day][i] = 1
	}
}

func (ag *Table) FindFreeTime() []string {
	var freeTimes []string
	for day := 1; day <= DayPerWeek; day++ {
		for hour := 1; hour <= HoursPerDay; hour++ {
			if ag.table[day][hour] == 0 {
				freeTimes = append(freeTimes, fmt.Sprintf("%d-%d", day, hour))
			}
		}
	}
	return freeTimes
}

// Exg 格式化输入
func Exg(time string) (int, int, int, int, int) {
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
	startWeek := res[0]
	endWeek := res[1]
	day := res[2]
	startTime := res[3]
	endTime := res[4]
	return startWeek, endWeek, day, startTime, endTime
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
				if hours[k]-hours[k-1] != 1 {
					fmt.Printf("第%d周：星期%d,%d-%d\n", week+1, day, start, hours[k-1])
					start = i
				}
			}
			if k == len(hours)-1 {
				fmt.Printf("第%d周：星期%d,%d-%d\n", week+1, day, start, hours[k])
			}
		}
	}
}

func main() {
	tables := make([]*Table, 0)

	for i := 0; i < TotalWeeks; i++ {
		table := NewTable()
		tables = append(tables, table)
	}
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
	for _, time := range s {
		startWeek, endWeek, day, startTime, endTime := Exg(time)
		for i := startWeek; i <= endWeek; i++ {
			tables[i-1].AddBusyTime(m[startTime], m[endTime], day)
		}

	}
	for i := 0; i < TotalWeeks; i++ {
		readData(i, tables[i].FindFreeTime())
	}
}
