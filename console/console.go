package console

import (
	"github.com/dfirebird/famigom/bus"
	"github.com/dfirebird/famigom/controller"
	"github.com/dfirebird/famigom/cpu"
	"github.com/dfirebird/famigom/cpu/ram"
	"github.com/dfirebird/famigom/mapper"
	"github.com/dfirebird/famigom/palette"
	"github.com/dfirebird/famigom/ppu"
	"github.com/dfirebird/famigom/ppu/nametable"
	"github.com/dfirebird/famigom/program"
)

const (
	CPU_CYCLE_DURATION_NS = 559
)

const (
	CONSOLE_BUTTON_A      = 1 << iota
	CONSOLE_BUTTON_B      = 1 << iota
	CONSOLE_BUTTON_SELECT = 1 << iota
	CONSOLE_BUTTON_START  = 1 << iota
	CONSOLE_BUTTON_UP     = 1 << iota
	CONSOLE_BUTTON_DOWN   = 1 << iota
	CONSOLE_BUTTON_LEFT   = 1 << iota
	CONSOLE_BUTTON_RIGHT  = 1 << iota
)

type Console struct {
	cpu         *cpu.CPU
	ppu         *ppu.PPU
	controllers *controller.Controllers

	mapper    *mapper.Mapper
	mapperNum byte
}

func CreateConsole(romData *[]byte) (*Console, error) {
	program, err := program.Parse(*romData)

	if err != nil {
		return nil, err
	}

	mainBus := bus.CreateMainBus()

	cpuRAM := ram.CreateRAM()
	mapper, err := mapper.GetMapper(program)
	if err != nil {
		return nil, err
	}

	mainBus.RegisterDevice(cpuRAM).RegisterDevice(mapper)

	cpu := cpu.New(&mainBus)

	dmaCallback := cpu.DMA
	dmaDevice := ram.CreateDMADevice(&dmaCallback)
	mainBus.RegisterDevice(dmaDevice)

	nmiCallback := cpu.NMI
	ppu := ppu.CreatePPU(&nmiCallback, nametable.FromNametableArrangement(program.NametableArrangement), mapper)
	mainBus.RegisterDevice(&ppu)
	mapper.SetMirroringUpdateCallback(ppu.UpdateMirroringCallback)

	controllers := controller.CreateControllers()
	mainBus.RegisterDevice(controllers)

	console := Console{
		cpu:         &cpu,
		ppu:         &ppu,
		controllers: controllers,
		mapper:      &mapper,
		mapperNum:   program.Mapper,
	}

	return &console, nil
}

func (c *Console) PowerUp() {
	c.cpu.PowerUp()
	c.ppu.PowerUp()
}

func (c *Console) Step() {
	c.cpu.Step()

	c.ppu.Step()
	c.ppu.Step()
	c.ppu.Step()
}

func (c *Console) LoadControllerButtons(port1, port2 byte) {
	c.controllers.LoadButtonData(port1, port2)
}

func (c *Console) GetPixelData() []byte {
	pixels := make([]byte, 0, 256*240*4)
	for _, colorIdx := range c.ppu.VirtualDisplay {
		color := palette.GetColor(colorIdx)
		pixels = append(pixels, color.R, color.G, color.B, 0xFF) // Last byte is Alpha
	}

	return pixels
}

func (c *Console) DrawNametable() {
	c.ppu.DrawNametable()
}
