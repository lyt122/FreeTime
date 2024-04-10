package model

import "fmt"

const (
	DayPerWeek  = 7  //一周7天
	HoursPerDay = 24 //一天24时
)

type Table struct {
	table [][]int
}

// NewTable 这里table其实是以一周的时间为单位的，
func NewTable() *Table {
	table := make([][]int, DayPerWeek+1)
	for i := range table {
		table[i] = make([]int, HoursPerDay+1)
	}
	return &Table{table}
}

// AddBusyTime 添加忙碌时间
func (ag *Table) AddBusyTime(startHour, endHour, day int) {
	for i := startHour; i <= endHour; i++ {
		ag.table[day][i] += 1
	}
}

// FindFreeTime 获取最终的空闲时间
func (ag *Table) FindFreeTime() []string {
	var freeTimes []string
	for day := 1; day <= DayPerWeek; day++ {
		for hour := 0; hour <= HoursPerDay; hour++ {
			if ag.table[day][hour] == 0 {
				freeTimes = append(freeTimes, fmt.Sprintf("%d-%d", day, hour))
			}
		}
	}
	return freeTimes
}

// Adjust 处理调课情况
func (ag *Table) Adjust(oStartHour, oEndHour, oDay, startHour, endHour, day int) {
	for i := oStartHour; i <= oEndHour; i++ {
		ag.table[oDay][i] -= 1
	}
	for i := startHour; i <= endHour; i++ {
		ag.table[day][i] += 1
	}
}
