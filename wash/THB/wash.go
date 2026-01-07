package THB

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Extracted struct {
	PayCoin float64
	PayTime int64
	Balance float64
}

// 工具函数：泰文月份转数字
func thaiMonthToNum(month string) int {
	arr := []string{"", "ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.", "ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."}
	for i, v := range arr {
		if month == v {
			return i
		}
	}
	return 0
}

// 工具函数：英文月份转数字
func engMonthToNum(month string) int {
	arr := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sept", "Oct", "Nov", "Dec"}
	for i, v := range arr {
		if month == v {
			return i + 1
		}
	}
	return 0
}

// SCB
func ExtractSCB(ramk string) Extracted {
	remarkArr := strings.Fields(ramk)
	var payCoin float64
	var payTime int64
	var balance float64
	if len(remarkArr) < 2 {
		return Extracted{}
	}
	if remarkArr[0] == "เงิน" {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[1], ",", ""), 64)
	} else if remarkArr[0] == "Transfer" && len(remarkArr) > 5 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[5], ",", ""), 64)
	} else if remarkArr[0] == "Cash/transfer" && len(remarkArr) > 4 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[4], ",", ""), 64)
	} else {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[1], ",", ""), 64)
	}
	// 时间
	if len(remarkArr) > 3 {
		timeArr := remarkArr[3]
		if len(timeArr) == 11 {
			year := time.Now().Year()
			payTimeStr := strconv.Itoa(year) + "-" + timeArr[3:5] + "-" + timeArr[0:2] + " " + timeArr[6:8] + ":" + timeArr[9:11]
			t, _ := time.ParseInLocation("2006-1-2 15:04", payTimeStr, time.Local)
			payTime = t.Unix()
		}
	}
	// 余额
	last := remarkArr[len(remarkArr)-1]
	re := regexp.MustCompile(`(\d+\.\d+)`)
	if m := re.FindStringSubmatch(last); len(m) > 1 {
		balance, _ = strconv.ParseFloat(m[1], 64)
	}
	return Extracted{PayCoin: payCoin, PayTime: payTime, Balance: balance}
}

// SCB读取
func ExtractSCBRead(ramk string) Extracted {
	return ExtractSCB(ramk)
}

// SCB通知
func ExtractSCBNotify(ramk string) Extracted {
	remarkArr := strings.Fields(ramk)
	var payCoin float64
	var payTime int64
	if len(remarkArr) == 15 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[1], ",", ""), 64)
		month := thaiMonthToNum(remarkArr[11])
		if month > 0 && month < 13 {
			year := time.Now().Year()
			payTimeStr := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-" + remarkArr[10] + " " + remarkArr[14]
			t, _ := time.ParseInLocation("2006-1-2 15:04", payTimeStr, time.Local)
			payTime = t.Unix()
		}
	}
	return Extracted{PayCoin: payCoin, PayTime: payTime}
}

// SCB流水
func ExtractSCBWater(msgTime string, coin float64) Extracted {
	// msgTime 例：01/05/2024 12:34
	parts := strings.Fields(msgTime)
	if len(parts) < 2 {
		return Extracted{}
	}
	dateParts := strings.Split(parts[0], "/")
	if len(dateParts) != 3 {
		return Extracted{}
	}
	payTimeStr := dateParts[2] + "-" + dateParts[1] + "-" + dateParts[0] + " " + parts[1]
	t, _ := time.ParseInLocation("2006-1-2 15:04", payTimeStr, time.Local)
	return Extracted{PayCoin: coin, PayTime: t.Unix()}
}

// TM流水
func ExtractTMWater(ramk string) Extracted {
	var obj struct {
		Time  string `json:"time"`
		Money string `json:"money"`
	}
	_ = json.Unmarshal([]byte(ramk), &obj)
	payCoin, _ := strconv.ParseFloat(strings.ReplaceAll(obj.Money, ",", ""), 64)
	timeArr := strings.Fields(obj.Time)
	if len(timeArr) < 4 {
		return Extracted{}
	}
	// 泰历转公历
	year, _ := strconv.Atoi(timeArr[2])
	year = year - 543
	month := thaiMonthToNum(timeArr[1])
	payTimeStr := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-" + timeArr[0] + " " + timeArr[3]
	t, _ := time.ParseInLocation("2006-1-2 15:04", payTimeStr, time.Local)
	return Extracted{PayCoin: payCoin, PayTime: t.Unix()}
}

