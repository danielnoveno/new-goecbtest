/*
    file:           app/types/application_settings.go
    description:    Model dan helper UI untuk ecbsetting
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

type ECBSetting struct {
	ServerIPAddress string `json:"server_ip_address"`
	Simo3IPAddress  string `json:"simo3_ip_address"`
	UseWLAN         string `json:"use_wlan"`
}
