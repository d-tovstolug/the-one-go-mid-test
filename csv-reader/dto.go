package main

import (
	"errors"
	"strconv"
)

type Entry struct {
	Year                      int32
	IndustryAggregationNZSIOC int32
	IndustryCodeNZSIOC        int32
	IndustryNameNZSIOC        string
	Units                     string
	VariableCode              string
	VariableName              string
	VariableCategory          string
	Value                     int64
	IndustryCodeANZSIC06      string
}

func parseCSVLine(data []string) (*Entry, error) {
	if len(data) < 10 {
		return nil, errors.New("not enough data to parse from csv")
	}

	year, _ := strconv.ParseInt(data[0], 10, 32)
	indAgg, _ := strconv.ParseInt(data[1], 10, 32)
	indCode, _ := strconv.ParseInt(data[2], 10, 32)
	value, _ := strconv.ParseInt(data[8], 10, 32)

	return &Entry{
		Year:                      int32(year),
		IndustryAggregationNZSIOC: int32(indAgg),
		IndustryCodeNZSIOC:        int32(indCode),
		IndustryNameNZSIOC:        data[3],
		Units:                     data[4],
		VariableCode:              data[5],
		VariableName:              data[6],
		VariableCategory:          data[7],
		Value:                     value,
		IndustryCodeANZSIC06:      data[9],
	}, nil
}
