//Только структуры

package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	name string `json:"name"` 	
	password string
	status string `json:"status"`
	role string `json:"role"`
	organization string `json:"organization"`
}