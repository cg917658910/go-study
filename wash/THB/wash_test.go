package THB

import (
	"fmt"
	"testing"
)

func TestWashBBL(t *testing.T) {
	msg := `ถอน/โอน/จ่ายเงินจากบ/ชX0280 ผ่านMB 1,042.00บ ใช้ได้36,447.43บ`
	msgTime := "2025-05-17 16:27:51"
	result := ExtractBBL(msg, msgTime)
	fmt.Println("res: ", result)
	if result.PayCoin != 41042.00 {
		t.Errorf("wash BBL expected 41042.00, but got %v", result.PayCoin)
	}
}
func TestWashGSB(t *testing.T) {
	//msg1 := `เงินเข้า: มีการฝาก/โอนเงิน 10.00 บาท จากบัญชี KBNK`
	msg := `เงินเข้า: มีการฝาก/โอนเงิน 74.76 บาท จากบัญชี KBNK 0013XXXX1945 เข้าบัญชี GSBA 0204XXXX0743 วันที่ 24 พ.ค. 2568 เวลา 15:04 น. คงเหลือ 2,702.93 บาท`
	result := ExtractGSB(msg)
	fmt.Println("res: ", result)
	if result.PayCoin != 10.00 {
		t.Errorf("wash GSB expected 10.00, but got %v", result.PayCoin)
	}
}
func TestWashTTB(t *testing.T) {
	msg := `27-05@16:25 บชX46746X:เงินเข้า 499.95บ ใช้ได้ 11,871.96บ`
	result := ExtractTTB(msg)
	fmt.Println("res: ", result)
	if result.PayCoin != 499.95 {
		t.Errorf("wash TTB expected 499.95, but got %v", result.PayCoin)
	}
}
