package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/opsmx/grpc-bidir/tunnel"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	connectedAgentsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "agents_connected",
		Help: "The currently connected agents",
	}, []string{"agent", "protocol"})
)

type Agent interface {
	Send(*httpMessage)
	CancelRequest(*cancelRequest)

	Endpoint() endpoint
	Session() string
	ConnectedAt() uint64
	LastPing() uint64
	LastUse() uint64
}

type Agents struct {
	sync.RWMutex
	m map[endpoint][]Agent
}

type httpMessage struct {
	out chan *tunnel.ASEventWrapper
	cmd *tunnel.HttpRequest
}

type cancelRequest struct {
	id string
}

type agentState struct {
	ep              endpoint
	session         string
	inHTTPRequest   chan *httpMessage
	inCancelRequest chan *cancelRequest
	connectedAt     uint64
	lastPing        uint64
	lastUse         uint64
}

func (s *agentState) Endpoint() endpoint {
	return s.ep
}

func (s *agentState) Session() string {
	return s.session
}

func (s *agentState) ConnectedAt() uint64 {
	return s.connectedAt
}

func (s *agentState) LastPing() uint64 {
	return s.lastPing
}

func (s *agentState) LastUse() uint64 {
	return s.lastUse
}

func (s *agentState) String() string {
	return fmt.Sprintf("(%s, %s, %s)", s.ep.name, s.ep.protocol, s.session)
}

type endpoint struct {
	name     string // The agent name
	protocol string // "kubernetes" or whatever API we are handling
}

func MakeAgents() *Agents {
	return &Agents{
		m: make(map[endpoint][]Agent),
	}
}

func sliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func (a *Agents) AddAgent(state *agentState) {
	agents.Lock()
	defer agents.Unlock()
	agentList, ok := agents.m[state.ep]
	if !ok {
		agentList = make([]Agent, 0)
	}
	agentList = append(agentList, state)
	agents.m[state.ep] = agentList
	log.Printf("Agent %s added, now at %d endpoints", state, len(agentList))
	connectedAgentsGauge.WithLabelValues(state.ep.name, state.ep.protocol).Inc()
}

func (a *Agents) RemoveAgent(state *agentState) {
	agents.Lock()
	defer agents.Unlock()
	agentList, ok := agents.m[state.ep]
	if !ok {
		log.Printf("Attempt to remove unknown agent %s", state)
		return
	}

	close(state.inHTTPRequest)
	close(state.inCancelRequest)

	// TODO: We should always find our entry...
	i := sliceIndex(len(agentList), func(i int) bool { return agentList[i] == state })
	if i != -1 {
		agentList[i] = agentList[len(agentList)-1]
		agentList[len(agentList)-1] = nil
		agentList = agentList[:len(agentList)-1]
		agents.m[state.ep] = agentList
		connectedAgentsGauge.WithLabelValues(state.ep.name, state.ep.protocol).Dec()
	} else {
		log.Printf("Attempt to remove unknown agent %s", state)
	}
	log.Printf("Agent %s removed, now at %d endpoints", state, len(agentList))
}

func (a *Agents) SendToAgent(ep endpoint, message *httpMessage) bool {
	agents.RLock()
	defer agents.RUnlock()
	agentList, ok := agents.m[ep]
	if !ok || len(agentList) == 0 {
		log.Printf("No agents connected for: %s", ep)
		return false
	}
	agent := agentList[rnd.Intn(len(agentList))]
	agent.Send(message)
	return true
}

func (a *Agents) CancelRequest(ep endpoint, message *cancelRequest) bool {
	agents.RLock()
	defer agents.RUnlock()
	agentList, ok := agents.m[ep]
	if !ok || len(agentList) == 0 {
		log.Printf("No agents connected for: %s", ep)
		return false
	}
	agent := agentList[rnd.Intn(len(agentList))]
	agent.CancelRequest(message)
	return true
}

func (s *agentState) Send(message *httpMessage) {
	s.inHTTPRequest <- message
}

func (s *agentState) CancelRequest(message *cancelRequest) {
	s.inCancelRequest <- message
}