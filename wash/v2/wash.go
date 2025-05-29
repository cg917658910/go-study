package v2

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TransactionData struct {
	PayCoin  float64
	PayTime  time.Time
	Balance  float64
	BankName string
}

// Main extraction function that routes to specific handlers
func ExtractTransactionData(bankType string, ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}
	var err error

	switch bankType {
	case "SCB", "SCB读取":
		data, err = handleSCB(ramk, msgTime)
	case "SCB通知":
		data, err = handleSCBNotify(ramk, msgTime)
	case "SCB流水":
		data, err = handleSCBWater(ramk, msgTime)
	case "KTB":
		data, err = handleKTB(ramk, msgTime)
	case "KTBLine":
		data, err = handleKTBLine(ramk, msgTime)
	case "KTB通知":
		data, err = handleKTBNotice(ramk, msgTime)
	case "KTB流水":
		data, err = handleKTBWater(ramk, msgTime)
	case "KBANK通知":
		data, err = handleKKRNotify(ramk, msgTime)
	case "KBANK读取":
		data, err = handleKKRRead(ramk, msgTime)
	case "KBANK流水":
		data, err = handleKKRWater(ramk, msgTime)
	case "BBL流水":
		data, err = handleBBLWater(ramk, msgTime)
	case "BBL":
		data, err = handleBBL(ramk, msgTime)
	case "BAAC":
		data, err = handleBAAC(ramk, msgTime)
	case "TTB", "TTB读取", "TTB通知":
		data, err = handleTTB(ramk, msgTime)
	default:
		return data, fmt.Errorf("unsupported bank type: %s", bankType)
	}

	return data, err
}

// SCB Bank handler
func handleSCB(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}
	remarkArr := strings.Fields(ramk)

	// Check if it's a withdrawal
	if strings.Contains(ramk, "ถอน") || strings.Contains(ramk, "โอนเงิน") ||
		strings.Contains(ramk, "ถอน/โอนเงิน") || strings.Contains(ramk, "True") {
		return data, fmt.Errorf("not a deposit transaction")
	}

	// Extract amount
	var amountStr string
	if remarkArr[0] == "เงิน" {
		amountStr = strings.ReplaceAll(remarkArr[1], ",", "")
	} else if remarkArr[0] == "Transfer" {
		amountStr = strings.ReplaceAll(remarkArr[5], ",", "")
	} else if remarkArr[0] == "Cash/transfer" {
		amountStr = strings.ReplaceAll(remarkArr[4], ",", "")
	} else {
		amountStr = strings.ReplaceAll(remarkArr[1], ",", "")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return data, fmt.Errorf("error parsing amount: %v", err)
	}
	data.PayCoin = amount

	// Extract time
	var payTime time.Time
	if remarkArr[0] == "เงิน" {
		if len(remarkArr) > 3 {
			timeStr := remarkArr[3]
			if len(timeStr) == 11 {
				// Format: DDMMYY HH:MM
				year := time.Now().Year()
				month := timeStr[3:5]
				day := timeStr[0:2]
				hour := timeStr[6:8]
				minute := timeStr[9:11]
				payTime, err = time.Parse("2006-01-02 15:04", fmt.Sprintf("%d-%s-%s %s:%s", year, month, day, hour, minute))
			}
		}
	} else if remarkArr[0] == "Transfer" && len(remarkArr) > 11 {
		timeStr := remarkArr[11]
		if len(timeStr) == 11 {
			year := time.Now().Year()
			month := timeStr[3:5]
			day := timeStr[0:2]
			hour := timeStr[6:8]
			minute := timeStr[9:11]
			payTime, err = time.Parse("2006-01-02 15:04", fmt.Sprintf("%d-%s-%s %s:%s", year, month, day, hour, minute))
		}
	}

	if payTime.IsZero() {
		payTime, err = time.Parse("2006-01-02 15:04:05", msgTime)
		if err != nil {
			return data, fmt.Errorf("error parsing fallback time: %v", err)
		}
	}
	data.PayTime = payTime

	// Extract balance (last element)
	if len(remarkArr) > 0 {
		balanceStr := strings.ReplaceAll(remarkArr[len(remarkArr)-1], ",", "")
		balance, err := strconv.ParseFloat(balanceStr, 64)
		if err == nil {
			data.Balance = balance
		}
	}

	return data, nil
}

