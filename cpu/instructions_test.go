package cpu_test

import "testing"

func Test_BRK(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/00.json"
	runInstructionTest(t, test_data_path)
}

func Test_ORA_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/01.json"
	runInstructionTest(t, test_data_path)
}

func Test_ORA_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/05.json"
	runInstructionTest(t, test_data_path)
}

func Test_ORA_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/09.json"
	runInstructionTest(t, test_data_path)
}

func Test_ORA_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/0d.json"
	runInstructionTest(t, test_data_path)
}

func Test_ORA_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/11.json"
	runInstructionTest(t, test_data_path)
}

func Test_ORA_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/15.json"
	runInstructionTest(t, test_data_path)
}

func Test_ORA_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/19.json"
	runInstructionTest(t, test_data_path)
}

func Test_ORA_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/1d.json"
	runInstructionTest(t, test_data_path)
}

func Test_ASL_A(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/0a.json"
	runInstructionTest(t, test_data_path)
}

func Test_ASL_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/06.json"
	runInstructionTest(t, test_data_path)
}

func Test_ASL_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/0e.json"
	runInstructionTest(t, test_data_path)
}

func Test_ASL_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/16.json"
	runInstructionTest(t, test_data_path)
}

func Test_ASL_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/1e.json"
	runInstructionTest(t, test_data_path)
}

func Test_PHP(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/08.json"
	runInstructionTest(t, test_data_path)
}

func Test_BPL_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/10.json"
	runInstructionTest(t, test_data_path)
}

func Test_CLC(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/18.json"
	runInstructionTest(t, test_data_path)
}

func Test_JSR_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/20.json"
	runInstructionTest(t, test_data_path)
}

func Test_AND_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/21.json"
	runInstructionTest(t, test_data_path)
}

func Test_AND_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/25.json"
	runInstructionTest(t, test_data_path)
}

func Test_AND_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/29.json"
	runInstructionTest(t, test_data_path)
}

func Test_AND_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/2d.json"
	runInstructionTest(t, test_data_path)
}

func Test_AND_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/31.json"
	runInstructionTest(t, test_data_path)
}

func Test_AND_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/35.json"
	runInstructionTest(t, test_data_path)
}

func Test_AND_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/39.json"
	runInstructionTest(t, test_data_path)
}

func Test_AND_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/3d.json"
	runInstructionTest(t, test_data_path)
}

func Test_ROL_A(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/2a.json"
	runInstructionTest(t, test_data_path)
}

func Test_ROL_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/26.json"
	runInstructionTest(t, test_data_path)
}

func Test_ROL_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/2e.json"
	runInstructionTest(t, test_data_path)
}

func Test_ROL_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/36.json"
	runInstructionTest(t, test_data_path)
}

func Test_ROL_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/3e.json"
	runInstructionTest(t, test_data_path)
}

func Test_PLP(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/28.json"
	runInstructionTest(t, test_data_path)
}

func Test_BMI_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/30.json"
	runInstructionTest(t, test_data_path)
}

func Test_SEC(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/38.json"
	runInstructionTest(t, test_data_path)
}

func Test_RTI(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/40.json"
	runInstructionTest(t, test_data_path)
}

func Test_EOR_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/41.json"
	runInstructionTest(t, test_data_path)
}

func Test_EOR_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/45.json"
	runInstructionTest(t, test_data_path)
}

func Test_EOR_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/49.json"
	runInstructionTest(t, test_data_path)
}

func Test_EOR_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/4d.json"
	runInstructionTest(t, test_data_path)
}

func Test_EOR_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/51.json"
	runInstructionTest(t, test_data_path)
}

func Test_EOR_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/55.json"
	runInstructionTest(t, test_data_path)
}

func Test_EOR_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/59.json"
	runInstructionTest(t, test_data_path)
}

func Test_EOR_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/5d.json"
	runInstructionTest(t, test_data_path)
}

func Test_LSR_A(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/4a.json"
	runInstructionTest(t, test_data_path)
}

func Test_LSR_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/46.json"
	runInstructionTest(t, test_data_path)
}

func Test_LSR_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/4e.json"
	runInstructionTest(t, test_data_path)
}

