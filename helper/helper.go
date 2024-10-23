package helper

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GGGMGN struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		History []struct {
			Maker                string   `json:"maker"`
			BaseAmount           float64  `json:"base_amount"`
			QuoteAmount          float64  `json:"quote_amount"`
			AmountUsd            float64  `json:"amount_usd"`
			Timestamp            int      `json:"timestamp"`
			Event                string   `json:"event"`
			TxHash               string   `json:"tx_hash"`
			PriceUsd             string   `json:"price_usd"`
			MakerTags            []any    `json:"maker_tags"`
			MakerTwitterUsername string   `json:"maker_twitter_username"`
			MakerTwitterName     string   `json:"maker_twitter_name"`
			MakerName            string   `json:"maker_name"`
			MakerAvatar          string   `json:"maker_avatar"`
			MakerEns             any      `json:"maker_ens"`
			MakerTokenTags       []string `json:"maker_token_tags"`
			TokenAddress         string   `json:"token_address"`
			QuoteAddress         string   `json:"quote_address"`
			TotalTrade           int      `json:"total_trade"`
			ID                   string   `json:"id"`
			IsFollowing          int      `json:"is_following"`
			IsOpenOrClose        int      `json:"is_open_or_close"`
			BuyCostUsd           float64  `json:"buy_cost_usd"`
			Balance              string   `json:"balance"`
			Cost                 float64  `json:"cost"`
			HistoryBoughtAmount  float64  `json:"history_bought_amount"`
			HistorySoldIncome    float64  `json:"history_sold_income"`
			HistorySoldAmount    float64  `json:"history_sold_amount"`
			UnrealizedProfit     int      `json:"unrealized_profit"`
			RealizedProfit       float64  `json:"realized_profit"`
		} `json:"history"`
	} `json:"data"`
}
type Result struct {
	UserAddress  string
	TokenAddress string
	DateTime     string
	Amount       float64
	TxHash       string
}

type PFResult struct {
	Signature    string `json:"signature"`
	Mint         string `json:"mint"`
	SolAmount    int    `json:"sol_amount"`
	TokenAmount  int64  `json:"token_amount"`
	IsBuy        bool   `json:"is_buy"`
	User         string `json:"user"`
	Timestamp    int    `json:"timestamp"`
	TxIndex      int    `json:"tx_index"`
	Username     string `json:"username"`
	ProfileImage string `json:"profile_image"`
	Slot         int    `json:"slot"`
}

const LamportsPerSol = 1000000000
const LimitPerRequest = 200

func FetchGMGN(token, address string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://gmgn.ai/defi/quotation/v1/trades/sol/"+token+"?limit=100&maker="+address, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("dnt", "1")
	req.Header.Set("referer", "https://gmgn.ai/sol/token/j4B2oyth_uSjv4A4CcXjjwSetmEsZEzxCPgp1rYtDoCVokwrpump")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="129", "Not=A?Brand";v="8", "Chromium";v="129"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	req.Header.Set("x-kl-saas-ajax-request", "Ajax_Request")
	req.Header.Set("cookie", "_ga=GA1.1.851833522.1729001920; _ga_0XM0LYXGC8=deleted; __cf_bm=vWBsSS9af4QkpNdM65Sr4qFCjiDfaVNroHb9lc9a0cE-1729479284-1.0.1.1-4HHqe.i45m3CUem.mRChkTy_gNqy34GD8ZWPaoyvGgdavuXgUzz6_m9B4t_7yv9br689lefjqA0ryyxyahIUAw; cf_clearance=sLc1P98PEAX8zumPkE3RFWmFoOyimfmaBFSbDZKMGxc-1729479287-1.2.1.1-A7ZKUqgd1qVbOxp14XcAPjtOBDUz1yxUXOTHL0tEnD_c2noDFnLDG0hHXtZ6c2imuh60rb7O6rkThSxu8NroCISEBkSOzyj7HsV28OwIPN3iYVYb8Rhv6A5V1vDiHlUfawa6exk84Ult5aXe.V2EQnVC6ICrps1_6bmm.9P86aBFNTRCoCZHqqX9lRSKfrUR1Ow3w89ZRxl8PoKjAkdxpH0zDD4TVV1Z2Rq1SwfcKQNx5ipjQH0FW7vWPgH04D6iCimVpmJ0vxiRT0vvhx5s9czooZWki.wSijiAaTFn_Cjps6GbSEAhGE2O8jsR.malh7oD9.KKeIIVAg.G3WDU49QlZ0iJboNfW6qZq3kIilqhxS.00S4OOIhCMVC8U2sJ; _ga_0XM0LYXGC8=GS1.1.1729479287.29.1.1729479414.0.0.0")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var result GGGMGN
	err = json.Unmarshal(bodyText, &result)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Got GMGN result with length:", len(result.Data.History))
	if len(result.Data.History) > 0 {
		return true
	}
	return false
}

