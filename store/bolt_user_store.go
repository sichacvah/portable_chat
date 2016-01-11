package store

import (
	"strings"

	"github.com/joyrexus/buckets"
	"github.com/sichacvah/portable_chat/model"
)

type BoltUserStore struct {
	usersBucket        *buckets.Bucket
	usersByLoginBucket *buckets.Bucket
}

const (
	USERS          = "users"
	USERS_BY_LOGIN = "users_by_login"
)

func NewBoltDbUserStore(boltStore *BoltDBStore) UserStore {
	us := &BoltUserStore{}
	us.usersBucket, _ = boltStore.db.New([]byte(USERS))
	us.usersByLoginBucket, _ = boltStore.db.New([]byte(USERS_BY_LOGIN))

	return us
}

func (us BoltUserStore) GetUsers() StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		items, err := us.usersBucket.Items()

		if err != nil {
			result.Err = model.NewAppError("BoltUserStore.GetUsers", err.Error(), "")
			storeChannel <- result
			close(storeChannel)
		}

		users := []string{}

		for _, item := range items {
			users = append(users, string(item.Value))
		}

		result.Data = users
		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (us BoltUserStore) Delete(userId string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		if len(userId) <= 0 {
			result.Err = model.NewAppError("BoltUserStore.Delete", "You must get userId in delete", "user_id = "+userId)
			storeChannel <- result
			close(storeChannel)
		}

		err := us.usersBucket.Delete([]byte(userId))
		if err != nil {
			result.Err = model.NewAppError("BoltUserStore.Delete", err.Error(), "")
			storeChannel <- result
			close(storeChannel)
		}

		result.Data = "ok"

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (us BoltUserStore) Update(user *model.User) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		if len(user.Password) > 0 {
			if user.PasswordConfirmation != user.Password {
				result.Err = model.NewAppError("BoltUserStore.Update", "When update password must be equal to passwordConfirmation", "user_id"+user.Id)
				storeChannel <- result
				close(storeChannel)
			}
		}
		if us.isLoginTaken(user.Login) {
			result.Err = model.NewAppError("BoltUserStore.Save", "User Login already taken", "user_login="+user.Login)
			storeChannel <- result
			close(storeChannel)
			return
		}

		user.PreSave()
		userJson := user.ToJson()
		us.usersBucket.Put([]byte(user.Id), []byte(userJson))
		us.usersByLoginBucket.Put([]byte(user.Login), []byte(user.Id))

		result.Data = user

		storeChannel <- result
		return
	}()

	return storeChannel
}

func (us BoltUserStore) Save(user *model.User) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		if len(user.Id) > 0 {
			result.Err = model.NewAppError("BoltUserStore.Save", "Must call update for exisiting user", "user_id="+user.Id)
			storeChannel <- result
			close(storeChannel)
			return
		}

		if us.isLoginTaken(user.Login) {
			result.Err = model.NewAppError("BoltUserStore.Save", "User Login already taken", "user_login="+user.Login)
			storeChannel <- result
			close(storeChannel)
			return
		}

		user.PreSave()
		userJSON := user.ToJson()
		us.usersBucket.Put([]byte(user.Id), []byte(userJSON))
		us.usersByLoginBucket.Put([]byte(user.Login), []byte(user.Id))

		result.Data = user

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (us BoltUserStore) getJson(id string) string {
	user, err := us.usersBucket.Get([]byte(id))
	if err != nil {
		panic(err)
	} else {
		return string(user)
	}
}

func (us BoltUserStore) get(id string) *model.User {
	user := string(us.getJson(id))
	return model.UserFromJson(strings.NewReader(user))
}

func (us BoltUserStore) Get(id string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}
		data := us.get(id)
		if data == nil {
			result.Err = model.NewAppError("BoltUserStore.Get", "User not found", "userId="+id)
			storeChannel <- result
			close(storeChannel)
			return
		}

		result.Data = data
		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}

func (us BoltUserStore) getIdByLogin(login string) string {
	userID, err := us.usersByLoginBucket.Get([]byte(login))
	if err != nil {
		panic(err)
	} else {
		return string(userID)
	}
}

func (us BoltUserStore) GetByLogin(login string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}
		data := us.getByLogin(login)
		if data == nil {
			result.Err = model.NewAppError("BoltUserStore.GetByLogin", "User not found", "userLogin="+login)
			storeChannel <- result
			close(storeChannel)
			return
		}

		result.Data = data
		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}

func (us BoltUserStore) getByLogin(login string) *model.User {
	userId := us.getIdByLogin(login)
	user := string(us.getJson(userId))
	return model.UserFromJson(strings.NewReader(user))
}

func (us BoltUserStore) isLoginTaken(login string) bool {
	return us.getIdByLogin(login) != ""
}
