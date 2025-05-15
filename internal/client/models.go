package client

type PinFileToIpfs struct {
	IPFSHash string `json:"IpfsHash"`
	ID       string `json:"ID"`
	Name     string `json:"Name"`
}

type PinById struct {
	Data struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		CID  string `json:"cid"`
	} `json:"data"`
}

type File struct {
	name string
	path string
}
