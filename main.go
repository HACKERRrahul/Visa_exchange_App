package main

import (
	"encoding/json"
	//"log"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	//"github.com/gorilla/mux"
)

func main() {
	// http.HandleFunc("/convert", CurrencyConversionHandler)
	// http.ListenAndServe(":9090", nil)
	// log.Fatal(http.ListenAndServe(":9090", nil))
	lambda.Start(CurrencyConversionHandler)
	//r := mux.NewRouter()
	//r.HandleFunc("/convert", CurrencyConversionHandler).Methods("POST")
	//http.ListenAndServe(":9090", r)

}

var (

	// THIS IS EXAMPLE ONLY how will user_id and password look like
	// userId = "1WM2TT4IHPXC8DQ5I3CH21n1rEBGK-Eyv_oLdzE2VZpDqRn_U";
	// password = "19JRVdej9";
	username = "XQ3DQ32CB0FI2COHC237217KbvBb1uhY_jUhl9xjj-uoe23oA"
	password = "7A7R7n2I6Np0DFrLZ5ZOP9imWBQCRh6itEMJv"

	// THIS IS EXAMPLE ONLY how will cert and key look like
	// clientCertificateFile = 'cert.pem'
	// clientCertificateKeyFile = 'key_83d11ea6-a22d-4e52-b310-e0558816727d.pem'
	// caCertificateFile = 'ca_bundle.pem'

	clientCertificateFile    = GetCurrentPath() + "cert_prod.pem" //"https://certifcates007.s3.ap-northeast-1.amazonaws.com/cert.pem"//arn:aws:lambda:ap-northeast-1:445192904874:layer:clientCertificateFile:1"//"/Users/rahul.kaushik/Downloads/cert.pem"
	clientCertificateKeyFile = GetCurrentPath() + "key_prod.pem"  //"https://certifcates007.s3.ap-northeast-1.amazonaws.com/key.pem"//"arn:aws:lambda:ap-northeast-1:445192904874:layer:clientCertificateKeyFile:1"///Users/rahul.kaushik/Downloads/key.pem"
	caCertificateFile        = GetCurrentPath() + "ca_prod.pem"   //"https://certifcates007.s3.ap-northeast-1.amazonaws.com/ca.pem"//arn:aws:lambda:ap-northeast-1:445192904874:layer:caCertificateFile:1"///Users/rahul.kaushik/Downloads/ca.pem"
)

func GetCurrentPath() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path + "/"
}

type IncRequest struct {
	FromCurrencyCode string `json:"fromCurrencyCode"`
	ToCurrencyCode   string `json:"toCurrencyCode"`
	FromAmount       string `json:"fromAmount"`
}
type Response struct {
	FromCurrencyCode  string `json:"fromCurrencyCode"`
	ToCurrencyCode    string `json:"toCurrencyCode"`
	FromAmount        string `json:"fromAmount"`
	DestinationAmount string `json:"destinationAmount"`
	ConversionRate string `json:"conversionRate"`
}

func CurrencyConversionHandler(incRequest IncRequest) (Response, error) {

	//headerContentTtype := r.Header.Get("Content-Type")
	//log.Println(headerContentTtype)

	//log.Println(incRequest, err)

	clientCACert, err := ioutil.ReadFile(caCertificateFile)
	if err != nil {
		panic(err)
	}

	//Load Client Key Pair
	clientKeyPair, err := tls.LoadX509KeyPair(clientCertificateFile, clientCertificateKeyFile)
	if err != nil {
		panic(err)
	}
	clientCertPool, _ := x509.SystemCertPool()
	if clientCertPool == nil {
		clientCertPool = x509.NewCertPool()
	}

	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientKeyPair},
		RootCAs:      clientCertPool,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	url := "https://sandbox.api.visa.com/forexrates/v2/foreignexchangerates"
	method := "POST"

	payload := &sampleRequest{}
	AcquirerDetail := &acquirerDetail{}
	AcquirerDetail.Bin = 408999
	AcquirerDetail.Settlement = *&Settlements{}
	AcquirerDetail.Settlement.CurrencyCode = "840"
	payload.AcquirerDetails = *AcquirerDetail
	payload.MarkupRate = "0.07"
	payload.DestinationCurrencyCode = incRequest.ToCurrencyCode
	payload.SourceAmount = incRequest.FromAmount
	payload.SourceCurrencyCode = incRequest.FromCurrencyCode
	payload.RateProductCode = "A"

	client := &http.Client{Transport: transport}
	stringpayload, err := json.Marshal(&payload)
	if err != nil {
		panic(err)
	}
	//fmt.Println(payload)
	Payload := strings.NewReader(string(stringpayload))

	req, err := http.NewRequest(method, url, Payload)

	if err != nil {
		fmt.Println(err)

	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic SzNQWUNLWUgwNU5ZVTBJQ1I0UEYyMXdKc0hjQVRDUkJYR2xOMWEzMEZWd2tfQjM2QTpmb01VTktnOFRzcUpJdmVVSWdBcGdETHhQVDBBM3Z4dzg=")
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}
	var resBody map[string]interface{}
	err = json.Unmarshal(body, &resBody)
	if err != nil {
		fmt.Println(err)

	}

	resp := &Response{}
	resp.FromAmount = incRequest.FromAmount
	resp.FromCurrencyCode = incRequest.FromCurrencyCode
	resp.ToCurrencyCode = incRequest.ToCurrencyCode

	destAmt, exists := resBody["destinationAmount"].(string)
	if exists {
		resp.DestinationAmount = destAmt
	}
	convrate,exists:=resBody["conversionRate"].(string)
	if exists{
        resp.ConversionRate=convrate
	}
	
	fmt.Println(string(body))
	return *resp, nil

}

type sampleRequest struct {
	SourceCurrencyCode      string         `json:"sourceCurrencyCode"`
	SourceAmount            string         `json:"sourceAmount"`
	AcquirerDetails         acquirerDetail `json:"acquirerDetails"`
	DestinationCurrencyCode string         `json:"destinationCurrencyCode"`
	MarkupRate              string         `json:"markupRate"`
	RateProductCode         string         `json:"rateProductCode"`
}

type acquirerDetail struct {
	Bin        int         `json:"bin"`
	Settlement Settlements `json:"settlement"`
}
type Settlements struct {
	CurrencyCode string `json:"currencyCode"`
}
