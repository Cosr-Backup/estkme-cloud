package lpac

type ChipInfo struct {
	EID                      string                   `json:"eidValue"`
	EuiccConfiguredAddresses EuiccConfiguredAddresses `json:"euiccConfiguredAddresses"`
	EUICCInfo2               EuiccInfo2               `json:"euiccInfo2"`
}

type EuiccConfiguredAddresses struct {
	DefaultDPAddress string `json:"defaultDpAddress"`
	RootDSAddress    string `json:"rootDsAddress"`
}

type EuiccInfo2 struct {
	SasAccreditationNumber string          `json:"sasAcreditationNumber"`
	ExtCardResource        ExtCardResource `json:"extCardResource"`
	PkiForSigning          []string        `json:"euiccCiPKIdListForSigning"`
}

type ExtCardResource struct {
	FreeNonVolatileMemory int `json:"freeNonVolatileMemory"`
	FreeVolatileMemory    int `json:"freeVolatileMemory"`
}

func (c *Cmd) Info() (*ChipInfo, error) {
	info := &ChipInfo{}
	if err := c.Run([]string{"chip", "info"}, &info, nil); err != nil {
		return info, err
	}
	return info, nil
}
