package wash

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	MsgTemplate interface {
		GetPayTimeRegex() string
		GetPayCoinRegex() string
	}
	WashResult struct {
		Msg      string
		PayTime  *time.Time
		PayCoin  float64
		Error    []error
		ErrorMsg string
	}
)

func (wr *WashResult) AddError(err error) {
	if err != nil {
		wr.Error = append(wr.Error, err)
	}
}

func (wr *WashResult) HasError() bool {
	return len(wr.Error) > 0
}
func (wr *WashResult) GetErrorMsg() string {
	if wr.ErrorMsg != "" {
		return wr.ErrorMsg
	}
	if len(wr.Error) > 0 {
		var sb strings.Builder
		for _, err := range wr.Error {
			sb.WriteString(err.Error() + "; ")
		}
		wr.ErrorMsg = sb.String()
	}
	return wr.ErrorMsg
}

func WashMsg(msg string, template MsgTemplate) (res *WashResult) {
	res = &WashResult{
		Msg: msg,
	}
	if msg == "" {
		res.AddError(errors.New("msg is empty"))
		return
	}
	if template == nil {
		res.AddError(errors.New("template is nil"))
		return
	}
	washpayTime(res, msg, template.GetPayTimeRegex())
	washpayCoin(res, msg, template.GetPayCoinRegex())

	return
}

func washpayTime(res *WashResult, msg string, timeRegex string) {

	if res == nil {
		return
	}
	timeMatch := extractPayTime(msg, timeRegex)
	payTime, err := parseTime(timeMatch)
	if err != nil {
		res.AddError(err)
		return
	}
	res.PayTime = payTime
}

func washpayCoin(res *WashResult, msg string, coinRegex string) {
	if res == nil {
		return
	}
	coinMatche := extractPayCoin(msg, coinRegex)
	payCoin, err := parseCoin(coinMatche)
	if err != nil {
		res.AddError(err)
		return
	}
	res.PayCoin = payCoin
}

func extractPayTime(msg string, payTimeRegex string) (payTime string) {
	timeRe := regexp.MustCompile(payTimeRegex)
	if timeRe.MatchString(msg) {
		payTime = timeRe.FindStringSubmatch(msg)[1]
	}
	return
}

func extractPayCoin(msg string, payCoinRegex string) (payCoin string) {
	cleanMsg := strings.ReplaceAll(msg, ",", "")
	coinRe := regexp.MustCompile(payCoinRegex)
	if coinRe.MatchString(cleanMsg) {
		payCoin = coinRe.FindStringSubmatch(cleanMsg)[1]
	}
	return
}

func parseTime(timeStr string) (*time.Time, error) {
	payTime, err := time.Parse(time.Layout, timeStr)
	if err != nil {
		return nil, err
	}
	return &payTime, nil
}
func parseCoin(coinStr string) (float64, error) {
	coin, err := strconv.ParseFloat(coinStr, 64)
	if err != nil {
		return 0, err
	}
	return coin, nil
}
