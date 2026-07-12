package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func (c *Client) GetPinById(ctx context.Context, id string) (*PinById, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v3/files/public/%s", c.HostURL, id), nil)
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

func (c *Client) PinFolder(ctx context.Context, files []string, name, version string) (*PinFileToIpfs, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer writer.Close()
		for _, file := range files {
			err := func() error {
				reader, err := os.Open(file)
				if err != nil {
					return err
				}
				defer reader.Close()
				part, err := writer.CreateFormFile("file", file)
				if err != nil {
					return err
				}
				_, err = io.Copy(part, reader)
				return err
			}()
			if err != nil {
				pw.CloseWithError(err)
				return
			}
		}
		metadata, err := writer.CreateFormField("pinataMetadata")
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		if _, err := metadata.Write(fmt.Appendf(nil, `{"name":"%s"}`, name)); err != nil {
			pw.CloseWithError(err)
			return
		}
		options, err := writer.CreateFormField("pinataOptions")
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		if _, err := options.Write(fmt.Appendf(nil, `{"cidVersion":%s}`, version)); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/pinning/pinFileToIPFS", c.HostURL), pr)
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

func (c *Client) Unpin(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/v3/files/public/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
