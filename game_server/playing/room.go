package playing

import (
	"mahjong/game_server/card"
	"mahjong/game_server/log"
	"time"
	"mahjong/game_server/util"
	"fmt"
)

type RoomStatusType int
const (
	RoomStatusWaitAllPlayerEnter	RoomStatusType = iota	// 等待玩家进入房间
	RoomStatusWaitStartPlayGame				// 等待游戏开始
	RoomStatusPlayingGame					// 正在进行游戏，结束后会进入RoomStatusEndPlayGame
	RoomStatusEndPlayGame					// 游戏结束后会回到等待游戏开始状态，或者进入结束房间状态
	RoomStatusRoomEnd							// 房间结束状态，比如东南西北风都打完了
)

type RoomObserver interface {
	OnRoomClosed(room *Room)
}

type Room struct {
	id				uint64					//房间id
	config 			*RoomConfig				//房间配置
	players 		[]*Player				//当前房间的玩家列表
	observers		[]RoomObserver			//房间观察者，需要实现OnRoomClose，房间close的时候会通知它

	roomStatus		RoomStatusType						//房间当前的状态

	firstMasterPlayer *Player				//第一个做东的玩家
	lastHuPlayer	*Player					//最后一次胡牌的玩家
	playedGameCnt	int						//已经玩了的游戏的次数

	//begin playingGameData, reset when start playing game
	cardPool		*card.Pool				//洗牌池
	magicCard		*card.Card				//当前的癞子牌
	masterPlayer	*Player					//做东的玩家，打筛子的玩家
	curOperator	*Player				//获得当前操作的玩家，可能是摸牌，碰牌，杠牌，吃牌，等他出牌
	prevOperator *Player
	quanFeng		*QuanFeng						//当前风圈
	otherPlayerOperate  []*PlayerOperate //当有玩家出牌，其它玩家的操作队列，依据优先级高低处理：胡 > 碰/杠 > 吃
	//end playingGameData, reset when start playing game

	playerOpCh		chan *PlayerOperate		//用户操作的channel

	stop bool
}

func NewRoom(config *RoomConfig) *Room {
	room := &Room{
		id:				util.UniqueId(),
		config:			config,
		players:		make([]*Player, 0),
		cardPool:		card.NewPool(),
		observers:		make([]RoomObserver, 0),
		roomStatus:		RoomStatusWaitAllPlayerEnter,
		quanFeng:		newQuanFeng(card.Feng_CardNo_Dong),
		playedGameCnt:	0,

		otherPlayerOperate:	make([]*PlayerOperate, 0),

		playerOpCh:		make(chan *PlayerOperate, 1024),
	}

	room.init()
	return room
}

func (room *Room) GetId() uint64 {
	return room.id
}

func (room *Room) PlayerOperate(op *PlayerOperate) {
	room.playerOpCh <- op
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
	case RoomStatusPlayingGame:
		room.playingGame()
	case RoomStatusEndPlayGame:
		room.endPlayGame()
	case RoomStatusRoomEnd:
		room.close()
	}
}

