package cpu_test

import "testing"

func Test_BRK(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/00.json"
	run_insturction_test(t, test_data_path)
}

func Test_ORA_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/01.json"
	run_insturction_test(t, test_data_path)
}

func Test_ORA_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/05.json"
	run_insturction_test(t, test_data_path)
}

func Test_ORA_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/09.json"
	run_insturction_test(t, test_data_path)
}

func Test_ORA_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/0D.json"
	run_insturction_test(t, test_data_path)
}

func Test_ORA_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/11.json"
	run_insturction_test(t, test_data_path)
}

func Test_ORA_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/15.json"
	run_insturction_test(t, test_data_path)
}

func Test_ORA_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/19.json"
	run_insturction_test(t, test_data_path)
}

func Test_ORA_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/1D.json"
	run_insturction_test(t, test_data_path)
}

func Test_ASL_A(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/0A.json"
	run_insturction_test(t, test_data_path)
}

func Test_ASL_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/06.json"
	run_insturction_test(t, test_data_path)
}

func Test_ASL_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/0E.json"
	run_insturction_test(t, test_data_path)
}

func Test_ASL_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/16.json"
	run_insturction_test(t, test_data_path)
}

func Test_ASL_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/1E.json"
	run_insturction_test(t, test_data_path)
}

func Test_PHP(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/08.json"
	run_insturction_test(t, test_data_path)
}

func Test_BPL_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/10.json"
	run_insturction_test(t, test_data_path)
}

func Test_CLC(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/18.json"
	run_insturction_test(t, test_data_path)
}

func Test_JSR_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/20.json"
	run_insturction_test(t, test_data_path)
}

func Test_AND_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/21.json"
	run_insturction_test(t, test_data_path)
}

func Test_AND_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/25.json"
	run_insturction_test(t, test_data_path)
}

func Test_AND_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/29.json"
	run_insturction_test(t, test_data_path)
}

func Test_AND_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/2D.json"
	run_insturction_test(t, test_data_path)
}

func Test_AND_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/31.json"
	run_insturction_test(t, test_data_path)
}

func Test_AND_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/35.json"
	run_insturction_test(t, test_data_path)
}

func Test_AND_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/39.json"
	run_insturction_test(t, test_data_path)
}

func Test_AND_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/3D.json"
	run_insturction_test(t, test_data_path)
}

func Test_ROL_A(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/2A.json"
	run_insturction_test(t, test_data_path)
}

func Test_ROL_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/26.json"
	run_insturction_test(t, test_data_path)
}

func Test_ROL_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/2E.json"
	run_insturction_test(t, test_data_path)
}

func Test_ROL_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/36.json"
	run_insturction_test(t, test_data_path)
}

func Test_ROL_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/3E.json"
	run_insturction_test(t, test_data_path)
}

func Test_PLP(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/28.json"
	run_insturction_test(t, test_data_path)
}

func Test_BMI_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/30.json"
	run_insturction_test(t, test_data_path)
}

func Test_SEC(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/38.json"
	run_insturction_test(t, test_data_path)
}

func Test_RTI(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/40.json"
	run_insturction_test(t, test_data_path)
}

func Test_EOR_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/41.json"
	run_insturction_test(t, test_data_path)
}

func Test_EOR_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/45.json"
	run_insturction_test(t, test_data_path)
}

func Test_EOR_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/49.json"
	run_insturction_test(t, test_data_path)
}

func Test_EOR_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/4D.json"
	run_insturction_test(t, test_data_path)
}

func Test_EOR_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/51.json"
	run_insturction_test(t, test_data_path)
}

func Test_EOR_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/55.json"
	run_insturction_test(t, test_data_path)
}

func Test_EOR_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/59.json"
	run_insturction_test(t, test_data_path)
}

func Test_EOR_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/5D.json"
	run_insturction_test(t, test_data_path)
}

func Test_LSR_A(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/4A.json"
	run_insturction_test(t, test_data_path)
}

func Test_LSR_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/46.json"
	run_insturction_test(t, test_data_path)
}

func Test_LSR_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/4E.json"
	run_insturction_test(t, test_data_path)
}

func Test_LSR_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/56.json"
	run_insturction_test(t, test_data_path)
}

