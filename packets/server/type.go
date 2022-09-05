package server

type PacketServerType uint16
type LoginResultType uint16

/*
	NULL = 0,
	HELLO = 1,
	ADMIN_INFO = 2,
	ADMIN_MESSAGE = 3,
	LOGIN = 4,
	INFO = 5,
	INFO_NET = 6,
	BALANCE = 7,
	BUY = 9,
	CONTEST_ITEMS = 10,
	ADMIN_ORDERS_INFO = 11,
	EVENTS = 12,
	REWARDS = 13,
	ADMIRERS = 14,
	GIFT = 16,
	BONUS = 17,
	LEADERBOARDS = 18,
	VIP = 20,
	ROOM_INVITE = 21,
	MOVE = 22,
	BOTTLE_PLAY_DENIED = 24,
	BOTTLE_ROOM = 25,
	BOTTLE_JOIN = 26,
	BOTTLE_LEAVE = 27,
	BOTTLE_LEADER = 28,
	BOTTLE_ROLL = 29,
	BOTTLE_KISS = 30,
	BOTTLE_ENTER = 35,
	CHAT_MESSAGE = 37,
	CHAT_WHISPER = 38,
	HISTORY_CONTACTS = 39,

	IGNORE_LIST = 41,

	BEST = 43,

	FRIENDS = 45,
	TOP = 46,
	LEAGUES_RATING = 47,
	LEAGUE_GROUP = 48,
	LEAGUE_INFO = 49,
	SEARCH = 50,
	LAST_MESSAGE = 51,
	LEAGUES_TIMEOUTS = 52,
	BOTTLE_LIVE_ROOMS = 54,

	MESSAGE_REACTION = 56,
	GIFTS_FOR_ACTIONS_STATS = 65,
	RATING_SIZE = 66,
	WEDDING_PROPOSAL_ANSWER = 67,
	WEDDING_PROPOSAL_CANCEL = 68,
	WEDDING_PROPOSAL_MAKE = 69,
	WEDDING_PROPOSAL_REFUSE = 70,
	WEDDING_PROPOSAL_INFO = 71,
	WEDDING_ADMISSIONS = 72,
	WEDDING_INFO = 73,
	WEDDING_ITEMS = 75,
	WEDDING_SETTLED = 76,
	WEDDING_GUESTS = 77,
	WEDDING_CONTEST = 80,
	WEDDING_JOIN = 81,
	WEDDING_KISS = 82,
	WEDDING_LEADER = 83,
	WEDDING_LEAVE = 84,
	WEDDING_PLAY_DENIED = 85,
	WEDDING_ROLL = 86,
	WEDDING_ROOM = 87,
	WEDDING_STATUS = 88,
	WEDDING_VOW = 89,
	WEDDING_GARTER = 90,
	WEDDING_BOUQUET = 91,
	WEDDING_HAPPY = 93,
	WEDDING_DIVORCE = 94,
	WEDDING_RATING_HAPPY = 95,
	WEDDING_CANCEL = 96,
	WEDDING_START_TIME = 97,		// For self wedding
	ACHIEVEMENT_GET = 106,
	CURIOS = 117,
	CURIOS_GIFT = 118,
	CHAT_HISTORY = 120,
	COLLECTIONS_ASSEMBLE = 121,
	COLLECTIONS_AWARD = 122,
	STATUS_GIFT_STATS = 127,

	COLLECTIONS_POINTS = 130,
	SELF_RICH_UPDATE = 131,
	POSTING_REWARDS = 135,
	PASSION_PASS = 150,
	OFFERS_INFO = 163,

	TIMEOUTS = 173,
	ADMIRE_SERIES = 200,
	PLAYER_ROOM_TYPE = 212,
	PHOTOS_INFO = 213,
	RATING_VIEWS = 216,

	PLAYERS_KISSES = 218,
	PLAYERS_VIEW = 230,
	UPDATE_INFO = 232,
	QUEST = 236,
	MODERATION_LIST = 238,
	OFFERS_BALANCE = 249,

	BIRTHDAY_NOTIFY = 253,
	POPULAR_GIFTS = 254,
	IGNORED = 255,
	HISTORY_MESSAGES = 256,

	BANS = 264,
	ADMIN_BUYINGS_INFO = 295,
	VIDEO_INFO = 296,
	VIDEO_ROOM = 297,
	VIDEO_LISTS = 298,
	VIDEO_QUEUE = 299,
	CHAT_OFFER = 300,
	WINK = 301,
	SPECTATOR_JOIN_LEAVE = 302,
	CAPTCHA = 304,
	CHAT_GIF = 305,
	KISS_PRIORITY = 307,
	KICK_KICKS = 308,
	KICK_SAVE = 309,
	BALANCE_ITEMS = 310,

	FRAMES = 314,
	REWARD_GOT = 315,
	POPULAR_VIDEOS = 317,
	ADMIN_REWARDS_INFO = 318,
	GIFT_BOXES = 319,
	BALLOONS = 320,

	CONTEST_WITH_ACTIONS = 322,
	CONTEST_ACTION = 323,
	LEAGUE_CURRENT_POINTS = 324,
	CONTEST_TETRIS = 325,
	MAX_TYPE
*/
