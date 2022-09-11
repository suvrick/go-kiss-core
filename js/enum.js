/* eslint-disable @typescript-eslint/indent */
export enum PlayType
{
	Play_new = 0,
	Play_change = 1,
	Play_private = 2
}

export enum WhisperMessageDirection
{
	Incomming = 0,	// General
	Outcomming = 1	// Confirm
}

export enum ChatReactionActionType
{
	Add = 0,
	Remove = 1
}

export enum GiftsChallengeState
{
	WaitingForPlayer = 0,
	Waiting = 1,
	Playing = 2
}

export enum SpectatorAction
{
	Join = 0,
	Leave = 1
}

export enum PlayerRoomType
{
	PLAYER_ROOM_BOTTLE = 1
}

export enum MoveFailReason
{
	OFFLINE = 0,
	NOT_IN_ROOM = 1,
	SAME_ROOM = 2,
	WRONG_AGE = 3,
	BOTTLE_NOT_FOUND = 4,
	BOTTLE_FULL = 5,
	BOTTLE_WAIT = 6,
	WEDDING_PRIVATE = 8
}

export enum PlayFailReason
{
	NO_AVAILABLE_ROOMS = 0,
	ROOM_FULL = 3,
	QUEUE_REMOVED = 6
}

export enum BuyResult
{
	BUY_SUCCESS = 0,
	// BUY_PRICE_CHANGED = 1,
	// BUY_NO_BALANCE = 2,
	// BUY_FAILED = 3
}

export enum ProposalInfoFields	// fields order in WEDDING_PROPOSAL_INFO(71)
{
	TargetId = 0,
	SenderId = 1,
	RingId = 2,
	Message = 3,
	Status = 4,
	Delay = 5
}

export enum ProposalStatus
{
	Proposal = 0,
	Removed = 1
}

export enum ServerProposalAnswer
{
	Accept,
	Reject,
	FailedAccept,
	FailedReject,
	AcceptUnadmire
}

export enum ProposalCancelStatus
{
	Success = 0,
	Failed = 1
}

export enum WeddingAdmissionStatus
{
	Own = 0,
	Invited = 1,
	Denied = 2
}

export enum WeddingStatus
{
	Expired = 0,
	Reject = 1,
	Divorce = 2,
	Canceled = 3,
	Created = 4,
	Ready = 5,
	Vow = 6,
	GarterReady = 7,
	GarterThrown = 8,
	BouquetReady = 9,
	BouquetThrown = 10,
	Completed = 13,
	Contest = 14
}

export enum WeddingItemType
{
	Ring = 0,
	Table,
	Bottle,
	Emcee,
	Garter,
	Bouquet,
	GroomSuit,
	BrideSuit,
	MaxType
}

export enum WeddingEventStatus	// for garter and bouquet
{
	Throw = 0,
	Catch = 1
}

export enum WeddingPlayDeniedReason	// for garter and bouquet
{
	Closed = 0,
	WaitTimeExpired = 1,
	NotActive = 2
}

export enum ProposalMakeStatus
{
	Success = 0,
	Engaged,
	Admire
}

export enum LoginDataIndex
{
	Status = 0,
	InnerId = 1,
	Balance = 2,
	InviterId = 3,
	LogoutTime = 4,
	Flags = 5,
	GamesCount = 6,
	KissesDaily = 7,
	LastPaymentTime = 8,
	SubscribeExpires = 9,
	Params = 10,
	SexSet = 11,
	TutorialStep = 12,
	Tag = 13,
	GlobalTime = 14,
	FirstLogin = 15,
	PhotosHash = 16
}

export enum LoginStatus
{
	Success = 0,
	Failed = 1,
	Exist = 2,
	Blocked = 3,
	WronngVersion = 4,
	NoSex = 5,
	Captcha = 6,
	BlockedByAge = 7,	// TODO договориться с сервером
	NeedVerify = 8,
	Deleted = 9
}

export enum MoveType
{
	MOVE_BOTTLE = 0
}

export enum CounterType
{
	BANK_OPEN = 26,
	BANK_PURCHASE = 27,
	ADVERTISING_AVAILABLE = 530,
	TRANSLATED_CHARS = 612
}

export enum TutorialStepType
{
	NotStarted = 0,
	Started = 1,
	Step1 = 2,
	Step2 = 3,
	Step3 = 4,
	Step4 = 5,
	Step5 = 6,
	MaxType
}

export enum KissAnswer
{
	No = 0,
	Yes = 1
}

export enum ClientProposalAnswer
{
	Yes = 0,
	No = 1
}

export enum VideoAdType
{
	DailyBonus = 0,
	CinemaGift = 1,
	ChatMessage = 2,
	Guest = 3
}

export enum ContestItemGift
{
	Fail = 0,
	Success = 1
}

/* eslint-enable @typescript-eslint/indent */