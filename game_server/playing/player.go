package playing

import (
	"peng_hu_zi/game_server/card"
	"peng_hu_zi/log"
	"time"
	"fmt"
)

type HuResult struct {
	IsHu 	bool
	IsZiMo	bool
	Desc	string
	Score	int
}

func (result *HuResult) String() string {
	if result == nil {
		return "{ishu=nil, desc=nil, score=nil}"
	}
	return fmt.Sprintf("{ishu=%v, desc=%s, score=%d}", result.IsHu, result.Desc, result.Score)
}

func newHuResult(isHu bool, desc string, score int) *HuResult {
	return &HuResult{
		IsHu: isHu,
		Desc: desc,
		Score: score,
	}
}

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

//begin operate

func (player *Player) OperateEnterRoom(room *Room) bool{
	log.Debug(player, "OperateEnterRoom room :", room)
	data := &OperateEnterRoomData{}
	op := NewOperateEnterRoom(player, data)
	room.PlayerOperate(op)
	result := player.waitResult(op.Result)
	if r, ok := result.(*OperateEnterRoomResult); ok {
		log.Debug("OperateEnterRoom succ, then set player.room")
		player.room = room
		return r.Ok
	}

	return false
}

func (player *Player) OperateLeaveRoom() bool{
	log.Debug(player, "OperateLeaveRoom", player.room)
	if player.room == nil {
		return true
	}

	data := &OperateLeaveRoomData{}
	op := NewOperateLeaveRoom(player, data)
	player.room.PlayerOperate(op)
	result := player.waitResult(op.Result)
	if _, ok := result.(*OperateLeaveRoomResult); ok {
		player.room = nil
	}

	return true
}

func (player *Player) OperateDropCard(card *card.Card) bool {
	log.Debug(player, "OperateDrop card :", card)
	data := &OperateDropCardData{
		Card: card,
	}
	op := NewOperateDropCard(player, data)
	player.room.PlayerOperate(op)
	result := player.waitResult(op.Result)
	if r, ok := result.(*OperateDropCardResult); ok {
		return r.OK
	}
	return false
}


func (player *Player) OperateChi(whatCard *card.Card, whatGroup *card.Cards) bool {
	log.Debug(player, "OperateChi, card :", whatCard.Name(), " group :", whatGroup.ToString())
	result := make(chan bool, 1)
	data := &PlayerOperateChiData{
		Card: whatCard,
		Group: whatGroup,
	}
	op := NewPlayerOperateChi(player, result, data)
	player.room.PlayerOperate(op)
	return player.waitResult(result)
}

func (player *Player) OperatePeng(card *card.Card) bool {
	log.Debug(player, "OperatePeng card :", card.Name())
	result := make(chan bool, 1)
	data := &PlayerOperatePengData{
		Card: card,
	}
	op := NewPlayerOperatePeng(player, result, data)
	player.room.PlayerOperate(op)
	return player.waitResult(result)
}

func (player *Player) OperateGang(card *card.Card) bool {
	log.Debug(player, "OperateGang card :", card.Name())
	result := make(chan bool, 1)
	data := &PlayerOperateGangData{
		Card: card,
	}
	op := NewPlayerOperateGang(player, result, data)
	player.room.PlayerOperate(op)
	return player.waitResult(result)
}

func (player *Player) OperateZiMo() bool {
	log.Debug(player, "OperateZiMo")
	result := make(chan bool, 1)
	data := &PlayerOperateZiMoData{}
	op := NewPlayerOperateZiMo(player, result, data)
	player.room.PlayerOperate(op)
	return player.waitResult(result)
}

func (player *Player) OperateDianPao(card *card.Card) bool {
	log.Debug(player, "OperateDianPao card :", card.Name())
	result := make(chan bool, 1)
	data := &PlayerOperateDianPaoData{
		Card: card,
	}
	op := NewPlayerOperateDianPao(player, result, data)
	player.room.PlayerOperate(op)
	return player.waitResult(result)
}

//end operate

func (player *Player) BeGangBy(gangPlayer *Player) {
	player.gangByPlayers = append(player.gangByPlayers, gangPlayer)
}

func (player *Player) Chi(whatCard *card.Card, whatGroup *card.Cards) bool {
	log.Debug(player, "Chi card :", whatCard.Name(), ", group :", whatGroup.ToString())
	return player.playingCards.Chi(whatCard, whatGroup)
}

func (player *Player) Peng(card *card.Card) bool {
	log.Debug(player, "Peng card :", card.Name())
	return player.playingCards.Peng(card)
}

