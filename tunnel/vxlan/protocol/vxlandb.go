// vxlandb.go
package vxlan

import (
	"net"
)

var vxlanDB map[uint32]*vxlanDbEntry

type vxlanDbEntry struct {
	VNI         uint32
	VlanId      uint16 // used to tag inner ethernet frame when egressing
	Group       net.IP // multicast group IP
	MTU         uint32 // MTU size for each VTEP
	VtepMembers []uint32
}

func NewVxlanDbEntry(c *VxlanConfig) *vxlanDbEntry {
	return &vxlanDbEntry{
		VNI:         c.VNI,
		VlanId:      c.VlanId,
		Group:       c.Group,
		MTU:         c.MTU,
		VtepMembers: make([]uint32, 0),
	}
}

func (s *VXLANServer) saveVxLanConfigData(c *VxlanConfig) {
	if _, ok := vxlanDB[c.VNI]; !ok {
		vxlan := NewVxlanDbEntry(c)
		vxlanDB[c.VNI] = vxlan
	}
}