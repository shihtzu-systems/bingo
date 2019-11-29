package bingosvc

import (
	"encoding/json"
	"github.com/shihtzu-systems/bingo/pkg/bingo"
	"github.com/shihtzu-systems/redix"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func NewBoard(id string, boxes bingo.Boxes) bingo.Board {
	if len(boxes) < 25 {
		log.Fatal("not enough boxes to build a bingo board")
	}

	// shuffle
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(boxes), func(i, j int) { boxes[i], boxes[j] = boxes[j], boxes[i] })

	return bingo.Board{
		Id: id,
		B:  boxes[:5],
		I:  boxes[5:10],
		N:  boxes[10:15],
		G:  boxes[15:20],
		O:  boxes[20:25],
	}
}

func SaveBoard(board bingo.Board, redis redix.Redis) {
	redis.Connect()
	defer redis.Disconnect()

	redis.Set(board.Id, board.PrettyJson())
}

func BoardExists(id string, redis redix.Redis) bool {
	redis.Connect()
	defer redis.Disconnect()

	return redis.Exists(id)
}

func GetBoard(id string, redis redix.Redis) (out bingo.Board) {
	redis.Connect()
	defer redis.Disconnect()

	bjson := redis.Get(id)
	if err := json.Unmarshal([]byte(bjson), &out); err != nil {
		log.Fatal(err)
	}
	return out
}
