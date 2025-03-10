package location

import (
	"errors"
	"strings"
)

type IndiaState struct {
	Name         string
	NameInitials string
}

type IndiaStatesMap map[int]IndiaState

type NameOrShortToIndiaStateMap map[string]IndiaState

func (ism IndiaStatesMap) ToNameOrShortNameToIndiaStateMap() NameOrShortToIndiaStateMap {
	mapp := make(NameOrShortToIndiaStateMap)
	for _, v := range ism {
		mapp[strings.ToLower(v.Name)] = v
		mapp[strings.ToLower(v.NameInitials)] = v
	}
	return mapp
}

var indiaStates IndiaStatesMap = map[int]IndiaState{
	0:  {"Unknown", "Unknown"},
	1:  {"Andaman & Nicobar", "AN"},
	2:  {"Andhra Pradesh", "AP,AD"},
	3:  {"Arunachal Pradesh", "AR"},
	4:  {"Assam", "AS"},
	5:  {"Bihar", "BR"},
	6:  {"Chandigarh", "CH"},
	7:  {"Chhatisgarh", "CG"},
	8:  {"Dadra and Nagar Haveli", "DH,DN"},
	9:  {"Daman & Diu", "DD"},
	10: {"Goa", "GA"},
	11: {"Gujarat", "GJ"},
	12: {"Haryana", "HR"},
	13: {"Himachal Pradesh", "HP"},
	14: {"Jammu & Kashmir", "JK"},
	15: {"Jharkhand", "JH"},
	16: {"Karnataka", "KA"},
	17: {"Kerala", "KL"},
	18: {"Lakshadweep", "LD"},
	19: {"Madhya Pradesh", "MP"},
	20: {"Manipur", "MN"},
	21: {"Mizoram", "MZ"},
	22: {"Nagaland", "NL"},
	23: {"Nagar Haveli", "Nagar Haveli"},
	24: {"Delhi", "DL,OR"},
	25: {"Orissa", "OD"},
	26: {"Missing 26", "Missing 26"},
	27: {"Punjab", "PB"},
	28: {"Sikkim", "SK"},
	29: {"Tamil Nadu", "TN"},
	30: {"Tripura", "TR"},
	31: {"Uttar Pradesh", "UP"},
	32: {"Uttarakhand", "UK"},
	33: {"West Bengal", "WB"},
	34: {"Maharashtra", "MH"},
	35: {"Meghalaya", "ML"},
	36: {"Pondicherry", "PY"},
	37: {"Rajasthan", "RJ"},
	38: {"Unspecified", "Unspecified"},
	39: {"Telangana", "TG,TS"},
	40: {"Restricted State", "Restricted"},
}

// errors
var (
	ErrInvalidState = errors.New("invalid state")
)