// KTB
func ExtractKTB(ramk string, msgTime string) Extracted {
	remark := strings.Fields(ramk)
	var payCoin float64
	var payTime int64
	var balance float64
	if len(remark) > 1 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remark[1], ",", ""), 64)
	}
	t, _ := time.ParseInLocation("2006-01-02 15:04", msgTime, time.Local)
	payTime = t.Unix()
	return Extracted{PayCoin: payCoin, PayTime: payTime, Balance: balance}
}

// KTBLine
func ExtractKTBLine(ramk string, msgTime string) Extracted {
	remark := strings.Fields(ramk)
	var payCoin float64
	var payTime int64
	var balance float64
	if len(remark) > 1 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remark[1], ",", ""), 64)
	}
	t, _ := time.ParseInLocation("2006-01-02 15:04", msgTime, time.Local)
	payTime = t.Unix()
	return Extracted{PayCoin: payCoin, PayTime: payTime, Balance: balance}
}

// KTB通知
func ExtractKTBNotice(ramk string, msgTime string) Extracted {
	return ExtractKTB(ramk, msgTime)
}

// KTB流水
func ExtractKTBWater(postMsg string) Extracted {
	var obj struct {
		MsgTime string  `json:"msg_time"`
		Coin    float64 `json:"coin"`
	}
	_ = json.Unmarshal([]byte(postMsg), &obj)
	parts := strings.Fields(obj.MsgTime)
	if len(parts) < 2 {
		return Extracted{}
	}
	dateParts := strings.Split(parts[0], "-")
	if len(dateParts) != 3 {
		return Extracted{}
	}
	payTimeStr := dateParts[2] + "-" + dateParts[1] + "-" + dateParts[0] + " " + parts[1]
	t, _ := time.ParseInLocation("2006-1-2 15:04", payTimeStr, time.Local)
	return Extracted{PayCoin: obj.Coin, PayTime: t.Unix()}
}

// KBANK通知
func ExtractKBANKNotify(ramk string) Extracted {
	remarkArr := strings.Fields(ramk)
	var payCoin float64
	var payTime int64
	//arr := []string{"", "ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.", "ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."}
	if len(remarkArr) == 14 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[4], ",", ""), 64)
		year, _ := strconv.Atoi(remarkArr[10])
		year = 2500 + year - 543
		month := thaiMonthToNum(remarkArr[9])
		day := remarkArr[8]
		payTimeStr := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-" + day + " " + remarkArr[12]
		t, _ := time.ParseInLocation("2006-1-2 15:04", payTimeStr, time.Local)
		payTime = t.Unix()
	}
	return Extracted{PayCoin: payCoin, PayTime: payTime}
}

// KBANK读取
func ExtractKBANKRead(ramk string) Extracted {
	remark := strings.ReplaceAll(ramk, "KBank:", "")
	remarkArr := strings.Fields(remark)
	var payCoin float64
	var payTime int64
	var balance float64
	if len(remarkArr) == 9 || len(remarkArr) == 14 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[5], ",", ""), 64)
	} else if len(remarkArr) == 10 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[6], ",", ""), 64)
		if !isNumeric(remarkArr[6]) {
			payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[5], ",", ""), 64)
		}
	} else if len(remarkArr) == 4 || len(remarkArr) == 5 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[3], ",", ""), 64)
	} else if len(remarkArr) == 6 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[4], ",", ""), 64)
	}
	// 余额
	if len(remarkArr) > 2 {
		balanceStr := remarkArr[len(remarkArr)-2]
		balance, _ = strconv.ParseFloat(strings.ReplaceAll(balanceStr, ",", ""), 64)
	}
	// 时间
	if len(remarkArr) > 1 {
		t := strings.Split(remarkArr[0], "/")
		if len(t) == 3 {
			year := time.Now().Year()
			payTimeStr := strconv.Itoa(year) + "-" + t[1] + "-" + t[0] + " " + remarkArr[1]
			tm, _ := time.ParseInLocation("2006-01-02 15:04", payTimeStr, time.Local)
			payTime = tm.Unix()
		}
	}
	return Extracted{PayCoin: payCoin, PayTime: payTime, Balance: balance}
}

