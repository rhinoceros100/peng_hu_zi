package playing

import (
	"peng_hu_zi/game_server/card"
	"peng_hu_zi/log"
	"time"
	"peng_hu_zi/util"
	"fmt"
)

type RoomStatusType int
const (
	RoomStatusWaitAllPlayerEnter	RoomStatusType = iota	// 等待玩家进入房间
	RoomStatusWaitStartPlayGame				// 等待游戏开始
	RoomStatusPlayGame					// 正在进行游戏，结束后会进入RoomStatusEndPlayGame
	RoomStatusEndPlayGame					// 游戏结束后会回到等待游戏开始状态，或者进入结束房间状态
	RoomStatusRoomEnd						// 房间结束状态
)

func (status RoomStatusType) String() string {
	switch status {
	case RoomStatusWaitAllPlayerEnter :
		return "RoomStatusWaitAllPlayerEnter"
	case RoomStatusWaitStartPlayGame:
		return "RoomStatusWaitStartPlayGame"
	case RoomStatusPlayGame:
		return "RoomStatusPlayGame"
	case RoomStatusEndPlayGame:
		return "RoomStatusEndPlayGame"
	case RoomStatusRoomEnd:
		return "RoomStatusRoomEnd"
	}
	return "unknow RoomStatus"
}

type RoomObserver interface {
	OnRoomClosed(room *Room)
}

type Room struct {
	id				uint64					//房间id
	config 			*RoomConfig				//房间配置
	players 		[]*Player				//当前房间的玩家列表

	observers		[]RoomObserver			//房间观察者，需要实现OnRoomClose，房间close的时候会通知它

	roomStatus		RoomStatusType						//房间当前的状态

	lastHuPlayer	*Player					//最后一次胡牌的玩家
	playedGameCnt	int						//已经玩了的游戏的次数

	//begin playingGameData, reset when start playing game
	cardPool		*card.Pool				//洗牌池
	masterPlayer	*Player					//做东的玩家，打筛子的玩家
	curOperator	*Player				//获得当前操作的玩家，可能是碰牌，跑牌，吃牌，等他出牌
	prevOperator *Player			//上一个操作的玩家
	lastDropCardOperator *Player	//最后一个出牌的人，可能是系统，可能是玩家，系统发牌时，该值为nil
	operateQue  []*Operate //当有玩家出牌或者系统排牌时，玩家的操作队列，依据优先级高低处理：胡 > 碰/跑 > 吃
	//end playingGameData, reset when start playing game

	//playerOpCh		chan *Operate		//用户操作的channel

	roomOperateCh				chan *Operate
	dropCardCh					chan *Operate
	masterChiCardCh				chan *Operate
	masterNextPlayerChiCardCh 	chan *Operate
	pengCardCh					chan *Operate
	paoCardCh					chan *Operate

	stop bool
}

func NewRoom(id uint64, config *RoomConfig) *Room {
	room := &Room{
		id:				id,
		config:			config,
		players:		make([]*Player, 0),
		cardPool:		card.NewPool(),
		observers:		make([]RoomObserver, 0),
		roomStatus:		RoomStatusWaitAllPlayerEnter,
		playedGameCnt:	0,

		operateQue:	make([]*Operate, 0),

		//playerOpCh:		make(chan *Operate, 1024),
		roomOperateCh: make(chan *Operate, 1024),
		dropCardCh: make(chan *Operate, 1024),
		masterChiCardCh: make(chan *Operate, 1024),
		masterNextPlayerChiCardCh: make(chan *Operate, 1024),
		pengCardCh: make(chan *Operate, 1024),
	}

	return room
}

func (room *Room) GetId() uint64 {
	return room.id
}

func (room *Room) PlayerOperate(op *Operate) {
	switch op.Op {
	case OperateEnterRoom, OperateLeaveRoom:
		room.roomOperateCh <- op
	case OperateDropCard:
		room.dropCardCh <- op
	case OperateChiCard:
		if op.Operator == room.curOperator {
			room.masterChiCardCh <- op
		} else if op.Operator == room.nextPlayer(room.curOperator) {
			room.masterNextPlayerChiCardCh <- op
		}
	case OperatePengCard:
		room.pengCardCh <- op
	}
}

