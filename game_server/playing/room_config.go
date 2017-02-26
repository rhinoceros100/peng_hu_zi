package playing

import "peng_hu_zi/util"

type RoomConfig struct {
	NeedPlayerNum				int        `json:"need_player_num"`
	WaitPlayerEnterRoomTimeout	int        `json:"wait_player_enter_room_timeout"`
	WaitPlayerOperateTimeout	int        `json:"wait_player_operate_timeout"`
	MaxPlayGameCnt			int            `json:"max_play_game_cnt"`	//不支持圈风的时候，最大的游戏局数
}

func NewRoomConfig() *RoomConfig {
	return &RoomConfig{}
}

func (config *RoomConfig) Init(file string) error {
	return util.InitJsonConfig(file, config)
}