package exchange

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Liqui struct {
	signer    Signer
	endpoint  LiquiEndpoint
	pairs     []common.TokenPair
	addresses map[string]ethereum.Address
}

func (self *Liqui) MarshalText() (text []byte, err error) {
	return []byte(self.ID()), nil
}

func (self *Liqui) Address(token common.Token) (ethereum.Address, bool) {
	addr, supported := self.addresses[token.ID]
	return addr, supported
}

func (self *Liqui) ID() common.ExchangeID {
	return common.ExchangeID("liqui")
}

func (self *Liqui) TokenPairs() []common.TokenPair {
	return self.pairs
}

func (self *Liqui) Name() string {
	return "liqui"
}

// map of token pair to map of asks/bids to array of [rate, amount]
type liqresp map[string]map[string][][]float64

type liqinfo struct {
	Success int `json:"success"`
	Return  map[string]map[string]float64
	Error   string `json:"error"`
}

func (self *Liqui) Trade(tradeType string, base common.Token, quote common.Token, rate float64, amount float64) (done float64, remaining float64, finished bool, err error) {
	return self.endpoint.Trade(
		self.signer.GetLiquiKey(),
		tradeType, base, quote, rate, amount, self.signer)
}

func (self *Liqui) Withdraw(token common.Token, amount *big.Int, address ethereum.Address) error {
	return self.endpoint.Withdraw(
		self.signer.GetLiquiKey(),
		token, amount, address, self.signer)
}

func (self *Liqui) FetchEBalanceData() (common.EBalanceEntry, error) {
	result := common.EBalanceEntry{}
	result.Timestamp = common.GetTimestamp()
	result.Valid = true
	resp_data, err := self.endpoint.GetInfo(
		self.signer.GetLiquiKey(),
		self.signer,
	)
	result.ReturnTime = common.GetTimestamp()
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	} else {
		if resp_data.Success == 1 {
			balances := resp_data.Return["funds"]
			result.Balance = map[string]float64{}
			for tokenID, _ := range common.SupportedTokens {
				result.Balance[tokenID] = balances[strings.ToLower(tokenID)]
			}
		} else {
			result.Valid = false
			result.Error = resp_data.Error
		}
	}
	return result, nil
}

func (self *Liqui) FetchPriceData() (map[common.TokenPairID]common.ExchangePrice, error) {
	result := map[common.TokenPairID]common.ExchangePrice{}
	pairs_str := []string{}
	for _, pair := range self.pairs {
		pairs_str = append(pairs_str, fmt.Sprintf("%s_%s", pair.Base.ID, pair.Quote.ID))
	}
	timestamp := common.GetTimestamp()
	resp_data, err := self.endpoint.Depth(
		strings.ToLower(strings.Join(pairs_str, "-")),
	)
	returnTime := common.GetTimestamp()
	if err != nil {
		for _, pair := range self.pairs {
			one_pair_result := common.ExchangePrice{}
			one_pair_result.Timestamp = timestamp
			one_pair_result.ReturnTime = returnTime
			one_pair_result.Valid = false
			one_pair_result.Error = err.Error()
			result[pair.PairID()] = one_pair_result
		}
	} else {
		for _, pair := range self.pairs {
			one_pair_result := common.ExchangePrice{}
			one_pair_result.Timestamp = timestamp
			one_pair_result.ReturnTime = returnTime
			one_pair_result.Valid = true
			one_data := resp_data[fmt.Sprintf(
				"%s_%s",
				strings.ToLower(pair.Base.ID),
				strings.ToLower(pair.Quote.ID),
			)]
			for _, buy := range one_data["bids"] {
				one_pair_result.BuyPrices = append(
					one_pair_result.BuyPrices,
					common.PriceEntry{
						buy[1],
						buy[0],
					},
				)
			}
			for _, sell := range one_data["asks"] {
				one_pair_result.SellPrices = append(
					one_pair_result.SellPrices,
					common.PriceEntry{
						sell[1],
						sell[0],
					},
				)
			}
			result[pair.PairID()] = one_pair_result
		}
	}
	return result, err
}

func NewLiqui(signer Signer, endpoint LiquiEndpoint) *Liqui {
	return &Liqui{
		signer,
		endpoint,
		[]common.TokenPair{
			common.MustCreateTokenPair("OMG", "ETH"),
			common.MustCreateTokenPair("DGD", "ETH"),
			common.MustCreateTokenPair("CVC", "ETH"),
			common.MustCreateTokenPair("MCO", "ETH"),
			common.MustCreateTokenPair("GNT", "ETH"),
			common.MustCreateTokenPair("ADX", "ETH"),
			common.MustCreateTokenPair("EOS", "ETH"),
			common.MustCreateTokenPair("PAY", "ETH"),
			common.MustCreateTokenPair("BAT", "ETH"),
			common.MustCreateTokenPair("KNC", "ETH"),
		},
		map[string]ethereum.Address{
			"ETH": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"OMG": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"DGD": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"CVC": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"MCO": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"GNT": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"ADX": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"EOS": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"PAY": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"BAT": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
			"KNC": ethereum.HexToAddress("0xce656971fe4fc43a0211b792d380900761b7862c"),
		},
	}
}