func (room *Room) addObserver(observer RoomObserver) {
	room.observers = append(room.observers, observer)
}

func (room *Room) Start() {
	go func() {
		for  {
			if !room.stop {
				room.checkStatus()
			}
		}
	}()
}

func (room *Room) checkStatus() {
	switch room.roomStatus {
	case RoomStatusWaitAllPlayerEnter:
		room.waitAllPlayerEnter()
	case RoomStatusWaitStartPlayGame:
		room.startPlayGame()
	case RoomStatusPlayGame:
		room.playGame()
	case RoomStatusEndPlayGame:
		room.endPlayGame()
	case RoomStatusRoomEnd:
		room.close()
	}
}

func (room *Room) isRoomEnd() bool {
	return room.playedGameCnt >= room.config.MaxPlayGameCnt
}

func (room *Room) close() {
	log.Debug(room, "Room.close")
	room.stop = true
	for _, observer := range room.observers {
		observer.OnRoomClosed(room)
	}

	for _, player := range room.players {
		player.OnRoomClosed()
	}
}

func (room *Room) isAllPlayerEnter() bool {
	length := len(room.players)
	log.Debug(room, "Room.isAllPlayerEnter, player num :", length, ", need :", room.config.NeedPlayerNum)
	return length >= room.config.NeedPlayerNum
}

func (room *Room) switchStatus(status RoomStatusType) {
	log.Debug(room, "room status switch,", room.roomStatus, " =>", status)
	room.roomStatus = status
}

func (room *Room) startPlayGame()  {
	log.Debug(room, "Room.startPlayGame")

	// 重置牌池, 洗牌
	room.cardPool.ReGenerate()

	// 选择东家
	room.masterPlayer = room.selectMasterPlayer()

	// 设定获得牌的玩家为东家
	room.switchOperator(room.masterPlayer)
	room.prevOperator = nil
	room.lastDropCardOperator = nil

	room.operateQue = room.operateQue[0:0]

	//发初始牌给所有玩家
	room.putInitCardsToPlayers()

	//通知所有玩家手上的牌
	for _, player := range room.players {
		player.OnGetInitCards()
	}

	//todo
	//检查天胡
	//计算所有玩家的提扫

	room.switchStatus(RoomStatusPlayGame)

	room.waitPlayerDrop()
}

func (room *Room) playGame() {
	// 发牌给玩家
	card := room.cardPool.PopFront()
	log.Debug(room, "Room.playingGame put card[ ", card, "]to", room.curOperator)
	if card == nil {//没有牌了，该局结束，流局
		room.switchStatus(RoomStatusEndPlayGame)
		return
	}

	room.broadcastDispatchCard(card)
	room.checkTestCard(card, true)
}

func (room *Room) endPlayGame() {
	room.playedGameCnt++
	log.Debug(room, "Room.endPlayGame cnt :", room.playedGameCnt)
	if room.isRoomEnd() {
		log.Debug(room, "Room.endPlayGame room end")
		room.switchStatus(RoomStatusRoomEnd)
	} else {
		for _, player := range room.players {
			player.OnEndPlayGame()
		}
		log.Debug(room, "Room.endPlayGame restart play game")
		room.switchStatus(RoomStatusWaitStartPlayGame)
	}
}

//等待玩家进入
func (room *Room) waitAllPlayerEnter() {
	log.Debug(room, "Room.waitAllPlayerEnter")
	breakTimerTime := time.Duration(0)
	timeout := time.Duration(room.config.WaitPlayerEnterRoomTimeout) * time.Second
	for {
		timer := timeout - breakTimerTime
		select {
		case <-time.After(timer):
			log.Debug(room, "waitAllPlayerEnter timeout", timeout)
			room.switchStatus(RoomStatusRoomEnd) //超时发现没有足够的玩家都进入房间了，则结束
			return
		case op := <-room.roomOperateCh:
			if op.Op == OperateEnterRoom || op.Op == OperateLeaveRoom {
				log.Debug(room, "Room.waitAllPlayerEnter catch operate:", op)
				room.dealPlayerOperate(op)
				if room.isAllPlayerEnter() {
					room.switchStatus(RoomStatusWaitStartPlayGame)
					return
				}
			}
		}
	}
}

