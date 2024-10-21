package service

import (
	"time"
)

type CreateEntity struct {
	Mark 			string	`json:"mark"`
	Model 			string	`json:"model"`
	Volume 			string	`json:"volume"`
	Price			string	`json:"price"`
	Year 			string	`json:"year"`
	Kilometers 		string	`json:"kilometers"`
	Power 			string	`json:"power"`
	Transmission 	string	`json:"transmission"`
	Fuel 			string	`json:"fuel"`
	Owners 			string	`json:"owners"`
	Drive 			string	`json:"drive"`
}

type Entity struct {
	ID				int			`json:"id"`
	Mark 			string		`json:"mark"`
	Model 			string		`json:"model"`
	Volume 			string		`json:"volume"`
	Price			string		`json:"price"`
	Year 			string		`json:"year"`
	Kilometers 		string		`json:"kilometers"`
	Power 			string		`json:"power"`
	Transmission 	string		`json:"transmission"`
	Fuel 			string		`json:"fuel"`
	Owners 			string		`json:"owners"`
	Drive 			string		`json:"drive"`
	CreatedAt		time.Time	`json:"creation time"`
	UpdatedAt		time.Time	`json:"update time"`
}