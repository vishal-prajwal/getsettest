package digilocker

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"bitbucket.org/junglee_games/getsetgo/downloader"
	"github.com/pkg/errors"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AccountStatusRequest struct {
	Mobile  string `json:"mobile"`
	Aadhaar string `json:"aadhaar"`
}

type AccountStatusDetails struct {
	Code                 string `json:"code"`
	Mobile               string `json:"mobile"`
	Aadhaar              string `json:"aadhaar"`
	MobileAadhaarLinkage string `json:"mobileAadhaarLinkage"`
}

type AccountStatusResponse struct {
	Status     string               `json:"status"`
	StatusCode string               `json:"statusCode"`
	Result     AccountStatusDetails `json:"result"`
	Error      Error                `json:"error"`
}

type KYCStartRequest struct {
	ReferenceId string `json:"referenceId"`
	RedirectURL string `json:"redirectURL"`
}
type KYCStartDetails struct {
	URL string `json:"url"`
}
type KYCStartResponse struct {
	Status     string          `json:"status"`
	StatusCode string          `json:"statusCode"`
	Result     KYCStartDetails `json:"result"`
	Error      Error           `json:"error"`
}

type EAadhaarDetailsRequest struct {
	ReferenceId string `json:"referenceId"`
	AadhaarFile string `json:"aadhaarFile"`
}

type AadhaarDetails struct {
	Name string `json:"name"`
	// format dd-mm-yy
	DOB                 string `json:"dob"`
	Address             string `json:"address"`
	Pincode             string
	MaskedAadhaarNumber string `json:"maskedAadhaarNumber"`
	Photo               string `json:"photo"`
	AadhaarFile         string `json:"aadhaarFile"`
	XMLAadhaarFile      string `json:"xmlAadhaarFile"`
}

func (this *AadhaarDetails) SethPinCodeFromXmlFile() error {

	path, err := downloader.DownloadFile(this.XMLAadhaarFile, "addhar.xml")
	if err != nil {
		return errors.Wrap(ErrDownloadFile, err.Error())
	}
	// Open our xmlFile
	xmlFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array
	var users DigilockerAddharXMLResponse
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	err = xml.Unmarshal(byteValue, &users)
	if err != nil {
		return err
	}

	this.Pincode = users.CertificateData.KycRes.UidData.LData.Pc
	return nil
}

type AadhaarDetailsResponse struct {
	Status     string         `json:"status"`
	StatusCode string         `json:"statusCode"`
	Result     AadhaarDetails `json:"result"`
	Error      Error          `json:"error"`
}

type HypervergeHealthcheckResponse struct {
	StatusCode string            `json:"statusCode"`
	Result     HealthcheckResult `json:"result"`
}

type HealthcheckResult struct {
	Message          string `json:"message"`
	Endpoint         string `json:"endpoint"`
	Severity         string `json:"severity"`
	MeanResponseTime string `json:"meanResponseTime"`
	UserErrors       string `json:"userErrors"`
}

