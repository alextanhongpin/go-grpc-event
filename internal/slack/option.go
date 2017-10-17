package slack

type Options struct {
	webhookURL string
	channel    string
	username   string
	iconEmoji  string
}

type Option func(*Options)

func Username(username string) Option {
	return func(o *Options) {
		o.username = username
	}
}

func WebhookURL(url string) Option {
	return func(o *Options) {
		o.webhookURL = url
	}
}

func Channel(channel string) Option {
	return func(o *Options) {
		o.channel = channel
	}
}

func IconEmoji(iconEmoji string) Option {
	return func(o *Options) {
		o.iconEmoji = iconEmoji
	}
}
