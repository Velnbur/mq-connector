package mqconnector

import "context"

type Runner interface {
	Run(ctx context.Context)
}
