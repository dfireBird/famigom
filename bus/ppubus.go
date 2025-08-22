package bus

import "github.com/dfirebird/famigom/types"

type PPUBusDevice interface {
	ReadCHRMemory(addr types.Word) (bool, byte)
	WriteCHRMemory(addr types.Word, value byte)
}

type PPUBus struct {
	devicesMap []*PPUBusDevice
}

func NewPPUBus() PPUBus {
	return PPUBus{
		devicesMap: []*PPUBusDevice{},
	}
}

func (b *PPUBus) ReadCHRMemory(addr types.Word) byte {
	for _, device := range b.devicesMap {
		if isRead, value := (*device).ReadCHRMemory(addr); isRead {
			return value
		}
	}

	return 0xFF
}

func (b *PPUBus) WriteCHRMemory(addr types.Word, value byte) {
	for _, device := range b.devicesMap {
		(*device).WriteCHRMemory(addr, value)
	}
}

func (b *PPUBus) RegisterDevice(deviceStruct PPUBusDevice) *PPUBus {
	b.devicesMap = append(b.devicesMap, &deviceStruct)
	return b
}
