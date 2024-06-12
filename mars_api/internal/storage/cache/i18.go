package cache

import (
	"context"

	"github.com/google/common/pkg/models"
	"github.com/rs/zerolog/log"
)

func PhraseError(lang string) string {
	tips := ""
	switch lang {
	case "zh_cn":
		tips = "获取词组失败。"
		break
	case "zh_hk":
		tips = "獲取片語失敗。"
		break
	case "ja_jp":
		tips = "フレーズの取得に失敗しました"
		break
	case "ko_kr":
		tips = "구문 가 져 오 는 데 실 패 했 습 니 다."
		break
	default:
		tips = "Failed to get phrase."
	}
	return tips
}

func (c *Cache) initI18() {
	// 获取上下文
	ctx := context.TODO()

	// 匹配以 "lang" 开头的键
	iter := c.conn.Scan(ctx, 0, "lang:*", 0).Iterator()

	// 遍历匹配的键并逐个删除
	for iter.Next(ctx) {
		err := c.conn.Del(ctx, iter.Val()).Err()
		if err != nil {
			log.Error().Msgf("Error deleting key %s: %s\n", iter.Val(), err)
		} else {
			log.Info().Msgf("Deleted key: %s\n", iter.Val())
		}
	}

	if err := iter.Err(); err != nil {
		log.Error().Msgf("Error iterating keys: %s\n", err)
	}

	langs := models.LangPhrase{}
	langsList := langs.GetByModular("api", c.db)
	for _, v := range langsList {
		key := "lang:" + v.LanguageKey
		c.conn.HSet(context.TODO(), key, "zh_cn", v.ZhCN)
		c.conn.HSet(context.TODO(), key, "zh_hk", v.ZhHk)
		c.conn.HSet(context.TODO(), key, "en_us", v.EnUs)
		c.conn.HSet(context.TODO(), key, "ko_kr", v.KoKr)
		c.conn.HSet(context.TODO(), key, "ja_jp", v.JaJp)
	}

	log.Info().Msg("I18 cache initialized")
}
