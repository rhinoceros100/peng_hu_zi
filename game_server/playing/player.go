package playing

import (
	"peng_hu_zi/game_server/card"
	"peng_hu_zi/log"
	"time"
	"fmt"
)

type PlayerObserver interface {
	OnMsg(player *Player, msg *Message)
}

type Player struct {
	id				uint64			// 玩家id
	room			*Room			// 玩家所在的房间

	playingCards 	*card.PlayingCards	//玩家手上的牌

	pengFromPlayer 	map[int64]*Player			//碰的牌来源于 int64:Card.MakeKey()
	paoFromPlayers  map[int64]*Player			//跑的牌来源于 int64:Card.MakeKey()


	observers	 []PlayerObserver
	operateCh		chan *Operate
}

func NewPlayer(id uint64) *Player {
	player :=  &Player{
		id:		id,
		playingCards:	card.NewPlayingCards(),
		pengFromPlayer: make(map[int64]*Player),
		paoFromPlayers: make(map[int64]*Player),
		observers:	make([]PlayerObserver, 0),
		operateCh: make(chan *Operate),
	}
	return player
}

func (player *Player) GetId() uint64 {
	return player.id
}

func (player *Player) Reset() {
	log.Debug(player,"Player.Reset")
	player.playingCards.Reset()
	for key := range player.paoFromPlayers {
		delete(player.paoFromPlayers, key)
	}

	for key := range player.pengFromPlayer {
		delete(player.pengFromPlayer, key)
	}
}

func (player *Player) AddObserver(ob PlayerObserver) {
	player.observers = append(player.observers, ob)
}

func (player *Player) AddCard(card *card.Card) {
	log.Debug(player, "Player.AddCard :", card)
	player.playingCards.AddCard(card)
}


/*	计算指定的牌可以吃牌的组合
*/
func (player *Player) ComputeChiGroup(card *card.Card) []*card.Cards {
	return player.playingCards.ComputeChiGroup(card)
}

func (player *Player) CanPeng(card *card.Card) bool {
	return player.playingCards.CanPeng(card)
}

//检查是否能吃
func (player *Player) CanChi(whatCard *card.Card, whatGroup *card.Cards) bool {
	return player.playingCards.CanChi(whatCard, whatGroup)
}

//检查是否能扫
func (player *Player) CanSao(whatCard *card.Card) bool {
	return player.playingCards.CanSao(whatCard)
}

//检查是否能跑
func (player *Player) CanPao(whatCard *card.Card) bool {
	return player.playingCards.CanPao(whatCard)
}

//检查是否能提龙
func (player *Player) CanTiLong(whatCard *card.Card) bool {
	return player.playingCards.CanTiLong(whatCard)
}

func (player *Player) TestCard(whatCard *card.Card) *card.TestCardResult{
	return player.playingCards.TestCard(whatCard)
}

//begin try operate