func Test_LSR_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/56.json"
	runInstructionTest(t, test_data_path)
}

func Test_LSR_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/5e.json"
	runInstructionTest(t, test_data_path)
}

func Test_PHA(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/48.json"
	runInstructionTest(t, test_data_path)
}

func Test_BVC_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/50.json"
	runInstructionTest(t, test_data_path)
}

func Test_CLI(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/58.json"
	runInstructionTest(t, test_data_path)
}

func Test_RTS(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/60.json"
	runInstructionTest(t, test_data_path)
}

func Test_ADC_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/61.json"
	runInstructionTest(t, test_data_path)
}

func Test_ADC_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/65.json"
	runInstructionTest(t, test_data_path)
}

func Test_ADC_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/69.json"
	runInstructionTest(t, test_data_path)
}

func Test_ADC_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/6d.json"
	runInstructionTest(t, test_data_path)
}

func Test_ADC_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/71.json"
	runInstructionTest(t, test_data_path)
}

func Test_ADC_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/75.json"
	runInstructionTest(t, test_data_path)
}

func Test_ADC_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/79.json"
	runInstructionTest(t, test_data_path)
}

func Test_ADC_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/7d.json"
	runInstructionTest(t, test_data_path)
}

func Test_ROR_A(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/6a.json"
	runInstructionTest(t, test_data_path)
}

func Test_ROR_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/66.json"
	runInstructionTest(t, test_data_path)
}

func Test_ROR_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/6e.json"
	runInstructionTest(t, test_data_path)
}

func Test_ROR_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/76.json"
	runInstructionTest(t, test_data_path)
}

func Test_ROR_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/7e.json"
	runInstructionTest(t, test_data_path)
}

func Test_PLA(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/68.json"
	runInstructionTest(t, test_data_path)
}

func Test_BVS_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/70.json"
	runInstructionTest(t, test_data_path)
}

func Test_SEI(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/78.json"
	runInstructionTest(t, test_data_path)
}

func Test_STA_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/81.json"
	runInstructionTest(t, test_data_path)
}

func Test_STA_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/85.json"
	runInstructionTest(t, test_data_path)
}

func Test_STA_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/8d.json"
	runInstructionTest(t, test_data_path)
}

func Test_STA_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/91.json"
	runInstructionTest(t, test_data_path)
}

func Test_STA_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/95.json"
	runInstructionTest(t, test_data_path)
}

func Test_STA_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/99.json"
	runInstructionTest(t, test_data_path)
}

func Test_STA_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/9d.json"
	runInstructionTest(t, test_data_path)
}

func Test_STY_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/84.json"
	runInstructionTest(t, test_data_path)
}

func Test_STY_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/8c.json"
	runInstructionTest(t, test_data_path)
}

func Test_STY_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/94.json"
	runInstructionTest(t, test_data_path)
}

func Test_STX_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/86.json"
	runInstructionTest(t, test_data_path)
}

func Test_STX_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/8e.json"
	runInstructionTest(t, test_data_path)
}

func Test_STX_zpY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/96.json"
	runInstructionTest(t, test_data_path)
}

func Test_DEY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/88.json"
	runInstructionTest(t, test_data_path)
}

func Test_TXA(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/8a.json"
	runInstructionTest(t, test_data_path)
}

func Test_BCC_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/90.json"
	runInstructionTest(t, test_data_path)
}

func Test_TYA(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/98.json"
	runInstructionTest(t, test_data_path)
}

func Test_TXS(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/9a.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDY_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/a0.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDY_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/a4.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDY_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/ac.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDY_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/b4.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDY_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/bc.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDA_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/a1.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDA_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/a5.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDA_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/a9.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDA_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/ad.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDA_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/b1.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDA_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/b5.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDA_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/b9.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDA_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/bd.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDX_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/a2.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDX_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/a6.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDX_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/ae.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDX_zpY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/b6.json"
	runInstructionTest(t, test_data_path)
}

func Test_LDX_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/be.json"
	runInstructionTest(t, test_data_path)
}

func Test_TAY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/a8.json"
	runInstructionTest(t, test_data_path)
}

