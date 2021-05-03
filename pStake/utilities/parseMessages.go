package utilities

import "errors"

type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type SendCoin struct {
	//Type        string   `json:"@type"`
	FromAddress string   `json:"from_address"`
	ToAddress   string   `json:"to_address"`
	Amounts     []Amount `json:"amounts"`
}

const SendCoinMsgType = "/cosmos.bank.v1beta1.MsgSend"

func ParseAmount(amountInterface interface{}) Amount {
	map1 := amountInterface.(map[string]interface{})
	amount := map1["amount"].(string)
	denom := map1["denom"].(string)
	amt := Amount{
		Denom:  denom,
		Amount: amount,
	}
	return amt
}

func ParseSendCoin(msg interface{}) (SendCoin, error) {
	map1 := msg.(map[string]interface{})
	if map1["@type"].(string) != SendCoinMsgType {
		return SendCoin{}, errors.New("Incorrect msg type")
	}
	fromAddress := map1["from_address"].(string)
	toAddress := map1["from_address"].(string)
	amountsInterface := map1["amount"].([]interface{})

	amounts := make([]Amount, len(amountsInterface))

	for i, amountInterface := range amountsInterface {
		amounts[i] = ParseAmount(amountInterface)
	}

	sendCoin := SendCoin{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amounts:     amounts,
	}

	return sendCoin, nil
}

func GetMsgType(msg interface{}) string {
	map1 := msg.(map[string]interface{})
	return map1["@type"].(string)
}
