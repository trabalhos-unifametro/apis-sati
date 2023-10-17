package emails

import (
	"apis-sati/utils"
	"github.com/keighl/mandrill"
)

func MountEmail(body, title, email string) (responses []*mandrill.Response, err error) {
	client := mandrill.ClientWithKey(utils.MANDRILL)
	message := &mandrill.Message{}
	message.AddRecipient(utils.CheckToSend(email), "S.A.T.I", "to")
	message.FromEmail = utils.FROM
	message.FromName = "S.A.T.I"
	message.Subject = title
	message.HTML = body
	return client.MessagesSend(message)
}