func Test_LSR_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/5E.json"
	run_insturction_test(t, test_data_path)
}

func Test_PHA(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/48.json"
	run_insturction_test(t, test_data_path)
}

func Test_BVC_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/50.json"
	run_insturction_test(t, test_data_path)
}

func Test_CLI(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/58.json"
	run_insturction_test(t, test_data_path)
}

func Test_RTS(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/60.json"
	run_insturction_test(t, test_data_path)
}

func Test_ADC_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/61.json"
	run_insturction_test(t, test_data_path)
}

func Test_ADC_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/65.json"
	run_insturction_test(t, test_data_path)
}

func Test_ADC_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/69.json"
	run_insturction_test(t, test_data_path)
}

func Test_ADC_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/6D.json"
	run_insturction_test(t, test_data_path)
}

func Test_ADC_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/71.json"
	run_insturction_test(t, test_data_path)
}

func Test_ADC_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/75.json"
	run_insturction_test(t, test_data_path)
}

func Test_ADC_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/79.json"
	run_insturction_test(t, test_data_path)
}

func Test_ADC_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/7D.json"
	run_insturction_test(t, test_data_path)
}

func Test_ROR_A(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/6A.json"
	run_insturction_test(t, test_data_path)
}

func Test_ROR_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/66.json"
	run_insturction_test(t, test_data_path)
}

func Test_ROR_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/6E.json"
	run_insturction_test(t, test_data_path)
}

func Test_ROR_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/76.json"
	run_insturction_test(t, test_data_path)
}

func Test_ROR_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/7E.json"
	run_insturction_test(t, test_data_path)
}

func Test_PLA(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/68.json"
	run_insturction_test(t, test_data_path)
}

func Test_BVS_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/70.json"
	run_insturction_test(t, test_data_path)
}

func Test_SEI(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/78.json"
	run_insturction_test(t, test_data_path)
}

func Test_STA_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/81.json"
	run_insturction_test(t, test_data_path)
}

func Test_STA_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/85.json"
	run_insturction_test(t, test_data_path)
}

func Test_STA_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/8D.json"
	run_insturction_test(t, test_data_path)
}

func Test_STA_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/91.json"
	run_insturction_test(t, test_data_path)
}

func Test_STA_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/95.json"
	run_insturction_test(t, test_data_path)
}

func Test_STA_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/99.json"
	run_insturction_test(t, test_data_path)
}

func Test_STA_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/9D.json"
	run_insturction_test(t, test_data_path)
}

func Test_STY_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/84.json"
	run_insturction_test(t, test_data_path)
}

func Test_STY_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/8C.json"
	run_insturction_test(t, test_data_path)
}

func Test_STY_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/94.json"
	run_insturction_test(t, test_data_path)
}

func Test_STX_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/86.json"
	run_insturction_test(t, test_data_path)
}

func Test_STX_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/8E.json"
	run_insturction_test(t, test_data_path)
}

func Test_STX_zpY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/96.json"
	run_insturction_test(t, test_data_path)
}

func Test_DEY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/88.json"
	run_insturction_test(t, test_data_path)
}

func Test_TXA(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/8A.json"
	run_insturction_test(t, test_data_path)
}

func Test_BCC_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/90.json"
	run_insturction_test(t, test_data_path)
}

func Test_TYA(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/98.json"
	run_insturction_test(t, test_data_path)
}

func Test_TXS(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/9A.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDY_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/A0.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDY_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/A4.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDY_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/AC.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDY_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/B4.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDY_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/BC.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDA_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/A1.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDA_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/A5.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDA_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/A9.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDA_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/AD.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDA_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/B1.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDA_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/B5.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDA_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/B9.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDA_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/BD.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDX_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/A2.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDX_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/A6.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDX_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/AE.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDX_zpY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/B6.json"
	run_insturction_test(t, test_data_path)
}

func Test_LDX_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/BE.json"
	run_insturction_test(t, test_data_path)
}

func Test_TAY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/A8.json"
	run_insturction_test(t, test_data_path)
}

func Test_TAX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/AA.json"
	run_insturction_test(t, test_data_path)
}

