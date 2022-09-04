package server

const COLLECTIONS_POINTS PacketServerType = 130

// COLLECTIONS_POINTS(130) "I"
type CollectionsPoints struct {
	Points uint16
}
