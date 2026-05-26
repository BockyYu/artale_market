package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func itemViewKey(date string) string {
	return fmt.Sprintf("item_views:%s", date)
}

// FrequentEntry Redis ZSet 查詢結果的單筆紀錄
type FrequentEntry struct {
	ItemID uint // 商品 ID
	Count  int  // 被查詢次數
}

type QueryRepository interface {
	RecordQuery(userID string, itemID uint) error
	GetFrequent(userID string, limit int) ([]FrequentEntry, error)
	RecordItemView(itemID uint, date string) error
	GetAllItemViews(date string) (map[uint]int, error)
}

type queryRepo struct {
	rdb *redis.Client
}

func NewQueryRepository(rdb *redis.Client) QueryRepository {
	return &queryRepo{rdb: rdb}
}

func (r *queryRepo) key(userID string) string {
	return fmt.Sprintf("user:%s:item_queries", userID)
}

func (r *queryRepo) RecordQuery(userID string, itemID uint) error {
	return r.rdb.ZIncrBy(context.Background(), r.key(userID), 1, strconv.Itoa(int(itemID))).Err()
}

func (r *queryRepo) RecordItemView(itemID uint, date string) error {
	return r.rdb.HIncrBy(context.Background(), itemViewKey(date), strconv.Itoa(int(itemID)), 1).Err()
}

func (r *queryRepo) GetAllItemViews(date string) (map[uint]int, error) {
	raw, err := r.rdb.HGetAll(context.Background(), itemViewKey(date)).Result()
	if err != nil {
		return nil, err
	}
	result := make(map[uint]int, len(raw))
	for k, v := range raw {
		id, _ := strconv.ParseUint(k, 10, 64)
		count, _ := strconv.Atoi(v)
		result[uint(id)] = count
	}
	return result, nil
}

func (r *queryRepo) GetFrequent(userID string, limit int) ([]FrequentEntry, error) {
	results, err := r.rdb.ZRevRangeWithScores(context.Background(), r.key(userID), 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}
	entries := make([]FrequentEntry, 0, len(results))
	for _, z := range results {
		id, _ := strconv.ParseUint(fmt.Sprintf("%v", z.Member), 10, 64)
		entries = append(entries, FrequentEntry{ItemID: uint(id), Count: int(z.Score)})
	}
	return entries, nil
}