// KBANK流水
func ExtractKBANKWater(postMsg string) Extracted {
	var obj struct {
		MsgTime string      `json:"msg_time"`
		Coin    interface{} `json:"coin"`
	}
	_ = json.Unmarshal([]byte(postMsg), &obj)
	parts := strings.Fields(obj.MsgTime)
	if len(parts) < 4 {
		return Extracted{}
	}
	year, _ := strconv.Atoi(parts[2])
	year = 2500 + year - 543
	month := thaiMonthToNum(parts[1])
	day := parts[0]
	payTimeStr := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-" + day + " " + parts[3]
	t, _ := time.ParseInLocation("2006-1-2 15:04", payTimeStr, time.Local)
	coinStr := ""
	switch v := obj.Coin.(type) {
	case string:
		coinStr = v
	case float64:
		coinStr = strconv.FormatFloat(v, 'f', 2, 64)
	}
	payCoin, _ := strconv.ParseFloat(strings.ReplaceAll(coinStr, ",", ""), 64)
	return Extracted{PayCoin: payCoin, PayTime: t.Unix()}
}

// BBL流水
func ExtractBBLWater(postMsg string) Extracted {
	var obj struct {
		MsgTime string  `json:"msg_time"`
		Coin    float64 `json:"coin"`
	}
	_ = json.Unmarshal([]byte(postMsg), &obj)
	parts := strings.Fields(obj.MsgTime)
	if len(parts) < 4 {
		return Extracted{}
	}
	month := engMonthToNum(parts[1])
	payTimeStr := parts[2] + "-" + strconv.Itoa(month) + "-" + parts[0] + " " + parts[3]
	t, _ := time.ParseInLocation("2006-1-2 15:04", payTimeStr, time.Local)
	return Extracted{PayCoin: obj.Coin, PayTime: t.Unix()}
}

// BBL
func ExtractBBL(remark string, msgTime string) Extracted {
	remarkArr := strings.Fields(remark)
	var payCoin float64
	var balance float64
	if len(remarkArr) > 0 && remarkArr[0] == "PromptPay" && len(remarkArr) > 8 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[8], ",", ""), 64)
	} else if len(remarkArr) > 2 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[2], ",", ""), 64)
	}
	if len(remarkArr) > 0 {
		last := remarkArr[len(remarkArr)-1]
		re := regexp.MustCompile(`(\d+\.\d+)`)
		if m := re.FindStringSubmatch(last); len(m) > 1 {
			balance, _ = strconv.ParseFloat(m[1], 64)
		}
	}
	t, _ := time.ParseInLocation("2006-01-02 15:04", msgTime, time.Local)
	return Extracted{PayCoin: payCoin, PayTime: t.Unix(), Balance: balance}
}

// BAAC
func ExtractBAAC(ramk string, msgTime string) Extracted {
	remark := strings.Fields(ramk)
	var payCoin float64
	var payTime int64
	var balance float64
	if len(remark) > 6 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remark[6], ",", ""), 64)
		timeArr := strings.Split(remark[0], "/")
		if len(timeArr) == 3 {
			year := time.Now().Year()
			payTimeStr := strconv.Itoa(year) + "-" + timeArr[1] + "-" + timeArr[0] + " " + remark[1]
			t, _ := time.ParseInLocation("2006-01-02 15:04", payTimeStr, time.Local)
			payTime = t.Unix()
		} else {
			t, _ := time.ParseInLocation("2006-01-02 15:04", msgTime, time.Local)
			payTime = t.Unix()
		}
		if len(remark) > 9 {
			balance, _ = strconv.ParseFloat(strings.ReplaceAll(remark[9], ",", ""), 64)
		}
	}
	return Extracted{PayCoin: payCoin, PayTime: payTime, Balance: balance}
}

