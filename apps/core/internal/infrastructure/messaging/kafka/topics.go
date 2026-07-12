package kafka

import appmessaging "github.com/didiegovieira/go-financial-core/apps/core/internal/application/port/messaging"

const (
	TopicAccountCreated = "core.account.created"
	TopicAccountDeleted = "core.account.deleted"
	TopicWalletCredited = "core.wallet.credited"
	TopicWalletDebited  = "core.wallet.debited"
)

var topicByEventName = map[appmessaging.EventName]string{
	appmessaging.EventAccountCreated: TopicAccountCreated,
	appmessaging.EventAccountDeleted: TopicAccountDeleted,
	appmessaging.EventWalletCredited: TopicWalletCredited,
	appmessaging.EventWalletDebited:  TopicWalletDebited,
}