// SCB Notify handler
func handleSCBNotify(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}
	remarkArr := strings.Fields(ramk)

	if len(remarkArr) != 15 {
		return data, fmt.Errorf("unexpected message format")
	}

	// Extract amount
	amountStr := strings.ReplaceAll(remarkArr[1], ",", "")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return data, fmt.Errorf("error parsing amount: %v", err)
	}
	data.PayCoin = amount

	// Extract time (Thai month needs conversion)
	thaiMonth := remarkArr[11]
	month := changeThaiMonthToNumber(thaiMonth)
	if month < 1 || month > 12 {
		return data, fmt.Errorf("invalid month")
	}

	year := time.Now().Year()
	day := remarkArr[10]
	timePart := remarkArr[14]

	payTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%d-%02d-%s %s", year, month, day, timePart))
	if err != nil {
		return data, fmt.Errorf("error parsing time: %v", err)
	}
	data.PayTime = payTime

	return data, nil
}

// SCB Water handler
func handleSCBWater(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}

	// Parse the message time (format: DD/MM/YYYY HH:MM:SS)
	payTime, err := time.Parse("02/01/2006 15:04:05", msgTime)
	if err != nil {
		return data, fmt.Errorf("error parsing time: %v", err)
	}
	data.PayTime = payTime

	// Extract amount from ramk (assuming it's in JSON format)
	// This is a simplified version - actual implementation would parse JSON
	re := regexp.MustCompile(`"amount":\s*"([\d,]+)`)
	matches := re.FindStringSubmatch(ramk)
	if len(matches) < 2 {
		return data, fmt.Errorf("amount not found in message")
	}

	amountStr := strings.ReplaceAll(matches[1], ",", "")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return data, fmt.Errorf("error parsing amount: %v", err)
	}
	data.PayCoin = amount

	return data, nil
}

// KTB Bank handler
func handleKTB(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}

	if strings.Contains(ramk, "OTP") {
		return data, fmt.Errorf("OTP message, skipping")
	}

	remarkArr := strings.Fields(ramk)

	// Check if it's a deposit or withdrawal
	if strings.Contains(ramk, "Withdraw") {
		return data, fmt.Errorf("withdrawal transaction, skipping")
	}

	// Extract amount
	var amountStr string
	if strings.Contains(ramk, "Deposit") {
		// English format
		amountStr = strings.ReplaceAll(remarkArr[4], ",", "")
		amountStr = strings.ReplaceAll(amountStr, "THB", "")
	} else if strings.Contains(remarkArr[0], "@") {
		amountStr = strings.ReplaceAll(remarkArr[2], ",", "")
	} else {
		amountStr = strings.ReplaceAll(remarkArr[5], ",", "")
		amountStr = strings.ReplaceAll(amountStr, "บ", "")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return data, fmt.Errorf("error parsing amount: %v", err)
	}
	data.PayCoin = amount

	// Extract time
	var payTime time.Time
	if strings.Contains(remarkArr[0], "@") {
		timeParts := strings.Split(remarkArr[0], "@")
		if len(timeParts) == 2 {
			dateParts := strings.Split(timeParts[0], "-")
			if len(dateParts) == 2 {
				year := time.Now().Year()
				month := dateParts[1]
				day := dateParts[0]
				payTime, err = time.Parse("2006-01-02 15:04", fmt.Sprintf("%d-%s-%s %s", year, month, day, timeParts[1]))
			}
		}
	} else if len(remarkArr) > 3 && strings.Contains(remarkArr[3], "@") {
		timeParts := strings.Split(remarkArr[3], "@")
		if len(timeParts) == 2 {
			dateParts := strings.Split(timeParts[0], "-")
			if len(dateParts) == 2 {
				year := time.Now().Year()
				month := dateParts[1]
				day := dateParts[0]
				payTime, err = time.Parse("2006-01-02 15:04", fmt.Sprintf("%d-%s-%s %s", year, month, day, timeParts[1]))
			}
		}
	}

	if payTime.IsZero() {
		payTime, err = time.Parse("2006-01-02 15:04:05", msgTime)
		if err != nil {
			return data, fmt.Errorf("error parsing fallback time: %v", err)
		}
	}
	data.PayTime = payTime

	// Extract balance (last element)
	if len(remarkArr) > 0 {
		balanceStr := strings.ReplaceAll(remarkArr[len(remarkArr)-1], ",", "")
		balanceStr = strings.ReplaceAll(balanceStr, "บ", "")
		balanceStr = strings.ReplaceAll(balanceStr, "THB", "")
		balance, err := strconv.ParseFloat(balanceStr, 64)
		if err == nil {
			data.Balance = balance
		}
	}

	return data, nil
}

