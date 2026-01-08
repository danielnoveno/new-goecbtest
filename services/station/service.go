/*
   file:           services/station/service.go
   description:    Service untuk logic EcbStation (Identity & Configuration)
   created:        optimization 05-01-2026
*/

package station

import (
	"fmt"
	"net"

	"go-ecb/app/types"
	"go-ecb/configs"
	"go-ecb/pkg/logging"
	"go-ecb/repository"
)

type StationService struct {
	repo *repository.EcbStationRepository
}

func NewStationService(repo *repository.EcbStationRepository) *StationService {
	return &StationService{repo: repo}
}

func (s *StationService) Initialize() (*types.EcbStation, error) {
	myIP := GetOutboundIP()
	logging.Logger().Infof("[Station] Detecting machine identity... IP: %s", myIP)

	station, err := s.repo.FindEcbStationByIP(myIP)
	if err != nil {
		return nil, fmt.Errorf("failed to check station identity: %w", err)
	}

	if station != nil {
		logging.Logger().Infof("[Station] Indentity found in DB. Location: %s, Line: %s", station.Location, station.Linetype)
		return station, nil
	}

	logging.Logger().Infof("[Station] New machine detected! Auto-registering from ENV...")
	
	simoConfig := configs.LoadSimoConfig()
	newStation := types.EcbStation{
		Ipaddress:   myIP,
		Location:    simoConfig.EcbLocation,
		Linetype:    simoConfig.EcbLineType,
		Lineids:     simoConfig.EcbLineIds,
		Workcenters: simoConfig.EcbWorkcenters,
		Mode:        simoConfig.EcbMode,
		Tacktime:    int(simoConfig.EcbTacktime),
		Status:      simoConfig.EcbStateDefault,
		Theme:       simoConfig.Theme,
		Lineactive:  1,
	}

	id, err := s.repo.CreateEcbStation(newStation)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-register station: %w", err)
	}
	
	newStation.ID = id
	logging.Logger().Infof("[Station] Registration success. ID: %d", id)
	return &newStation, nil
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
