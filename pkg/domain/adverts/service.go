package adverts

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"

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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.addr, bytes.NewReader(jsonBody))
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

	if resp.StatusCode != 200 {
		return nil, nil
	}

	var result []model.Advertisement
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