type DigilockerAddharXMLResponse struct {
	XMLName         xml.Name `xml:"Certificate"`
	Text            string   `xml:",chardata"`
	CertificateData struct {
		Text   string `xml:",chardata"`
		KycRes struct {
			Text    string `xml:",chardata"`
			Code    string `xml:"code,attr"`
			Ret     string `xml:"ret,attr"`
			Ts      string `xml:"ts,attr"`
			Ttl     string `xml:"ttl,attr"`
			Txn     string `xml:"txn,attr"`
			UidData struct {
				Text string `xml:",chardata"`
				Tkn  string `xml:"tkn,attr"`
				Uid  string `xml:"uid,attr"`
				Poi  struct {
					Text   string `xml:",chardata"`
					Dob    string `xml:"dob,attr"`
					Gender string `xml:"gender,attr"`
					Name   string `xml:"name,attr"`
				} `xml:"Poi"`
				Poa struct {
					Text    string `xml:",chardata"`
					Co      string `xml:"co,attr"`
					Country string `xml:"country,attr"`
					Dist    string `xml:"dist,attr"`
					Lm      string `xml:"lm,attr"`
					Loc     string `xml:"loc,attr"`
					Pc      string `xml:"pc,attr"`
					State   string `xml:"state,attr"`
					Vtc     string `xml:"vtc,attr"`
				} `xml:"Poa"`
				LData struct {
					Text    string `xml:",chardata"`
					Co      string `xml:"co,attr"`
					Country string `xml:"country,attr"`
					Dist    string `xml:"dist,attr"`
					Lang    string `xml:"lang,attr"`
					Lm      string `xml:"lm,attr"`
					Loc     string `xml:"loc,attr"`
					Name    string `xml:"name,attr"`
					Pc      string `xml:"pc,attr"`
					State   string `xml:"state,attr"`
					Vtc     string `xml:"vtc,attr"`
				} `xml:"LData"`
				Pht string `xml:"Pht"`
			} `xml:"UidData"`
		} `xml:"KycRes"`
	} `xml:"CertificateData"`
	Signature struct {
		Text       string `xml:",chardata"`
		Xmlns      string `xml:"xmlns,attr"`
		SignedInfo struct {
			Text                   string `xml:",chardata"`
			CanonicalizationMethod struct {
				Text      string `xml:",chardata"`
				Algorithm string `xml:"Algorithm,attr"`
			} `xml:"CanonicalizationMethod"`
			SignatureMethod struct {
				Text      string `xml:",chardata"`
				Algorithm string `xml:"Algorithm,attr"`
			} `xml:"SignatureMethod"`
			Reference struct {
				Text       string `xml:",chardata"`
				URI        string `xml:"URI,attr"`
				Transforms struct {
					Text      string `xml:",chardata"`
					Transform struct {
						Text      string `xml:",chardata"`
						Algorithm string `xml:"Algorithm,attr"`
					} `xml:"Transform"`
				} `xml:"Transforms"`
				DigestMethod struct {
					Text      string `xml:",chardata"`
					Algorithm string `xml:"Algorithm,attr"`
				} `xml:"DigestMethod"`
				DigestValue string `xml:"DigestValue"`
			} `xml:"Reference"`
		} `xml:"SignedInfo"`
		SignatureValue string `xml:"SignatureValue"`
		KeyInfo        struct {
			Text     string `xml:",chardata"`
			X509Data struct {
				Text            string   `xml:",chardata"`
				X509SubjectName []string `xml:"X509SubjectName"`
				X509Certificate []string `xml:"X509Certificate"`
			} `xml:"X509Data"`
		} `xml:"KeyInfo"`
	} `xml:"Signature"`
}

type PanDetails struct {
	PanNumber   string
	Name        string
	DOB         string
	Address     string
	DocImageUrl string
}

type DigilockerPanResponse struct {
	Status     string              `json:"status"`
	StatusCode string              `json:"statusCode"`
	Result     DigilockerPanResult `json:"result"`
	Error      Error               `json:"error"`
}

type DigilockerPanResult struct {
	Code      string                        `json:"code"`
	DocsFound []string                      `json:"docs_found"`
	Details   []PanDigilockerRespnseDetails `json:"details"`
}

type PanDigilockerRespnseDetails struct {
	PAN     string `json:"pan"`
	Name    string `json:"name"`
	DOB     string `json:"dob"`
	FileUrl string `json:"file"`
}
type KYCResult struct {
	Status     string    `json:"status"`
	StatusCode string    `json:"statusCode"`
	Result     []KYCInfo `json:"result"`
	Error      Error     `json:"error"`
}

type KYCInfo struct {
	URI  string  `json:"uri"`
	Data KYCData `json:"data"`
}

type KYCData struct {
	Number  string `json:"number"`
	Name    string `json:"name"`
	DOB     string `json:"dob"`
	Address string `json:"address"`
	Photo   string `json:"photo"`
	File    string `json:"file"`
}

type PanDetailsRequestPayload struct {
	ReferenceID string         `json:"referenceId"`
	Docs        []DocumentInfo `json:"docs"`
}

type DigilockerDocumentsRequetsPayload struct {
	ReferenceID string `json:"referenceId"`
	PAN         string `json:"pan"`
	PANFile     string `json:"panFile"`
}

type DocumentInfo struct {
	DocID  string         `json:"docId"`
	File   string         `json:"file"`
	Params DocumentParams `json:"params"`
}

type DocumentParams struct {
	PANNo       string `json:"panno"`
	PANFullName string `json:"PANFullName"`
}