func Test_TAX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/aa.json"
	runInstructionTest(t, test_data_path)
}

func Test_BCS_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/b0.json"
	runInstructionTest(t, test_data_path)
}

func Test_CLV(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/b8.json"
	runInstructionTest(t, test_data_path)
}

func Test_TSX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/ba.json"
	runInstructionTest(t, test_data_path)
}

func Test_CPY_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/c0.json"
	runInstructionTest(t, test_data_path)
}

func Test_CPY_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/c4.json"
	runInstructionTest(t, test_data_path)
}

func Test_CPY_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/cc.json"
	runInstructionTest(t, test_data_path)
}

func Test_CMP_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/c1.json"
	runInstructionTest(t, test_data_path)
}

func Test_CMP_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/c5.json"
	runInstructionTest(t, test_data_path)
}

func Test_CMP_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/c9.json"
	runInstructionTest(t, test_data_path)
}

func Test_CMP_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/cd.json"
	runInstructionTest(t, test_data_path)
}

func Test_CMP_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/d1.json"
	runInstructionTest(t, test_data_path)
}

func Test_CMP_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/d5.json"
	runInstructionTest(t, test_data_path)
}

func Test_CMP_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/d9.json"
	runInstructionTest(t, test_data_path)
}

func Test_CMP_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/dd.json"
	runInstructionTest(t, test_data_path)
}

func Test_DEC_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/c6.json"
	runInstructionTest(t, test_data_path)
}

func Test_DEC_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/ce.json"
	runInstructionTest(t, test_data_path)
}

func Test_DEC_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/d6.json"
	runInstructionTest(t, test_data_path)
}

func Test_DEC_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/de.json"
	runInstructionTest(t, test_data_path)
}

func Test_INY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/c8.json"
	runInstructionTest(t, test_data_path)
}

func Test_DEX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/ca.json"
	runInstructionTest(t, test_data_path)
}

func Test_BNE_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/d0.json"
	runInstructionTest(t, test_data_path)
}

func Test_CLD(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/d8.json"
	runInstructionTest(t, test_data_path)
}

func Test_CPX_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/e0.json"
	runInstructionTest(t, test_data_path)
}

func Test_CPX_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/e4.json"
	runInstructionTest(t, test_data_path)
}

func Test_CPX_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/ec.json"
	runInstructionTest(t, test_data_path)
}

func Test_SBC_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/e1.json"
	runInstructionTest(t, test_data_path)
}

func Test_SBC_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/e5.json"
	runInstructionTest(t, test_data_path)
}

func Test_SBC_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/e9.json"
	runInstructionTest(t, test_data_path)
}

func Test_SBC_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/ed.json"
	runInstructionTest(t, test_data_path)
}

func Test_SBC_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/f1.json"
	runInstructionTest(t, test_data_path)
}

func Test_SBC_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/f5.json"
	runInstructionTest(t, test_data_path)
}

func Test_SBC_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/f9.json"
	runInstructionTest(t, test_data_path)
}

func Test_SBC_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/fd.json"
	runInstructionTest(t, test_data_path)
}

func Test_INC_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/e6.json"
	runInstructionTest(t, test_data_path)
}

func Test_INC_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/ee.json"
	runInstructionTest(t, test_data_path)
}

func Test_INC_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/f6.json"
	runInstructionTest(t, test_data_path)
}

func Test_INC_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/fe.json"
	runInstructionTest(t, test_data_path)
}

func Test_INX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/e8.json"
	runInstructionTest(t, test_data_path)
}

func Test_NOP(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/ea.json"
	runInstructionTest(t, test_data_path)
}

func Test_BEQ_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/f0.json"
	runInstructionTest(t, test_data_path)
}

func Test_SED(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/f8.json"
	runInstructionTest(t, test_data_path)
}

func Test_JMP_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/4c.json"
	runInstructionTest(t, test_data_path)
}

func Test_JMP_ind(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/6c.json"
	runInstructionTest(t, test_data_path)
}

func Test_BIT_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/24.json"
	runInstructionTest(t, test_data_path)
}

func Test_BIT_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/2c.json"
	runInstructionTest(t, test_data_path)
}