//给所有玩家发初始化的14张牌, 东家15张
func (room *Room) putInitCardsToPlayers() {
	log.Debug(room, "Room.initAllPlayer")
	for _, player := range room.players {
		player.Reset()
		for num := 0; num < 14; num++ {
			room.putCardToPlayer(player)
		}
	}
	room.putCardToPlayer(room.masterPlayer)
}

//添加玩家
func (room *Room) addPlayer(player *Player) bool {
	if room.roomStatus != RoomStatusWaitAllPlayerEnter {
		return false
	}
	room.players = append(room.players, player)
	return true
}

func (room *Room) delPlayer(player *Player)  {
	for idx, p := range room.players {
		if p == player {
			room.players = append(room.players[0:idx], room.players[idx+1:]...)
			return
		}
	}
}

//发牌给指定玩家
func (room *Room) putCardToPlayer(player *Player) *card.Card {
	card := room.cardPool.PopFront()
	if card == nil {
		return nil
	}
	player.AddCard(card)
	return card
}

//选择东家
func (room *Room) selectMasterPlayer() *Player {
	idx := util.RandomN(len(room.players))
	log.Debug(room, "Room.selectMasterPlayer", room.players[idx])
	return room.players[idx]

}

//等待碰、吃牌的玩家出牌，超时的话，自动帮他出一张牌
func (room *Room) waitPlayerDrop() {
	log.Debug(room, "Room.waitPlayerDrop")
	breakTimerTime := time.Duration(0)
	timeout := time.Duration(room.config.WaitPlayerOperateTimeout) * time.Second
	var dropOp *Operate
	for  {
		timer := timeout - breakTimerTime
		select {
		case <- time.After(timer):
			dropOp = room.makeDropCardOperate(room.curOperator, room.curOperator.GetTailCard())
			log.Debug(room, "Room.waitPlayerDrop ", room.curOperator, "auto drop", " op :", dropOp)
			room.dealPlayerOperate(dropOp)
			return
		case dropOp = <-room.dropCardCh :
			log.Debug(room, "Room.waitPlayerDrop operate :", dropOp)
			if room.curOperator != dropOp.Operator {
				continue
			}
			room.dealPlayerOperate(dropOp)
			return
		}
	}
}

//等待玩家碰牌
func (room *Room) waitPlayerPeng(player *Player) (isTimeout bool){
	log.Debug(room, "Room.waitPlayerPeng")
	breakTimerTime := time.Duration(0)
	timeout := time.Duration(room.config.WaitPlayerOperateTimeout) * time.Second
	for  {
		timer := timeout - breakTimerTime
		select {
		case <- time.After(timer):
			return true
		case op := <-room.pengCardCh :
			log.Debug(room, "Room.waitPlayerPeng operate :", op)
			if room.curOperator != op.Operator {
				continue
			}
			if op.Op != OperatePengCard {
				continue
			}
			room.dealPlayerOperate(op)
			return false
		}
	}

	return false
}

func (room *Room) waitMasterChi() (isTimeout bool){
	log.Debug(room, "Room.waitMasterChi")
	breakTimerTime := time.Duration(0)
	timeout := time.Duration(room.config.WaitPlayerOperateTimeout) * time.Second
	for  {
		timer := timeout - breakTimerTime
		select {
		case <- time.After(timer):
			return true
		case op := <-room.masterChiCardCh :
			log.Debug(room, "Room.waitMasterChi operate :", op)
			if op.Operator != room.curOperator {
				continue
			}
			room.dealPlayerOperate(op)
			return false
		}
	}
	return false
}

func (room *Room) waitMasterNextPlayerChi(){
	log.Debug(room, "Room.waitMasterNextPlayerChi")
	breakTimerTime := time.Duration(0)
	timeout := time.Duration(room.config.WaitPlayerOperateTimeout) * time.Second
	for  {
		timer := timeout - breakTimerTime
		select {
		case <- time.After(timer):
			return
		case op := <-room.masterNextPlayerChiCardCh:
			log.Debug(room, "Room.waitMasterNextPlayerChi operate :", op)
			if op.Operator != room.nextPlayer(room.curOperator) {
				continue
			}
			room.dealPlayerOperate(op)
			return
		}
	}
}