func Test_BCS_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/B0.json"
	run_insturction_test(t, test_data_path)
}

func Test_CLV(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/B8.json"
	run_insturction_test(t, test_data_path)
}

func Test_TSX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/BA.json"
	run_insturction_test(t, test_data_path)
}

func Test_CPY_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/C0.json"
	run_insturction_test(t, test_data_path)
}

func Test_CPY_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/C4.json"
	run_insturction_test(t, test_data_path)
}

func Test_CPY_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/CC.json"
	run_insturction_test(t, test_data_path)
}

func Test_CMP_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/C1.json"
	run_insturction_test(t, test_data_path)
}

func Test_CMP_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/C5.json"
	run_insturction_test(t, test_data_path)
}

func Test_CMP_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/C9.json"
	run_insturction_test(t, test_data_path)
}

func Test_CMP_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/CD.json"
	run_insturction_test(t, test_data_path)
}

func Test_CMP_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/D1.json"
	run_insturction_test(t, test_data_path)
}

func Test_CMP_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/D5.json"
	run_insturction_test(t, test_data_path)
}

func Test_CMP_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/D9.json"
	run_insturction_test(t, test_data_path)
}

func Test_CMP_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/DD.json"
	run_insturction_test(t, test_data_path)
}

func Test_DEC_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/C6.json"
	run_insturction_test(t, test_data_path)
}

func Test_DEC_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/CE.json"
	run_insturction_test(t, test_data_path)
}

func Test_DEC_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/D6.json"
	run_insturction_test(t, test_data_path)
}

func Test_DEC_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/DE.json"
	run_insturction_test(t, test_data_path)
}

func Test_INY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/C8.json"
	run_insturction_test(t, test_data_path)
}

func Test_DEX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/CA.json"
	run_insturction_test(t, test_data_path)
}

func Test_BNE_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/D0.json"
	run_insturction_test(t, test_data_path)
}

func Test_CLD(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/D8.json"
	run_insturction_test(t, test_data_path)
}

func Test_CPX_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/E0.json"
	run_insturction_test(t, test_data_path)
}

func Test_CPX_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/E4.json"
	run_insturction_test(t, test_data_path)
}

func Test_CPX_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/EC.json"
	run_insturction_test(t, test_data_path)
}

func Test_SBC_Xin(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/E1.json"
	run_insturction_test(t, test_data_path)
}

func Test_SBC_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/E5.json"
	run_insturction_test(t, test_data_path)
}

func Test_SBC_imd(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/E9.json"
	run_insturction_test(t, test_data_path)
}

func Test_SBC_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/ED.json"
	run_insturction_test(t, test_data_path)
}

func Test_SBC_inY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/F1.json"
	run_insturction_test(t, test_data_path)
}

func Test_SBC_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/F5.json"
	run_insturction_test(t, test_data_path)
}

func Test_SBC_abY(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/F9.json"
	run_insturction_test(t, test_data_path)
}

func Test_SBC_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/FD.json"
	run_insturction_test(t, test_data_path)
}

func Test_INC_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/E6.json"
	run_insturction_test(t, test_data_path)
}

func Test_INC_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/EE.json"
	run_insturction_test(t, test_data_path)
}

func Test_INC_zpX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/F6.json"
	run_insturction_test(t, test_data_path)
}

func Test_INC_abX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/FE.json"
	run_insturction_test(t, test_data_path)
}

func Test_INX(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/E8.json"
	run_insturction_test(t, test_data_path)
}

func Test_NOP(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/EA.json"
	run_insturction_test(t, test_data_path)
}

func Test_BEQ_rel(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/F0.json"
	run_insturction_test(t, test_data_path)
}

func Test_SED(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/F8.json"
	run_insturction_test(t, test_data_path)
}

func Test_JMP_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/4C.json"
	run_insturction_test(t, test_data_path)
}

func Test_JMP_ind(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/6C.json"
	run_insturction_test(t, test_data_path)
}

func Test_BIT_zpg(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/24.json"
	run_insturction_test(t, test_data_path)
}

func Test_BIT_abs(t *testing.T) {
	test_data_path := "../test_data/single_step_tests/2C.json"
	run_insturction_test(t, test_data_path)
}
