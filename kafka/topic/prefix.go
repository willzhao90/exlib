package topic

const (
	// EngineInPrefix gives all messages to match engine
	// Including: NewOrder, UpdateOrder, CancelOrder
	EngineInPrefix = "engine-in-"

	// EngineOutPrefix gives all messages from match engine
	// Including: NewOrder, UpdateOrder, CancelOrder
	EngineOutPrefix = "engine-out-"

	MemberInPrefic  = "member-in-"
	MemberOutPrefic = "member-out-"

	// SavedTopicPrefix add prefix to saved topic to create a new topic
	SavedTopicPrefix = "save-topic-"
)