// TTB
func ExtractTTB(remark string) Extracted {
	var payCoin float64
	var payTime int64
	var balance float64
	reCoin := regexp.MustCompile(`([\d,]+\.\d{2})บ\.`)
	if m := reCoin.FindStringSubmatch(remark); len(m) > 1 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(m[1], ",", ""), 64)
	}
	reTime := regexp.MustCompile(`(\d{2}/\d{2}/\d{2})@(\d{2}:\d{2})`)
	if m := reTime.FindStringSubmatch(remark); len(m) == 3 {
		ymd := strings.Split(m[1], "/")
		if len(ymd) == 3 {
			payTimeStr := "20" + ymd[2] + "-" + ymd[1] + "-" + ymd[0] + " " + m[2]
			t, _ := time.ParseInLocation("2006-01-02 15:04", payTimeStr, time.Local)
			payTime = t.Unix()
		}
	}
	reBal := regexp.MustCompile(`เหลือ([\d,]+\.\d{2})บ`)
	if m := reBal.FindStringSubmatch(remark); len(m) > 1 {
		balance, _ = strconv.ParseFloat(strings.ReplaceAll(m[1], ",", ""), 64)
	}
	return Extracted{PayCoin: payCoin, PayTime: payTime, Balance: balance}
}

// TTB读取/TTB通知
func ExtractTTBRead(remark string) Extracted {
	return ExtractTTB(remark)
}
func ExtractTTBNotify(remark string) Extracted {
	return ExtractTTB(remark)
}

// KKRLSCLI
func ExtractKKRLSCLI(postMsg string) Extracted {
	var obj struct {
		Coin float64 `json:"coin"`
		Time int64   `json:"time"`
	}
	_ = json.Unmarshal([]byte(postMsg), &obj)
	return Extracted{PayCoin: obj.Coin, PayTime: obj.Time}
}

// BAY
func ExtractBAY(remark string) Extracted {
	arr := strings.Fields(remark)
	var payCoin, balance float64
	var payTime int64
	if len(arr) > 0 && arr[0] == "โอนเข้า" && len(arr) > 6 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(arr[3], ",", ""), 64)
		re := regexp.MustCompile(`\((\d{2})/(\d{2})/(\d{2}),(\d{2}):(\d{2})\)`)
		m := re.FindStringSubmatch(remark)
		if len(m) == 6 {
			y, _ := strconv.Atoi(m[3])
			y += 2000
			payTimeStr := strconv.Itoa(y) + "-" + m[2] + "-" + m[1] + " " + m[4] + ":" + m[5]
			t, _ := time.ParseInLocation("2006-01-02 15:04", payTimeStr, time.Local)
			payTime = t.Unix()
		}
		if len(arr) > 13 {
			balance, _ = strconv.ParseFloat(strings.ReplaceAll(arr[13], ",", ""), 64)
		}
	}
	return Extracted{PayCoin: payCoin, PayTime: payTime, Balance: balance}
}

// TM流水ios
func ExtractTMIosWater(ramk string) Extracted {
	var obj struct {
		Time  string `json:"time"`
		Money string `json:"money"`
	}
	_ = json.Unmarshal([]byte(ramk), &obj)
	payCoin, _ := strconv.ParseFloat(strings.ReplaceAll(obj.Money, ",", ""), 64)
	timeArr := strings.Fields(obj.Time)
	if len(timeArr) < 4 {
		return Extracted{}
	}
	year, _ := strconv.Atoi(timeArr[2])
	year = year - 543
	month := thaiMonthToNum(timeArr[1])
	payTimeStr := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-" + timeArr[0] + " " + timeArr[3]
	t, _ := time.ParseInLocation("2006-01-02 15:04", payTimeStr, time.Local)
	return Extracted{PayCoin: payCoin, PayTime: t.Unix()}
}

