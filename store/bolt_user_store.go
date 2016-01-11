package store

import (
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
	us.usersBucket, _ = BoltDBStore.db.New([]byte(USERS))
	us.usersByLoginBucket, _ = BoltDBStore.db.New([]byte(USERS_BY_LOGIN))

	return us
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

		result.Date = user

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
		return user
	}
}

func (us BoltUserStore) Get(id string) *User {
	user := getJson(id)
	return model.UserFromJson(string(user))
}

func (us BoltUserStore) getIdByLogin(login string) string {
	userID, err := us.usersByLoginBucket.Get([]byte(login))
	if err != nil {
		panic(err)
	} else {
		return userID
	}
}

func (us BoltUserStore) GetByLogin(login string) *User {
	userId := us.getIdByLogin(login)
	user := getJson(userId)
	return model.UserFromJson(string(user))
}

func (us BoltUserStore) isLoginTaken(login string) {
	return getByLoginJSON(login) != nil
}
