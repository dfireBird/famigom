package mapper

import (
	"github.com/dfirebird/famigom/bus"
	"github.com/dfirebird/famigom/mapper/mappernrom"
	"github.com/dfirebird/famigom/program"
)

type Mapper interface {
	GetMapperNum() byte
	bus.MainBusDevice
}

func GetMapper(program *program.Program) Mapper {
	switch program.Mapper {
	case 0x00:
		return mappernrom.CreateMapperNRom(program.PrgRom, [8192]byte(program.ChrRom), program.NametableArrangement)
	default:
		return nil
	}
}
