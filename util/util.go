package util

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"net/url"
)

func GetBigIntForDecimal(decimal int) *big.Int {
	floatDecimal := big.NewFloat(math.Pow10(decimal))
	bigIntDecimal := new(big.Int)
	floatDecimal.Int(bigIntDecimal)

	return bigIntDecimal
}

func QuoBigInt(a *big.Int, b *big.Int) *big.Float {
	fl := new(big.Float).SetInt(a)
	fl.Quo(fl, new(big.Float).SetInt(b))
	return fl
}

func CalcActualOutAmount(amount *big.Int, ratio *big.Float, fixedFee *big.Int) *big.Int {
	res := new(big.Float).SetInt(amount)
	res.Mul(res, ratio)

	amountInt := new(big.Int)
	res.Int(amountInt)

	amountInt.Sub(amountInt, fixedFee)
	return amountInt
}

func SendTelegramMessage(botId string, chatId string, msg string) {
	if botId == "" || chatId == "" || msg == "" {
		return
	}

	endPoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botId)
	formData := url.Values{
		"chat_id":    {chatId},
		"parse_mode": {"html"},
		"text":       {msg},
	}
	Logger.Infof("send tg message, bot_id=%s, chat_id=%s, msg=%s", botId, chatId, msg)
	res, err := http.PostForm(endPoint, formData)
	if err != nil {
		Logger.Errorf("send telegram message error, bot_id=%s, chat_id=%s, msg=%s, err=%s", botId, chatId, msg, err.Error())
		return
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		Logger.Errorf("read http response error, err=%s", err.Error())
		return
	}
	Logger.Infof("tg response: %s", string(bodyBytes))
}
