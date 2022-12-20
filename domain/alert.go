package domain

type SendAlert interface {
	Send() error
}

type DiscordField struct {
	Name   *string `json:"name,omitempty"`
	Value  *string `json:"value,omitempty"`
	Inline *bool   `json:"inline,omitempty"`
}

type DiscordFooter struct {
	Value   *string `json:"text,omitempty"`
	IconUrl *string `json:"icon_url,omitempty"`
}

type DiscordAuthor struct {
	Name    *string `json:"name,omitempty"`
	Url     *string `json:"url,omitempty"`
	IconUrl *string `json:"icon_url,omitempty"`
}

type DiscordImage struct {
	Url *string `json:"url,omitempty"`
}

type DiscordEmbed struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Url         *string `json:"url,omitempty"`
	Color       *int32  `json:"color,omitempty"`
	//Timestamp   time.Time      `json:"timestamp,omitempty"`
	Fields    []*DiscordField `json:"fields,omitempty"`
	Author    DiscordAuthor   `json:"author,omitempty"`
	Image     DiscordImage    `json:"image,omitempty"`
	Thumbnail DiscordImage    `json:"thumbnail,omitempty"`
	Footer    *DiscordFooter  `json:"footer,omitempty"`
}

// Root object for Discord Webhook messages
type DiscordMessage struct {
	Username *string         `json:"username,omitempty"`
	Content  *string         `json:"content,omitempty"`
	Embeds   *[]DiscordEmbed `json:"embeds,omitempty"`
}
