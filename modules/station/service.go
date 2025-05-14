package station

import (
	"encoding/json"
	"errors"
	"github.com/aldngrha/mrt-schedules.git/modules/common/client"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	GetAllStations() (response []StationResponse, err error)
	CheckSchedulesByStation(id string) (response []ScheduleResponse, err error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *service) GetAllStations() (response []StationResponse, err error) {
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"

	bodyResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var stations []Station
	err = json.Unmarshal(bodyResponse, &stations)

	for _, item := range stations {
		response = append(response, StationResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	return
}

func (s *service) CheckSchedulesByStation(id string) (response []ScheduleResponse, err error) {
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"
	bodyResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var schedule []Schedule
	err = json.Unmarshal(bodyResponse, &schedule)
	if err != nil {
		return
	}

	// schedule selected by id station
	var scheduleSelected Schedule
	for _, item := range schedule {
		if item.StationId == id {
			scheduleSelected = item
			break
		}
	}

	if scheduleSelected.StationId == "" {
		err = errors.New("station not found")
		return
	}
	response, err = ConvertDataToResponses(scheduleSelected)
	if err != nil {
		return
	}

	return
}

func ConvertDataToResponses(schedule Schedule) (response []ScheduleResponse, err error) {
	var (
		LebakBulusTripName = "Stasiun Lebak Bulus Grab"
		BundaranHITripName = "Stasiun Bundaran HI Bank DKI"
	)

	scheduleLebakBulus := schedule.ScheduleLebakBulus
	scheduleBundaranHI := schedule.ScheduleBundaranHI

	scheduleLebakBulusParsed, err := ConvertShceduleToTimeFormat(scheduleLebakBulus)
	if err != nil {
		return
	}
	scheduleBundaranHIParsed, err := ConvertShceduleToTimeFormat(scheduleBundaranHI)
	if err != nil {
		return
	}

	//convert to response
	for _, item := range scheduleLebakBulusParsed {
		if item.Format("15:04") > time.Now().Format("1504") {
			response = append(response, ScheduleResponse{
				StationName: LebakBulusTripName,
				Time:        item.Format("15:04"),
			})
		}
	}

	for _, item := range scheduleBundaranHIParsed {
		if item.Format("15:04") > time.Now().Format("1504") {
			response = append(response, ScheduleResponse{
				StationName: BundaranHITripName,
				Time:        item.Format("15:04"),
			})
		}
	}
	return
}

func ConvertShceduleToTimeFormat(schedule string) (response []time.Time, err error) {
	var (
		parsedTime time.Time
		schedules  = strings.Split(schedule, ",")
	)

	for _, item := range schedules {
		trimmedTime := strings.TrimSpace(item)
		if trimmedTime == "" {
			continue
		}

		parsedTime, err = time.Parse("15:04", trimmedTime)
		if err != nil {
			err = errors.New("time format error " + trimmedTime)
			return
		}

		response = append(response, parsedTime)
	}
	return
}
