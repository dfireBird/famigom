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

	window, renderer, err := sdl.CreateWindowAndRenderer("Famigom", 1024, 960, sdl.WINDOW_BORDERLESS)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer renderer.Destroy()

	nesScreenTex, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA32, sdl.TEXTUREACCESS_STREAMING, 256, 240)
	if err != nil {
		panic(err)
	}
	err = nesScreenTex.SetScaleMode(sdl.SCALEMODE_NEAREST)
	if err != nil {
		panic(err)
	}

	cycleTimer := sdl.TicksNS()
	elapsed := cycleTimer - cycleTimer
	sdl.RunLoop(func() error {
		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type {
			case sdl.EVENT_QUIT:
				return sdl.EndLoop
			case sdl.EVENT_KEY_DOWN:
				if event.KeyboardEvent().Scancode == sdl.SCANCODE_ESCAPE {
					return sdl.EndLoop
				} else if event.KeyboardEvent().Scancode == sdl.SCANCODE_Z {
					konsole.DrawNametable()
				}
			}
		}

		now := sdl.TicksNS()
		elapsed += now - cycleTimer
		cycleTimer = now

		for elapsed > console.CPU_CYCLE_DURATION_NS {
			konsole.Step()

			elapsed -= console.CPU_CYCLE_DURATION_NS
		}

		pixels := konsole.GetPixelData()

		nesScreenTex.Update(nil, pixels, 256*4)
		renderer.RenderTexture(nesScreenTex, nil, nil)
		renderer.Present()
		return nil
	})
}
