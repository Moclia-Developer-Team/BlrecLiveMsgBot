package bot

import zero "github.com/wdvxdr1123/ZeroBot"

func GetBotList() []int64 {
	var bots []int64
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		bots = append(bots, id)
		return true
	})
	return bots
}
