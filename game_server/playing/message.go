package playing

import (
	"fmt"
	"peng_hu_zi/game_server/card"
)

type MsgType	int

const  (
	MsgGetInitCards	MsgType = iota + 1
	MsgDispatchCard
	MsgDropCard
	MsgChiCard
	MsgPengCard
	MsgSaoCard
	MsgPaoCard
	MsgTiLongCard
	MsgHu
	MsgEnterRoom
	MsgLeaveRoom
	MsgGameEnd
	MsgRoomClosed
)

type Message struct {
	Type		MsgType
	Owner 	*Player
	Data 	interface{}
}

func (data *Message) String() string {
	if data == nil {
		return "{nil Message}"
	}
	return fmt.Sprintf("{type=%v, Owner=%v}", data.Type, data.Owner)
}

func newMsg(t MsgType, owner *Player, data interface{}) *Message {
	return &Message{
		Owner:	owner,
		Type: t,
		Data: data,
	}
}

//玩家获得初始牌的消息
type GetInitCardsMsgData struct {
	playingCards	*card.PlayingCards
}
func NewGetInitCardsMsg(owner *Player, data *GetInitCardsMsgData) *Message {
	return newMsg(MsgGetInitCards, owner, data)
}


//玩家获得牌的消息
type DispatchCardMsgData struct {
	Card *card.Card
}

func NewDispatchCardMsg(owner *Player, data *DispatchCardMsgData) *Message {
	return newMsg(MsgDispatchCard, owner, data)
}

//玩家出牌的消息
type DropCardMsgData struct {
	Card *card.Card
}
func NewDropCardMsg(owner *Player, data *DropCardMsgData) *Message {
	return newMsg(MsgDropCard, owner, data)
}

//玩家吃牌的消息
type ChiCardMsgData struct {
	ChiPlayer		*Player
	FromPlayer		*Player
	WhatCard		*card.Card
	WhatGroup		*card.Cards
}
func NewChiCardMsg(owner *Player, data *ChiCardMsgData) *Message {
	return newMsg(MsgChiCard, owner, data)
}

//玩家碰牌的消息
type PengCardMsgData struct {
	PengPlayer		*Player
	FromPlayer		*Player
	WhatCard		*card.Card
}
func NewPengCardMsg(owner *Player, data *PengCardMsgData) *Message {
	return newMsg(MsgPengCard, owner, data)
}

//玩家扫牌的消息
type SaoCardMsgData struct {
	SaoPlayer		*Player
	FromPlayer		*Player
	WhatCard		*card.Card
}
func NewSaoCardMsg(owner *Player, data *SaoCardMsgData) *Message {
	return newMsg(MsgSaoCard, owner, data)
}

//玩家跑牌的消息
type PaoCardMsgData struct {
	PaoPlayer		*Player
	FromPlayer		*Player
	WhatCard		*card.Card
}
func NewPaoCardMsg(owner *Player, data *PaoCardMsgData) *Message {
	return newMsg(MsgPaoCard, owner, data)
}


//玩家提龙牌的消息
type TiLongCardMsgData struct {
	TiLongPlayer		*Player
	FromPlayer		*Player
	WhatCard		*card.Card
}
func NewTiLongCardMsg(owner *Player, data *TiLongCardMsgData) *Message {
	return newMsg(MsgTiLongCard, owner, data)
}


type PlayerScore struct {
	P *Player
	Score int
}
//玩家自摸的消息
type HuMsgData struct {
	HuPlayer		*Player			// 胡牌的玩家
	FromPlayer		*Player			// 牌来源于哪个玩家
	Desc			string			// 胡的描述
	PlayerScore 	[]*PlayerScore	// 玩家的分数
}
func NewHuMsg(owner *Player, data *HuMsgData) *Message {
	return newMsg(MsgHu, owner, data)
}

//玩家进入房间的消息
type EnterRoomMsgData struct {
	EnterPlayer *Player
	AllPlayer 	[]*Player
}
func NewEnterRoomMsg(owner *Player, data *EnterRoomMsgData) *Message {
	return newMsg(MsgEnterRoom, owner, data)
}

//玩家离开房间的消息
type LeaveRoomMsgData struct {
	LeavePlayer *Player
	AllPlayer 	[]*Player
}
func NewLeaveRoomMsg(owner *Player, data *LeaveRoomMsgData) *Message {
	return newMsg(MsgLeaveRoom, owner, data)
}

//一盘游戏结束的消息
type GameEndMsgData struct {}
func NewGameEndMsg(owner *Player, data *GameEndMsgData) *Message{
	return newMsg(MsgGameEnd, owner, data)
}

//房间结束的消息
type RoomClosedMsgData struct {}
func NewRoomClosedMsg(owner *Player, data *RoomClosedMsgData) *Message{
	return newMsg(MsgRoomClosed, owner, data)
}