// KTB Line handler
func handleKTBLine(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}

	if strings.Contains(ramk, "OTP") {
		return data, fmt.Errorf("OTP message, skipping")
	}

	if !strings.Contains(ramk, "เงินเข้า") {
		return data, fmt.Errorf("not a deposit transaction")
	}

	remarkArr := strings.Fields(ramk)

	// Extract amount
	amountStr := strings.ReplaceAll(remarkArr[1], ",", "")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return data, fmt.Errorf("error parsing amount: %v", err)
	}
	data.PayCoin = amount

	// Extract time
	var payTime time.Time
	if len(remarkArr) > 9 {
		var timeParts []string
		if !strings.Contains(remarkArr[7], "/") {
			timeParts = strings.Split(remarkArr[6], "/")
		} else {
			timeParts = strings.Split(remarkArr[7], "/")
		}

		if len(timeParts) == 2 {
			year := time.Now().Year()
			month := timeParts[1]
			day := timeParts[0]
			payTime, err = time.Parse("2006-01-02 15:04", fmt.Sprintf("%d-%s-%s %s", year, month, day, remarkArr[7]))
		}
	}

	if payTime.IsZero() {
		payTime, err = time.Parse("2006-01-02 15:04:05", msgTime)
		if err != nil {
			return data, fmt.Errorf("error parsing fallback time: %v", err)
		}
	}
	data.PayTime = payTime

	// Extract balance
	if len(remarkArr) > 9 {
		balanceStr := strings.ReplaceAll(remarkArr[9], ",", "")
		balance, err := strconv.ParseFloat(balanceStr, 64)
		if err == nil {
			data.Balance = balance
		}
	}

	return data, nil
}

// KTB Notice handler
func handleKTBNotice(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}

	if strings.Contains(ramk, "OTP") {
		return data, fmt.Errorf("OTP message, skipping")
	}

	if !strings.Contains(ramk, "ได้รับ") {
		return data, fmt.Errorf("not a deposit transaction")
	}

	remarkArr := strings.Fields(ramk)

	// Extract amount
	amountStr := strings.ReplaceAll(remarkArr[1], ",", "")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return data, fmt.Errorf("error parsing amount: %v", err)
	}
	data.PayCoin = amount

	// Use message time as fallback
	payTime, err := time.Parse("2006-01-02 15:04:05", msgTime)
	if err != nil {
		return data, fmt.Errorf("error parsing time: %v", err)
	}
	data.PayTime = payTime

	return data, nil
}

// KTB Water handler
func handleKTBWater(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}

	// Parse the message time (format: DD-MM-YYYY HH:MM:SS)
	timeParts := strings.Fields(msgTime)
	if len(timeParts) < 2 {
		return data, fmt.Errorf("invalid time format")
	}

	dateParts := strings.Split(timeParts[0], "-")
	if len(dateParts) != 3 {
		return data, fmt.Errorf("invalid date format")
	}

	year := dateParts[2]
	month := dateParts[1]
	day := dateParts[0]
	timeStr := timeParts[1]

	payTime, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s-%s-%s %s", year, month, day, timeStr))
	if err != nil {
		return data, fmt.Errorf("error parsing time: %v", err)
	}
	data.PayTime = payTime

	// Extract amount from ramk (assuming it's in JSON format)
	re := regexp.MustCompile(`"amount":\s*"([\d,]+)`)
	matches := re.FindStringSubmatch(ramk)
	if len(matches) < 2 {
		return data, fmt.Errorf("amount not found in message")
	}

	amountStr := strings.ReplaceAll(matches[1], ",", "")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return data, fmt.Errorf("error parsing amount: %v", err)
	}
	data.PayCoin = amount

	return data, nil
}

