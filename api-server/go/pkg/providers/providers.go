package providers

import (
	"github.com/caellach/shorturl/api-server/go/pkg/providers/discord"
)

func init() {
	discord.Init()
}

func GetProvidersInfo() map[string]interface{} {
	return map[string]interface{}{
		"discord": discord.GetInfo(),
	}
}
