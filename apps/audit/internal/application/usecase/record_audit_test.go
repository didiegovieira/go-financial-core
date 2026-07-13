package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/didiegovieira/go-financial-core/apps/audit/internal/application/dto"
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/application/port/repository"
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/domain/entity"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestRecordAuditSavesMappedAuditLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	auditLogRepository := repository.NewMockAuditLogRepository(ctrl)
	usecase := NewRecordAudit(auditLogRepository)
	input := dto.RecordAuditInput{
		EventName:   "payment.created",
		AggregateID: "account-123",
		Payload:     json.RawMessage(`{"amount":1000}`),
	}
	before := time.Now().UTC()

	auditLogRepository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, auditLog entity.AuditLog) error {
			if auditLog.ID == uuid.Nil {
				t.Fatal("expected audit log ID to be generated")
			}
			if auditLog.EventName != input.EventName {
				t.Fatalf("expected event name %q, got %q", input.EventName, auditLog.EventName)
			}
			if auditLog.AggregateID != input.AggregateID {
				t.Fatalf("expected aggregate ID %q, got %q", input.AggregateID, auditLog.AggregateID)
			}
			if string(auditLog.Payload) != string(input.Payload) {
				t.Fatalf("expected payload %q, got %q", string(input.Payload), string(auditLog.Payload))
			}
			if auditLog.OccurredAt.Before(before) || auditLog.OccurredAt.After(time.Now().UTC()) {
				t.Fatalf("expected occurred at between %s and now, got %s", before, auditLog.OccurredAt)
			}

			return nil
		})

	err := usecase.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRecordAuditReturnsRepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	auditLogRepository := repository.NewMockAuditLogRepository(ctrl)
	usecase := NewRecordAudit(auditLogRepository)
	expectedErr := errors.New("save audit failed")

	auditLogRepository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(expectedErr)

	err := usecase.Execute(context.Background(), dto.RecordAuditInput{
		EventName:   "payment.created",
		AggregateID: "account-123",
		Payload:     json.RawMessage(`{"amount":1000}`),
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
