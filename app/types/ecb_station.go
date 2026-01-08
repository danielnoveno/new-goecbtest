/*
   file:           app/types/station_definition.go
   description:    Model dan helper UI untuk ecbstation
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import (
	"time"
)

type EcbStation struct {
	ID          int       `db:"id"`
	Ipaddress   string    `db:"ipaddress"`
	Location    string    `db:"location"`
	Mode        string    `db:"mode"`
	Linetype    string    `db:"linetype"`
	Lineids     string    `db:"lineids"`
	Lineactive  int       `db:"lineactive"`
	Ecbstate    string    `db:"ecbstate"`
	Theme       string    `db:"theme"`
	Tacktime    int       `db:"tacktime"`
	Workcenters string    `db:"workcenters"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