// SwooleTM
func ExtractSwooleTM(coin string, msgTime int64) Extracted {
	coin = strings.ReplaceAll(coin, " ", "")
	coin = strings.ReplaceAll(coin, "฿", "")
	coin = strings.ReplaceAll(coin, ",", "")
	payCoin, _ := strconv.ParseFloat(coin, 64)
	return Extracted{PayCoin: payCoin, PayTime: msgTime}
}

// GSB
func ExtractGSB(ramk string) Extracted {
	remark := strings.Fields(ramk)
	var payCoin float64
	var payTime int64
	var balance float64
	if len(remark) == 0 {
		return Extracted{}
	}
	if remark[0] == "คุณได้รับเงิน" {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remark[1], ",", ""), 64)
	} else if len(remark) > 2 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(remark[2], ",", ""), "฿", ""), 64)
	}
	re := regexp.MustCompile(`วันที่ (\d+) (.*?) (\d+) เวลา (\d+):(\d+) น`)
	m := re.FindStringSubmatch(ramk)
	if len(m) == 6 {
		month := thaiMonthToNum(m[2])
		year := time.Now().Year()
		payTimeStr := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-" + m[1] + " " + m[4] + ":" + m[5] + ":00"
		t, err := time.Parse(time.DateTime, payTimeStr)
		if err != nil {
			fmt.Printf("parse time err：%v", err)
		}
		payTime = t.Unix()
		fmt.Println("time", payTimeStr, t.String(), payTime)
	}
	if len(remark) > 18 {
		balance, _ = strconv.ParseFloat(strings.ReplaceAll(remark[18], ",", ""), 64)
	}
	return Extracted{PayCoin: payCoin, PayTime: payTime, Balance: balance}
}

// GSB读取/GSB通知/GSBLine
func ExtractGSBRead(ramk string) Extracted   { return ExtractGSB(ramk) }
func ExtractGSBNotify(ramk string) Extracted { return ExtractGSB(ramk) }
func ExtractGSBLine(ramk string) Extracted   { return ExtractGSB(ramk) }

// python-SCBGH
func ExtractSCBGH(postMsg string) Extracted {
	var obj struct {
		Money string `json:"money"`
		Time  string `json:"time"`
	}
	_ = json.Unmarshal([]byte(postMsg), &obj)
	payCoin, _ := strconv.ParseFloat(strings.ReplaceAll(obj.Money, ",", ""), 64)
	t, _ := time.ParseInLocation("2006-01-02 15:04", obj.Time, time.Local)
	return Extracted{PayCoin: payCoin, PayTime: t.Unix()}
}

// python-KTBGH
func ExtractKTBGH(postMsg string) Extracted {
	return ExtractSCBGH(postMsg)
}

// python-Kbankgh
func ExtractKbankGH(postMsg string) Extracted {
	return ExtractSCBGH(postMsg)
}

// python-Ttbgh
func ExtractTtbGH(postMsg string) Extracted {
	return ExtractSCBGH(postMsg)
}

// python-BAYGH
func ExtractBayGH(postMsg string) Extracted {
	return ExtractSCBGH(postMsg)
}

// tm_protocol_water
func ExtractTMProtocolWater(postMsg string) Extracted {
	var obj struct {
		DateTime string `json:"date_time"`
		Amount   string `json:"amount"`
	}
	_ = json.Unmarshal([]byte(postMsg), &obj)
	payCoin, _ := strconv.ParseFloat(strings.ReplaceAll(obj.Amount, ",", ""), 64)
	t, _ := time.ParseInLocation("2006-01-02 15:04", obj.DateTime, time.Local)
	return Extracted{PayCoin: payCoin, PayTime: t.Unix()}
}

// scb_protocol_water
func ExtractSCBProtocolWater(postMsg string) Extracted {
	return ExtractTMProtocolWater(postMsg)
}

// ttb_protocol_water
func ExtractTTBProtocolWater(postMsg string) Extracted {
	return ExtractTMProtocolWater(postMsg)
}

// 工具
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
