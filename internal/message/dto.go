package message

type GetSendMessagesQuery struct {
	Limit  int `json:"limit" validate:"required,min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}