//取指定玩家的下一个玩家
func (room *Room) nextPlayer(player *Player) *Player {
	idx := player.idxOfRoom
	if idx == len(room.players) - 1 {
		log.Debug(room, "Room.nextPlayer :", room.players[0])
		return room.players[0]
	}
	log.Debug(room, "Room.nextPlayer :", room.players[idx+1])
	return room.players[idx+1]
}

//获取上一次出牌的玩家
func (room *Room) getDropCardOperator() *Player {
	return room.lastDropCardOperator
}

//处理玩家操作
func (room *Room) dealPlayerOperate(op *Operate) bool{
	log.Debug(room, "Room.dealPlayerOperate :", op)
	switch op.Op {
	case OperateEnterRoom:
		if _, ok := op.Data.(*OperateEnterRoomData); ok {
			if room.addPlayer(op.Operator) { //	玩家进入成功
				op.Operator.EnterRoom(room, len(room.players)-1)
				log.Debug(room, "Room.dealPlayerOperate player enter :", op.Operator)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				return true
			}
		}

	case OperateLeaveRoom:
		if _, ok := op.Data.(*OperateLeaveRoomData); ok {
			log.Debug(room, "Room.dealPlayerOperate player leave :", op.Operator)
			room.delPlayer(op.Operator)
			op.Operator.LeaveRoom()
			op.ResultCh <- true
			room.broadcastPlayerSuccessOperated(op)
			return true
		}

	case OperateDropCard:
		if data, ok := op.Data.(*OperateDropCardData); ok {
			if op.Operator.Drop(data.Card) { //出牌
				log.Debug(room, "Room.dealPlayerOperate Drop card :", data.Card)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				room.checkTestCard(data.Card, false)
				return true
			}
		}

	case OperateChiCard:
		if data, ok := op.Data.(*OperateChiCardData); ok {
			if op.Operator.Chi(data.Card, data.Group) {
				log.Debug(room, "Room.dealPlayerOperate chi card :", data.Card, ", group :", data.Group)
				//吃成功了，设定当前玩家为吃牌者，并等待他出牌
				room.switchOperator(op.Operator)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				room.waitPlayerDrop()
				return true
			}
		}

	case OperatePengCard:
		if data, ok := op.Data.(*OperatePengCardData); ok {
			if op.Operator.Peng(data.Card, room.lastDropCardOperator) {
				//碰成功了，设定当前玩家为碰牌者，并等待他出牌
				log.Debug(room, "Room.dealPlayerOperate peng :", data.Card)
				room.switchOperator(op.Operator)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				room.waitPlayerDrop()
				return true
			}
		}

	case OperateSaoCard:
		if data, ok := op.Data.(*OperateSaoCardData); ok {
			if op.Operator.Sao(data.Card) {
				log.Debug(room, "Room.dealPlayerOperate gang :", data.Card)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				return true
			}
		}

	case OperatePaoCard:
		if data, ok := op.Data.(*OperatePaoCardData); ok {
			if op.Operator.Pao(data.Card, room.lastDropCardOperator) {
				//跑成功了，设定当前玩家为跑牌者，并等待他出牌
				log.Debug(room, "Room.dealPlayerOperate Pao :", data.Card)
				room.switchOperator(op.Operator)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				if op.Operator.GetPaoAndTiLongNum() < 2 {
					room.waitPlayerDrop()
				}
				return true
			}
		}

	case OperateTiLongCard :
		if data, ok := op.Data.(*OperateSaoCardData); ok {
			if op.Operator.TiLong(data.Card) {
				log.Debug(room, "Room.dealPlayerOperate TiLong :", data.Card)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				return true
			}
		}

	case OperateHu :
		if _, ok := op.Data.(*OperateHuData); ok {
			room.broadcastPlayerSuccessOperated(op)
			op.ResultCh <- true
			return true
		}

	}


	op.ResultCh <- false
	return false
}

