package slack

// New returns a new pointer to the slack webhook struct
func New(opts ...Option) *SlackWebhook {
	options := Options{
		webhookURL: "",
		channel:    "#test",
		username:   "bot",
		iconEmoji:  ":ghost:",
	}
	for _, o := range opts {
		o(&options)
	}
	slack := SlackWebhook{
		Data: Slack{
			Channel:   options.channel,
			IconEmoji: options.iconEmoji,
			Username:  options.username,
		},
		WebhookURL: options.webhookURL,
	}
	return &slack
}
