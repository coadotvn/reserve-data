package data

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type ReserveData struct {
	storage Storage
	fetcher Fetcher
}

func (self ReserveData) CurrentPriceVersion() (common.Version, error) {
	return self.storage.CurrentPriceVersion()
}

func (self ReserveData) GetAllPrices() (common.AllPriceResponse, error) {
	timestamp := common.GetTimestamp()
	version, err := self.storage.CurrentPriceVersion()
	if err != nil {
		return common.AllPriceResponse{}, err
	} else {
		result := common.AllPriceResponse{}
		data, err := self.storage.GetAllPrices(version)
		returnTime := common.GetTimestamp()
		result.Version = version
		result.Timestamp = timestamp
		result.ReturnTime = returnTime
		result.Data = data
		return result, err
	}
}

func (self ReserveData) GetOnePrice(pairID common.TokenPairID) (common.OnePriceResponse, error) {
	timestamp := common.GetTimestamp()
	version, err := self.storage.CurrentPriceVersion()
	if err != nil {
		return common.OnePriceResponse{}, err
	} else {
		result := common.OnePriceResponse{}
		data, err := self.storage.GetOnePrice(pairID, version)
		returnTime := common.GetTimestamp()
		result.Version = version
		result.Timestamp = timestamp
		result.ReturnTime = returnTime
		result.Data = data
		return result, err
	}
}

func (self ReserveData) CurrentEBalanceVersion() (common.Version, error) {
	return self.storage.CurrentEBalanceVersion()
}

func (self ReserveData) GetAllEBalances() (common.AllEBalanceResponse, error) {
	timestamp := common.GetTimestamp()
	version, err := self.storage.CurrentEBalanceVersion()
	if err != nil {
		return common.AllEBalanceResponse{}, err
	} else {
		result := common.AllEBalanceResponse{}
		data, err := self.storage.GetAllEBalances(version)
		returnTime := common.GetTimestamp()
		result.Version = version
		result.Timestamp = timestamp
		result.ReturnTime = returnTime
		result.Data = data
		return result, err
	}
}

func (self ReserveData) CurrentBalanceVersion() (common.Version, error) {
	return self.storage.CurrentBalanceVersion()
}

func (self ReserveData) GetAllBalances() (common.AllBalanceResponse, error) {
	timestamp := common.GetTimestamp()
	version, err := self.storage.CurrentBalanceVersion()
	if err != nil {
		return common.AllBalanceResponse{}, err
	} else {
		result := common.AllBalanceResponse{}
		balances, err := self.storage.GetAllBalances(version)
		returnTime := common.GetTimestamp()
		result.Version = version
		result.Timestamp = timestamp
		result.ReturnTime = returnTime
		data := map[string]common.BalanceResponse{}
		for tokenID, balance := range balances {
			data[tokenID] = balance.ToBalanceResponse(
				common.MustGetToken(tokenID).Decimal,
			)
		}
		result.Data = data
		return result, err
	}
}

func (self ReserveData) CurrentRateVersion() (common.Version, error) {
	return self.storage.CurrentRateVersion()
}

func (self ReserveData) GetAllRates() (common.AllRateResponse, error) {
	timestamp := common.GetTimestamp()
	version, err := self.storage.CurrentRateVersion()
	if err != nil {
		return common.AllRateResponse{}, err
	} else {
		result := common.AllRateResponse{}
		rates, err := self.storage.GetAllRates(version)
		returnTime := common.GetTimestamp()
		result.Version = version
		result.Timestamp = timestamp
		result.ReturnTime = returnTime
		data := map[common.TokenPairID]common.RateResponse{}
		for tokenPairID, rate := range rates.Data {
			data[tokenPairID] = common.RateResponse{
				Valid:       rates.Valid,
				Error:       rates.Error,
				Timestamp:   rates.Timestamp,
				ReturnTime:  rates.ReturnTime,
				Rate:        common.BigToFloat(rate.Rate, 18),
				ExpiryBlock: rate.ExpiryBlock.Int64(),
				Balance:     common.BigToFloat(rate.Balance, 18),
			}
		}
		result.Data = data
		return result, err
	}
}

func (self ReserveData) Run() error {
	return self.fetcher.Run()
}

func NewReserveData(storage Storage, fetcher Fetcher) *ReserveData {
	return &ReserveData{storage, fetcher}
}
