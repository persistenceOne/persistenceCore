package responses

type Response struct {
	LastBlockHeight string `json:"last_block_height"`
}

type ABCIResult struct {
	Response Response `json:"response"`
}

type ABCIResponse struct {
	Result ABCIResult `json:"result"`
}
