package services

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func ListResourceReactions(tx *gorm.DB) (map[string]int64, error) {
	var reactions []struct {
		Symbol string
		Count  int64
	}

	if err := tx.Model(&models.Reaction{}).
		Select("symbol, COUNT(id) as count").
		Group("symbol").
		Scan(&reactions).Error; err != nil {
		return map[string]int64{}, err
	}

	return lo.SliceToMap(reactions, func(item struct {
		Symbol string
		Count  int64
	},
	) (string, int64) {
		return item.Symbol, item.Count
	}), nil
}

func BatchListResourceReactions(tx *gorm.DB) (map[uint]map[string]int64, error) {
	var reactions []struct {
		ArticleID uint
		Symbol    string
		Count     int64
	}

	reactInfo := map[uint]map[string]int64{}
	if err := tx.Model(&models.Reaction{}).
		Select("article_id, symbol, COUNT(id) as count").
		Group("article_id, symbol").
		Scan(&reactions).Error; err != nil {
		return reactInfo, err
	}

	for _, info := range reactions {
		if _, ok := reactInfo[info.ArticleID]; !ok {
			reactInfo[info.ArticleID] = make(map[string]int64)
		}
		reactInfo[info.ArticleID][info.Symbol] = info.Count
	}

	return reactInfo, nil
}
