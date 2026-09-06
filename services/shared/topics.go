package shared

const (
	PaymentTopic = "payment_topic"
	AccountTopic = "account_topic"
)

var PartitionAlias = map[string]int32{
	"starting": 0,
	"verified": 1,
}