// KBANK Notify handler
func handleKKRNotify(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}
	remarkArr := strings.Fields(ramk)

	if len(remarkArr) != 15 && len(remarkArr) != 14 {
		return data, fmt.Errorf("unexpected message format")
	}

	// Extract amount
	amountStr := strings.ReplaceAll(remarkArr[4], ",", "")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return data, fmt.Errorf("error parsing amount: %v", err)
	}
	data.PayCoin = amount

	// Extract time (Thai month needs conversion)
	thaiMonth := remarkArr[9]
	month := changeThaiMonthToNumber(thaiMonth)
	if month < 1 || month > 12 {
		return data, fmt.Errorf("invalid month")
	}

	year := 2500 + atoi(remarkArr[10]) - 543
	day := remarkArr[8]
	timePart := remarkArr[12]
	if len(remarkArr) == 15 {
		timePart = remarkArr[14]
	}

	payTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%d-%02d-%s %s", year, month, day, timePart))
	if err != nil {
		return data, fmt.Errorf("error parsing time: %v", err)
	}
	data.PayTime = payTime

	return data, nil
}

// KBANK Read handler
func handleKKRRead(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}
	ramk = strings.ReplaceAll(ramk, "KBank:", "")
	remarkArr := strings.Fields(ramk)

	// Extract amount
	var amountStr string
	var isWithdrawal bool

	switch len(remarkArr) {
	case 9, 14:
		amountStr = strings.ReplaceAll(remarkArr[5], ",", "")
		for _, word := range remarkArr {
			if strings.Contains(word, "เงินออก") || strings.Contains(word, "หักบช") {
				isWithdrawal = true
				break
			}
		}
	case 10:
		amountStr = strings.ReplaceAll(remarkArr[6], ",", "")
		if !isNumeric(amountStr) {
			amountStr = strings.ReplaceAll(remarkArr[5], ",", "")
		}
		for _, word := range remarkArr {
			if strings.Contains(word, "เงินออก") || strings.Contains(word, "หักบช") {
				isWithdrawal = true
				break
			}
		}
	case 4, 5:
		amountStr = strings.ReplaceAll(remarkArr[3], ",", "")
		if strings.Contains(remarkArr[3], "เงินออก") || strings.Contains(remarkArr[3], "หักบช") {
			isWithdrawal = true
		}
	case 6:
		amountStr = strings.ReplaceAll(remarkArr[4], ",", "")
		if strings.Contains(remarkArr[3], "เงินออก") || strings.Contains(remarkArr[3], "หักบช") {
			isWithdrawal = true
		}
	default:
		return data, fmt.Errorf("unexpected message format")
	}

	amountStr = strings.ReplaceAll(amountStr, "เงินออก", "")
	amountStr = strings.ReplaceAll(amountStr, "เงินเข้า", "")
	amountStr = strings.ReplaceAll(amountStr, "บ", "")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return data, fmt.Errorf("error parsing amount: %v", err)
	}

	if isWithdrawal {
		data.PayCoin = -amount
	} else {
		data.PayCoin = amount
	}

	// Extract time
	var payTime time.Time
	if len(remarkArr) >= 9 || len(remarkArr) == 10 || len(remarkArr) == 14 {
		dateParts := strings.Split(remarkArr[0], "/")
		if len(dateParts) >= 2 {
			year := time.Now().Year()
			month := dateParts[1]
			day := dateParts[0]
			payTime, err = time.Parse("2006-01-02 15:04", fmt.Sprintf("%d-%s-%s %s", year, month, day, remarkArr[1]))
		}
	} else if len(remarkArr) >= 4 && len(remarkArr) <= 6 {
		dateParts := strings.Split(remarkArr[0], "/")
		if len(dateParts) >= 3 {
			year := 2500 + atoi(dateParts[2]) - 543
			month := dateParts[1]
			day := dateParts[0]
			payTime, err = time.Parse("2006-01-02 15:04", fmt.Sprintf("%d-%s-%s %s", year, month, day, remarkArr[1]))
		}
	}

	if payTime.IsZero() {
		payTime, err = time.Parse("2006-01-02 15:04:05", msgTime)
		if err != nil {
			return data, fmt.Errorf("error parsing fallback time: %v", err)
		}
	}
	data.PayTime = payTime

	// Extract balance (second last element)
	if len(remarkArr) >= 2 {
		balanceStr := strings.ReplaceAll(remarkArr[len(remarkArr)-2], ",", "")
		balance, err := strconv.ParseFloat(balanceStr, 64)
		if err == nil {
			data.Balance = balance
		}
	}

	return data, nil
}

