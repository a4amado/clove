// we have 3 plans available
// Free, Standard, Pro
// each plan has its own limits and features
// free has a message limit of 4KB per message, 10000 per month  and  50 concurrent connections
// standard has a message limit of 16KB per message, unlimited monthly messages, and 200 concurrent connections
// pro has a message limit of 32KB per message, unlimited monthly messages, and 1000 concurrent connections
package plans

import (
	"clove/internals/repository"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

const KBtoBytes = 1024

type limit struct {
	MessageSizeLimit           uint32 // in KB
	MonthlyMessageLimit        int32  // -1 means unlimited
	ConcurrentConnectionsLimit uint32
}

var planLimits = map[repository.AppType]limit{
	repository.AppTypeFree: {
		MessageSizeLimit:           4 * KBtoBytes,
		MonthlyMessageLimit:        10000,
		ConcurrentConnectionsLimit: 50,
	},
	repository.AppTypeStandard: {
		MessageSizeLimit:           16 * KBtoBytes,
		MonthlyMessageLimit:        -1, // -1 means unlimited
		ConcurrentConnectionsLimit: 200,
	},
	repository.AppTypePro: {
		MessageSizeLimit:           32 * KBtoBytes,
		MonthlyMessageLimit:        -1, // -1 means unlimited
		ConcurrentConnectionsLimit: 1000,
	},
}

var (
	ErrMessageSizeLimitExceeded           = errors.New("message size limit exceeded, please upgrade your plan")
	ErrMonthlyMessageLimitExceeded        = errors.New("monthly message limit exceeded, please upgrade your plan")
	ErrConcurrentConnectionsLimitExceeded = errors.New("concurrent connections limit exceeded, please upgrade your plan")
)

type planOptions struct {
	// if appId is nil, only MessageSizeLimit is enforced
	appId       *pgtype.UUID
	messageSize uint32
}
type Option func(*planOptions)

// WithAppID sets the application ID for plan limit checking.
// When appId is provided, MonthlyMessageLimit and ConcurrentConnectionsLimit are enforced.
// WithAppID returns an Option that sets the application ID used during plan validation.
// When an app ID is provided, monthly message and concurrent connection limits will be enforced.
func WithAppID(appId pgtype.UUID) Option {
	return func(options *planOptions) {
		options.appId = &appId
	}
}

// WithMessageSize returns an Option that sets the message size (in bytes) on a planOptions instance.
func WithMessageSize(size uint32) Option {
	return func(options *planOptions) {
		options.messageSize = size
	}
}

// ValidatePlan validates the given plan against configured limits using the provided options.
// It enforces the plan's message size limit and, when an application ID is supplied and the plan
// defines a finite monthly limit, is intended to also enforce monthly message and concurrent
// connections limits (those checks are not implemented yet).
//
// The function panics if no options are provided.
//
// Returns ErrMessageSizeLimitExceeded when the supplied message size exceeds the plan's limit,
// or nil when validation passes.
func ValidatePlan(ctx context.Context, plan repository.AppType, options ...Option) error {

	if len(options) == 0 {
		panic("at least one option is required")
	}
	message := planOptions{}
	for _, option := range options {
		option(&message)
	}

	if message.messageSize > planLimits[plan].MessageSizeLimit {
		return ErrMessageSizeLimitExceeded
	}

	if planLimits[plan].MonthlyMessageLimit != -1 && message.appId != nil {
		// TODO: check monthly message limit

		// TODO: check concurrent connections limit
	}

	return nil
}
