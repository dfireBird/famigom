package mapper

import (
	"fmt"

	"github.com/dfirebird/famigom/bus"
	"github.com/dfirebird/famigom/mapper/mapperaxrom"
	"github.com/dfirebird/famigom/mapper/mappercnrom"
	"github.com/dfirebird/famigom/mapper/mappermmc1"
	"github.com/dfirebird/famigom/mapper/mappernrom"
	"github.com/dfirebird/famigom/mapper/mapperuxrom"
	"github.com/dfirebird/famigom/ppu/nametable"
	"github.com/dfirebird/famigom/program"
	"github.com/dfirebird/famigom/types"
)

type Mapper interface {
	CPUStep()
	PPUStep()
	GetMapperNum() byte
	SetMirroringUpdateCallback(func(nametable.NametableMirroring))
	bus.MainBusDevice
	bus.PPUBusDevice
}

var (
	ErrUnsupported = func(mNo types.Word) error { return fmt.Errorf("mapper is not supported yet %d", mNo) }
)

func GetMapper(program *program.Program) (Mapper, error) {
	switch program.Mapper {
	case 0x00:
		return mappernrom.CreateMapperNRom(program), nil
	case 0x01:
		return mappermmc1.CreateMapperMMC1(program), nil
	case 0x02:
		return mapperuxrom.CreateMapperUxROM(program), nil
	case 0x03:
		return mappercnrom.CreateMapperCNROM(program), nil
	case 0x07:
		return mapperaxrom.CreateMapperAxROM(program), nil
	default:
		return nil, ErrUnsupported(program.Mapper)
	}
}
