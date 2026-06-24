package payment

import (
	"context"

	dto "github.com/didiegovieira/go-payments-core/apps/payments/internal/application/dto/payment"
	"github.com/didiegovieira/go-payments-core/pkg/base"
	"github.com/didiegovieira/go-payments-core/pkg/metrics"
)

type CreateUseCase = base.UseCase[dto.CreateInput, *dto.CreateOutput]

type Create struct{}

func NewCreate() *Create {
	return &Create{}
}

func (c *Create) Execute(ctx context.Context, input dto.CreateInput) (out *dto.CreateOutput, err error) {
	ctx, finish := metrics.TraceWithError(ctx, "UseCase.Payment.Create", &err)
	defer finish()

	return out, nil
}
