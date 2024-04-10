package model

import "fmt"

const (
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
		ag.table[day][i] += 1
	}
}

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
