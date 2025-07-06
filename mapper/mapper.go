package mapper

import (
	"github.com/dfirebird/famigom/bus"
	"github.com/dfirebird/famigom/mapper/mappernrom"
	"github.com/dfirebird/famigom/program"
	. "github.com/dfirebird/famigom/types"
)

type Mapper interface {
    GetMapperNum() byte
	bus.MemoryBus
}

var defaultAddrRange = AddrRange {
	LowAddr: 0xFFFF,
	HighAddr: 0xFFFF,
}

func GetMapper(program *program.Program) (Mapper, AddrRange) {
	switch program.Mapper {
	case 0x00:
		return mappernrom.CreateMapperNRom(program.PrgRom, [8192]byte(program.ChrRom), program.NametableArrangement)
	default:
		return nil, defaultAddrRange
	}
}
