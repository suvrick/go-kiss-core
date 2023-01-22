package server

import (
	"log"

	"github.com/suvrick/go-kiss-core/interfaces"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/types"
)

const INFO types.PacketServerType = 5

// INFO(5) "BB"
type Info struct {
	//ArrLen types.I
	//ArrLen2 types.I
	Players []PlayerInfo
}

const INFOMASK types.I = 328588

// ISBSBBIIBIBIIBBIII
type PlayerInfo struct {
	GameID types.I
	//NetType  types.B
	Name types.S
	Sex  types.B
	//Tag      types.I
	//Referrer types.I
	//Ddate    types.I
	Avatar  models.Avatar
	Profile types.S
	Status  types.S
	Vip     types.B
	Kisses  models.Kiss
}

func (packet *Info) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {

	player := models.Player{}

	if len(packet.Players) > 0 {
		player.PlayerID = packet.Players[0].GameID
		player.Name = packet.Players[0].Name
		player.Avatar = packet.Players[0].Avatar.Avatar
		player.Profile = packet.Players[0].Profile
		player.Sex = packet.Players[0].Sex
		player.Vip = packet.Players[0].Vip
		player.Kissed = packet.Players[0].Kisses.Kissed
		player.KissedDay = packet.Players[0].Kisses.KissedDay
	}

	for _, v := range room.Players {
		if v.PlayerID == player.PlayerID {
			room.Players[player.PlayerID] = &player
		}
	}

	if hiro.ID == player.PlayerID {
		hiro.Info = &player
		log.Printf("I`m %s, ID: %d\n", hiro.Info.Name, hiro.ID)
	}

	return nil
}

/*
0		.defineField("I", "nid")
1		.defineField("B", "type")
2		.defineField("S", "name")
3		.defineField("B", "sex")
4		.defineField("I", "tag")
5		.defineField("I", "referrer")
6		.defineField("I", "bdate")
7		.defineField("SB", [ "avatar", "avatar_status" ])
8		.defineField("S", "profile")
9		.defineField("S", "status")
10		.defineField("B", "countryId")
11		.defineField("B", "online")
12		.defineField("I", "admirer_id")
13		.defineField("I", "admirer_price")
14		.defineField("I", "admirer_time_finish")
15		.defineField("I", "views")
16		.defineField("B", "vip")
17		.defineField("B", "color")
18		.defineField("II", [ "kisses", "kisses_today" ])
19		.defineField("II", [ "gifts", "gifts_today" ])
20		.defineField("[III]", "lastGifts") //[source_id:I, gift_id:I, time:I]
		.defineField("B", "device")
		.defineField("I", "wedding_id")
		.defineField("[III]", "achievements")
		.defineField("[BI]", [ "collections" ])
		.defineField("B", "avatar_id")
		.defineField("B", "rights")
		.defineField("I", "register_time")
		.defineField("I", "logout_time")
		.defineField("[S][B]", [ "photos", "photos_statuses" ])
		.defineField("IIBII", [ "bridals_place", "wedlocks_place", "is_popular", "rich_place", "views_place" ])
		.defineField("B", "frame_id");
*/