// KBANK Water handler
func handleKKRWater(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}

	// Parse the message time (format: DD Month(Thai) YYYY HH:MM)
	timeParts := strings.Fields(msgTime)
	if len(timeParts) < 4 {
		return data, fmt.Errorf("invalid time format")
	}

	thaiMonths := []string{"", "ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.",
		"ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."}

	day := timeParts[0]
	thaiMonth := timeParts[1]
	year := 2500 + atoi(timeParts[2]) - 543
	timeStr := timeParts[3]

	var month int
	for i, m := range thaiMonths {
		if m == thaiMonth {
			month = i
			break
		}
	}
	if month == 0 {
		return data, fmt.Errorf("invalid Thai month")
	}

	payTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%d-%02d-%s %s", year, month, day, timeStr))
	if err != nil {
		return data, fmt.Errorf("error parsing time: %v", err)
	}
	data.PayTime = payTime

	// Extract amount from ramk (assuming it's in JSON format)
	re := regexp.MustCompile(`"amount":\s*"([\d,]+)`)
	matches := re.FindStringSubmatch(ramk)
	if len(matches) < 2 {
		return data, fmt.Errorf("amount not found in message")
	}

	amountStr := strings.ReplaceAll(matches[1], ",", "")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return data, fmt.Errorf("error parsing amount: %v", err)
	}
	data.PayCoin = amount

	return data, nil
}

// BBL Water handler
func handleBBLWater(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}

	// Parse the message time (format: Month(English) DD YYYY HH:MM:SS)
	timeParts := strings.Fields(msgTime)
	if len(timeParts) < 4 {
		return data, fmt.Errorf("invalid time format")
	}

	englishMonths := []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sept", "Oct", "Nov", "Dec"}

	monthStr := timeParts[1]
	day := timeParts[2]
	year := timeParts[3]
	timeStr := timeParts[4]

	var month int
	for i, m := range englishMonths {
		if m == monthStr {
			month = i
			break
		}
	}
	if month == 0 {
		return data, fmt.Errorf("invalid English month")
	}

	payTime, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s-%02d-%s %s", year, month, day, timeStr))
	if err != nil {
		return data, fmt.Errorf("error parsing time: %v", err)
	}
	data.PayTime = payTime

	// Extract amount from ramk (assuming it's in JSON format)
	re := regexp.MustCompile(`"amount":\s*"([\d,]+)`)
	matches := re.FindStringSubmatch(ramk)
	if len(matches) < 2 {
		return data, fmt.Errorf("amount not found in message")
	}

	amountStr := strings.ReplaceAll(matches[1], ",", "")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return data, fmt.Errorf("error parsing amount: %v", err)
	}
	data.PayCoin = amount

	return data, nil
}

// BBL handler
func handleBBL(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}

	// Check if it's a withdrawal
	if strings.Contains(ramk, "ถอน/โอน") || strings.Contains(ramk, "Money Withdrawal") {
		return data, fmt.Errorf("withdrawal transaction, skipping")
	}

	// Check if it's a deposit
	if !strings.Contains(ramk, "PromptPay") && !strings.Contains(ramk, "Deposit") {
		return data, fmt.Errorf("not a deposit transaction")
	}

	remarkArr := strings.Fields(ramk)

	// Extract amount
	var amountStr string
	if strings.Contains(ramk, "PromptPay") {
		amountStr = strings.ReplaceAll(remarkArr[8], ",", "")
	} else {
		amountStr = strings.ReplaceAll(remarkArr[2], ",", "")
		if amountStr == "0" {
			// Try alternative pattern
			re := regexp.MustCompile(`MB\s+([\d,]+)\s+บ`)
			matches := re.FindStringSubmatch(ramk)
			if len(matches) >= 2 {
				amountStr = strings.ReplaceAll(matches[1], ",", "")
			}
		}
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return data, fmt.Errorf("error parsing amount: %v", err)
	}
	data.PayCoin = amount

	// Extract time (use message time as fallback)
	payTime, err := time.Parse("2006-01-02 15:04:05", msgTime)
	if err != nil {
		return data, fmt.Errorf("error parsing time: %v", err)
	}
	data.PayTime = payTime

	// Extract balance (last element)
	if len(remarkArr) > 0 {
		balanceStr := remarkArr[len(remarkArr)-1]
		// Extract numeric part from balance string
		re := regexp.MustCompile(`(\d+\.\d+)`)
		matches := re.FindStringSubmatch(balanceStr)
		if len(matches) >= 2 {
			balance, err := strconv.ParseFloat(matches[1], 64)
			if err == nil {
				data.Balance = balance
			}
		}
	}

	return data, nil
}

