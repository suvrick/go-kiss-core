package server

const INFO PacketServerType = 5

// INFO(5) "BB"
type Info struct {
	ArrLen    int
	ItemCount uint16
	GameID    uint64
	LoginID   uint64
	NetType   byte
	Name      string
	Sex       byte
	Tag       int
	Referrer  int
	Ddate     int
	Avatar    string
	AvatarID  byte
	Profile   string
	Status    string
}

type Avatar struct {
	Avatar   string
	AvatarID byte
}

/*
	nid?: string; // PlayerInfoParser.NET_ID; [I;
	abstract type: NetType; // PlayerInfoParser.TYPE; [B;
	name?: string; // PlayerInfoParser.NAME; [S;
	sex?: Gender; // PlayerInfoParser.SEX; [B;
	tag?: number; // PlayerInfoParser.TAG; [I;
	referrer?: number; // PlayerInfoParser.REFERRER; [I;
	bdate?: number; // PlayerInfoParser.BDAY; [I;
	avatar?: string; // PlayerInfoParser.PHOTO; [SB;
	avatar_status?: number;
	profile?: string; // PlayerInfoParser.PROFILE; [S;
	status?: string; // PlayerInfoParser.STATUS; [S;
	countryId?: number; // PlayerInfoParser.COUNTRY_ID; [B;
	online?: boolean; // PlayerInfoParser.ONLINE; [B;
	admirer_id?: number; // PlayerInfoParser.ADMIRER_ID; [I;
	admirer_price?: number; // PlayerInfoParser.ADMIRER_PRICE; [I;
	admirer_time_finish?: number; // PlayerInfoParser.ADMIRER_LEFT; [I; it is timestamp
	views?: number; // PlayerInfoParser.VIEWS; [I;
	vip?: number; // PlayerInfoParser.IS_VIP; [B;
	color?: number; // PlayerInfoParser.COLOR; [B;
	kisses?: number; // PlayerInfoParser.KISSES; [II;
	gifts?: number; // PlayerInfoParser.GIFTS; [II;
	kisses_today?: number;
	gifts_today?: number;
	lastGifts?: [ number, number, number ][]; // PlayerInfoParser.PLAYER_GIFTS; ["[IIBI]"; // [source_id:I; gift_id:I; time:I]
	device?: DeviceType; // PlayerInfoParser.DEVICE; [B;
	wedding_id?: number; // PlayerInfoParser.WEDDING_ID; [I;
	achievements?: [ number, number, number ][]; // PlayerInfoParser.ACHIEVEMENTS; ["[III]";
	collections?: [ number, number ][]; // PlayerInfoParser.COLLECTIONS_SETS; ["[BI]";
	avatar_id?: number; // PlayerInfoParser.AVATAR_ID; [B;
	rights?: number; // PlayerInfoParser.RIGHTS; [B;
	register_time?: number; // PlayerInfoParser.REGISTER_TIME; [I;
	logout_time?: number; // PlayerInfoParser.LOGOUT_TIME; [I;
	photos?: string[]; // PlayerInfoParser.PHOTOS; ["[S][B]";
	photos_statuses?: PhotoStatus[];
	frame_id?: number; // PlayerInfoParser.VIP_FRAM
	mask 2
	bridals_place?: number;
	wedlocks_place?: number;
	is_popular?: number;
	rich_place?: number;
	views_place?: number; // PlayerInfoParser.RATING_PLACES; [IIBBI;
	level?: number; // PlayerInfoParser.LEVEL; [I;
	ability_expire?: number;
	ability_type?: number;	// PlayerInfoParser.ABILITY_EXPIRE, ["IB", "ability_expire, ability_type"]]
	rolls_rewarded?: number; // PlayerInfoParser.ROLLS_REWARDED, ["B", "rolls_rewarded"]]
	subscribe_past_days?: number; // PlayerInfoParser.SUBSCRIBE_PAST_DAYS, ["I", "subscribe_past_days"]]
	vip_days_left?: number;
	vip_trial_used?: number;
	league?: number;
	league_common_points?: number;
	league_group_id?: number;
	league_points?: number; // PlayerInfoParser.LEAGUES, ["BI"
	last_complaint_date?: number; // PlayerInfoParser.LAST_COMPLAINT_DATE, "I"
	admire_reward_timestamp?: number;
	interests?: number;
	professions?: number;
	deleted?: boolean;
*/
