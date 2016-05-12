// ribdRouteServiceApis.go
package rpc

import (
	"fmt"
	"l3/rib/server"
	"ribd"
	"ribdInt"
)

func (m RIBDServicesHandler) CreateIPv4Route(cfg *ribd.IPv4Route) (val bool, err error) {
	logger.Info(fmt.Sprintln("Received create route request for ip", cfg.DestinationNw, " mask ", cfg.NetworkMask, "cfg.NextHopIntRef: ", cfg.NextHop[0].NextHopIntRef))
	err = m.server.RouteConfigValidationCheck(cfg, "add")
	if err != nil {
		logger.Err(fmt.Sprintln("validation check failed with error ", err))
		return false, err
	}
	m.server.RouteCreateConfCh <- cfg
	return true, nil
}
func (m RIBDServicesHandler) OnewayCreateIPv4Route(cfg *ribd.IPv4Route) (err error) {
	logger.Info(fmt.Sprintln("OnewayCreateIPv4Route - Received create route request for ip", cfg.DestinationNw, " mask ", cfg.NetworkMask, "cfg.NextHopIntRef: ", cfg.NextHop[0].NextHopIntRef))
	m.CreateIPv4Route(cfg)
	return err
}
func (m RIBDServicesHandler) OnewayCreateBulkIPv4Route(cfg []*ribdInt.IPv4Route) (err error) {
	//logger.Info(fmt.Sprintln("OnewayCreateIPv4Route - Received create route request for ip", cfg.DestinationNw, " mask ", cfg.NetworkMask, "cfg.OutgoingIntfType: ", cfg.OutgoingIntfType, "cfg.OutgoingInterface: ", cfg.OutgoingInterface))
	logger.Info(fmt.Sprintln("OnewayCreateBulkIPv4Route for ", len(cfg), " routes"))
	for i := 0; i < len(cfg); i++ {
		newCfg := ribd.IPv4Route{
			DestinationNw: cfg[i].DestinationNw,
			NetworkMask:   cfg[i].NetworkMask,
			Cost:          cfg[i].Cost,
			Protocol:      cfg[i].Protocol,
		}
		newCfg.NextHop = make([]*ribd.NextHopInfo, 0)
		nextHop := ribd.NextHopInfo{
			NextHopIp:     cfg[i].NextHopIp,
			NextHopIntRef: cfg[i].NextHopIntRef,
			Weight:        cfg[i].Weight,
		}
		newCfg.NextHop = append(newCfg.NextHop, &nextHop)
		m.CreateIPv4Route(&newCfg)
	}
	return err
}
func (m RIBDServicesHandler) DeleteIPv4Route(cfg *ribd.IPv4Route) (val bool, err error) {
	logger.Info(fmt.Sprintln("DeleteIPv4:RouteReceived Route Delete request for ", cfg.DestinationNw, ":", cfg.NetworkMask, "nextHopIP:", cfg.NextHop[0].NextHopIp, "Protocol ", cfg.Protocol))
	err = m.server.RouteConfigValidationCheck(cfg, "del")
	if err != nil {
		logger.Err(fmt.Sprintln("validation check failed with error ", err))
		return false, err
	}
	m.server.RouteDeleteConfCh <- cfg
	return true, nil
}
func (m RIBDServicesHandler) OnewayDeleteIPv4Route(cfg *ribd.IPv4Route) (err error) {
	logger.Info(fmt.Sprintln("OnewayDeleteIPv4Route:RouteReceived Route Delete request for ", cfg.DestinationNw, ":", cfg.NetworkMask, "nextHopIP:", cfg.NextHop[0].NextHopIp, "Protocol ", cfg.Protocol))
	m.DeleteIPv4Route(cfg)
	return err
}
func (m RIBDServicesHandler) UpdateIPv4Route(origconfig *ribd.IPv4Route, newconfig *ribd.IPv4Route, attrset []bool, op string) (val bool, err error) {
	logger.Println("UpdateIPv4Route: Received update route request")
	err = m.server.RouteConfigValidationCheckForUpdate(newconfig, attrset)
	if err != nil {
		logger.Err(fmt.Sprintln("validation check failed with error ", err))
		return false, err
	}
	routeUpdateConfig := server.UpdateRouteInfo{origconfig, newconfig, attrset, op}
	m.server.RouteUpdateConfCh <- routeUpdateConfig
	return true, nil
}
func (m RIBDServicesHandler) OnewayUpdateIPv4Route(origconfig *ribd.IPv4Route, newconfig *ribd.IPv4Route, attrset []bool) (err error) {
	logger.Println("OneWayUpdateIPv4Route: Received update route request")
	m.UpdateIPv4Route(origconfig, newconfig, attrset, "replace")
	return err
}
func (m RIBDServicesHandler) GetIPv4RouteState(destNw string) (*ribd.IPv4RouteState, error) {
	logger.Info("Get state for IPv4Route")
	route := ribd.NewIPv4RouteState()
	return route, nil
}

func (m RIBDServicesHandler) GetBulkIPv4RouteState(fromIndex ribd.Int, rcount ribd.Int) (routes *ribd.IPv4RouteStateGetInfo, err error) {
	ret, err := m.server.GetBulkIPv4RouteState(fromIndex, rcount)
	return ret, err
}

func (m RIBDServicesHandler) GetIPv4EventState(index int32) (*ribd.IPv4EventState, error) {
	logger.Info("Get state for IPv4EventState")
	route := ribd.NewIPv4EventState()
	return route, nil
}

func (m RIBDServicesHandler) GetBulkIPv4EventState(fromIndex ribd.Int, rcount ribd.Int) (routes *ribd.IPv4EventStateGetInfo, err error) {
	ret, err := m.server.GetBulkIPv4EventState(fromIndex, rcount)
	return ret, err
}

func (m RIBDServicesHandler) GetBulkRoutesForProtocol(srcProtocol string, fromIndex ribdInt.Int, rcount ribdInt.Int) (routes *ribdInt.RoutesGetInfo, err error) {
	ret, err := m.server.GetBulkRoutesForProtocol(srcProtocol, fromIndex, rcount)
	return ret, err
}

func (m RIBDServicesHandler) GetBulkRouteDistanceState(fromIndex ribd.Int, rcount ribd.Int) (routeDistanceStates *ribd.RouteDistanceStateGetInfo, err error) {
	ret, err := m.server.GetBulkRouteDistanceState(fromIndex, rcount)
	return ret, err
}
func (m RIBDServicesHandler) GetRouteDistanceState(protocol string) (*ribd.RouteDistanceState, error) {
	logger.Info("Get state for RouteDistanceState")
	route := ribd.NewRouteDistanceState()
	return route, nil
}
func (m RIBDServicesHandler) GetNextHopIfTypeStr(nextHopIfType ribdInt.Int) (nextHopIfTypeStr string, err error) {
	nhStr, err := m.server.GetNextHopIfTypeStr(nextHopIfType)
	return nhStr, err
}
func (m RIBDServicesHandler) GetRoute(destNetIp string, networkMask string) (route *ribdInt.Routes, err error) {
	ret, err := m.server.GetRoute(destNetIp, networkMask)
	return ret, err
}
func (m RIBDServicesHandler) GetRouteReachabilityInfo(destNet string) (nextHopIntf *ribdInt.NextHopInfo, err error) {
	nh, err := m.server.GetRouteReachabilityInfo(destNet)
	return nh, err
}
func (m RIBDServicesHandler) TrackReachabilityStatus(ipAddr string, protocol string, op string) (err error) {
	m.server.TrackReachabilityCh <- server.TrackReachabilityInfo{ipAddr, protocol, op}
	return nil
}
