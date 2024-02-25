package server

import (
	"bytes"
	"math/big"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/types"
)

const INFO types.PacketServerType = 5

// INFO(5) "BB"
type Info struct {
	//ArrLen types.I
	//ArrLen2 types.I
	Players map[uint64]PlayerInfo
}

func (p Info) String() string {
	return "INFO(5)"
}

func (info *Info) Unmarshal(r *bytes.Reader) error {

	info.Players = map[uint64]PlayerInfo{}

	len, err := leb128.ReadUInt64(r)
	if err != nil {
		return err
	}

	buf := make([]byte, len)
	_, err = r.Read(buf)
	if err != nil {
		return err
	}

	mask, err := leb128.ReadInt64(r)
	if err != nil {
		return err
	}

	_ = mask

	rp := bytes.NewReader(buf)

	playerCount, err := leb128.ReadUInt64(rp)
	if err != nil {
		return err
	}

	_ = playerCount

	playerID, err := leb128.ReadUInt64(rp)
	if err != nil {
		return err
	}

	player := PlayerInfo{}

	use := needParse(mask)

	if use() {
		player.NetID, err = leb128.ReadBigNumber(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.NetType, err = leb128.ReadByte(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Sex, err = leb128.ReadByte(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.CountryId, err = leb128.ReadByte(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Online, err = leb128.ReadByte(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Vip, err = leb128.ReadByte(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Color, err = leb128.ReadByte(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Device, err = leb128.ReadByte(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.AvatarId, err = leb128.ReadByte(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Rights, err = leb128.ReadByte(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.FrameId, err = leb128.ReadByte(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Deleted, err = leb128.ReadByte(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Tag, err = leb128.ReadUInt64(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.BDate, err = leb128.ReadUInt64(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Name, err = leb128.ReadString(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Avatar, err = leb128.ReadString(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Profile, err = leb128.ReadString(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Status, err = leb128.ReadString(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.AdmirerId, err = leb128.ReadUInt64(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.AdmirerPrice, err = leb128.ReadUInt64(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.AdmirerTimeFinish, err = leb128.ReadUInt64(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.AdmireRewardTimestamp, err = leb128.ReadUInt64(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Kisses, err = leb128.ReadUInt64(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.KissesDaily, err = leb128.ReadUInt64(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.Gifts, err = leb128.ReadUInt64(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.GiftsDaily, err = leb128.ReadUInt64(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.WeddingId, err = leb128.ReadByte(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		player.ClubId, err = leb128.ReadByte(rp)
		if err != nil {
			return err
		}
	}

	if use() {
		achievementCount, err := leb128.ReadUInt64(rp)
		if err != nil {
			return nil
		}

		player.Achievements = make([]Achievement, 0)

		for achievementCount > 0 {

			achievement := Achievement{}

			achievement.AchievementID, err = leb128.ReadUInt64(rp)
			if err != nil {
				return err
			}

			achievement.NumberField1, err = leb128.ReadUInt64(rp)
			if err != nil {
				return err
			}

			achievement.NumberField2, err = leb128.ReadUInt64(rp)
			if err != nil {
				return err
			}

			player.Achievements = append(player.Achievements, achievement)

			achievementCount--
		}
	}

	if use() {
		collectionCount, err := leb128.ReadUInt64(rp)
		if err != nil {
			return nil
		}

		player.Collections = make([]Collection, 0)

		for collectionCount > 0 {

			collection := Collection{}

			collection.ElementID, err = leb128.ReadUInt64(rp)
			if err != nil {
				return err
			}

			collection.Count, err = leb128.ReadUInt64(rp)
			if err != nil {
				return err
			}

			player.Collections = append(player.Collections, collection)

			collectionCount--
		}
	}

	info.Players[playerID] = player

	return nil
}

func needParse(mask int64) func() bool {
	bits := big.NewInt(mask)
	l := bits.BitLen()
	i := -1
	return func() bool {
		i++

		if i > l-1 {
			return false
		}

		return bits.Bit(i) == 1
	}
}

// const INFOMASK types.I = 328588

// // ISBSBBIIBIBIIBBIII
type PlayerInfo struct {
	NetID     string
	NetType   byte
	Sex       byte
	CountryId byte
	Online    byte
	Vip       byte
	Color     byte
	Device    byte
	AvatarId  byte
	Rights    byte
	FrameId   byte
	Deleted   byte

	Tag   uint64
	BDate uint64

	Name    string
	Avatar  string
	Profile string
	Status  string

	AdmirerId             uint64
	AdmirerPrice          uint64
	AdmirerTimeFinish     uint64
	AdmireRewardTimestamp uint64

	Kisses      uint64
	KissesDaily uint64
	Gifts       uint64
	GiftsDaily  uint64

	WeddingId byte
	ClubId    byte

	Achievements []Achievement
	Collections  []Collection
}

type Achievement struct {
	AchievementID uint64
	NumberField1  uint64
	NumberField2  uint64
}
type Collection struct {
	ElementID uint64
	Count     uint64
}

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

/*
/.defineField("IBBBUUBBBBBU", [ "netId", "netType", "sex", "countryId", "online", "vip", "color", "device", "avatarId", "rights", "frameId", "deleted" ])
.defineField("I", "tag")
.defineField("I", "bdate")
.defineField("S", "name")
.defineField("S", "avatar")
.defineField("S", "profile")
.defineField("S", "status")
.defineField("IIII", [ "admirerId", "admirerPrice", "admirerTimeFinish", "admireRewardTimestamp" ])
.defineField("IIII", [ "kisses", "kissesDaily", "gifts", "giftsDaily" ])
.defineField("II", [ "weddingId", "clubId" ])
.defineField("[III]", "achievements")
.defineField("[BI]", [ "collections" ])
.defineField("II", [ "registerTime", "logoutTime" ])
.defineField("[S][B]II", [ "photos", "photosStatuses", "totalCountLikesDaily", "totalCountLikes" ])
.defineField("IIIII", [ "bridalsPlace", "wedlocksPlace", "popularPlace", "forbesPlace", "viewsPlace" ])
.defineField("BI", [ "abilityType", "abilityExpire" ])
.defineField("UI", [ "vipTrialUsed", "subscribePastDays" ])
.defineField("BIII", [ "league", "leagueCommonPoints", "leagueGroupId", "leaguePoints" ])
.defineField("I", "lastComplaintDate")
.defineField("I", "cascadeItems");
*/
