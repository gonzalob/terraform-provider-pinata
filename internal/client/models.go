package client

type PinFileToIpfs struct {
	IPFSHash string `json:"IpfsHash"`
	ID       string `json:"ID"`
	Name     string `json:"Name"`
}

type PinById struct {
	Data struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		CID  string `json:"cid"`
	} `json:"data"`
}
