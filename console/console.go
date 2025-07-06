package console

import (
	"github.com/dfirebird/famigom/bus"
	"github.com/dfirebird/famigom/cpu"
	"github.com/dfirebird/famigom/cpu/ram"
	"github.com/dfirebird/famigom/mapper"
	"github.com/dfirebird/famigom/program"
)

type Console struct {
	cpu *cpu.CPU

	mapper    *mapper.Mapper
	mapperNum byte
}

func CreateConsole(program *program.Program) Console {
	mainBus := bus.CreateMainBus()

	ram, ramAddrRange := ram.CreateRAM()
	mapper, mapperAddrRange := mapper.GetMapper(program)

	mainBus.RegisterDevice(ramAddrRange, ram).RegisterDevice(mapperAddrRange, mapper)

	cpu := cpu.New(&mainBus)

	return Console{
		cpu:       &cpu,
		mapper:    &mapper,
		mapperNum: program.Mapper,
	}
}
