package tokenprovider

import (
	"context"
	"errors"
	"fmt"

	"github.com/90poe/service-chassis/m2m/v2"
)

const (
	Ahoy   = "c0bf8da5-8fd0-48b5-ae3a-c13c23a630de"
	Zodiac = "8a396ee1-a16b-4cc4-8b54-d3de465b8fc8"

	Dev = iota
	Int
)

type M2MTokenProvider struct {
	provider *m2m.GRPCTokenProvider
}

func Build(ctx context.Context, env int) (*M2MTokenProvider, error) {
	config := m2m.Config{
		Target: fmt.Sprintf("%s:%d", "localhost", 8090),
	}

	switch env {
	case Dev:
		config.ClientID = "EilCjNZMWNLfEZzNBkyoV1mg"
		config.ClientSecret = "Sv45fstXWIutTcQ8oGoXkI05g6wITmG49JzY"
	case Int:
		config.ClientID = "EilCjNZMWNLfEZzNBkyoVint"
		config.ClientSecret = "Sv45fstXWIutTcQ8oGoXkI05g6wITmG49int"
	default:
		return nil, errors.New("unknown env")
	}

	provider, err := m2m.New(ctx, config)
	if err != nil {
		return nil, err
	}

	return &M2MTokenProvider{provider: provider}, nil
}

func (p *M2MTokenProvider) PrepareAuthContext(ctx context.Context, accountID string) (context.Context, error) {
	return p.provider.PrepareAuthContext(ctx, accountID)
}

func (p *M2MTokenProvider) Close() error {
	return p.provider.Close()
}
