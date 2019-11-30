package bingosvc

import (
	"encoding/json"
	"github.com/shihtzu-systems/bingo/pkg/bingo"
	"github.com/shihtzu-systems/redix"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func SaveBoardId(sessionId, boardId string, redis redix.Redis) {
	redis.Connect()
	defer redis.Disconnect()

	redis.Set(sessionId+":board:id", []byte(boardId))
}

func GetBoardId(sessionId string, redis redix.Redis) string {
	redis.Connect()
	defer redis.Disconnect()

	return redis.Get(sessionId + ":board:id")
}

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

func CheckForBingo(board *bingo.Board) bool {
	board.Bingos = []bingo.Bingo{}
	for _, b := range bingos {
		log.Debugf("[%s] [%s]", b.Type, b.Id)
		if checkBoxes(b.B, board.B) &&
			checkBoxes(b.I, board.I) &&
			checkBoxes(b.N, board.N) &&
			checkBoxes(b.G, board.G) &&
			checkBoxes(b.O, board.O) {
			log.Info("Bingo!")
			board.Bingos = append(board.Bingos, b)
		}
	}
	log.Debug("bingos: ", len(board.Bingos))
	board.Bingoed = len(board.Bingos) > 0
	return board.Bingoed
}

func checkBoxes(bingos []bool, boxes bingo.Boxes) bool {
	for i, b := range bingos {
		if b == false {
			continue
		}
		if !boxes[i].Marked {
			return false
		}
	}
	return true
}
