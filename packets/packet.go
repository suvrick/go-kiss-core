package packets

import (
	"github.com/suvrick/go-kiss-core/leb128"
)

type Packet struct {
	Len    int
	ID     int
	Type   uint16
	Format string
	Name   string
	Error  error
	Buffer []byte
	Params []interface{}
}

type Reward struct {
	ID    int64
	Count int64
}

type Bot struct {
	Result         int64
	GameID         int64
	Balance        int64
	BalanceHistory []int64
	Name           string
	CanCollect     int64
	BonusDay       int64
	Rewards        []Reward
	RewardsHistory []Reward
	RewardsGot     []int64
}

type Mask struct {
	PacketID uint16
	Index    int
	Type     string
	Name     string
}

// LOGIN(4);  status:B, inner_id:I, balance:I
// BALANCE(7); bottles:I, reason:B
// GAME_REWARDS(13); [id:I, count:I]
// BONUS(17); can_collect:B, day:B
// REWARD_GOT(315); owner_id:I, reward_id:I
var masks = []Mask{
	{Name: "Result", PacketID: 4, Index: 0},
	{Name: "GameID", PacketID: 4, Index: 1},
	{Name: "Balance", PacketID: 4, Index: 2},
	{Name: "Balance", PacketID: 7, Index: 0},
	{Name: "BalanceHistory", PacketID: 4, Index: 2},
	{Name: "BalanceHistory", PacketID: 7, Index: 0},
	{Name: "Rewards", PacketID: 13, Index: 0},
	{Name: "RewardsGot", PacketID: 315, Index: 1},
	{Name: "CanCollect", PacketID: 17, Index: 0},
	{Name: "BonusDay", PacketID: 17, Index: 1},
}

func (p *Packet) GetBuffer(msgID int64) {

	a, err := leb128.Compress(msgID)
	if err != nil {
		p.Error = err
		return
	}

	b := len(p.Buffer) + len(a)

	c, err := leb128.Compress(b)
	if err != nil {
		p.Error = err
		return
	}

	data := make([]byte, 0)
	data = append(data, c...)        // итоговая длина пакета
	data = append(data, a...)        // ID сообщения
	data = append(data, p.Buffer...) // данные

	p.Buffer = data

	//fmt.Println(data)
}
