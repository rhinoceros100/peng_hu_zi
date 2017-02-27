package main

import (
	"bufio"
	"os"
	"peng_hu_zi/log"
	"strconv"
	"peng_hu_zi/game_server/playing"
	"peng_hu_zi/util"
	"peng_hu_zi/game_server/card"
	"strings"
)

func help() {
	/*
	OperateEnterRoom	OperateType = iota + 1
	OperateLeaveRoom

	OperateDropCard
	OperateChiCard
	OperatePengCard
	OperateSaoCard
	OperatePaoCard
	OperateTiLongCard
	OperateHu
	 */
	log.Debug("h")
	log.Debug("exit")
	log.Debug("mycards")
	log.Debug(playing.OperateEnterRoom, int(playing.OperateEnterRoom))
	log.Debug(playing.OperateLeaveRoom, int(playing.OperateLeaveRoom))
	log.Debug(playing.OperateDropCard, int(playing.OperateDropCard), "1(small) 7(cardno)")
	log.Debug(playing.OperateChiCard, int(playing.OperateChiCard), "1(small) 2(cardno_1) (cardno_2)")
	log.Debug(playing.OperatePengCard, int(playing.OperatePengCard))
	//log.Debug(playing.OperateSaoCard, " : OperateSaoCard")
	//log.Debug(playing.OperatePaoCard, " : OperatePaoCard")
	//log.Debug(playing.OperateTiLongCard, " : OperateTiLongCard")
	//log.Debug(playing.OperateHu, " : OperateHu")
}

type PlayerObserver struct {}
func (ob *PlayerObserver) OnMsg(player *playing.Player, msg *playing.Message) {
	log.Debug(player, "receive msg", msg)
	log.Debug(player, "playingcards :", player.GetPlayingCards())
}

func main() {
	running := true

	//init room
	conf := playing.NewRoomConfig()
	err := conf.Init("./playing/room_config.json")
	if err != nil {
		log.Debug("room config init", err)
		return
	}
	room := playing.NewRoom(util.UniqueId(), conf)
	room.Start()

	robots := []*playing.Player{
		playing.NewPlayer(1),
		playing.NewPlayer(2),
		playing.NewPlayer(3),
	}

	for _, robot := range robots {
		robot.OperateEnterRoom(room)
		robot.AddObserver(&PlayerObserver{})
	}

	curPlayer := playing.NewPlayer(4)
	curPlayer.AddObserver(&PlayerObserver{})

	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		cmd := string(data)
		if cmd == "h" {
			help()
		} else if cmd == "exit" {
			return
		} else if cmd == "mycards" {
			log.Debug(curPlayer.GetPlayingCards())
		}
		splits := strings.Split(cmd, " ")
		c, _ := strconv.Atoi(splits[0])
		switch playing.OperateType(c) {
		case playing.OperateEnterRoom:
			curPlayer.OperateEnterRoom(room)
		case playing.OperateLeaveRoom:
			curPlayer.OperateLeaveRoom()
		case playing.OperateDropCard:
			card := &card.Card{}
			card.CardType, _ = strconv.Atoi(splits[1])
			card.CardNo, _ = strconv.Atoi(splits[2])
			curPlayer.OperateDropCard(card)
		case playing.OperateChiCard:
			card1 := &card.Card{}
			card1.CardType, _ = strconv.Atoi(splits[1])
			card1.CardNo, _ = strconv.Atoi(splits[2])
			card2 := &card.Card{CardType:card1.CardType}
			card2.CardNo, _ = strconv.Atoi(splits[3])
			cards := card.NewCards(card1, card2)
			curPlayer.OperateChiCard(cards)
		case playing.OperatePengCard:
			curPlayer.OperatePengCard()
		}
	}
}