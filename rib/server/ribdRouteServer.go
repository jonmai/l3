//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

// ribdRouteServer.go
package server

import (
	"ribd"
)

type RouteConfigInfo struct {
	OrigRoute *ribd.IPv4Route
	NewRoute  *ribd.IPv4Route
	Attrset   []bool
	Op        string //"add"/"del"/"update"
}
type TrackReachabilityInfo struct {
	IpAddr   string
	Protocol string
	Op       string
}
type NextHopInfoKey struct {
	nextHopIp string
}
type NextHopInfo struct {
	refCount int //number of routes using this as a next hop
}
type PerProtocolRouteInfo struct {
	routeMap   map[string]int
	totalcount int
}

var ProtocolRouteMap map[string]PerProtocolRouteInfo //map[string]int

func UpdateProtocolRouteMap(protocol string, op string, value string) {
	if ProtocolRouteMap == nil {
		if op == "del" {
			return
		}
		ProtocolRouteMap = make(map[string]PerProtocolRouteInfo) //map[string]int)
	}
	val, ok := ProtocolRouteMap[protocol]
	if !ok {
		if op == "del" {
			return
		}
		ProtocolRouteMap[protocol].routeMap = make(map[string]int)
	}
	protocolroutemap, ok := ProtocolRouteMap[protocol].routeMap
	if !ok {
		return
	}
	totalcount := val.totalcount
	count, ok := protocolroutemap[value]
	if !ok {
		if op == "del" {
			return
		}
	}
	if op == "add" {
		count++
		totalcount++
	} else if op == "del" {
		count--
		totalcount--
	}
	protocolroutemap[value] = count
	ProtocolRouteMap[protocol].routeMap = protocolroutemap
	ProtocolRouteMap[protocol].totalcount = totalcount

}
func (ribdServiceHandler *RIBDServer) StartRouteProcessServer() {
	logger.Info("Starting the routeserver loop")
	ProtocolRouteMap = make(map[string]map[string]int)
	for {
		select {
		case routeConf := <-ribdServiceHandler.RouteConfCh:
			//logger.Debug(fmt.Sprintln("received message on RouteConfCh channel, op: ", routeConf.Op)
			if routeConf.Op == "add" {
				ribdServiceHandler.ProcessV4RouteCreateConfig(routeConf.OrigConfigObject.(*ribd.IPv4Route))
			} else if routeConf.Op == "addBulk" {
				ribdServiceHandler.ProcessBulkRouteCreateConfig(routeConf.OrigBulkRouteConfigObject) //.([]*ribd.IPv4Route))
			} else if routeConf.Op == "del" {
				ribdServiceHandler.ProcessV4RouteDeleteConfig(routeConf.OrigConfigObject.(*ribd.IPv4Route))
			} else if routeConf.Op == "update" {
				if routeConf.PatchOp == nil || len(routeConf.PatchOp) == 0 {
					ribdServiceHandler.Processv4RouteUpdateConfig(routeConf.OrigConfigObject.(*ribd.IPv4Route), routeConf.NewConfigObject.(*ribd.IPv4Route), routeConf.AttrSet)
				} else {
					ribdServiceHandler.Processv4RoutePatchUpdateConfig(routeConf.OrigConfigObject.(*ribd.IPv4Route), routeConf.NewConfigObject.(*ribd.IPv4Route), routeConf.PatchOp)
				}
			} else if routeConf.Op == "addv6" {
				//create ipv6 route
				ribdServiceHandler.ProcessV6RouteCreateConfig(routeConf.OrigConfigObject.(*ribd.IPv6Route))
			} else if routeConf.Op == "delv6" {
				//delete ipv6 route
				ribdServiceHandler.ProcessV6RouteDeleteConfig(routeConf.OrigConfigObject.(*ribd.IPv6Route))
			} else if routeConf.Op == "updatev6" {
				//update ipv6 route
				if routeConf.PatchOp == nil || len(routeConf.PatchOp) == 0 {
					ribdServiceHandler.Processv6RouteUpdateConfig(routeConf.OrigConfigObject.(*ribd.IPv6Route), routeConf.NewConfigObject.(*ribd.IPv6Route), routeConf.AttrSet)
				} else {
					ribdServiceHandler.Processv6RoutePatchUpdateConfig(routeConf.OrigConfigObject.(*ribd.IPv6Route), routeConf.NewConfigObject.(*ribd.IPv6Route), routeConf.PatchOp)
				}
			}
		}
	}
}