func (player *Player) OperateEnterRoom(room *Room) bool{
	log.Debug(player, "OperateEnterRoom room :", room)
	data := &OperateEnterRoomData{}
	op := NewOperateEnterRoom(player, data)
	room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperateLeaveRoom() bool{
	log.Debug(player, "OperateLeaveRoom", player.room)
	if player.room == nil {
		return true
	}

	data := &OperateLeaveRoomData{}
	op := NewOperateLeaveRoom(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperateDropCard(card *card.Card) bool {
	log.Debug(player, "OperateDrop card :", card)
	data := &OperateDropCardData{
		Card: card,
	}
	op := NewOperateDropCard(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}


func (player *Player) OperateChiCard(whatCard *card.Card, whatGroup *card.Cards) bool {
	log.Debug(player, "OperateChi, card :", whatCard, " group :", whatGroup)
	data := &OperateChiCardData{
		Card: whatCard,
		Group: whatGroup,
	}
	op := NewOperateChiCard(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperatePengCard(card *card.Card) bool {
	log.Debug(player, "OperatePeng card :", card)
	data := &OperatePengCardData{
		Card: card,
	}
	op := NewOperatePengCard(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperateSaoCard(card *card.Card) bool {
	log.Debug(player, "OperateSao card :", card)
	data := &OperateSaoCardData{
		Card: card,
	}
	op := NewOperateSaoCard(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperatePaoCard(card *card.Card) bool {
	log.Debug(player, "OperatePaoCard", card)
	data := &OperatePaoCardData{
		Card: card,
	}
	op := NewOperatePaoCard(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperateTiLongCard(card *card.Card) bool {
	log.Debug(player, "OperateTiLongCard card :", card)
	data := &OperateTiLongCardData{
		Card: card,
	}
	op := NewOperateTiLongCard(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

//end try operate

func (player *Player) EnterRoom(room *Room) {
	log.Debug(player, "enter", room)
	player.room = room
}

func (player *Player) LeaveRoom() {
	log.Debug(player, "leave", player.room)
	player.room = nil
}

func (player *Player) Drop(card *card.Card) bool {
	log.Debug(player, "Drop card :", card)
	return player.playingCards.DropCard(card)
}

func (player *Player) AutoDrop() *card.Card {
	log.Debug(player, "AutoDrop")
	return player.playingCards.DropTail()
}

func (player *Player) Chi(whatCard *card.Card, whatGroup *card.Cards) bool {
	log.Debug(player, "Chi card :", whatCard, ", group :", whatGroup)
	return player.playingCards.Chi(whatCard, whatGroup)
}

func (player *Player) Peng(card *card.Card, fromPlayer *Player) bool {
	log.Debug(player, "Peng card :", card)
	ok := player.playingCards.Peng(card)
	if ok {
		player.pengFromPlayer[card.MakeKey()] = fromPlayer
	}
	return ok
}

func (player *Player) Sao(card *card.Card) bool {
	log.Debug(player, "Sao", card)
	return player.playingCards.Sao(card)
}

func (player *Player) Pao(card *card.Card, fromPlayer *Player) bool {
	log.Debug(player, "Pao", card)
	ok := player.playingCards.Pao(card)
	if ok {
		player.paoFromPlayers[card.MakeKey()] = fromPlayer
	}
	return ok
}

func (player *Player) TiLong(card *card.Card) bool {
	log.Debug(player, "TiLong", card)
	return player.playingCards.TiLong(card)
}

func (player *Player) IsHu() bool {
	//todo
}

func (player *Player) OnPlayingGameEnd() {
	log.Debug(player, "OnPlayingGameEnd")
	player.Reset()
	data := &GameEndMsgData{}
	player.notifyObserver(NewGameEndMsg(player, data))
}

func (player *Player) OnRoomClosed() {
	log.Debug(player, "OnRoomClosed")
	player.room = nil
	player.Reset()

	data := &RoomClosedMsgData{}
	player.notifyObserver(NewRoomClosedMsg(player, data))
}

func (player *Player) OnDispatchCard(card *card.Card) *card.TestCardResult{
	return player.TestCard(card)
}

func (player *Player) OnGetInitCards() {
	log.Debug(player, "OnGetInitCards")

	data := &GetInitCardsMsgData{
		PlayingCards: player.playingCards,
	}
	player.notifyObserver(NewGetInitCardsMsg(player, data))
}

func (player *Player) String() string{
	if player == nil {
		return "{player=nil}"
	}
	return fmt.Sprintf("{player=%v}", player.id)
}

//玩家成功操作的通知
func (player *Player) OnPlayerSuccessOperated(op *Operate) {
	log.Debug(player, "OnPlayerSuccessOperated", op)
	switch op.Op {
	case OperateEnterRoom:
		player.onPlayerEnterRoom(op)
	case OperateLeaveRoom:
		player.onPlayerLeaveRoom(op)
	case OperateHu:
		//todo
	case OperateDropCard:
		player.onPlayerDrop(op)
	case OperateChiCard:
		player.onPlayerChi(op)
	case OperatePengCard:
		player.onPlayerPeng(op)
	case OperateSaoCard :
		player.onPlayerGang(op)
	case OperatePaoCard :
		player.onPlayerZiMo(op)
	case OperateTiLongCard :
		player.onPlayerDianPao(op)
	}
}

func (player *Player) GetRoom() *Room {
	return player.room
}

func (player *Player) waitResult(resultCh chan bool) bool{
	log.Debug(player, "Player.waitResult")
	select {
	case <- time.After(time.Second * 10):
		log.Debug(player, "Player.waitResult timeout")
		return false
	case result := <- resultCh:
		log.Debug(player, "Player.waitResult result :", result)
		return result
	}
	return false
}

func (player *Player) notifyObserver(msg *Message) {
	log.Debug(player, "notifyObserverMsg", msg)
	for _, ob := range player.observers {
		ob.OnMsg(player, msg)
	}
}

func (player *Player) calcScore(huPlayer *Player) int {
}

//begin player operate event handler

func (player *Player) onPlayerEnterRoom(op *Operate) {
	if _, ok := op.Data.(*OperateEnterRoomData); ok {
		if player.room == nil {
			return
		}
		msgData := &EnterRoomMsgData{
			EnterPlayer : op.Operator,
			AllPlayer: player.room.players,
		}
		player.notifyObserver(NewEnterRoomMsg(player, msgData))
	}
}

func (player *Player) onPlayerLeaveRoom(op *Operate) {
	if _, ok := op.Data.(*OperateLeaveRoomData); ok {
		if op.Operator == player {
			return
		}
		if player.room == nil {
			return
		}
		msgData := &LeaveRoomMsgData{
			LeavePlayer : op.Operator,
			AllPlayer: player.room.players,
		}
		player.notifyObserver(NewLeaveRoomMsg(player, msgData))
	}
}

func (player *Player) onPlyaerGetInitCards(op *PlayerOperate) {
	if data, ok := op.Data.(*PlayerOperateGetInitCardsData); ok {
		msg := &GetInitCardsMsg{
			CardsInHand: data.CardsInHand,
			MagicCards: data.MagicCards,
		}
		player.notifyObserver(NewPlayerGetInitCardsMsg(player, msg))
	}
}

func (player *Player) onPlayerGet(op *PlayerOperate) {
	if data, ok := op.Data.(*PlayerOperateGetData); ok {
		msg := &GetCardMsg{
			canZiMo : player.CanZiMo(),
		}
		if op.Operator == player {//拿到牌只告诉自己是什么牌
			msg.Card = data.Card
		}
		player.notifyObserver(NewPlayerGetCardMsg(player, msg))
	}
}

func (player *Player) onPlayerDrop(op *PlayerOperate) {
	if data, ok := op.Data.(*PlayerOperateDropData); ok {
		if op.Operator == player {//自己出牌，不用通知自己
			return
		}
		msg := &DropCardMsg{
			WhatCard : data.Card,
			canChiGroup : player.ComputeChiGroup(data.Card),
			canPeng	: player.CanPeng(data.Card),
			canGang : player.CanGang(data.Card),
			canDianPao : player.CanDianPao(data.Card),
		}
		player.notifyObserver(NewPlayerDropCardMsg(player, msg))
	}
}

func (player *Player) onPlayerChi(op *PlayerOperate) {
	if data, ok := op.Data.(*PlayerOperateChiData); ok {
		if op.Operator == player {//自己吃牌，不用通知自己
			return
		}
		msg := &ChiCardMsg{
			ChiPlayer: op.Operator,
			FromPlayer: player.room.getPrevOperator(),
			WhatCard: data.Card,
			WhatGroup: data.Group,
		}
		player.notifyObserver(NewPlayerChiCardMsg(player, msg))
	}
}

func (player *Player) onPlayerPeng(op *PlayerOperate) {
	if data, ok := op.Data.(*PlayerOperatePengData); ok {
		if op.Operator == player {//自己碰牌，不用通知自己
			return
		}
		msg := &PengCardMsg{
			PengPlayer: op.Operator,
			FromPlayer: player.room.getPrevOperator(),
			WhatCard: data.Card,
		}
		player.notifyObserver(NewPlayerPengCardMsg(player, msg))
	}
}

func (player *Player) onPlayerGang(op *PlayerOperate) {
	if data, ok := op.Data.(*PlayerOperateGangData); ok {
		if op.Operator == player {//自己杠牌，不用通知自己
			return
		}
		msg := &GangCardMsg{
			GangPlayer: op.Operator,
			FromPlayer: player.room.getPrevOperator(),
			WhatCard: data.Card,
		}
		player.notifyObserver(NewPlayerGangCardMsg(player, msg))
	}
}

func (player *Player) onPlayerZiMo(op *PlayerOperate) {
	if _, ok := op.Data.(*PlayerOperateZiMoData); ok {
		msg := &ZiMoMsg{
			HuPlayer: op.Operator,
			WhatCard: op.Operator.lastGetCard,
			Desc: op.Operator.result.Desc,
			PlayerScore: make([]*PlayerScore, 0),
		}
		for _, tmpPlayer := range player.room.players {
			score := &PlayerScore{
				P : tmpPlayer,
				Score: tmpPlayer.calcScore(op.Operator),
			}
			msg.PlayerScore = append(msg.PlayerScore, score)
		}
		player.notifyObserver(NewPlayerZiMoMsg(player, msg))
	}
}

func (player *Player) onPlayerDianPao(op *PlayerOperate) {
	if data, ok := op.Data.(*PlayerOperateDianPaoData); ok {
		msg := &DianPaoMsg{
			HuPlayer: op.Operator,
			FromPlayer: player.room.getPrevOperator(),
			WhatCard: data.Card,
			Desc: op.Operator.result.Desc,
			PlayerScore: make([]*PlayerScore, 0),
		}
		for _, tmpPlayer := range player.room.players {
			score := &PlayerScore{
				P : tmpPlayer,
				Score: tmpPlayer.calcScore(op.Operator),
			}
			msg.PlayerScore = append(msg.PlayerScore, score)
		}
		player.notifyObserver(NewPlayerDianPaoMsg(player, msg))
	}
}

//end player operate event handler
