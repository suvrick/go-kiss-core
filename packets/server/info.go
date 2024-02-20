package server

import (
	"bytes"

	"github.com/suvrick/go-kiss-core/types"
)

const INFO types.PacketServerType = 5

// INFO(5) "BB"
type Info struct {
	//ArrLen types.I
	//ArrLen2 types.I
	// Players []PlayerInfo
}

func (p Info) String() string {
	return "INFO(5)"
}

func (info *Info) Unmarshal(r *bytes.Reader) error {
	return nil
}

// const INFOMASK types.I = 328588

// // ISBSBBIIBIBIIBBIII
// type PlayerInfo struct {
// 	GameID types.I
// 	//NetType  types.B
// 	Name types.S
// 	Sex  types.B
// 	//Tag      types.I
// 	//Referrer types.I
// 	//Ddate    types.I
// 	Avatar  models.Avatar
// 	Profile types.S
// 	Status  types.S
// 	Vip     types.B
// 	Kisses  models.Kiss
// }

// func (packet *Info) Use(hiro *models.Hiro, room *models.Room, game interfaces.IGame) error {

// 	player := models.Player{}

// 	if len(packet.Players) > 0 {
// 		player.PlayerID = packet.Players[0].GameID
// 		player.Name = packet.Players[0].Name
// 		player.Avatar = packet.Players[0].Avatar.Avatar
// 		player.Profile = packet.Players[0].Profile
// 		player.Sex = packet.Players[0].Sex
// 		player.Vip = packet.Players[0].Vip
// 		player.Kissed = packet.Players[0].Kisses.Kissed
// 		player.KissedDay = packet.Players[0].Kisses.KissedDay
// 	}

// 	for _, v := range room.Players {
// 		if v.PlayerID == player.PlayerID {
// 			room.Players[player.PlayerID] = &player
// 		}
// 	}

// 	if hiro.ID == player.PlayerID {
// 		hiro.Info = &player
// 		log.Printf("I`m %s, ID: %d\n", hiro.Info.Name, hiro.ID)
// 	}

// 	return nil
// }

/*
	public netId?: string;
	public abstract netType: NetType;
	public sex?: Gender;
	public countryId?: CountryID;
	public online?: boolean;
	public vip?: boolean;
	public color?: number;
	public device?: DeviceType;
	public avatarId?: PhotoID;
	public rights?: number;
	public frameId?: FrameID;
	public deleted?: boolean;

	public tag?: number;
	public bdate?: number;

	public name?: string;
	public avatar?: string;
	public profile?: string;
	public status?: string;

	public admirerId?: PlayerID;
	public admirerPrice?: number;
	public admirerTimeFinish?: number;
	public admireRewardTimestamp?: number;

	public kisses?: number;
	public kissesDaily?: number;
	public gifts?: number;
	public giftsDaily?: number;

	public weddingId?: WeddingID;
	public clubId?: ClubID;

	public achievements?: [ AchievementID, number, number ][];
	public collections?: [ number, number ][];

	public registerTime?: number;
	public logoutTime?: number;

	public photos?: string[];
	public photosStatuses?: PhotoStatus[];

	public bridalsPlace?: number;
	public wedlocksPlace?: number;
	public popularPlace?: number;
	public forbesPlace?: number;
	public viewsPlace?: number;

	public abilityType?: number;
	public abilityExpire?: number;

	public vipTrialUsed?: boolean;
	public subscribePastDays?: number;

	public league?: LeagueType;
	public leagueCommonPoints?: number;
	public leagueGroupId?: number;
	public leaguePoints?: number;

	public lastComplaintDate?: number;
*/
