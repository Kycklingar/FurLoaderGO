package inkbunny

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func httpError(res *http.Response) error {
	if res.StatusCode != 200 {
		return fmt.Errorf("Http error [%d]", res.StatusCode)
	}
	return nil
}

type ibJsonError struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type ibJsonLogin struct {
	ibJsonError
	Sid string `json:"sid"`
}

type ibJsonSearch struct {
	ibJsonError
	Submissions []ibJsonSub `json:"submissions"`
}

type ibJsonSub struct {
	SubID     string `json:"submission_id"`
	Username  string `json:"username"`
	FileName  string `json:"file_name"`
	FileURL   string `json:"file_url"`
	PageCount string `json:"pagecount"`

	Files []ibJsonFiles `json:"files"`
}

type ibJsonFiles struct {
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url_full"`
}

type ibJsonWatchlist struct {
	ibJsonError
	Watches []ibJsonUser `json:"watches"`
}

type ibJsonUser struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

func (i *ibJsonLogin) decode(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(&i); err != nil {
		return err
	}

	if len(i.ErrorMessage) > 0 {
		return fmt.Errorf("%d %s", i.ErrorCode, i.ErrorMessage)
	}
	return nil
}

func (i *ibJsonSearch) decode(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(&i); err != nil {
		return err
	}

	if len(i.ErrorMessage) > 0 {
		return fmt.Errorf("%d %s", i.ErrorCode, i.ErrorMessage)
	}
	return nil
}

func (i *ibJsonWatchlist) decode(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(&i); err != nil {
		return err
	}

	if len(i.ErrorMessage) > 0 {
		return fmt.Errorf("%d %s", i.ErrorCode, i.ErrorMessage)
	}
	return nil
}
