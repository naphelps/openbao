package mongodb

import (
	"context"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigLease(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/lease",
		Fields: map[string]*framework.FieldSchema{
			"ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Default ttl for credentials.",
			},

			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Maximum time a set of credentials can be valid for.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigLeaseRead,
			logical.UpdateOperation: b.pathConfigLeaseWrite,
		},

		HelpSynopsis:    pathConfigLeaseHelpSyn,
		HelpDescription: pathConfigLeaseHelpDesc,
	}
}

func (b *backend) pathConfigLeaseWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config/lease", &configLease{
		TTL:    time.Second * time.Duration(d.Get("ttl").(int)),
		MaxTTL: time.Second * time.Duration(d.Get("max_ttl").(int)),
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigLeaseRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	leaseConfig, err := b.LeaseConfig(req.Storage)

	if err != nil {
		return nil, err
	}
	if leaseConfig == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"ttl":     leaseConfig.TTL.Seconds(),
			"max_ttl": leaseConfig.MaxTTL.Seconds(),
		},
	}, nil
}

type configLease struct {
	TTL    time.Duration
	MaxTTL time.Duration
}

const pathConfigLeaseHelpSyn = `
Configure the default lease TTL settings for credentials
generated by the mongodb backend.
`

const pathConfigLeaseHelpDesc = `
This configures the default lease TTL settings used for
credentials generated by this backend. The ttl specifies the
duration that a set of credentials will be valid for before
the lease must be renewed (if it is renewable), while the
max_ttl specifies the overall maximum duration that the
credentials will be valid regardless of lease renewals.

The format for the TTL values is an integer and then unit. For
example, the value "1h" specifies a 1-hour TTL. The longest
supported unit is hours.
`
