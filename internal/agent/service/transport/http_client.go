package transport

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/dontagr/pandora/internal/agent/config"
	store "github.com/dontagr/pandora/internal/agent/store/models"
	"github.com/dontagr/pandora/internal/models"
)

type HTTPManager struct {
	client       *http.Client
	log          *zap.SugaredLogger
	User         *store.User
	cnf          *config.Config
	serverHost   string
	waitForRetry int
}

func NewHTTPManager(log *zap.SugaredLogger, cnf *config.Config) (*HTTPManager, error) {
	manager := HTTPManager{
		log:          log,
		client:       &http.Client{Timeout: 3 * time.Second},
		waitForRetry: cnf.Transport.WaitForRetry,
		serverHost:   cnf.HTTPServer.Host,
		cnf:          cnf,
	}

	return &manager, nil
}

func (h *HTTPManager) NewRequest(method string, body *bytes.Buffer, url string, auth bool) (*models.CommonResponce, error) {
	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequest(method, fmt.Sprintf("%s%s", h.serverHost, url), nil)
	} else {
		req, err = http.NewRequest(method, fmt.Sprintf("%s%s", h.serverHost, url), body)
	}
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if h.cnf.HTTPServer.Gzip {
		req.Header.Set("Content-Encoding", "gzip")
	}
	req.Header.Set("Content-Type", "application/json")
	if auth {
		req.Header.Add("Authorization", h.User.Token)
	}

	var resp *http.Response
	var netErr *net.OpError
	var errSend error
	for i := 0; i < 3; i++ {
		resp, errSend = h.client.Do(req) // nolint
		if errSend == nil {
			responce, err := getResponce(resp)
			if err != nil {
				return nil, err
			}
			if responce.Status >= 200 && responce.Status < 300 {
				h.log.Debugf("Sent request was successfully with status code: %d", responce.Status)
			} else {
				h.log.Warnf("Received non-2xx status code: %d", responce.Status)
			}

			return responce, nil
		}
		if errors.As(errSend, &netErr) {
			h.log.Warnf("Connection error we try №%d", i+1)
			if i < 2 {
				time.Sleep(time.Duration(h.waitForRetry) * time.Second)
			}
		} else {
			return nil, fmt.Errorf("sending data: %w", errSend)
		}
	}

	return nil, fmt.Errorf("failed to send request after retrying: %w", errSend)
}

func getResponce(resp *http.Response) (*models.CommonResponce, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	out := models.CommonResponce{
		Status: resp.StatusCode,
		Body:   bodyBytes,
		Meta:   models.CommonMeta{Authorization: resp.Header.Get("Authorization")},
	}

	return &out, nil
}
