package piefed

import (
	piefedResponse "LemmyBeProxy/dto/response/piefed"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	goHttp "net/http"
)

// UploadImage forwards a file's raw bytes to Piefed's /upload/image endpoint.
// This bypasses defaultHandler/sendRequest entirely because those are built
// for JSON request bodies only — image upload needs a real multipart body,
// matching what mlmym itself sends when talking to a real Lemmy pict-rs.
func (receiver *Piefed) UploadImage(fileBytes []byte, filename string, bearerToken string) (*piefedResponse.ImageUploadResponse, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(fileBytes); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := goHttp.NewRequest(
		"POST",
		receiver.url()+"/upload/image",
		body,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	res, err := goHttp.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != goHttp.StatusOK {
		return nil, fmt.Errorf("piefed image upload failed with status %d: %s", res.StatusCode, string(respBody))
	}

	var response piefedResponse.ImageUploadResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
