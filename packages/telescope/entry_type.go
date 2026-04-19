package telescope

// Entry type constants mirror Laravel's EntryType class. They identify the
// category of a recorded Telescope entry.
const (
	EntryTypeBatch         = "batch"
	EntryTypeCache         = "cache"
	EntryTypeCommand       = "command"
	EntryTypeDump          = "dump"
	EntryTypeEvent         = "event"
	EntryTypeException     = "exception"
	EntryTypeJob           = "job"
	EntryTypeLog           = "log"
	EntryTypeMail          = "mail"
	EntryTypeModel         = "model"
	EntryTypeNotification  = "notification"
	EntryTypeQuery         = "query"
	EntryTypeRedis         = "redis"
	EntryTypeRequest       = "request"
	EntryTypeScheduledTask = "scheduled_task"
	EntryTypeGate          = "gate"
	EntryTypeView          = "view"
	EntryTypeClientRequest = "client_request"
)
