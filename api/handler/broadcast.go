package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/pion/webrtc/v3"
	"github.com/uCibar/bootcamp-radio/api/response"
	"github.com/uCibar/bootcamp-radio/entity"
	"github.com/uCibar/bootcamp-radio/peerconn"
	"github.com/uCibar/bootcamp-radio/signal"
	"github.com/uCibar/bootcamp-radio/stream"
	"log"
	"net/http"
)

type BroadcastRepository interface {
	All() []*stream.Broadcast
	Add(b *stream.Broadcast) error
	Get(id string) (*stream.Broadcast, error)
}

type BroadcastCreateReq struct {
	BroadcastTitle string `json:"broadcast_title"`
	Offer          string `json:"offer"`
}

type BroadcastJoinReq struct {
	BroadcastID string `json:"broadcast_id"`
	Offer       string `json:"offer"`
}

type BroadcastCreateRes struct {
	BroadcastID string `json:"broadcast_id"`
	Answer      string `json:"answer"`
}

type BroadcastJoinRes struct {
	BroadcastID    string `json:"broadcast_id"`
	BroadcastTitle string `json:"broadcast_title"`
	Username       string `json:"username"`
	Answer         string `json:"answer"`
}

type BroadcastRes struct {
	BroadcastID string `json:"broadcast_id"`
	Title       string `json:"title"`
	Username    string `json:"username"`
}

type BroadcastListRes struct {
	Broadcasts []*BroadcastRes `json:"broadcasts"`
}

type BroadcastHandler struct {
	*handler
	broadcastRepository BroadcastRepository
}

func NewBroadcastHandler(broadcastRepository BroadcastRepository, logger *log.Logger) *BroadcastHandler {
	return &BroadcastHandler{broadcastRepository: broadcastRepository, handler: newHandler(logger)}
}

func (h *BroadcastHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := r.Context().Value("user").(*entity.User)

	var createReq BroadcastCreateReq

	err := h.readBody(r, &createReq)
	if err != nil {
		h.writeResponse(w, response.InvalidBody(err.Error()))
		return
	}

	offer := webrtc.SessionDescription{}
	signal.Decode(createReq.Offer, &offer)

	pc, err := peerconn.NewPublisherPeer()
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	publisher := peerconn.NewPublisher(uuid.New().String(), pc)

	err = h.initPeer(pc, offer)
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	broadcastId := uuid.New().String()
	broadcast := stream.NewBroadcast(broadcastId, publisher, fmt.Sprintf("%s's Broadcast", user.Username), user.Username)

	err = h.broadcastRepository.Add(broadcast)
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	h.writeResponse(w, response.Success(201, BroadcastCreateRes{
		BroadcastID: broadcast.ID,
		Answer:      signal.Encode(pc.LocalDescription()),
	}))
}

func (h *BroadcastHandler) Join(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var joinReq BroadcastJoinReq

	err := h.readBody(r, &joinReq)
	if err != nil {
		h.writeResponse(w, response.InvalidBody(err.Error()))
		return
	}

	broadcast, err := h.broadcastRepository.Get(joinReq.BroadcastID)
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	offer := webrtc.SessionDescription{}
	signal.Decode(joinReq.Offer, &offer)

	pc, err := peerconn.NewSubscriberPeer()
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	subscriber := peerconn.NewSubscriber("test", pc)

	err = broadcast.AddSubscriber(subscriber)
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	err = h.initPeer(pc, offer)
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	h.writeResponse(w, response.Success(200, BroadcastJoinRes{
		BroadcastID:    broadcast.ID,
		BroadcastTitle: broadcast.Title,
		Username:       broadcast.Username,
		Answer:         signal.Encode(pc.LocalDescription()),
	}))
}

func (h *BroadcastHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	broadcasts := h.broadcastRepository.All()

	res := BroadcastListRes{Broadcasts: make([]*BroadcastRes, 0, len(broadcasts))}

	for _, broadcast := range broadcasts {
		res.Broadcasts = append(res.Broadcasts, &BroadcastRes{
			BroadcastID: broadcast.ID,
			Title:       broadcast.Title,
			Username:    broadcast.Username,
		})
	}

	h.writeResponse(w, response.Success(200, res))
}

func (h *BroadcastHandler) initPeer(pc *webrtc.PeerConnection, offer webrtc.SessionDescription) error {
	err := pc.SetRemoteDescription(offer)
	if err != nil {
		return err
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return err
	}

	gc := webrtc.GatheringCompletePromise(pc)

	err = pc.SetLocalDescription(answer)
	if err != nil {
		return err
	}

	<-gc
	return nil
}
