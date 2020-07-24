// Code generated by sysl DO NOT EDIT.
package dbendpoints

import (
	"context"
	"database/sql"
	"time"
)

// DefaultDbEndpointsImpl ...
type DefaultDbEndpointsImpl struct {
}

// NewDefaultDbEndpointsImpl for DbEndpoints
func NewDefaultDbEndpointsImpl() *DefaultDbEndpointsImpl {
	return &DefaultDbEndpointsImpl{}
}

// GetCompanyLocationList Client
type GetCompanyLocationListClient struct {
	conn                         *sql.Conn
	retrievebycompanyandlocation *sql.Stmt
}

// ServiceInterface for DbEndpoints
type ServiceInterface struct {
	GetCompanyLocationList func(ctx context.Context, req *GetCompanyLocationListRequest, client GetCompanyLocationListClient) (*GetCompanyLocationResponse, error)
}

// DownstreamConfig for DbEndpoints
type DownstreamConfig struct {
	ContextTimeout time.Duration `yaml:"contextTimeout"`
}