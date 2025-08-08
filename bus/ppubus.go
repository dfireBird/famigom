package bus

import "github.com/dfirebird/famigom/types"

type PPUBusDevice interface {
	ReadPRGMemory(addr types.Word) (bool, byte)
	WritePRGMemory(addr types.Word, value byte)
}

type PPUBus struct {
	devicesMap []*PPUBusDevice
}

func NewPPUBus() PPUBus {
	return PPUBus{
		devicesMap: []*PPUBusDevice{},
	}
}

func (b *PPUBus) ReadPRGMemory(addr types.Word) byte {
	for _, device := range b.devicesMap {
		if isRead, value := (*device).ReadPRGMemory(addr); isRead {
			return value
		}
	}

	return 0xFF
}

func (b *PPUBus) WritePRGMemory(addr types.Word, value byte) {
	for _, device := range b.devicesMap {
		(*device).WritePRGMemory(addr, value)
	}
}

func (b *PPUBus) RegisterDevice(deviceStruct PPUBusDevice) *PPUBus {
	b.devicesMap = append(b.devicesMap, &deviceStruct)
	return b
}
