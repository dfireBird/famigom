package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/dfirebird/famigom/console"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/sdl"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

func main() {
	verbose := flag.Bool("v", false, "Enables verbose logging")
	flag.Parse()

	if len(flag.Args()) != 1 {
		logger.Printf("ERROR: NES ROM file path is not passed as an argument")
		os.Exit(1)
	}

	romPath, _ := filepath.Abs(flag.Arg(0))
	logger.Printf("Reading ROM/Program file from path: %s", romPath)
	romData, err := os.ReadFile(romPath)
	if err != nil {
		logger.Printf("ERROR: Reading the ROM file %v", err)
		os.Exit(1)
	}

	konsole, err := console.CreateConsole(&romData, *verbose)
	if err != nil {
		logger.Printf("ERROR: Creating Emulator failed: %v", err)
		os.Exit(1)
	}

	logger.Printf("Powering up the console")
	konsole.PowerUp()

	defer binsdl.Load().Unload()
	defer sdl.Quit()

	if err := sdl.SetHint(sdl.HINT_RENDER_VSYNC, "1"); err != nil {
		panic(err)
	}

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}

	window, renderer, err := sdl.CreateWindowAndRenderer("Famigom", 256, 240, sdl.WINDOW_BORDERLESS)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer renderer.Destroy()

	origin := sdl.FPoint{10.0, 100.0}
	scale := float32(1.5)
	var scaledPoints []sdl.FPoint
	for _, point := range FAMIGOM_POINTS {
		scaledPoints = append(scaledPoints, sdl.FPoint{origin.X + point.X*scale, origin.Y + point.Y*scale})
	}

	renderer.SetDrawColor(255, 255, 255, sdl.ALPHA_OPAQUE)

	cycleTimer := sdl.TicksNS()
	elapsed := cycleTimer - cycleTimer
	sdl.RunLoop(func() error {
		now := sdl.TicksNS()
		elapsed += now - cycleTimer
		cycleTimer = now

		for elapsed > console.CPU_CYCLE_DURATION_NS {
			var event sdl.Event
			for sdl.PollEvent(&event) {
				switch event.Type {
				case sdl.EVENT_QUIT:
					return sdl.EndLoop
				case sdl.EVENT_KEY_DOWN:
					if event.KeyboardEvent().Scancode == sdl.SCANCODE_ESCAPE {
						return sdl.EndLoop
					}
				}
			}

			konsole.Step()

			renderer.RenderPoints(scaledPoints)
			renderer.Present()
			elapsed -= console.CPU_CYCLE_DURATION_NS
		}

		return nil
	})
}

var FAMIGOM_POINTS []sdl.FPoint = []sdl.FPoint{
	// F (starts at x=0)
	{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}, {6, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {1, 3}, {2, 3}, {3, 3}, {4, 3},
	// A (starts at x=25)
	{25 + 3, 0}, {25 + 2, 1}, {25 + 4, 1}, {25 + 1, 2}, {25 + 5, 2}, {25 + 0, 3}, {25 + 6, 3}, {25 + 0, 4}, {25 + 6, 4}, {25 + 1, 5}, {25 + 2, 5}, {25 + 3, 5}, {25 + 4, 5}, {25 + 5, 5}, {25 + 0, 6}, {25 + 6, 6},
	// M (starts at x=50)
	{50 + 0, 0}, {50 + 6, 0}, {50 + 0, 1}, {50 + 6, 1}, {50 + 0, 2}, {50 + 1, 2}, {50 + 5, 2}, {50 + 6, 2}, {50 + 0, 3}, {50 + 2, 3}, {50 + 4, 3}, {50 + 6, 3}, {50 + 0, 4}, {50 + 3, 4}, {50 + 6, 4}, {50 + 0, 5}, {50 + 6, 5},
	// I (starts at x=75)
	{75 + 0, 0}, {75 + 1, 0}, {75 + 2, 0}, {75 + 3, 0}, {75 + 4, 0}, {75 + 5, 0}, {75 + 6, 0},
	{75 + 3, 1}, {75 + 3, 2}, {75 + 3, 3}, {75 + 3, 4}, {75 + 3, 5},
	{75 + 0, 6}, {75 + 1, 6}, {75 + 2, 6}, {75 + 3, 6}, {75 + 4, 6}, {75 + 5, 6}, {75 + 6, 6},
	// G (starts at x=100)
	{100 + 1, 0}, {100 + 2, 0}, {100 + 3, 0}, {100 + 4, 0}, {100 + 5, 0}, {100 + 6, 1}, {100 + 6, 2}, {100 + 6, 3}, {100 + 4, 3}, {100 + 5, 3}, {100 + 3, 4}, {100 + 6, 5},
	{100 + 1, 6}, {100 + 2, 6}, {100 + 3, 6}, {100 + 4, 6}, {100 + 5, 6}, {100 + 0, 1}, {100 + 0, 2}, {100 + 0, 3}, {100 + 0, 4}, {100 + 0, 5},
	// O (starts at x=125)
	{125 + 1, 0}, {125 + 2, 0}, {125 + 3, 0}, {125 + 4, 0}, {125 + 5, 0}, {125 + 0, 1}, {125 + 6, 1}, {125 + 0, 2}, {125 + 6, 2},
	{125 + 0, 3}, {125 + 6, 3}, {125 + 0, 4}, {125 + 6, 4}, {125 + 0, 5}, {125 + 6, 5},
	{125 + 1, 6}, {125 + 2, 6}, {125 + 3, 6}, {125 + 4, 6}, {125 + 5, 6},
	// M (starts at x=150)
	{150 + 0, 0}, {150 + 6, 0}, {150 + 0, 1}, {150 + 6, 1}, {150 + 0, 2}, {150 + 1, 2}, {150 + 5, 2}, {150 + 6, 2}, {150 + 0, 3}, {150 + 2, 3}, {150 + 4, 3}, {150 + 6, 3}, {150 + 0, 4}, {150 + 3, 4}, {150 + 6, 4}, {150 + 0, 5}, {150 + 6, 5},
}
