package market

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ychi/coinbase-pro-go/common"
	"github.com/ychi/coinbase-pro-go/restclient"
)

type Product struct {
	Id				string 	`json:"id"`
	BaseCurrency	string 	`json:"base_currency"`
	QuoteCurrency	string 	`json:"quote_currency"`
	BaseMinSize		string	`json:"base_min_size"`
	BaseMaxSize		string	`json:"base_max_size"`
	QuoteIncrement	string	`json:"quote_increment"`
	BaseIncrement	string	`json:"base_increment"`
	DisplayName		string	`json:"display_name"`
	MinMarketFunds	string	`json:"min_market_funds"`
	MaxMarketFunds	string	`json:"max_market_funds"`
	MarginEnabled	bool	`json:"margin_enabled"`
	PostOnly		bool	`json:"post_only"`
	LimitOnly		bool	`json:"limit_only"`
	CancelOnly		bool	`json:"cancel_only"`
	TradingDisabled	bool	`json:"trading_disabled"`
	Status 			string	`json:"status"`
	StatusMessage	string	`json:"status_message"`
}

type ProductService interface {
	List()
	Find()
	GetBook()
	GetTicker()
	GetStats()
	// Trades, Historical Rates
}

type productService struct {
	rc restclient.RestClient
}


func NewProductService(rc restclient.RestClient) (*productService) {
	return &productService{
		rc: rc,
	}
}

func (svc *productService) List() ([]Product, error){
	var ret []Product
	_, err := svc.rc.Request("GET", "/products", nil, nil, &ret)
	return ret, err
}

func (svc *productService) Find(id string) (Product, error){
	var ret Product
	path := fmt.Sprintf("%s%s", "/products/", id)
	_, err := svc.rc.Request("GET", path, nil, nil, &ret)
	return ret, err
}


type BookEntry struct {
	Price		string
	Size 		string
	OrderId		string
	numOrders	int
}

func (e *BookEntry) UnmarshalJSON(b []byte) error {
	var els []interface{}
	if err := json.Unmarshal(b, &els); err != nil {
		return err
	}
	
	priceStr, ok := els[0].(string)
	if !ok {
		return errors.New("expect bookEntry[0] (price) to be string")
	}

	sizeStr, ok := els[1].(string)
	if !ok {
		return errors.New("expect bookEntry[1] (size) to be string")
	}

	*e = BookEntry{
		Price: priceStr,
		Size: sizeStr,
	}

	if numOrdersInt, ok := els[2].(float64); ok {
		e.numOrders = int(numOrdersInt)
	} else if orderId, ok := els[2].(string); ok {
		e.OrderId = orderId
	} else {
		return errors.New("expect bookEntry[2] to be float64 (numOrders) or string (orderId)")
	}

	return nil
}

type Book struct {
	Sequence 	int			`json:"sequence"`
	Bids		[]BookEntry	`json:"bids"`
	Asks		[]BookEntry `json:"asks"`
}

func (svc *productService) GetBook(id string, level int)(Book, error){
	var ret Book
	path := fmt.Sprintf("/products/%s/book?level=%d", id, level)
	_, err := svc.rc.Request("GET", path, nil, nil, &ret)
	return ret, err
}

type Ticker struct {
	TradeId int					`json:"trade_id,number"`
	Price 	string				`json:"price"`
	Size	string				`json:"size"`
	Bid		string				`json:"bid"`
	Ask		string				`json:"ask"`
	Volume	common.DecimalStr	`json:"volume"`
	Time 	common.Time			`json:"time,string"`

}

func (svc *productService) GetTicker(id string) (Ticker, error) {
	var ret Ticker

	path := fmt.Sprintf("/products/%s/ticker", id)
	_, err := svc.rc.Request("GET", path, nil, nil, &ret)
	return ret, err
}