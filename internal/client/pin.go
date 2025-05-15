package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func (c *Client) GetPinById(id string) (*PinById, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v3/files/public/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	pin := PinById{}
	err = json.Unmarshal(body, &pin)
	if err != nil {
		return nil, err
	}

	return &pin, nil
}

func (c *Client) PinFolder(files []string, name, version string) (*PinFileToIpfs, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, file := range files {
		reader, _ := os.Open(file)
		defer reader.Close()
		part, _ := writer.CreateFormFile("file", file)
		io.Copy(part, reader)
	}
	metadata, _ := writer.CreateFormField("pinataMetadata")
	metadata.Write(fmt.Appendf(nil, `{"name":"%s"}`, name))
	options, _ := writer.CreateFormField("pinataOptions")
	options.Write(fmt.Appendf(nil, `{"cidVersion":%s}`, version))
	writer.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/pinning/pinFileToIPFS", c.HostURL), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	raw, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := PinFileToIpfs{}
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) Unpin(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v3/files/public/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