func (room *Room) broadcastPlayerSuccessOperated(op *Operate) {
	log.Debug(room, "Room.broadcastPlayerSuccessOperated :", op)
	for _, player := range room.players {
		player.OnPlayerSuccessOperated(op)
	}
}

func (room *Room) broadcastDispatchCard(card *card.Card) {
	log.Debug(room, "Room.broadcastDispatchCard", card)
	for _, player := range room.players {
		player.OnDispatchCard(card)
	}
}

func (room *Room) broadcastShowDispatchCard(card *card.Card)  {
	log.Debug(room, "broadcastShowDispatchCard", card)
	for _, player := range room.players {
		player.ShowDispatchCard(card)
	}
}

func (room *Room) checkTestCard(card *card.Card, isSystemDisptach bool)  {

	//检查有没有人胡，有就自动胡了
	results := make([]*TestCardResult, len(room.players))
	tmpPlayer := room.curOperator
	for {
		idx := tmpPlayer.idxOfRoom
		results[idx] =  tmpPlayer.TestCard(card, room.lastDropCardOperator)
		if results[idx].CanHu {
			room.switchStatus(RoomStatusEndPlayGame)
			room.lastHuPlayer = room.curOperator
			room.dealPlayerOperate(room.makeHuCardOperate())
			return
		}

		tmpPlayer = room.nextPlayer(tmpPlayer)
		if tmpPlayer == room.curOperator{
			break
		}
	}

	//没有人胡，检查提龙，跑, 扫
	for idx, result := range results {
		if result.CanTiLong {
			op := room.makeTiLongOperate(room.players[idx], card)
			room.dealPlayerOperate(op)
			close(op.ResultCh)
			return
		}
		if result.CanPao {
			op := room.makePaoOperate(room.players[idx], card)
			room.dealPlayerOperate(op)
			close(op.ResultCh)
			return
		}
		if result.CanSao{
			op := room.makeSaoOperate(room.players[idx], card)
			room.dealPlayerOperate(op)
			close(op.ResultCh)
			return
		}
	}

	if isSystemDisptach {//如果是系统派发的牌，没有人胡，也没有跑，扫，提龙，把牌显示给其它所有玩家
		room.broadcastShowDispatchCard(card)
	}

	//检查有没有人能碰
	for idx, result := range results {
		if result.CanPeng {
			timeout := room.waitPlayerPeng(room.players[idx])
			if !timeout {
				return
			}
		}
	}

	//检查当前玩家吃牌
	if len (results[room.curOperator.idxOfRoom].ChiGroup) > 0 {
		waitTimeout := room.waitMasterChi()
		if !waitTimeout {
			return
		}
	}

	nextPlayer := room.nextPlayer(room.curOperator)
	if len(results[nextPlayer.idxOfRoom].ChiGroup) > 0 {
		room.waitMasterNextPlayerChi()
	}
}

func (room *Room) makeDropCardOperate(operator *Player, card *card.Card) *Operate {
	data := &OperateDropCardData{
		Card: card,
	}
	return NewOperateDropCard(operator, data)
}

func (room *Room) makeHuCardOperate() *Operate{
	return NewOperateHu(
		room.curOperator,
		&OperateHuData{
			HuPlayer: room.curOperator,
			FromPlayer: room.lastDropCardOperator,
			Desc: "",
		},
	)
}

func (room *Room) makeTiLongOperate(operator *Player, card *card.Card) *Operate {
	return NewOperateTiLongCard(operator, &OperateTiLongCardData{Card:card})
}

func (room *Room) makePaoOperate(operator *Player, card *card.Card) *Operate {
	return NewOperatePaoCard(operator, &OperatePaoCardData{Card:card})
}

func (room *Room) makeSaoOperate(operator *Player, card *card.Card) *Operate {
	return NewOperateSaoCard(operator, &OperateSaoCardData{Card:card})
}

func (room *Room) switchOperator(player *Player) {
	log.Debug(room, "switchOperator", room.curOperator, "=>", player)
	room.prevOperator = room.curOperator
	room.curOperator = player
}

func (room *Room) String() string {
	if room == nil {
		return "{room=nil}"
	}
	return fmt.Sprintf("{room=%v}", room.GetId())
}