func (room *Room) isRoomEnd() bool {
	if room.config.WithQuanFeng {
		if !room.quanFeng.isLastQuanFeng() {//不是最后一圈，肯定没结束
			return false
		}
		return room.computeQuanFeng().isFirstQuanFeng()
	}

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

func (room *Room) checkAllPlayerEnter() {
	length := len(room.players)
	log.Debug(room, "Room.checkAllPlayerEnter, player num :", length, ", need :", room.config.NeedPlayerNum)
	if length >= room.config.NeedPlayerNum {
		room.switchStatus(RoomStatusWaitStartPlayGame)
	}
}

func (room *Room) switchStatus(status RoomStatusType) {
	log.Debug(room, "room status switch,", room.roomStatus, " =>", status)
	room.roomStatus = status
}

func (room *Room) startPlayGame()  {
	log.Debug(room, "Room.startPlayerGame")

	// 重置牌池, 洗牌
	room.cardPool.ReGenerate()

	// 计算癞子牌，如果有的话
	room.computeMagicCard()

	// 选择东家
	room.masterPlayer = room.selectMasterPlayer()
	if room.firstMasterPlayer == nil {
		room.firstMasterPlayer = room.masterPlayer
	}

	//选完东家后，计算圈风
	if room.config.WithQuanFeng {
		room.quanFeng = room.computeQuanFeng()
	}

	// 设定获得牌的玩家的索引为东家
	room.curOperator = room.masterPlayer
	room.prevOperator = nil

	room.otherPlayerOperate = room.otherPlayerOperate[0:0]

	//发初始牌给所有玩家
	room.putInitCardsToPlayers()

	//通知所有玩家手上的牌
	for _, player := range room.players {
		data := &PlayerOperateGetInitCardsData{
			CardsInHand: player.playingCards.GetCardsInHand(),
			MagicCards: player.playingCards.GetMagicCards(),
		}
		player.OnPlayerSuccessOperated(NewPlayerOperateGetInitCards(player, nil, data))
	}

	room.switchStatus(RoomStatusPlayingGame)
}

func (room *Room) playingGame() {
	// 发牌给玩家
	card := room.putCardToPlayer(room.curOperator)
	log.Debug(room, "Room.playingGame put card[ ", card.Name(), "]to", room.curOperator)
	if card == nil {
		room.switchStatus(RoomStatusEndPlayGame)
	} else {
		data := &PlayerOperateGetData{
			Card: card,
		}
		operate := NewPlayerOperateGet(room.curOperator, data)
		room.broadcastPlayerSuccessOperated(operate)
		room.waitCurPlayerOperate()
	}
}

func (room *Room) endPlayGame() {
	room.playedGameCnt++
	log.Debug(room, "Room.endPlayGame cnt :", room.playedGameCnt)
	if room.isRoomEnd() {
		log.Debug(room, "Room.endPlayGame room end")
		room.switchStatus(RoomStatusRoomEnd)
	} else {
		for _, player := range room.players {
			player.OnPlayingGameEnd()
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
			room.switchStatus(RoomStatusRoomEnd) //超时发现没有全部玩家都进入房间了，则结束
			return
		case op := <-room.playerOpCh:
			if op.Op == PlayerOperateEnterRoom || op.Op == PlayerOperateLeaveRoom {
				log.Debug(room, "Room.waitAllPlayerEnter catch operate:", op.String())
				room.dealPlayerOperate(op)
			}
		}
	}
}

//给所有玩家发初始化的13张牌
func (room *Room) putInitCardsToPlayers() {
	log.Debug(room, "Room.initAllPlayer")
	for _, player := range room.players {
		player.Reset()
		for num := 0; num < 13; num++ {
			room.putCardToPlayer(player)
		}
	}
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

//初始化cardool
func (room *Room) init() {
	config := room.config
	if config.WithFengCard {
		room.cardPool.AddFengGenerater()
	}
	if config.WithJianCard {
		room.cardPool.AddJianGenerater()
	}
	if config.WithHuaCard {
		room.cardPool.AddHuaGenerater()
	}
	if config.WithWanCard {
		room.cardPool.AddWanGenerater()
	}
	if config.WithTiaoCard {
		room.cardPool.AddTiaoGenerater()
	}
	if config.WithTongCard {
		room.cardPool.AddTongGenerater()
	}
}

//计算癞子牌
func (room *Room) computeMagicCard() {
	log.Debug(room, "Room.computeMagicCard, HasMagicCard:", room.config.HasMagicCard)
	if !room.config.HasMagicCard {
		return
	}
	cardIdx := room.config.NeedPlayerNum * 13
	card := room.cardPool.At(cardIdx)
	room.magicCard = card.Next()
	log.Debug(room, "Room.computeMagicCard, MagicCard :", room.magicCard.Name())
}

//是否癞子牌
func (room *Room) isMagicCard(card *card.Card) bool {
	if !room.config.HasMagicCard {
		return false
	}
	return card.SameAs(room.magicCard)
}

//发牌给指定玩家
func (room *Room) putCardToPlayer(player *Player) *card.Card {
	card := room.cardPool.PopFront()
	if card == nil {
		return nil
	}
	if room.isMagicCard(card) {
		player.AddMagicCard(card)
	} else {
		player.AddCard(card)
	}
	return card
}

//选择东家
func (room *Room) selectMasterPlayer() *Player {
	log.Debug(room, "Room.selectMasterPlayer")
	if room.playedGameCnt == 0 { //第一盘，随机一个做东
		idx := util.RandomN(len(room.players))
		log.Debug(room, "Room.selectMasterPlayer", room.players[idx])
		return room.players[idx]
	}

	if room.lastHuPlayer == nil {//流局，上一盘没有人胡牌
		log.Debug(room, "Room.selectMasterPlayer", room.masterPlayer)
		return room.masterPlayer
	}

	if !room.config.WithQuanFeng { //不支持圈风，那就谁胡谁做东
		log.Debug(room, "Room.selectMasterPlayer", room.lastHuPlayer)
		return room.lastHuPlayer
	}

	//支持圈风，那就如果他胡了就继续做东，否则他的下一个玩家做东
	if room.masterPlayer == room.lastHuPlayer { //上一次做东的人最后一次胡牌了，继续他做东
		log.Debug(room, "Room.selectMasterPlayer", room.masterPlayer)
		return room.masterPlayer
	}
	next := room.nextPlayer(room.masterPlayer)
	log.Debug(room, "Room.selectMasterPlayer", next)
	return next
}

//等待获得牌的玩家操作, 胡？杠？出牌？如果没有任何操作，超时的话，自动帮他出一张牌
func (room *Room) waitCurPlayerOperate() {
	log.Debug(room, "Room.waitCurPlayerOperate")
	breakTimerTime := time.Duration(0)
	timeout := time.Duration(room.config.WaitPlayerOperateTimeout) * time.Second
	for  {
		timer := timeout - breakTimerTime
		select {
		case <-time.After(timer):
			//超时没有操作，自动帮他出一张牌
			card := room.curOperator.AutoDrop()
			log.Debug(room, "Room.waitCurPlayerOperate auto drop :", card.Name(), ", timeout :", timeout)
			room.curOperator.OperateDrop(card)
		case op := <-room.playerOpCh:
			log.Debug(room, "Room.waitCurPlayerOperate catch operate :", op.String())
			if op.Operator != room.curOperator {
				//不是当前玩家的操作，直接无视
				return
			}
			if op.Op == PlayerOperateChi || op.Op == PlayerOperatePeng || op.Op == PlayerOperateDianPao {
				//当前玩家不可能吃牌、碰牌、点炮胡别人
				return
			}
			room.dealPlayerOperate(op)
		}
	}
}

//当玩家出牌后，等待其他玩家操作
func (room *Room) waitOtherPlayerOperateAfterDrop() {
	log.Debug(room, "Room.waitOtherPlayerOperateAfterDrop")
	breakTimerTime := time.Duration(0)
	timeout := time.Duration(room.config.WaitPlayerOperateTimeout) * time.Second
	for  {
		timer := timeout - breakTimerTime
		select {
		case <- time.After(timer):
			//超时没有其它玩家有任何操作, 设置下一个操作者，继续
			room.curOperator = room.nextPlayer(room.curOperator)
			log.Debug(room, "Room.waitOtherPlayerOperateAfterDrop timeout :", timeout, ", so set next operator :", room.curOperator)
		case op := <-room.playerOpCh :
			log.Debug(room, "Room.waitOtherPlayerOperateAfterDrop operate:", op.String())
			if op.Operator == room.curOperator {//操作者不可能是出牌者，直接无视
				return
			}
			room.dealPlayerOperate(op)
		}
	}
}

//等待碰、吃牌的玩家出牌，超时的话，自动帮他出一张牌
func (room *Room) waitPlayerDrop() {
	log.Debug(room, "Room.waitPlayerDrop")
	breakTimerTime := time.Duration(0)
	timeout := time.Duration(room.config.WaitPlayerOperateTimeout) * time.Second
	for  {
		timer := timeout - breakTimerTime
		select {
		case <- time.After(timer):
			card := room.curOperator.AutoDrop()
			log.Debug(room, "Room.waitPlayerDrop auto drop :", card.Name(), ", timeout :", timeout)
			room.curOperator.OperateDrop(card)
			return
		case op := <-room.playerOpCh :
			log.Debug(room, "Room.waitPlayerDrop operate :", op.String())
			if room.curOperator != op.Operator {
				return
			}
			if op.Op != PlayerOperateDrop {
				return
			}
			room.dealPlayerOperate(op)
		}
	}
}

//取指定玩家的下一个玩家
func (room *Room) nextPlayer(player *Player) *Player {
	idx := 0
	for i, p := range room.players {
		if p == player {
			idx = i
			break
		}
	}
	if idx == len(room.players) - 1 {
		log.Debug(room, "Room.nextPlayer :", room.players[0])
		return room.players[0]
	}
	log.Debug(room, "Room.nextPlayer :", room.players[idx+1])
	return room.players[idx+1]
}

//获取上一次操作的玩家
func (room *Room) getPrevOperator() *Player {
	return room.prevOperator
}

//处理玩家操作
func (room *Room) dealPlayerOperate(op *PlayerOperate) {
	log.Debug(room, "Room.dealPlayerOperate :", op.String())
	switch op.Op {
	case PlayerOperateDianPao :
		if room.config.OnlyZiMo {
			//只能自摸，则不允许非摸牌者胡
			op.Notify <- false
			return
		}
		if data, ok := op.Data.(*PlayerOperateDianPaoData); ok {
			result := op.Operator.DianPao(data.Card)
			log.Debug(room, "Room.dealPlayerOperate DianPao result :", result)
			if result.IsHu {
				room.lastHuPlayer = op.Operator
				room.switchStatus(RoomStatusEndPlayGame)
				op.Notify <- true
				room.broadcastPlayerSuccessOperated(op)
				return
			}
		}
		op.Notify <- false
		return

	case PlayerOperateZiMo :
		if op.Operator != room.curOperator {
			op.Notify <- false
			return
		}
		result := op.Operator.ZiMo()
		log.Debug(room, "Room.dealPlayerOperate ZiMo result :", result)
		if result.IsHu {
			room.lastHuPlayer = op.Operator
			room.switchStatus(RoomStatusEndPlayGame)
			op.Notify <- true
			room.broadcastPlayerSuccessOperated(op)
			return
		}
		op.Notify <- false
		return

	case PlayerOperateDrop:
		if op.Operator != room.curOperator {
			op.Notify <- false
			return
		}
		if data, ok := op.Data.(*PlayerOperateDropData); ok {
			if op.Operator.Drop(data.Card) {//出牌
				log.Debug(room, "Room.dealPlayerOperate Drop card :", data.Card.Name())
				op.Notify <- true
				room.broadcastPlayerSuccessOperated(op)
				//出牌后等待其他玩家操作
				room.waitOtherPlayerOperateAfterDrop()
				return
			}
		}
		op.Notify <- false
		return

	case PlayerOperateChi:
		if !room.config.WithChi {
			op.Notify <- false
			return
		}
		nextPlayer := room.nextPlayer(room.curOperator)
		if nextPlayer != op.Operator {
			op.Notify <- false
			return
		}
		//先加入等待队列，然后选一个优先级最高的操作来处理， todo
		//room.otherPlayerOperate = append(room.otherPlayerOperate, op)
		if chiData, ok := op.Data.(*PlayerOperateChiData); ok {
			if op.Operator.Chi(chiData.Card, chiData.Group) {
				log.Debug(room, "Room.dealPlayerOperate chi card :", chiData.Card.Name(), ", group :", chiData.Group)
				//吃成功了，设定当前玩家为吃牌者，并等待他出牌
				room.prevOperator = room.curOperator
				room.curOperator = op.Operator
				op.Notify <- true
				room.broadcastPlayerSuccessOperated(op)
				room.waitPlayerDrop()
				return
			}
		}
		op.Notify <- false
		return

	case PlayerOperatePeng:
		if !room.config.WithPeng {
			op.Notify <- false
			return
		}
		if data, ok := op.Data.(*PlayerOperatePengData); ok {
			if op.Operator.Peng(data.Card) {
				//碰成功了，设定当前玩家为碰牌者，并等待他出牌
				log.Debug(room, "Room.dealPlayerOperate peng :", data.Card.Name())
				room.prevOperator = room.curOperator
				room.curOperator = op.Operator
				op.Notify <- true
				room.broadcastPlayerSuccessOperated(op)
				room.waitPlayerDrop()
				return
			}
		}
		op.Notify <- false
		return

	case PlayerOperateGang :
		if !room.config.WithGang {
			op.Notify <- false
			return
		}
		if room.cardPool.GetCardNum() <= len(room.players) {//牌数少于人数，不允许杠牌了
			op.Notify <- false
			return
		}
		if data, ok := op.Data.(*PlayerOperateGangData); ok {
			if op.Operator.Gang(data.Card, room.curOperator) {
				room.curOperator.BeGangBy(op.Operator)//设定当前操作者被杠了

				//杠成功了，设定当前玩家为杠牌者
				log.Debug(room, "Room.dealPlayerOperate gang :", data.Card.Name())
				room.prevOperator = room.curOperator
				room.curOperator = op.Operator
				op.Notify <- true
				room.broadcastPlayerSuccessOperated(op)
				return
			}
		}
		op.Notify <- false
		return

	case PlayerOperateEnterRoom:
		if data, ok := op.Data.(*PlayerOperateEnterRoomData); ok {
			if room.addPlayer(op.Operator) { //	玩家进入成功
				log.Debug(room, "Room.dealPlayerOperate player enter :", op.Operator)
				op.Notify <- true
				data.Players = room.players
				room.checkAllPlayerEnter()
				room.broadcastPlayerSuccessOperated(op)
			} else {
				op.Notify <- false //玩家进入失败
			}
			return
		}

	case PlayerOperateLeaveRoom:
		if data, ok := op.Data.(*PlayerOperateLeaveRoomData); ok {
			log.Debug(room, "Room.dealPlayerOperate player leave :", op.Operator)
			room.delPlayer(op.Operator)
			op.Notify <- false
			data.Players = room.players
			room.broadcastPlayerSuccessOperated(op)
			return
		}
	}
}

//计算圈风
func (room *Room) computeQuanFeng() *QuanFeng{
	if !room.config.WithQuanFeng {
		log.Debug(room, "Room.computeQuanFeng", room.quanFeng)
		return room.quanFeng
	}

	if room.lastHuPlayer == nil {//上一盘没有人胡, 圈风不变
		log.Debug(room, "Room.computeQuanFeng", room.quanFeng)
		return room.quanFeng
	}

	if room.lastHuPlayer == room.masterPlayer {//上一次胡牌的玩家是东家，圈风不变
		log.Debug(room, "Room.computeQuanFeng", room.quanFeng)
		return room.quanFeng
	}

	if room.masterPlayer != room.firstMasterPlayer {//没有回到起点，还是同一个圈风
		log.Debug(room, "Room.computeQuanFeng", room.quanFeng)
		return room.quanFeng
	}

	next := room.quanFeng.next()
	log.Debug(room, "Room.computeQuanFeng", next)
	return next
}

func (room *Room) broadcastPlayerSuccessOperated(op *PlayerOperate) {
	log.Debug(room, "Room.broadcastPlayerSuccessOperated :", op.String())
	for _, player := range room.players {
		player.OnPlayerSuccessOperated(op)
	}
}

func (room *Room) String() string {
	if room == nil {
		return "{room=nil}"
	}
	return fmt.Sprintf("{room=%v}", room.GetId())
}