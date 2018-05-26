package rate_test

import (
	"errors"
	"github.com/fossapps/starter/rate"
	"github.com/fossapps/starter/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLimiter_CountErrorWhileRemovingDecayed(t *testing.T) {
	expect := assert.New(t)
	ctrl := gomock.NewController(t)
	mockRedis := mocks.NewMockRedisClient(ctrl)
	mockRedis.EXPECT().ZRemRangeByScore("key", "0", gomock.Any()).Return(int64(0), errors.New("error"))

	limiter := rate.Limiter{
		Decay:       5 * time.Second,
		Limit:       5,
		RedisClient: mockRedis,
	}
	card, err := limiter.Count("key")
	expect.Equal(int64(-1), card)
	expect.NotNil(err)
}

func TestLimiter_Count(t *testing.T) {
	expect := assert.New(t)
	ctrl := gomock.NewController(t)
	mockRedis := mocks.NewMockRedisClient(ctrl)
	mockRedis.EXPECT().ZCard("key").Return(int64(2), nil)
	mockRedis.EXPECT().ZRemRangeByScore("key", "0", gomock.Any()).Return(int64(0), nil)
	limiter := rate.Limiter{
		Decay:       5 * time.Second,
		Limit:       5,
		RedisClient: mockRedis,
	}
	card, err := limiter.Count("key")
	expect.Equal(card, int64(2))
	expect.Nil(err)
}

func TestLimiter_HitErrorWhileAdding(t *testing.T) {
	expect := assert.New(t)
	ctrl := gomock.NewController(t)
	mockRedis := mocks.NewMockRedisClient(ctrl)
	mockRedis.EXPECT().ZAdd("key", gomock.Any()).Return(int64(0), errors.New("error"))

	limiter := rate.Limiter{
		Decay:       5 * time.Second,
		Limit:       5,
		RedisClient: mockRedis,
	}
	card, err := limiter.Hit("key")
	expect.Equal(int64(-1), card)
	expect.NotNil(err)
}

func TestLimiter_HitErrorSettingExpiry(t *testing.T) {
	expect := assert.New(t)
	ctrl := gomock.NewController(t)
	mockRedis := mocks.NewMockRedisClient(ctrl)
	mockRedis.EXPECT().ZAdd("key", gomock.Any()).Return(int64(0), nil)
	mockRedis.EXPECT().Expire("key", gomock.Any()).Return(false, errors.New("error"))

	limiter := rate.Limiter{
		Decay:       5 * time.Second,
		Limit:       5,
		RedisClient: mockRedis,
	}
	card, err := limiter.Hit("key")
	expect.Equal(int64(-1), card)
	expect.NotNil(err)
}

func TestLimiter_Hit(t *testing.T) {
	expect := assert.New(t)
	ctrl := gomock.NewController(t)
	mockRedis := mocks.NewMockRedisClient(ctrl)
	mockRedis.EXPECT().ZAdd("key", gomock.Any()).Return(int64(0), nil)
	mockRedis.EXPECT().ZRemRangeByScore("key", "0", gomock.Any()).Return(int64(0), nil)
	mockRedis.EXPECT().Expire("key", gomock.Any()).Return(true, nil)
	mockRedis.EXPECT().ZCard("key").Return(int64(2), nil)

	limiter := rate.Limiter{
		Decay:       5 * time.Second,
		Limit:       5,
		RedisClient: mockRedis,
	}
	card, err := limiter.Hit("key")
	expect.Equal(int64(2), card)
	expect.Nil(err)
}
