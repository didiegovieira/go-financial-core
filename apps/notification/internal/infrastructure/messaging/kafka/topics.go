package kafka

const (
	TopicAccountCreated = "core.account.created"
	TopicAccountDeleted = "core.account.deleted"
	TopicWalletCredited = "core.wallet.credited"
	TopicWalletDebited  = "core.wallet.debited"
)

var Topics = []string{
	TopicAccountCreated,
	TopicAccountDeleted,
	TopicWalletCredited,
	TopicWalletDebited,
}

var eventNameByTopic = map[string]string{
	TopicAccountCreated: "account.created",
	TopicAccountDeleted: "account.deleted",
	TopicWalletCredited: "wallet.credited",
	TopicWalletDebited:  "wallet.debited",
}
