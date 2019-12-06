package bingosvc

import (
	"encoding/json"
	"errors"
	"github.com/shihtzu-systems/bingo/pkg/bingo"
	"github.com/shihtzu-systems/redix"
	"math/rand"
	"time"
)

func SaveBoardId(sessionId, boardId string, redis redix.Redis) error {
	if err := redis.Connect(); err != nil {
		return err
	}
	defer func() { _ = redis.Disconnect() }()

	_, err := redis.Set(sessionId+":board:id", []byte(boardId))
	return err
}

func GetBoardId(sessionId string, redis redix.Redis) (string, error) {
	if err := redis.Connect(); err != nil {
		return "", err
	}
	defer func() { _ = redis.Disconnect() }()

	return redis.Get(sessionId + ":board:id")
}

func NewBoard(id string, boxes bingo.Boxes) (out bingo.Board, err error) {
	if len(boxes) < 25 {
		return bingo.Board{}, errors.New("not enough boxes to build a bingo board")
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
	}, nil
}

func SaveBoard(board bingo.Board, redis redix.Redis) error {
	if err := redis.Connect(); err != nil {
		return err
	}
	defer func() { _ = redis.Disconnect() }()

	_, err := redis.Set(board.Id, board.PrettyJson())
	return err

}

func BoardExists(id string, redis redix.Redis) (bool, error) {
	if err := redis.Connect(); err != nil {
		return false, err
	}
	defer func() { _ = redis.Disconnect() }()

	return redis.Exists(id)
}

func GetBoard(id string, redis redix.Redis) (out bingo.Board, err error) {
	if err := redis.Connect(); err != nil {
		return bingo.Board{}, err
	}
	defer func() { _ = redis.Disconnect() }()

	bjson, err := redis.Get(id)
	if err != nil {
		return bingo.Board{}, err
	}
	if err := json.Unmarshal([]byte(bjson), &out); err != nil {
		return bingo.Board{}, err
	}
	return out, nil
}

func CheckForBingo(board *bingo.Board) bool {
	board.Bingos = []bingo.Bingo{}
	for _, b := range bingos {
		if checkBoxes(b.B, board.B) &&
			checkBoxes(b.I, board.I) &&
			checkBoxes(b.N, board.N) &&
			checkBoxes(b.G, board.G) &&
			checkBoxes(b.O, board.O) {
			board.Bingos = append(board.Bingos, b)
		}
	}
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
