package cloudflare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrMissingLiveID = errors.New("required live id missing")
)

/**
{
	"result": {
		"uid": "cb9c5e6d9bcf60308dd60a3e64cab64f",
		"rtmps": {
			"url": "rtmps://live.cloudflare.com:443/live/",
			"streamKey": "a3301eb9d763c45360d938985c6e0505kcb9c5e6d9bcf60308dd60a3e64cab64f"
		},
		"rtmpsPlayback": {
			"url": "rtmps://live.cloudflare.com:443/live/",
			"streamKey": "e044ee5e1d744fc204fa64cf2fc9dae8kcb9c5e6d9bcf60308dd60a3e64cab64f"
		},
		"srt": {
			"url": "srt://live.cloudflare.com:778",
			"streamId": "cb9c5e6d9bcf60308dd60a3e64cab64f",
			"passphrase": "7b5195e50e63c4292cef35ae3c796d4ekcb9c5e6d9bcf60308dd60a3e64cab64f"
		},
		"srtPlayback": {
			"url": "srt://live.cloudflare.com:778",
			"streamId": "playcb9c5e6d9bcf60308dd60a3e64cab64f",
			"passphrase": "479d0ce7f7499c46c536be799d873ff1kcb9c5e6d9bcf60308dd60a3e64cab64f"
		},
		"created": "2022-09-08T07:51:53.221284Z",
		"modified": "2022-09-13T12:42:19.813903Z",
		"meta": {
			"name": "\"2ETdSxa723jCNcmKN1VjFapBH11\""
		},
		"status": {
			"current": {
				"reason": "connected",
				"state": "connected",
				"statusEnteredAt": "2022-09-13T12:42:36.679Z",
				"statusLastSeen": "2022-09-13T12:42:36.679Z"
			},
			"history": []
		},
		"recording": {
			"mode": "automatic",
			"requireSignedURLs": false,
			"allowedOrigins": null
		}
	},
	"success": true,
	"errors": [],
	"messages": []
}
*/

type (
	StreamLiveParams struct {
		AccountID string
		LiveID    string
	}

	RTMPs struct {
		URL       string `json:"url,omitempty"`
		StreamKey string `json:"streamKey,omitempty"`
	}

	SRT struct {
		URL        string `json:"url,omitempty"`
		StreamID   string `json:"streamId,omitempty"`
		Passphrase string `json:"passphrase,omitempty"`
	}

	StreamLiveStatusEntry struct {
		Reason          string    `json:"reason,omitempty"`
		State           string    `json:"state,omitempty"`
		StatusEnteredAt time.Time `json:"statusEnteredAt,omitempty"`
		StatusLastSeen  time.Time `json:"statusLastSeen,omitempty"`
	}

	StreamLiveStatus struct {
		Current StreamLiveStatusEntry   `json:"current,omitempty"`
		History []StreamLiveStatusEntry `json:"history,omitempty"`
	}

	StreamLiveRecording struct {
		Mode              string   `json:"mode,omitempty"`
		RequireSignedURLs bool     `json:"requireSignedURLs,omitempty"`
		AllowedOrigins    []string `json:"allowedOrigins,omitempty"`
	}

	StreamLive struct {
		UID           string                 `json:"uid,omitempty"`
		RTMPS         RTMPs                  `json:"rtmps,omitempty"`
		RTMPSPlayback RTMPs                  `json:"rtmpsPlayback,omitempty"`
		SRT           SRT                    `json:"srt,omitempty"`
		SRTPlayback   SRT                    `json:"srtPlayback,omitempty"`
		Created       *time.Time             `json:"created,omitempty"`
		Modified      *time.Time             `json:"modified,omitempty"`
		Meta          map[string]interface{} `json:"meta,omitempty"`
		Status        StreamLiveStatus       `json:"status,omitempty"`
		Recording     StreamLiveRecording    `json:"recording,omitempty"`
	}

	StreamLiveResponse struct {
		Response
		Result StreamLive `json:"result,omitempty"`
	}
)

func (api *API) StreamLiveGet(ctx context.Context, options StreamLiveParams) (StreamLive, error) {
	if options.AccountID == "" {
		return StreamLive{}, ErrMissingAccountID
	}

	if options.LiveID == "" {
		return StreamLive{}, ErrMissingLiveID
	}

	uri := fmt.Sprintf("/accounts/%s/stream/live_inputs/%s", options.AccountID, options.LiveID)
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return StreamLive{}, err
	}

	var streamLiveResponse StreamLiveResponse
	if err := json.Unmarshal(res, &streamLiveResponse); err != nil {
		return StreamLive{}, err
	}

	return streamLiveResponse.Result, nil
}
