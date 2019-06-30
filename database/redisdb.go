/*
	This is very simple implementation of finite state machine with additional data storage
	Everything stored in fast mem-cache database REDIS https://redis.io

	All data you store must be able to be converted to json datatype (string, slice, etc)

	Data stored that way:
		{"RECORD-PREFIX:user_id:chat_id" : `{"__state__": "some state", "some_key": "some_val"}`

	Keys are unique per two

	Notes:
		- avoid using __state__ as a key for any data you store
		- redis is not persistent
		- check out default constants
*/

package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var Client *redis.Client

const (
	InitialState          string = "~"         // initial state will be "~"
	DefaultExpirationTime int64  = 0           // records won't expire
	RecordPrefix          string = "user_data" // prefix for redis records
)

// redis db setup
func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:         ":6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
}

func Scan(prefix string, count int64) []string {
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = Client.Scan(cursor, fmt.Sprintf("%s*", prefix), count).Result()
		if err != nil {
			panic(err)
		}
		if cursor == 0 {
			var res []string
			for key := range keys {
				res = append(res, keys[key])
			}
			return res
		}
	}
}

// base set reimplementation with err handling
func Set(key string, value string) {
	err := Client.Set(key, value, time.Duration(DefaultExpirationTime)).Err()
	if err != nil {
		panic(err)
	}
}

// get data from key
// use `create=true` case you want to create new bucket with `key: {"__state__": INITIAL_VALUE}`
// if it was not found by `key`
func Get(key string, create bool) map[string]interface{} {
	val, err := Client.Get(key).Result()

	var result map[string]interface{}

	if err == redis.Nil {
		if create {
			result = map[string]interface{}{
				"__state__": InitialState,
			}
			Set(key,
				fmt.Sprintf(`{"__state__": "%s"}`, InitialState)) // noqa

		} else {
			return result
		}
	} else if err != nil {
		panic(err)
	} else {
		err := json.Unmarshal([]byte(val), &result)
		if err != nil {
			panic(err)
		}
	}

	return result
}

// will create new record if can't get one with particular key
func GetOrCreate(key string) map[string]interface{} {
	return Get(key, true)
}

// get user identifier by most usable telegram update
// TODO: make all updates identifier
func Key(update tgbotapi.Update) string {
	var FID int
	var CHID int64

	switch {
	case update.Message != nil:
		FID = update.Message.From.ID
		CHID = update.Message.Chat.ID
	case update.CallbackQuery != nil:
		FID = update.CallbackQuery.From.ID
		CHID = update.CallbackQuery.Message.Chat.ID
	case update.ChannelPost != nil:
		CHID = update.ChannelPost.Chat.ID
		FID = int(CHID)
	case update.ChosenInlineResult != nil:
		FID = update.ChosenInlineResult.From.ID
		CHID = int64(CHID)
	case update.EditedMessage != nil:
		FID = update.EditedMessage.From.ID
		CHID = update.EditedMessage.Chat.ID
	}

	return fmt.Sprintf(RecordPrefix+":%d:%d", FID, CHID)
}

// update data bucket
func UpdateData(update tgbotapi.Update, data map[string]interface{}) map[string]interface{} {

	key := Key(update)
	oldData := GetOrCreate(key)

	for k, v := range data {
		oldData[k] = v
	}

	body, err := json.Marshal(oldData)
	if err != nil {
		panic(err)
	}

	Set(key, string(body))
	return oldData
}

// sets new state
// TODO: discuss if __state__ key must be array of database and renamed to __states__
func UpdateState(update tgbotapi.Update, state string) map[string]interface{} {
	key := Key(update)
	oldData := GetOrCreate(key)

	oldData["__state__"] = state

	body, err := json.Marshal(oldData)
	if err != nil {
		panic(err)
	}

	Set(key, string(body))
	return oldData
}

func GetCurrentState(update tgbotapi.Update) string {
	/*
		GetCurrentState used to get user's particular state
		:param: update - telegramBotUpdate
		:return: string state if declared
	*/
	key := Key(update)
	oldData := GetOrCreate(key)
	log.Println(oldData)
	return oldData["__state__"].(string)
}

// get user's saved data if no user found - will create new bucket
func GetData(update tgbotapi.Update) map[string]interface{} {
	key := Key(update)
	data := GetOrCreate(key)
	return data
}

// get all stored records
func GetAllKeys(limit int64) []string {
	scanned := Scan(RecordPrefix+":", limit)
	return scanned
}

// remove all stored data and database (flush db)
func Flush() {
	Client.FlushDB()
}
