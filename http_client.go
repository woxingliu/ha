package rice

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

var defautlHttpClient = &http.Client{}

func Request[R any](request *http.Request) (R, error) {
	var result R

	response, err := defautlHttpClient.Do(request)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	if len(b) == 0 {
		return result, errors.New("empty response body")
	}
	err = json.Unmarshal(b, &result)
	return result, err
}

func Get[R any](url string) (R, error) {
	var result R
	response, err := defautlHttpClient.Get(url)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	if len(b) == 0 {
		return result, errors.New("empty response body")
	}
	err = json.Unmarshal(b, &result)
	return result, err
}

// PostJson application/json
func PostJson[T, R any](url string, data T) (R, error) {
	var result R
	jsonData, err := json.Marshal(data)
	if err != nil {
		return result, err
	}
	response, err := defautlHttpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return result, err
	}
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	if len(b) == 0 {
		return result, errors.New("empty response body")
	}
	err = json.Unmarshal(b, &result)
	return result, err
}

// PostForm application/x-www-form-urlencoded
func PostForm[R any](url string, data url.Values) (R, error) {
	var result R
	response, err := defautlHttpClient.PostForm(url, data)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	if len(b) != 0 {
		if err := json.Unmarshal(b, &result); err != nil {
			return result, err
		}
	}
	return result, nil
}

type WriteMutipartFunc func(w *multipart.Writer) error

// PostMultipartForm multipart/form-data
func PostMultipartForm[R any](url string, data url.Values, wFunc WriteMutipartFunc) (R, error) {
	var result R
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for k, v := range data {
		for _, v2 := range v {
			w.WriteField(k, v2)
		}
	}
	if err := wFunc(w); err != nil {
		return result, err
	}

	response, err := defautlHttpClient.Post(url, w.FormDataContentType(), body)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	if len(b) != 0 {
		if err := json.Unmarshal(b, &result); err != nil {
			return result, err
		}
	}
	return result, nil
}
