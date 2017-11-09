package main

import (
	"github.com/gizak/termui"
	"strconv"
	"time"
)

var weekDays = map[string]time.Weekday{
	"Mon": time.Monday,
	"Tue": time.Tuesday,
	"Wed": time.Wednesday,
	"Thu": time.Thursday,
	"Fri": time.Friday,
	"Sat": time.Saturday,
	"Sun": time.Sunday,
}

// months represent all the month choices available in a year
var months = map[string]time.Month{
	"Jan": time.January,
	"Feb": time.February,
	"Mar": time.March,
	"Apr": time.April,
	"May": time.May,
	"Jun": time.June,
	"Jul": time.July,
	"Aug": time.August,
	"Sep": time.September,
	"Oct": time.October,
	"Nov": time.November,
	"Dec": time.December,
}

var defaultHourlyProfile = map[string]float64{
	"0":  1,
	"1":  1,
	"2":  1,
	"3":  1,
	"4":  1,
	"5":  1,
	"6":  1,
	"7":  1,
	"8":  1,
	"9":  1,
	"10": 1,
	"11": 1,
	"12": 1,
	"13": 1,
	"14": 1,
	"15": 1,
	"16": 1,
	"17": 1,
	"18": 1,
	"19": 1,
	"20": 1,
	"21": 1,
	"22": 1,
	"23": 1,
}

var defaultWeeklyProfile = map[string]float64{
	"Sun": 1,
	"Mon": 1,
	"Tue": 1,
	"Wed": 1,
	"Thu": 1,
	"Fri": 1,
	"Sat": 1,
}

var defaultMonthlyProfile = map[string]float64{
	"Jan": 1,
	"Feb": 1,
	"Mar": 1,
	"Apr": 1,
	"May": 1,
	"Jun": 1,
	"Jul": 1,
	"Aug": 1,
	"Sep": 1,
	"Oct": 1,
	"Nov": 1,
	"Dec": 1,
}

func PlotBarChart(keys []int, labels []string, header string) {
	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()
	bc := termui.NewBarChart()
	bc.BorderLabel = header
	bc.Data = keys
	bc.Width = 170
	bc.Height = 15
	bc.DataLabels = labels
	bc.BarGap = 1
	bc.BarWidth = 6
	bc.TextColor = termui.ColorGreen
	bc.BarColor = termui.ColorRed
	bc.NumColor = termui.ColorYellow

	termui.Render(bc)

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Loop()
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func GetStringSliceFromInt(intKeys []int) []string {
	var stringKeys = []string{}
	for _, i := range intKeys {
		j := strconv.Itoa(i)
		stringKeys = append(stringKeys, j)
	}
	return stringKeys
}

func GetIntSliceFromFloat(floatKeys []float64) []int {
	var intKeys = []int{}
	for _, i := range floatKeys {
		j := int(i)
		intKeys = append(intKeys, j)
	}
	return intKeys
}
