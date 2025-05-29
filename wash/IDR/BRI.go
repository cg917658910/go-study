package IDR

//var BRIMsgTemplate = `11/04/2025 16:26:10 -  Transfer dari XXXXXXXXXXX2504 dengan nomor rekening tujuan XXXXXX0336 sebesar Rp10.012,00 BERHASIL. Info lebih lanjut hubungi Call Center BRI 1500017`

var (
	BRIMsgTemplate *BRIMsg
)

type (
	BRIMsg struct {
		payTimeRegex string
		payCoinRegex string
	}
)

func init() {
	BRIMsgTemplate = NewBRIMsg()
}

func NewBRIMsg() *BRIMsg {
	// TODO: get from config
	return &BRIMsg{
		payTimeRegex: `(\d{2}/\d{2}/\d{4} \d{2}:\d{2}:\d{2})`,
		payCoinRegex: `Rp(\d+\.\d{3})`,
	}
}

func (m *BRIMsg) GetPayTimeRegex() string {
	return m.payTimeRegex
}

func (m *BRIMsg) GetPayCoinRegex() string {
	return m.payCoinRegex
}
