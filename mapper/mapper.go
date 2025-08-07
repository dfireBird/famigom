package mapper

import (
	"fmt"

	"github.com/dfirebird/famigom/bus"
	"github.com/dfirebird/famigom/mapper/mappernrom"
	"github.com/dfirebird/famigom/program"
)

type Mapper interface {
	GetMapperNum() byte
	bus.MainBusDevice
	bus.PPUBusDevice
}

var (
	ErrUnsupported = func(mNo byte) error { return fmt.Errorf("mapper is not supported yet %d", mNo) }
)

func GetMapper(program *program.Program) (Mapper, error) {
	switch program.Mapper {
	case 0x00:
		return mappernrom.CreateMapperNRom(program.PrgRom, [8192]byte(program.ChrRom), program.NametableArrangement), nil
	default:
		return nil, ErrUnsupported(program.Mapper)
	}
}
