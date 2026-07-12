package risk

import "testing"

func TestPolicyEvaluateApprovesTransactionBelowReviewThreshold(t *testing.T) {
	policy := NewDefaultPolicy()

	decision := policy.Evaluate(100_00)

	if decision.Status != DecisionApproved {
		t.Fatalf("expected status %q, got %q", DecisionApproved, decision.Status)
	}
}

func TestPolicyEvaluateRequiresReviewAtReviewThreshold(t *testing.T) {
	policy := NewDefaultPolicy()

	decision := policy.Evaluate(defaultReviewAmountInCents)

	if decision.Status != DecisionReview {
		t.Fatalf("expected status %q, got %q", DecisionReview, decision.Status)
	}
}

func TestPolicyEvaluateRejectsAtRejectThreshold(t *testing.T) {
	policy := NewDefaultPolicy()

	decision := policy.Evaluate(defaultRejectAmountInCents)

	if decision.Status != DecisionRejected {
		t.Fatalf("expected status %q, got %q", DecisionRejected, decision.Status)
	}
}
