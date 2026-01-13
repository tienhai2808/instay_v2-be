package consumer

import (
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"go.uber.org/zap"
)

type Consumer struct {
	log     *zap.Logger
	mqPro   port.MessageQueueProvider
	smtpPro port.SMTPProvider
}

func NewConsumer(
	log *zap.Logger,
	mqPro port.MessageQueueProvider,
	smtpPro port.SMTPProvider,
) *Consumer {
	return &Consumer{
		log,
		mqPro,
		smtpPro,
	}
}

func (c *Consumer) Start() {
	c.startEmailConsumer()
}
