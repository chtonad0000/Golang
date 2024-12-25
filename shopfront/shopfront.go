package shopfront

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type RedisCounters struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) Counters {
	return &RedisCounters{
		rdb: rdb,
	}
}

func (r *RedisCounters) RecordView(ctx context.Context, id ItemID, userID UserID) error {
	key := r.generateKey(id, userID)

	err := r.rdb.SetNX(ctx, key, true, 0).Err()
	if err != nil {
		return err
	}

	viewCountKey := "view_count:" + strconv.FormatInt(int64(id), 10)
	_, err = r.rdb.Incr(ctx, viewCountKey).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisCounters) GetItems(ctx context.Context, ids []ItemID, userID UserID) ([]Item, error) {
	var items []Item
	pipe := r.rdb.Pipeline()

	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = r.generateKey(id, userID)
	}
	cmd := pipe.MGet(ctx, keys...)

	viewCountKeys := make([]string, len(ids))
	for i, id := range ids {
		viewCountKeys[i] = "view_count:" + strconv.FormatInt(int64(id), 10)
	}
	viewCountCmd := pipe.MGet(ctx, viewCountKeys...)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	viewedResults := cmd.Val()
	viewCountResults := viewCountCmd.Val()

	for i, result := range viewedResults {
		viewed := false
		if result != nil {
			viewed = true
		}

		viewCount := 0
		if viewCountResults[i] != nil {
			viewCount, _ = strconv.Atoi(viewCountResults[i].(string))
		}

		items = append(items, Item{
			ViewCount: viewCount,
			Viewed:    viewed,
		})
	}

	return items, nil
}

func (r *RedisCounters) generateKey(id ItemID, userID UserID) string {
	return "item:" + strconv.FormatInt(int64(id), 10) + ":user:" + strconv.FormatInt(int64(userID), 10)
}
