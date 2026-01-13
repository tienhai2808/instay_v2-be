package consumer

import (
	"encoding/json"

	"github.com/InstaySystem/is_v2-be/internal/application/dto"
	"github.com/InstaySystem/is_v2-be/pkg/constants"
	"go.uber.org/zap"
)

func (c *Consumer) startEmailConsumer() {
	go c.startSendAuthEmail()
}

func (c *Consumer) startSendAuthEmail() {
	if err := c.mqPro.ConsumeMessage(constants.QueueNameAuthEmail, constants.ExchangeEmail, constants.RoutingKeyAuthEmail, func(body []byte) error {
		var emailMsg dto.AuthEmailMessage
		if err := json.Unmarshal(body, &emailMsg); err != nil {
			c.log.Error("json unmarshal auth email message failed", zap.Error(err))
			return err
		}

		if err := c.smtpPro.AuthEmail(emailMsg.To, emailMsg.Subject, emailMsg.Otp); err != nil {
			c.log.Error("send auth email failed", zap.Error(err))
			return err
		}

		return nil
	}); err != nil {
		c.log.Error("start consumer send auth email failed", zap.Error(err))
	}
}
