package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/redis/go-redis/v9"
	"github.com/suricat89/rinha-2024-q1/src/interfaces"
)

type CacheRepository struct {
	Rdb *redis.Client
}

func NewCacheRepository(Rdb *redis.Client) interfaces.CacheRepository {
	return &CacheRepository{Rdb}
}

func (r *CacheRepository) WaitForCustomerLock(customerId int, reqUuid string) error {
	customerKey := fmt.Sprintf("customerTransaction::%d", customerId)
	ctx := context.Background()

	log.Debugf("[%s] Waiting to lock customer '%d'", reqUuid, customerId)
	defer log.Debugf("[%s] Customer '%d' locked with success", reqUuid, customerId)

	chFinished := make(chan bool)
	chErr := make(chan error)
	go func(chFinished chan bool, chErr chan error) {
		for {
			after := time.After(5 * time.Millisecond)

			isReady, err := r.checkCustomerLock(ctx, reqUuid, customerKey)
			if err != nil {
				chErr <- err
				return
			}
			if isReady {
				chFinished <- true
				return
			}

			<-after
		}
	}(chFinished, chErr)

  for {
    select {
    case err := <-chErr:
      return err
    case <-chFinished:
      return nil
    }  
  }
}

func (r *CacheRepository) getCustomerKey(customerId int) string {
	return fmt.Sprintf("customerTransaction::%d", customerId)
}

// Returns `true` if the customer is available, locked for this host
// and ready for use, or `false` if it is locked for another host
func (r *CacheRepository) checkCustomerLock(ctx context.Context, reqUuid string, customerKey string) (bool, error) {
	val, err := r.Rdb.Get(ctx, customerKey).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}

	if val == "" || err == redis.Nil {
		err = r.Rdb.Set(ctx, customerKey, reqUuid, 0).Err()
		if err != nil {
			return false, err
		}
		log.Debugf("[%s] Set key '%s'", reqUuid, customerKey)
		time.Sleep(time.Millisecond * 5)
		return r.checkCustomerLock(ctx, reqUuid, customerKey)
	}
	if val == reqUuid {
		log.Debugf("[%s] Got value '%s' from key '%s'", reqUuid, val, customerKey)
		return true, nil
	}
	return false, nil
}

func (r *CacheRepository) UnlockCustomer(customerId int) error {
	customerKey := r.getCustomerKey(customerId)
	ctx := context.Background()

	return r.Rdb.Set(ctx, customerKey, "", 0).Err()
}
