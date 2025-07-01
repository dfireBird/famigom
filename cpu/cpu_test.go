package cpu_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/dfirebird/famigom/cpu"
	. "github.com/dfirebird/famigom/types"
)

type ramValue struct {
	addr  Word
	value byte
}

type cycleState struct {
	addr  Word
	value byte
	type_ string
}

type testCPUState struct {
	pc  Word
	s   byte
	a   byte
	x   byte
	y   byte
	p   byte
	ram []ramValue
}

type testScenario struct {
	name    string
	initial testCPUState
	final   testCPUState
	cycles  []cycleState
}

func run_insturction_test(t *testing.T, jsonFilePath string) {
	testScenarios := parse_test_data(jsonFilePath)

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			testCPU := createCPU(scenario.initial)

			testCPU.Step()

			for _, ramV := range scenario.final.ram {
				if cpuValue := testCPU.ReadMemory(ramV.addr); cpuValue != ramV.value {
					t.Fatalf("Memory value mismatch. Expect %d at %d but got %d", ramV.value, ramV.addr, cpuValue)
				}
			}

			if testCPU.A != scenario.final.a {
				t.Fatalf("Register value mismatch. Expect %d at %s but got %d", scenario.final.a, "A", testCPU.A)
			}
			if testCPU.X != scenario.final.x {
				t.Fatalf("Register value mismatch. Expect %d at %s but got %d", scenario.final.x, "X", testCPU.X)
			}
			if testCPU.Y != scenario.final.y {
				t.Fatalf("Register value mismatch. Expect %d at %s but got %d", scenario.final.y, "Y", testCPU.Y)
			}
			if testCPU.SP != scenario.final.s {
				t.Fatalf("Register value mismatch. Expect %d at %s but got %d", scenario.final.s, "SP", testCPU.SP)
			}
			if byte(testCPU.Flags) != scenario.final.p {
				t.Fatalf("Register value mismatch. Expect %d at %s but got %d", scenario.final.p, "SR", testCPU.Flags)
			}
			if testCPU.PC != scenario.final.pc {
				t.Fatalf("Register value mismatch. Expect %d at %s but got %d", scenario.final.pc, "PC", testCPU.PC)
			}

		})
	}
}

func parse_test_data(jsonFilePath string) []testScenario {
	parseCPUState := func(data map[string]any) testCPUState {
		pc := Word(data["pc"].(float64))
		s := byte(data["s"].(float64))
		a := byte(data["a"].(float64))
		x := byte(data["x"].(float64))
		y := byte(data["y"].(float64))
		p := byte(data["p"].(float64))

		var ram []ramValue
		for _, data := range data["ram"].([]any) {
			ramData := data.([]any)
			addr := Word(ramData[0].(float64))
			value := byte(ramData[1].(float64))

			ram = append(ram, ramValue{
				addr:  addr,
				value: value,
			})
		}

		return testCPUState{
			pc:  pc,
			s:   s,
			a:   a,
			x:   x,
			y:   y,
			p:   p,
			ram: ram,
		}
	}
	parseCycle := func(data []any) cycleState {
		return cycleState{
			addr:  Word(data[0].(float64)),
			value: byte(data[1].(float64)),
			type_: data[2].(string),
		}
	}
	var testScenarios []testScenario

	var parsedScenarios []any

	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return testScenarios
	}

	err = json.Unmarshal(jsonData, &parsedScenarios)
	for _, data := range parsedScenarios {
		scenarioData := data.(map[string]any)
		name := scenarioData["name"].(string)
		initial := parseCPUState(scenarioData["initial"].(map[string]any))
		final := parseCPUState(scenarioData["final"].(map[string]any))

		var cycles []cycleState
		for _, cycle := range scenarioData["cycles"].([]any) {
			cycles = append(cycles, parseCycle(cycle.([]any)))
		}

		testScenarios = append(testScenarios, testScenario{
			name:    name,
			initial: initial,
			final:   final,
			cycles:  cycles,
		})
	}

	return testScenarios
}

const maxMemory = (1 << 16)
func createCPU(cpuState testCPUState) cpu.CPU {
	testCPU := cpu.CPU{
		X:      cpuState.x,
		Y:      cpuState.y,
		A:      cpuState.a,
		Flags:  cpu.Status(cpuState.p),
		SP:     cpuState.s,
		PC:     cpuState.pc,
		Memory: [maxMemory]byte{},
	}

	for _, ramV := range cpuState.ram {
		testCPU.WriteMemory(ramV.addr, ramV.value)
	}

	return testCPU
}
