package playing

import (
	"peng_hu_zi/game_server/card"
	"peng_hu_zi/log"
	"time"
	"fmt"
)

type TestCardResult struct {
	CanHu		bool
	CanTiLong	bool
	CanSao		bool
	CanPao		bool
	CanPeng		bool
	ChiGroup	[]*card.Cards
}

func (result *TestCardResult) String() string {
	return fmt.Sprintf("{CanPao=%v, CanSao=%v, CanHu=%v, CanPeng=%v, ChiGroup=%v}",
		result.CanPao,
		result.CanSao,
		result.CanHu,
		result.CanPeng,
		result.ChiGroup,
	)
}

type PlayerObserver interface {
	OnMsg(player *Player, msg *Message)
}

type Player struct {
	id				uint64			// 玩家id
	idxOfRoom		int				//	玩家在房间的位置
	room			*Room			// 玩家所在的房间

	playingCards 	*card.PlayingCards	//玩家手上的牌

	pengFromPlayer 	map[int64]*Player			//碰的牌来源于 int64:Card.MakeKey()
	paoFromPlayers  map[int64]*Player			//跑的牌来源于 int64:Card.MakeKey()


	observers	 []PlayerObserver
}

func NewPlayer(id uint64) *Player {
	player :=  &Player{
		id:		id,
		playingCards:	card.NewPlayingCards(),
		pengFromPlayer: make(map[int64]*Player),
		paoFromPlayers: make(map[int64]*Player),
		observers:	make([]PlayerObserver, 0),
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

func (player *Player) TestCard(whatCard *card.Card, fromPlayer *Player) *TestCardResult{
	result := &TestCardResult{}
	if player == fromPlayer {//牌来源于自己，检查是否扫/提龙
		result.CanTiLong = player.playingCards.CanTiLong(whatCard)
		result.CanSao = player.playingCards.CanSao(whatCard)
	} else {
		result.CanPao = player.playingCards.CanPao(whatCard)
		result.CanPeng = player.playingCards.CanPeng(whatCard)
	}

	if player == player.room.curOperator {//当前操作者需要计算吃
		result.ChiGroup = player.playingCards.ComputeChiGroup(whatCard)
	} else if player == player.room.nextPlayer(player.room.curOperator) {//当前操作者的下家需要计算吃
		result.ChiGroup = player.playingCards.ComputeChiGroup(whatCard)
	}

	result.CanHu = player.playingCards.IsHuThisCard(whatCard)

	return result
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


func (player *Player) OperateChiCard(whatGroup *card.Cards) bool {
	card := player.room.lastDropCard
	log.Debug(player, "OperateChi, card :", card, " group :", whatGroup)
	data := &OperateChiCardData{
		Card: card,
		Group: whatGroup,
	}
	op := NewOperateChiCard(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperatePengCard() bool {
	card := player.room.lastDropCard
	//player.room.lastDropCardOperator
	log.Debug(player, "OperatePeng card :", card)
	data := &OperatePengCardData{
		Card: card,
	}
	op := NewOperatePengCard(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperateSaoCard() bool {
	card := player.room.lastDropCard
	log.Debug(player, "OperateSao card :", card)
	data := &OperateSaoCardData{
		Card: card,
	}
	op := NewOperateSaoCard(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperatePaoCard() bool {
	card := player.room.lastDropCard
	log.Debug(player, "OperatePaoCard", card)
	data := &OperatePaoCardData{
		Card: card,
	}
	op := NewOperatePaoCard(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperateTiLongCard() bool {
	card := player.room.lastDropCard
	log.Debug(player, "OperateTiLongCard card :", card)
	data := &OperateTiLongCardData{
		Card: card,
	}
	op := NewOperateTiLongCard(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperateGuo() bool {
	log.Debug(player, "OperateGuo")
	data := &OperateGuoData{}
	op := NewOperateGuo(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

//end try operate

func (player *Player) EnterRoom(room *Room, idxOfRoom int) {
	log.Debug(player, "enter", room)
	player.room = room
	player.idxOfRoom = idxOfRoom
}

func (player *Player) LeaveRoom() {
	log.Debug(player, "leave", player.room)
	player.room = nil
	player.idxOfRoom = -1
}

func (player *Player) Drop(card *card.Card) bool {
	log.Debug(player, "Drop card :", card)
	return player.playingCards.DropCard(card)
}

func (player *Player) GetTailCard() *card.Card {
	log.Debug(player, "GetTailCard")
	return player.playingCards.Tail()
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
	player.playingCards.IsHu()
	return false
}

func (player *Player) OnEndPlayGame() {
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

func (player *Player) OnGetInitCards() {
	log.Debug(player, "OnGetInitCards", player.playingCards)

	data := &GetInitCardsMsgData{
		PlayingCards: player.playingCards,
	}
	player.notifyObserver(NewGetInitCardsMsg(player, data))
}

func (player *Player) OnDispatchCard(card *card.Card) {
	log.Debug(player, "player.OnDispatchCard", card)
	player.notifyObserver(NewDispatchCardMsg(player, &DispatchCardMsgData{}))
}

func (player *Player) ShowDispatchCard(card *card.Card) {
	log.Debug(player, "player.ShowDispatchCard", card)
	player.notifyObserver(NewShowDispatchCardMsg(player, &ShowDispatchCardMsgData{Card:card}))
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
		player.onPlayerHu(op)
	case OperateDropCard:
		player.onPlayerDrop(op)
	case OperateChiCard:
		player.onPlayerChi(op)
	case OperatePengCard:
		player.onPlayerPeng(op)
	case OperateSaoCard :
		player.onPlayerSao(op)
	case OperatePaoCard :
		player.onPlayerPao(op)
	case OperateTiLongCard :
		player.onPlayerTiLong(op)
	}
}

func (player *Player) GetRoom() *Room {
	return player.room
}

func (player *Player) waitResult(resultCh chan bool) bool{
	select {
	case <- time.After(time.Second * 10):
		log.Debug(player, "Player.waitResult timeout")
		return false
	case result := <- resultCh:
		log.Debug(player, "Player.waitResult result :", result)
		return result
	}
	log.Debug(player, "Player.waitResult fasle")
	return false
}

func (player *Player) notifyObserver(msg *Message) {
	log.Debug(player, "notifyObserverMsg", msg)
	for _, ob := range player.observers {
		ob.OnMsg(player, msg)
	}
}

func (player *Player) calcScore(huPlayer *Player) int {
	return 0
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

func (player *Player) onPlayerHu(op *Operate) {
	if opData, ok := op.Data.(*OperateHuData); ok {
		msgData := &HuMsgData{
			HuPlayer : opData.HuPlayer,
			FromPlayer : opData.FromPlayer,
			Desc : opData.Desc,
			PlayerScore: opData.PlayerScore,
		}
		player.notifyObserver(NewHuMsg(player, msgData))
	}
}


func (player *Player) onPlayerDrop(op *Operate) {
	if opData, ok := op.Data.(*OperateDropCardData); ok {
		if op.Operator == player {//自己出牌，不用通知自己
			return
		}
		msgData := &DropCardMsgData{
			Card: opData.Card,
			CanPeng: player.CanPeng(opData.Card),
			ChiGroup: player.ComputeChiGroup(opData.Card),
		}
		player.notifyObserver(NewDropCardMsg(player, msgData))
	}
}

func (player *Player) onPlayerChi(op *Operate) {
	if opData, ok := op.Data.(*OperateChiCardData); ok {
		if op.Operator == player {//自己吃牌，不用通知自己
			return
		}
		msgData := &ChiCardMsgData{
			ChiPlayer: op.Operator,
			FromPlayer: player.room.getDropCardOperator(),
			WhatCard: opData.Card,
			WhatGroup: opData.Group,
		}
		player.notifyObserver(NewChiCardMsg(player, msgData))
	}
}

func (player *Player) onPlayerPeng(op *Operate) {
	if opData, ok := op.Data.(*OperatePengCardData); ok {
		if op.Operator == player {//自己碰牌，不用通知自己
			return
		}
		msgData := &PengCardMsgData{
			PengPlayer: op.Operator,
			FromPlayer: player.room.getDropCardOperator(),
			WhatCard: opData.Card,
		}
		player.notifyObserver(NewPengCardMsg(player, msgData))
	}
}

func (player *Player) onPlayerSao(op *Operate) {
	if opData, ok := op.Data.(*OperateSaoCardData); ok {
		if op.Operator == player {//自己扫牌，不用通知自己
			return
		}
		msgData := &SaoCardMsgData{
			SaoPlayer: op.Operator,
			WhatCard: opData.Card,
		}
		player.notifyObserver(NewSaoCardMsg(player, msgData))
	}
}

func (player *Player) onPlayerPao(op *Operate) {
	if opData, ok := op.Data.(*OperatePaoCardData); ok {
		if op.Operator == player {//自己跑牌，不用通知自己
			return
		}
		msgData := &PaoCardMsgData{
			PaoPlayer: op.Operator,
			WhatCard: opData.Card,
		}
		player.notifyObserver(NewPaoCardMsg(player, msgData))
	}
}

func (player *Player) onPlayerTiLong(op *Operate) {
	if opData, ok := op.Data.(*OperateTiLongCardData); ok {
		if op.Operator == player {//自己跑牌，不用通知自己
			return
		}
		msgData := &TiLongCardMsgData{
			TiLongPlayer: op.Operator,
			WhatCard: opData.Card,
		}
		player.notifyObserver(NewTiLongCardMsg(player, msgData))
	}
}
//end player operate event handler

func (player *Player) ComputeTiLong() {
	player.playingCards.ComputeTiLong()
}

func (player *Player) ComputeSao() {
	player.playingCards.ComputeSao()
}

func (player *Player) GetPaoAndTiLongNum() int{
	return player.playingCards.GetPaoAndTiLongNum()
}

func (player *Player) GetPlayingCards() *card.PlayingCards {
	return player.playingCards
}