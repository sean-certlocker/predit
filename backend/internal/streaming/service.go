package streaming

import (
	"time"

	"github.com/livekit/protocol/auth"
)

type LiveKitService struct {
	apiKey    string
	apiSecret string
}

func NewLiveKitService(apiKey, apiSecret string) *LiveKitService {
	return &LiveKitService{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

func (s *LiveKitService) CreateJoinToken(roomName, identity string) (string, error) {
	at := auth.NewAccessToken(s.apiKey, s.apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomName,
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	return at.ToJWT()
}
