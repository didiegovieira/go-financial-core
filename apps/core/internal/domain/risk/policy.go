package risk

const (
	defaultReviewAmountInCents = int64(500_00)
	defaultRejectAmountInCents = int64(5_000_00)
)

type DecisionStatus string

const (
	DecisionApproved DecisionStatus = "approved"
	DecisionReview   DecisionStatus = "review"
	DecisionRejected DecisionStatus = "rejected"
)

type Policy struct {
	ReviewAmountInCents int64
	RejectAmountInCents int64
}

type Decision struct {
	Status DecisionStatus
	Score  int
	Reason string
}

func NewDefaultPolicy() Policy {
	return Policy{
		ReviewAmountInCents: defaultReviewAmountInCents,
		RejectAmountInCents: defaultRejectAmountInCents,
	}
}

func (p Policy) Evaluate(amountInCents int64) Decision {
	reviewAmount := p.ReviewAmountInCents
	if reviewAmount == 0 {
		reviewAmount = defaultReviewAmountInCents
	}

	rejectAmount := p.RejectAmountInCents
	if rejectAmount == 0 {
		rejectAmount = defaultRejectAmountInCents
	}

	decision := Decision{
		Status: DecisionApproved,
		Score:  10,
		Reason: "transaction within configured risk policy",
	}

	if amountInCents >= rejectAmount {
		decision.Status = DecisionRejected
		decision.Score = 95
		decision.Reason = "transaction amount exceeds rejection threshold"
		return decision
	}

	if amountInCents >= reviewAmount {
		decision.Status = DecisionReview
		decision.Score = 65
		decision.Reason = "transaction amount requires manual review"
	}

	return decision
}