func (player *Player) Gang(card *card.Card, fromPlayer *Player) bool {
	log.Debug(player, "Gang card :", card.Name())
	ok := player.playingCards.Gang(card)
	if ok {
		player.gangFromPlayers = append(player.gangFromPlayers, fromPlayer)
	}
	return ok
}

func (player *Player) Drop(card *card.Card) bool {
	log.Debug(player, "Drop card :", card.Name())
	return player.playingCards.DropCard(card)
}

func (player *Player) AutoDrop() *card.Card {
	log.Debug(player, "AutoDrop")
	return player.playingCards.DropTail()
}

func (player *Player) IsHu() *HuResult {
	magicLen := player.playingCards.GetMagicCards().Len()
	log.Debug(player, "IsHu magicLen :", magicLen)
	for _, checker := range player.huChecker {
		if magicLen == 0 {
			isHu, conf := checker.IsHu(player.playingCards)
			if isHu {
				log.Debug(player, "checker :", checker.GetConfig().ToString(), "succ")
				return newHuResult(isHu, conf.Desc, conf.Score)
			}
		} else {
			//支持赖子牌的胡牌计算, 暴力穷举法，把赖子牌的所有候选集一个个试，胜在够简单
			candidate := player.computeMagicCandidate()
			tryCnt := 0
			for _, cards := range candidate {
				tryCnt++
				player.playingCards.AddCards(cards)
				isHu, conf := checker.IsHu(player.playingCards)
				if isHu {
					log.Debug(player, "checker :", checker.GetConfig().ToString(), "succ, tryMagicCnt :", tryCnt, ",cards:", cards.ToString())
					log.Debug(player, "tryCnt :", tryCnt, ", cards :", cards.ToString())
					return newHuResult(isHu, conf.Desc, conf.Score)
				} else {
					player.playingCards.DropCards(cards)
				}
			}
		}
		log.Debug(player, "checker :", checker.GetConfig().ToString(), "failed")
	}

	return newHuResult(false, "", 0)
}

func (player *Player) OnPlayingGameEnd() {
	log.Debug(player, "OnPlayingGameEnd")
	msg := &PlayingGameEndMsg{}
	player.notifyObserver(NewPlayerPlayingGameEndMsg(player, msg))
}

func (player *Player) OnRoomClosed() {
	log.Debug(player, "OnRoomClosed")
	player.room = nil
	player.Reset()

	msg := &RoomClosedMsg{}
	player.notifyObserver(NewPlayerRoomClosedMsg(player, msg))
}

func (player *Player) OnDispatchCard() {
	//TODO
}

func (player *Player) OnGetInitCards() {

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
	case PlayerOperateGet:
		player.onPlayerGet(op)
	case PlayerOperateDrop:
		player.onPlayerDrop(op)
	case PlayerOperateChi:
		player.onPlayerChi(op)
	case PlayerOperatePeng:
		player.onPlayerPeng(op)
	case PlayerOperateGang :
		player.onPlayerGang(op)
	case PlayerOperateZiMo :
		player.onPlayerZiMo(op)
	case PlayerOperateDianPao :
		player.onPlayerDianPao(op)
	case PlayerOperateEnterRoom:
		player.onPlayerEnterRoom(op)
	case PlayerOperateLeaveRoom:
		player.onPlayerLeaveRoom(op)
	}
}

func (player *Player) GetRoom() *Room {
	return player.room
}

func (player *Player) waitResult(resultCh chan interface{}) interface{}{
	log.Debug(player, "Player.waitResult")
	select {
	case <- time.After(time.Second * 10):
		log.Debug(player, "Player.waitResult timeout")
		return nil
	case result := <- resultCh:
		log.Debug(player, "Player.waitResult result :", result)
		return result
	}
	return nil
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

func (player *Player) onPlayerEnterRoom(op *PlayerOperate) {
	if data, ok := op.Data.(*PlayerOperateEnterRoomData); ok {
		msg := &EnterRoomMsg{
			EnterPlayer : op.Operator,
			AllPlayer: data.Players,
		}
		log.Debug(player, "onPlayerEnterRoom, msg :", msg)
		player.notifyObserver(NewPlayerEnterRoomMsg(player, msg))
	}
}

func (player *Player) onPlayerLeaveRoom(op *PlayerOperate) {
	if data, ok := op.Data.(*PlayerOperateLeaveRoomData); ok {
		msg := &LeaveRoomMsg{
			LeavePlayer : op.Operator,
			AllPlayer: data.Players,
		}
		player.notifyObserver(NewPlayerLeaveRoomMsg(player, msg))
	}
}
//end player operate event handler
