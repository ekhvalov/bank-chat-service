package keycloakclient

//go:generate options-gen -out-filename=client_options.gen.go -from-struct=Options

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Options struct {
	basePath  string `option:"mandatory" validate:"required,http_url"`
	realm     string `option:"mandatory" validate:"required"`
	username  string `option:"mandatory" validate:"required"`
	password  string `option:"mandatory" validate:"required"`
	userAgent string `validate:"required"`
	debugMode bool
}

// Client is a tiny client to the KeyCloak realm operations. UMA configuration:
// http://localhost:3010/realms/Bank/.well-known/uma2-configuration
type Client struct {
	username string
	password string
	realm    string
	cli      *resty.Client
}

func New(opts Options) (*Client, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	cli := resty.New()
	cli.SetDebug(opts.debugMode)
	cli.SetBaseURL(opts.basePath)
	if opts.userAgent != "" {
		cli.Header.Set("User-Agent", opts.userAgent)
	}

	return &Client{
		username: opts.username,
		password: opts.password,
		realm:    opts.realm,
		cli:      cli,
	}, nil
}
