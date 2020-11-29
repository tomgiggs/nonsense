package common

import (
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"time"
)

var ErrTimeOut = errors.New("get uid timeout")

type Uid struct {
	db         *sql.DB    // 数据库连接
	appId string     // 业务id
	ch         chan int64 // id缓冲池
	min, max   int64      // id段最小值，最大值
	step  int64
}

// 从数据库分配uid到缓存用于自增
func NewUid(db *sql.DB, appId string, size int64) *Uid {
	uid := &Uid{
		db:         db,
		appId: appId,
		ch:         make(chan int64, size),
		step: size,
	}
	go uid.productId()
	return uid
}

func (u *Uid) GetDeviceUid() (int64, error) {
	select {
	case <-time.After(1 * time.Second):
		return 0, ErrTimeOut
	case uid := <-u.ch:
		return uid, nil
	}
}

// 生产id，当ch达到最大容量时，这个方法会阻塞，直到ch中的id被消费
func (u *Uid) productId() {
	u.reLoad()

	for {
		if u.min >= u.max {
			u.reLoad()
		}
		u.min++
		u.ch <- u.min
	}
}

func (u *Uid) reLoad() (err error) {
	for {
		err = u.AllocateFromDB()
		if err == nil {
			return
		}
		Logger.Error("deviceUid",zap.Any("reload error", err))
		time.Sleep(time.Second)// 查询异常，一秒后重试
	}
}

// 从数据库获取id段
func (u *Uid) AllocateFromDB() error {

	tx, err := u.db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	row := tx.QueryRow("SELECT max_id,step FROM uid WHERE app_id = ? FOR UPDATE", u.appId)
	err = row.Scan(&u.max, &u.step)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE uid SET max_id = ? WHERE app_id = ?", u.max+u.step, u.appId)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	u.min = u.max
	u.max = u.max + u.step
	return nil
}
