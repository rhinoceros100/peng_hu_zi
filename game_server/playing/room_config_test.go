package playing

import "testing"

func TestRoomConfig_Init(t *testing.T) {
	config := NewRoomConfig()
	err := config.Init("./room_config.json")
	t.Log(err, config)
}
