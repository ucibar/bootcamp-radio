package handler

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/pion/webrtc/v3"
	"github.com/uCibar/bootcamp-radio/api/response"
	"github.com/uCibar/bootcamp-radio/entity"
	"github.com/uCibar/bootcamp-radio/peerconn"
	"github.com/uCibar/bootcamp-radio/signal"
	"log"
	"net/http"
)

type SessionRepository interface {
	All() []*peerconn.Session
	Add(b *peerconn.Session) error
	Get(id string) (*peerconn.Session, error)
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
	BroadcastID string   `json:"broadcast_id"`
	Title       string   `json:"title"`
	Owner       string   `json:"owner"`
	Subscribers []string `json:"subscribers,omitempty"`
}

type BroadcastListRes struct {
	Broadcasts []*BroadcastRes `json:"broadcasts"`
}

type BroadcastHandler struct {
	*handler
	sessionRepository SessionRepository
}

func NewBroadcastHandler(sessionRepository SessionRepository, logger *log.Logger) *BroadcastHandler {
	return &BroadcastHandler{sessionRepository: sessionRepository, handler: newHandler(logger)}
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

	peer, err := peerconn.NewPublisherPeer(uuid.New().String(), user)
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	err = h.initPeer(peer.PC(), offer)
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	if createReq.BroadcastTitle == "" {
		createReq.BroadcastTitle = fmt.Sprintf("%s's broadcast", user.Username)
	}

	session := peerconn.NewSession(uuid.New().String(), createReq.BroadcastTitle, peer)

	err = h.sessionRepository.Add(session)
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	h.writeResponse(w, response.Success(201, BroadcastCreateRes{
		BroadcastID: session.ID,
		Answer:      signal.Encode(peer.PC().LocalDescription()),
	}))
}

func (h *BroadcastHandler) Join(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := r.Context().Value("user").(*entity.User)
	var joinReq BroadcastJoinReq

	err := h.readBody(r, &joinReq)
	if err != nil {
		h.writeResponse(w, response.InvalidBody(err.Error()))
		return
	}

	session, err := h.sessionRepository.Get(joinReq.BroadcastID)
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	offer := webrtc.SessionDescription{}
	signal.Decode(joinReq.Offer, &offer)

	peer, err := peerconn.NewSubscriberPeer(uuid.New().String(), user)
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	err = session.AddParticipant(peer)
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	err = h.initPeer(peer.PC(), offer)
	if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	h.writeResponse(w, response.Success(200, BroadcastJoinRes{
		BroadcastID:    session.ID,
		BroadcastTitle: session.Title,
		Username:       session.Owner().Username,
		Answer:         signal.Encode(peer.PC().LocalDescription()),
	}))
}

func (h *BroadcastHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	broadcasts := h.sessionRepository.All()

	res := BroadcastListRes{Broadcasts: make([]*BroadcastRes, 0, len(broadcasts))}

	for _, broadcast := range broadcasts {
		res.Broadcasts = append(res.Broadcasts, &BroadcastRes{
			BroadcastID: broadcast.ID,
			Title:       broadcast.Title,
			Owner:       broadcast.Owner().Username,
		})
	}

	h.writeResponse(w, response.Success(200, res))
}

func (h *BroadcastHandler) Info(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	broadcastId := p.ByName("broadcast_id")

	if broadcastId == "" {
		h.writeResponse(w, response.InvalidParam("broadcast_id not provided"))
		return
	}

	broadcast, err := h.sessionRepository.Get(broadcastId)
	if errors.Is(err, entity.ErrBroadcastNotFound) {
		h.writeResponse(w, response.NotFound("broadcast not found"))
		return
	} else if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	subscribers := broadcast.Participants()

	res := BroadcastRes{
		BroadcastID: broadcast.ID,
		Title:       broadcast.Title,
		Owner:       broadcast.Owner().Username,
		Subscribers: make([]string, 0, len(subscribers)),
	}

	for _, user := range subscribers {
		res.Subscribers = append(res.Subscribers, user.Username)
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
