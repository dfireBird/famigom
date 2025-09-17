package main

import (
	"flag"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/dfirebird/famigom/console"
	"github.com/dfirebird/famigom/log"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/sdl"
)

const PerFrameScreenTicks = 1000 / 60

func main() {
	verbose := flag.Bool("v", false, "Enables verbose logging")
	profile := flag.Bool("p", false, "Enables profile logging")
	flag.Parse()

	if err := log.InitLogger(*verbose); err != nil {
		panic("Could not initialize loggers")
	}
	defer log.FlushLoggers()

	logger := log.Logger()
	if *profile {
		f, err := os.Create("cpu.prof")
		if err != nil {
			logger.Fatal("Could not create CPU profile file: ", err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			logger.Fatal("Could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if len(flag.Args()) != 1 {
		logger.Errorf("NES ROM file path is not passed as an argument")
		os.Exit(1)
	}

	romPath, _ := filepath.Abs(flag.Arg(0))
	logger.Infof("Reading ROM/Program file from path: %s", romPath)
	romData, err := os.ReadFile(romPath)
	if err != nil {
		logger.Errorf("Reading the ROM file %v", err)
		os.Exit(1)
	}

	konsole, err := console.CreateConsole(&romData)
	if err != nil {
		logger.Errorf("Creating Emulator failed: %v", err)
		os.Exit(1)
	}

	konsole.PowerUp()

	defer binsdl.Load().Unload()
	defer sdl.Quit()

	if err := sdl.SetHint(sdl.HINT_RENDER_VSYNC, "0"); err != nil {
		panic(err)
	}

	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_GAMEPAD); err != nil {
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
	logger.Info("SDL Window Initialization complete")

	step := false
	pause := false

	cycleTimer := sdl.TicksNS()
	elapsed := cycleTimer - cycleTimer

	var player1, player2 byte

	sdl.RunLoop(func() error {
		frameStart := sdl.Ticks()
		var event sdl.Event

		for sdl.PollEvent(&event) {
			switch event.Type {
			case sdl.EVENT_QUIT:
				return sdl.EndLoop
			case sdl.EVENT_GAMEPAD_ADDED:
				_, err := event.GamepadDeviceEvent().Which.OpenGamepad()
				if err != nil {
					logger.Errorln(err)
				}
			case sdl.EVENT_GAMEPAD_REMOVED:
				gamepad, err := event.GamepadDeviceEvent().Which.Gamepad()
				if err == nil {
					gamepad.Close()
				} else {
					logger.Errorln(err)
				}
			case sdl.EVENT_GAMEPAD_BUTTON_DOWN:
				button := sdl.GamepadButton(event.GamepadButtonEvent().Button)
				if button == sdl.GAMEPAD_BUTTON_RIGHT_SHOULDER {
					pause = !pause
					logger.Info("Pause: ", pause)
				} else if button == sdl.GAMEPAD_BUTTON_LEFT_SHOULDER {
					step = true
				} else if button == sdl.GAMEPAD_BUTTON_RIGHT_STICK {
					log.IsTrace = !log.IsTrace
					logger.Info("Trace: ", log.IsTrace)
				}
			case sdl.EVENT_KEY_DOWN:
				if event.KeyboardEvent().Scancode == sdl.SCANCODE_ESCAPE {
					return sdl.EndLoop
				} else if event.KeyboardEvent().Scancode == sdl.SCANCODE_1 {
					konsole.DrawNametable()
				} else if event.KeyboardEvent().Scancode == sdl.SCANCODE_2 {
					pause = !pause
					logger.Info("Pause: ", pause)
				} else if event.KeyboardEvent().Scancode == sdl.SCANCODE_3 {
					step = true
				} else if event.KeyboardEvent().Scancode == sdl.SCANCODE_4 {
					log.IsTrace = !log.IsTrace
					logger.Info("Trace Enable: ", log.IsTrace)
				}
			}
			ConvertInputEventsForConsole(event, &player1, &player2)
		}

		if pause && step {
			for range 29780 {
				konsole.LoadControllerButtons(0, 0)
				konsole.Step()
			}
			step = false
		}

		if !pause {
			now := sdl.TicksNS()
			elapsed += now - cycleTimer
			cycleTimer = now
		} else {
			elapsed = 0
			cycleTimer = sdl.TicksNS()
		}

		if !pause {
			for elapsed > console.CPU_CYCLE_DURATION_NS {
				konsole.LoadControllerButtons(player1, player2)
				konsole.Step()

				elapsed -= console.CPU_CYCLE_DURATION_NS
			}
		}

		renderer.Clear()
		pixels := konsole.GetPixelData()

		nesScreenTex.Update(nil, pixels, 256*4)
		renderer.RenderTexture(nesScreenTex, nil, nil)
		renderer.Present()

		if !pause {
			deltaFrameTime := sdl.Ticks() - frameStart
			for deltaFrameTime < PerFrameScreenTicks {
				deltaFrameTime = sdl.Ticks() - frameStart
			}
		}
		return nil
	})
}

func ConvertInputEventsForConsole(event sdl.Event, port1, port2 *byte) {
	if event.Type == sdl.EVENT_GAMEPAD_BUTTON_DOWN || event.Type == sdl.EVENT_GAMEPAD_BUTTON_UP {
		event := event.GamepadButtonEvent()
		playerIdx := event.Which.GamepadPlayerIndex()

		var currentPort *byte
		if playerIdx == 1 {
			currentPort = port2
		} else {
			currentPort = port1
		}

		button := sdl.GamepadButton(event.Button)
		switch button {
		case sdl.GAMEPAD_BUTTON_SOUTH, sdl.GAMEPAD_BUTTON_NORTH:
			handleChangeWithEvent(currentPort, console.CONSOLE_BUTTON_A, event.Type)
		case sdl.GAMEPAD_BUTTON_EAST, sdl.GAMEPAD_BUTTON_WEST:
			handleChangeWithEvent(currentPort, console.CONSOLE_BUTTON_B, event.Type)
		case sdl.GAMEPAD_BUTTON_BACK:
			handleChangeWithEvent(currentPort, console.CONSOLE_BUTTON_SELECT, event.Type)
		case sdl.GAMEPAD_BUTTON_START:
			handleChangeWithEvent(currentPort, console.CONSOLE_BUTTON_START, event.Type)
		case sdl.GAMEPAD_BUTTON_DPAD_UP:
			handleChangeWithEvent(currentPort, console.CONSOLE_BUTTON_UP, event.Type)
		case sdl.GAMEPAD_BUTTON_DPAD_DOWN:
			handleChangeWithEvent(currentPort, console.CONSOLE_BUTTON_DOWN, event.Type)
		case sdl.GAMEPAD_BUTTON_DPAD_LEFT:
			handleChangeWithEvent(currentPort, console.CONSOLE_BUTTON_LEFT, event.Type)
		case sdl.GAMEPAD_BUTTON_DPAD_RIGHT:
			handleChangeWithEvent(currentPort, console.CONSOLE_BUTTON_RIGHT, event.Type)
		}
	}
}

func handleChangeWithEvent(port *byte, buttonValue byte, eventType sdl.EventType) {
	if eventType == sdl.EVENT_GAMEPAD_BUTTON_DOWN {
		*port |= buttonValue
	} else {
		*port &^= buttonValue
	}
}
