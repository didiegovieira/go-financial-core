package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/didiegovieira/go-financial-core/apps/notification/internal/application/dto"
	"github.com/didiegovieira/go-financial-core/apps/notification/internal/application/port/client"
	"github.com/didiegovieira/go-financial-core/apps/notification/internal/domain/entity"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestSendNotificationSendsMappedNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	notificationClient := client.NewMockNotificationClient(ctrl)
	usecase := NewSendNotification(notificationClient)
	input := dto.SendNotificationInput{
		EventName:   "payment.created",
		AggregateID: "account-123",
		Payload:     json.RawMessage(`{"amount":1000}`),
	}
	before := time.Now().UTC()

	notificationClient.EXPECT().
		Send(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, notification entity.Notification) error {
			if _, err := uuid.Parse(notification.ID); err != nil {
				t.Fatalf("expected valid notification ID, got %q", notification.ID)
			}
			if notification.EventName != input.EventName {
				t.Fatalf("expected event name %q, got %q", input.EventName, notification.EventName)
			}
			if notification.AggregateID != input.AggregateID {
				t.Fatalf("expected aggregate ID %q, got %q", input.AggregateID, notification.AggregateID)
			}
			if notification.Channel != "external" {
				t.Fatalf("expected channel %q, got %q", "external", notification.Channel)
			}
			if notification.Recipient != input.AggregateID {
				t.Fatalf("expected recipient %q, got %q", input.AggregateID, notification.Recipient)
			}
			if notification.Subject != input.EventName {
				t.Fatalf("expected subject %q, got %q", input.EventName, notification.Subject)
			}
			if notification.Body != string(input.Payload) {
				t.Fatalf("expected body %q, got %q", string(input.Payload), notification.Body)
			}
			if notification.Status != "pending" {
				t.Fatalf("expected status %q, got %q", "pending", notification.Status)
			}
			if notification.CreatedAt.Before(before) || notification.CreatedAt.After(time.Now().UTC()) {
				t.Fatalf("expected created at between %s and now, got %s", before, notification.CreatedAt)
			}

			return nil
		})

	err := usecase.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSendNotificationReturnsClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	notificationClient := client.NewMockNotificationClient(ctrl)
	usecase := NewSendNotification(notificationClient)
	expectedErr := errors.New("send notification failed")

	notificationClient.EXPECT().
		Send(gomock.Any(), gomock.Any()).
		Return(expectedErr)

	err := usecase.Execute(context.Background(), dto.SendNotificationInput{
		EventName:   "payment.created",
		AggregateID: "account-123",
		Payload:     json.RawMessage(`{"amount":1000}`),
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
