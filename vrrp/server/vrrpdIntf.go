package vrrpServer

import (
	"asicd/asicdConstDefs"
	"encoding/json"
	"fmt"
	nanomsg "github.com/op/go-nanomsg"
)

func VrrpAsicdSubscriber() {
	for {
		logger.Info("VRRP: Read on Asic Subscriber socket....")
		rxBuf, err := asicdSubSocket.Recv(0)
		if err != nil {
			logger.Err(fmt.Sprintln("VRRP: Recv on asicd Subscriber socket failed with error:", err))
			continue
		}
		//logger.Info(fmt.Sprintln("VRRP: asicd Subscriber recv returned:", rxBuf))
		var msg asicdConstDefs.AsicdNotification
		err = json.Unmarshal(rxBuf, &msg)
		if err != nil {
			logger.Err(fmt.Sprintln("VRRP: Unable to Unmarshal asicd msg:", msg.Msg))
			continue
		}
		if msg.MsgType == asicdConstDefs.NOTIFY_VLAN_CREATE ||
			msg.MsgType == asicdConstDefs.NOTIFY_VLAN_DELETE {
			//Vlan Create Msg
			var vlanNotifyMsg asicdConstDefs.VlanNotifyMsg
			err = json.Unmarshal(msg.Msg, &vlanNotifyMsg)
			if err != nil {
				logger.Err(fmt.Sprintln("VRRP: Unable to unmashal vlanNotifyMsg:", msg.Msg))
				return
			}
			//DhcpRelayAgentUpdateVlanInfo(vlanNotifyMsg, msg.MsgType)
		} else if msg.MsgType == asicdConstDefs.NOTIFY_IPV4INTF_CREATE ||
			msg.MsgType == asicdConstDefs.NOTIFY_IPV4INTF_DELETE {
			var ipv4IntfNotifyMsg asicdConstDefs.IPv4IntfNotifyMsg
			err = json.Unmarshal(msg.Msg, &ipv4IntfNotifyMsg)
			if err != nil {
				logger.Err(fmt.Sprintln("VRRP: Unable to Unmarshal ipv4IntfNotifyMsg:", msg.Msg))
				continue
			}
			//DhcpRelayAgentUpdateIntfPortMap(ipv4IntfNotifyMsg, msg.MsgType)
		} else if msg.MsgType == asicdConstDefs.NOTIFY_L3INTF_STATE_CHANGE {
			//INTF_STATE_CHANGE
			var l3IntfStateNotifyMsg asicdConstDefs.L3IntfStateNotifyMsg
			err = json.Unmarshal(msg.Msg, &l3IntfStateNotifyMsg)
			if err != nil {
				logger.Err(fmt.Sprintln("VRRP: unable to Unmarshal l3 intf state change:", msg.Msg))
				continue
			}
			//DhcpRelayAgentUpdateL3IntfStateChange(l3IntfStateNotifyMsg)
		}
	}
}

func VrrpRegisterWithAsicdUpdates(address string) error {
	var err error
	logger.Info("VRRP: setting up asicd update listener")
	if asicdSubSocket, err = nanomsg.NewSubSocket(); err != nil {
		logger.Err(fmt.Sprintln("VRRP: Failed to create ASIC subscribe socket, error:", err))
		return err
	}

	if err = asicdSubSocket.Subscribe(""); err != nil {
		logger.Err(fmt.Sprintln("VRRP:Failed to subscribe to \"\" on ASIC subscribe socket, error:",
			err))
		return err
	}

	if _, err = asicdSubSocket.Connect(address); err != nil {
		logger.Err(fmt.Sprintln("VRRP: Failed to connect to ASIC publisher socket, address:",
			address, "error:", err))
		return err
	}

	if err = asicdSubSocket.SetRecvBuffer(1024 * 1024); err != nil {
		logger.Err(fmt.Sprintln("VRRP: Failed to set the buffer size for ",
			"ASIC publisher socket, error:", err))
		return err
	}
	logger.Info("VRRP: asicd update listener is set")
	return nil
}

func VrrpGetPortInfoFromAsicd() error {
	logger.Info("VRRP: Calling Asicd to initialize port properties")
	err := VrrpRegisterWithAsicdUpdates(asicdConstDefs.PUB_SOCKET_ADDR)
	if err == nil {
		// Asicd subscriber thread
		go VrrpAsicdSubscriber()
	}
	return nil
}