func FindMatchingUsers(firstResult, nextResult []Result) []Result {
	userMap := make(map[string]Result)
	for _, result := range firstResult {
		userMap[result.UserAddress] = result
	}
	var matchedResults []Result
	for _, result := range nextResult {
		if _, exists := userMap[result.UserAddress]; exists {
			matchedResults = append(matchedResults, result)
		}
		// else {
		// 	FetchGMGN := FetchGMGN(result.TokenAddress, result.UserAddress)
		// 	if FetchGMGN {
		// 		matchedResults = append(matchedResults, result)
		// 	}
		// }
	}
	return matchedResults
}
func MainFetch(token, minbuy string) ([]Result, error) {
	fmt.Println("Processing token", token)
	var wg sync.WaitGroup
	var AllResult []Result
	totalTxn, err := FetchTotalData(token, minbuy)
	if err != nil {
		return nil, err
	}
	totalTxnInt, _ := strconv.Atoi(totalTxn)
	fmt.Printf("Total transaction: %d With Min Buy: %s\n", totalTxnInt, minbuy)
	GenerateOffset := GenerateLimitOffset(totalTxnInt, 200)
	for _, v := range GenerateOffset {
		wg.Add(1)
		go func(v string) {
			defer wg.Done()
			offspl := strings.Split(v, " ")
			limit := offspl[0]
			offset := offspl[1]
			fmt.Println("Fetching data with limit:", limit, "offset:", offset)
			pfres, err := FetchData(token, minbuy, limit, offset)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Got Result with length:", len(pfres))

			for _, v := range pfres {
				timestamp := int64(v.Timestamp)
				t := time.Unix(timestamp, 0)
				formattedTime := t.Format("Jan 2 3:04:05 pm")
				AllResult = append(AllResult, Result{
					TokenAddress: v.Mint,
					UserAddress:  v.User,
					DateTime:     formattedTime,
					Amount:       LamportsToSol(int64(v.SolAmount)),
					TxHash:       v.Signature,
				})
			}
		}(v)
		wg.Wait()
	}
	fmt.Println("Total Result:", len(AllResult))
	return AllResult, nil
}

func GenerateLimitOffset(totalResult, limitPerRequest int) []string {
	var result []string

	// Iterasi hingga totalResult tercakup sepenuhnya
	for offset := 0; offset < totalResult; offset += limitPerRequest {
		// Set limit untuk setiap request
		currentLimit := offset + limitPerRequest
		result = append(result, fmt.Sprintf("%d %d", currentLimit, offset))
	}

	return result
}
func FetchData(token, minbuy, limit, offset string) ([]PFResult, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", "https://frontend-api.pump.fun/trades/all/"+token+"?limit="+limit+"&offset="+offset+"&minimumSize="+minbuy, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:132.0) Gecko/20100101 Firefox/132.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "https://pump.fun/")
	req.Header.Set("Origin", "https://pump.fun")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err

	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result []PFResult
	err = json.Unmarshal(bodyText, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func LamportsToSol(lamports int64) float64 {
	return float64(lamports) / float64(LamportsPerSol)
}
func FetchTotalData(token, minbuy string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", "https://frontend-api.pump.fun/trades/count/"+token+"?minimumSize="+minbuy, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:132.0) Gecko/20100101 Firefox/132.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "https://pump.fun/")
	req.Header.Set("Origin", "https://pump.fun")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyText), nil

}

func SolToLamports(sol float64) string {
	lamports := new(big.Int).Mul(
		big.NewInt(int64(sol*LamportsPerSol)),
		big.NewInt(1),
	)

	return lamports.String()
}
