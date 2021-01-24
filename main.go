package main
import (
	"io"
	"os"
	"github.com/ychi/cbpgo/conf"
	"github.com/ychi/cbpgo/restclient"
)


// pure get
func view(rest restclient.RestClient) {
	res, _ := rest.Request("GET", "/coinbase-accounts", nil, nil, nil)
	io.Copy(os.Stdout, res.Body)
}

// get + query params(s)
func orders(rest restclient.RestClient) {
	q := make(map[string]string)
	q["status"] = "pending"
	res, _ := rest.Request("GET", "/orders", q, nil, nil)
	io.Copy(os.Stdout, res.Body)	
}

// post + body
func withdraw(rest restclient.RestClient) {
	withdrawal := WithdrawalCrypto{
		Amount: "100.00",
		Currency: "BTC",
		CryptoAddress: "0x5ad5769cd04681FeD900BCE3DDc877B50E83d469",	
	}
	res, _ := rest.Request("POST", "/withdrawals/crypto", nil, withdrawal, nil)
	io.Copy(os.Stdout, res.Body)

}


type WithdrawalCrypto struct {
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	CryptoAddress string `json:"crypto_address"`
}



func main() {
	rest := restclient.
	NewRestClient().
	SetBaseURL(conf.API_URL).
	SetApiKey(conf.ACCESS_KEY).
	SetApiSecret(conf.API_SECRET).
	SetApiPassphrase(conf.PASS_PHRASE)
	
	view(rest)
	orders(rest)
	withdraw(rest)
}
