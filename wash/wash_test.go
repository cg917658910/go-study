package wash

import (
	"testing"

	"github.com/cg917658910/go-study/wash/IDR"
)

func TestWashBRIMsg(t *testing.T) {
	var msg = "11/04/2025 16:26:10 -  Transfer dari XXXXXXXXXXX2504 dengan nomor rekening tujuan XXXXXX0336 sebesar Rp10.012,00 BERHASIL. Info lebih lanjut hubungi Call Center BRI 1500017"
	res := WashMsg(msg, IDR.BRIMsgTemplate)
	if res.Error != nil {
		t.Errorf("WashMsg expected nil, got error: %v", res.Error)
		return
	}
	if res.PayCoin != 10.01200 {
		t.Errorf("WashMsg expected PayCoin 10.01200, got: %f", res.PayCoin)
		return
	}
}
func TestWashMsgParseTime(t *testing.T) {
	var timeStr = "11/04/2025 16:26:10"
	res, err := parseTime(timeStr)
	if err != nil {
		t.Errorf("parseTime expected nil, got error: %v", err)
		return
	}
	if res == nil {
		t.Errorf("parseTime result expected not nil, got: %v", res)
		return
	}
}
