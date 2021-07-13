package storage

type QueueAttributes struct {
	DelaySeconds                  uint `validate:"min=0,max=900"`
	MaximumMessageSize            uint `validate:"min=1024,max=262144"`
	MessageRetentionPeriod        uint `validate:"min=60,max=1209600"`
	ReceiveMessageWaitTimeSeconds uint `validate:"min=0,max=20"`
	VisibilityTimeout             uint `validate:"min=0,max=43200"`
}

type QueueTags map[string]string

type CreateQueueInput struct {
	QueueName  string          `validate:"required,min=5,max=50"`
	Attributes QueueAttributes `validate:"required"`
	Tags       QueueTags       ``
}
