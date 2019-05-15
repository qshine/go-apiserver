package service

import (
	"fmt"
	"go-apiserver/model"
	"go-apiserver/util"
	"sync" // 锁
)

func ListUser(username string, offset, limit int) ([]*model.UserInfo, uint64, error) {
	infos := make([]*model.UserInfo, 0)
	users, count, err := model.ListUser(username, offset, limit)
	if err != nil {
		return nil, count, err
	}

	ids := []uint64{}
	for _, user := range users {
		ids = append(ids, user.Id)
	}

	wg := sync.WaitGroup{}
	userList := model.UserList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint64]*model.UserInfo, len(users)),
	}

	// 分别创建error和bool通道
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// 开启并行处理, 降低api延迟
	for _, u := range users {
		// 添加任务个数
		wg.Add(1)
		go func(u *model.UserModel) {
			// 任务做完后任务数需要 -1 操作
			defer wg.Done()

			// 获取短id
			shortId, err := util.GenShortId()
			if err != nil {
				errChan <- err
				return
			}

			userList.Lock.Lock()
			defer userList.Lock.Unlock()

			// 业务处理, 返回Hello
			userList.IdMap[u.Id] = &model.UserInfo{
				Id:        u.Id,
				Username:  u.Username,
				SayHello:  fmt.Sprintf("Hello %s", shortId),
				Password:  u.Password,
				CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
		}(u)
	}

	go func() {
		// 等待所有任务完成, 最后关闭finished通道
		wg.Wait()
		close(finished)
	}()

	// select 只能接收一次
	select {
	case <-finished:
	case err := <-errChan:
		return nil, count, err
	}

	for _, id := range ids {
		infos = append(infos, userList.IdMap[id])
	}

	return infos, count, nil
}
