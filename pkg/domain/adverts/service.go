package adverts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/opencars/seedwork/logger"

	"github.com/opencars/core/pkg/config"
	"github.com/opencars/core/pkg/domain/model"
)

type Service struct {
	addr   string
	secret string
	token  string

	c *http.Client
}

// TODO: Configure timeout and http client settings.
func NewService(cfg *config.ServiceHTTP) (*Service, error) {
	return &Service{
		addr:   cfg.Address(),
		secret: cfg.Secret,
		token:  cfg.Token,

		c: &http.Client{
			Timeout: time.Second,
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return net.DialTimeout(network, addr, time.Second)
				},
			},
		},
	}, nil
}

func (s *Service) FindByVINs(ctx context.Context, vins ...string) ([]model.Advertisement, error) {
	body := newRequestBody(vins...)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	logger.Infof("request: %s", string(jsonBody))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.addr+"/api/v1/data/adverts", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("statisfy-token", s.secret)
	req.Header.Set("Authorization", "Token "+s.token)

	resp, err := s.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logger.Infof("status: %s", resp.Status)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed, status: %s", resp.Status)
	}

	var result []model.Advertisement

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	logger.Infof("add: %d", result)

	return result, nil
}
