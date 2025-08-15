package console

import (
	"github.com/dfirebird/famigom/bus"
	"github.com/dfirebird/famigom/cpu"
	"github.com/dfirebird/famigom/cpu/ram"
	"github.com/dfirebird/famigom/mapper"
	"github.com/dfirebird/famigom/palette"
	"github.com/dfirebird/famigom/ppu"
	"github.com/dfirebird/famigom/program"
)

const (
	CPU_CYCLE_DURATION_NS = 559
)

type Console struct {
	cpu *cpu.CPU
	ppu *ppu.PPU

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
	ppu := ppu.CreatePPU(&nmiCallback, program.NametableArrangement.GetMirroring())
	mainBus.RegisterDevice(&ppu)
	ppu.RegisterDevice(mapper)

	console := Console{
		cpu:       &cpu,
		ppu:       &ppu,
		mapper:    &mapper,
		mapperNum: program.Mapper,
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

func (c *Console) GetPixelData() []byte {
	pixels := make([]byte, 0, 256 * 240 * 4)
	for _, colorIdx := range c.ppu.VirtualDisplay {
		color := palette.GetColor(colorIdx)
		pixels = append(pixels, color.R, color.G, color.B, 0xFF) // Last byte is Alpha
	}

	return pixels
}

func (c *Console) DrawNametable() {
	c.ppu.DrawNametable()
}
