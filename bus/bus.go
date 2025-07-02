package bus

import (. "github.com/dfirebird/famigom/types")

type MemoryBus interface {
	ReadMemory(addr Word) byte
	WriteMemory(addr Word, value byte)
}