// BAAC handler
func handleBAAC(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}

	if strings.Contains(ramk, "OTP") {
		return data, fmt.Errorf("OTP message, skipping")
	}

	if !strings.Contains(ramk, "รับโอนพร้อมเพย์") {
		return data, fmt.Errorf("not a deposit transaction")
	}

	remarkArr := strings.Fields(ramk)

	// Extract amount
	amountStr := strings.ReplaceAll(remarkArr[6], ",", "")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return data, fmt.Errorf("error parsing amount: %v", err)
	}
	data.PayCoin = amount

	// Extract time
	var payTime time.Time
	if len(remarkArr) > 0 {
		dateParts := strings.Split(remarkArr[0], "/")
		if len(dateParts) >= 3 {
			year := time.Now().Year()
			month := dateParts[1]
			day := dateParts[0]
			payTime, err = time.Parse("2006-01-02 15:04", fmt.Sprintf("%d-%s-%s %s", year, month, day, remarkArr[1]))
		}
	}

	if payTime.IsZero() {
		payTime, err = time.Parse("2006-01-02 15:04:05", msgTime)
		if err != nil {
			return data, fmt.Errorf("error parsing fallback time: %v", err)
		}
	}
	data.PayTime = payTime

	// Extract balance
	if len(remarkArr) > 9 {
		balanceStr := strings.ReplaceAll(remarkArr[9], ",", "")
		balance, err := strconv.ParseFloat(balanceStr, 64)
		if err == nil {
			data.Balance = balance
		}
	}

	return data, nil
}

// TTB handler
func handleTTB(ramk string, msgTime string) (TransactionData, error) {
	data := TransactionData{}

	// Check if it's a transfer (not deposit)
	if strings.Contains(ramk, "โอนเงิน") || strings.Contains(ramk, "transferred") {
		return data, fmt.Errorf("transfer transaction, skipping")
	}

	// Extract amount
	re := regexp.MustCompile(`(?:\d{1,3}(?:,\d{3})*)\.\d{2}(?=บ\.)`)
	matches := re.FindStringSubmatch(ramk)
	if len(matches) == 0 {
		// Try English format
		remarkArr := strings.Fields(ramk)
		if len(remarkArr) >= 3 {
			amountStr := strings.ReplaceAll(remarkArr[2], ",", "")
			amount, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				return data, fmt.Errorf("error parsing amount: %v", err)
			}
			data.PayCoin = amount
		} else {
			return data, fmt.Errorf("amount not found in message")
		}
	} else {
		amountStr := strings.ReplaceAll(matches[0], ",", "")
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return data, fmt.Errorf("error parsing amount: %v", err)
		}
		data.PayCoin = amount
	}

	// Extract time
	re = regexp.MustCompile(`\d{2}/\d{2}/\d{2}@\d{2}:\d{2}`)
	matches = re.FindStringSubmatch(ramk)
	if len(matches) == 0 {
		return data, fmt.Errorf("time not found in message")
	}

	dateStr := strings.Split(matches[0], "@")
	if len(dateStr) != 2 {
		return data, fmt.Errorf("invalid time format")
	}

	dateParts := strings.Split(dateStr[0], "/")
	if len(dateParts) != 3 {
		return data, fmt.Errorf("invalid date format")
	}

	year := "20" + dateParts[2] // Assuming 21st century
	month := dateParts[1]
	day := dateParts[0]
	timePart := dateStr[1]

	payTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s-%s-%s %s", year, month, day, timePart))
	if err != nil {
		return data, fmt.Errorf("error parsing time: %v", err)
	}
	data.PayTime = payTime

	// Extract balance
	re = regexp.MustCompile(`เหลือ([\d,]+\.\d{2})บ`)
	matches = re.FindStringSubmatch(ramk)
	if len(matches) >= 2 {
		balanceStr := strings.ReplaceAll(matches[1], ",", "")
		balance, err := strconv.ParseFloat(balanceStr, 64)
		if err == nil {
			data.Balance = balance
		}
	}

	return data, nil
}
