package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

// StickersUpdate is dispatched when a guild's stickers are updated.
// This event does not depend on a cache like StickerCreate, StickerUpdate or StickerDelete.
type StickersUpdate struct {
	*GenericEvent
	gateway.EventGuildStickersUpdate
}

// GenericSticker is called upon receiving StickerCreate , StickerUpdate or StickerDelete (requires gateway.IntentGuildEmojisAndStickers)
type GenericSticker struct {
	*GenericEvent
	GuildID snowflake.ID
	Sticker discord.Sticker
}

// StickerCreate indicates that a new discord.Sticker got created in a discord.Guild (requires gateway.IntentGuildEmojisAndStickers)
type StickerCreate struct {
	*GenericSticker
}

// StickerUpdate indicates that a discord.Sticker got updated in a discord.Guild (requires gateway.IntentGuildEmojisAndStickers)
type StickerUpdate struct {
	*GenericSticker
	OldSticker discord.Sticker
}

// StickerDelete indicates that a discord.Sticker got deleted in a discord.Guild (requires gateway.IntentGuildEmojisAndStickers)
type StickerDelete struct {
	*GenericSticker
}
