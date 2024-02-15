package paxos

import (
	"coms4113/hw5/pkg/base"
	"fmt"
)

const (
	Propose = "propose"
	Accept  = "accept"
	Decide  = "decide"
)

type Proposer struct {
	N             int
	Phase         string
	N_a_max       int
	V             interface{}
	SuccessCount  int
	ResponseCount int
	// To indicate if response from peer is received, should be initialized as []bool of len(server.peers)
	Responses []bool
	// Use this field to check if a message is latest.
	SessionId int

	// in case node will propose again - restore initial value
	InitialValue interface{}
}

type ServerAttribute struct {
	peers []base.Address
	me    int

	// Paxos parameter
	n_p int
	n_a int
	v_a interface{}

	// final result
	agreedValue interface{}

	// Propose parameter
	proposer Proposer

	// retry
	timeout *TimeoutTimer
}

type Server struct {
	base.CoreNode
	ServerAttribute
}

func NewServer(peers []base.Address, me int, proposedValue interface{}) *Server {
	response := make([]bool, len(peers))
	return &Server{
		CoreNode: base.CoreNode{},
		ServerAttribute: ServerAttribute{
			peers: peers,
			me:    me,
			proposer: Proposer{
				InitialValue: proposedValue,
				Responses:    response,
			},
			timeout: &TimeoutTimer{},
		},
	}
}

func (server *Server) MessageHandler(message base.Message) []base.Node {
	//TODO: implement it
	// panic("implement me")

	newNodes := []base.Node{}
	// fmt.Println("server attr", serverAttr)
	switch messageType := message.(type) {
	case *ProposeRequest:
		// we should call the acceptor handler of the server and demonstrate
		// the expected behaviors when a Paxos server receives a Propose request from a proposer.
		x := server.handleProposeRequest(message.(*ProposeRequest))
		return x
	case *AcceptRequest:
		return server.handleAcceptRequest(message.(*AcceptRequest))
	case *DecideRequest:
		return server.handleDecideRequest(message.(*DecideRequest))
	case *ProposeResponse:
		// fmt.Println("ProposeResponse")
		return server.handleProposeResponse(message.(*ProposeResponse))
	case *AcceptResponse:
		// fmt.Println("AcceptResponse case")
		return server.handleAcceptResponse(message.(*AcceptResponse))
	case *DecideResponse:
		// fmt.Println("DecideResponse")
		fmt.Println(messageType)

	}

	return newNodes
}

// one 1 new state, reject or accept
func (server *Server) handleProposeRequest(req *ProposeRequest) []base.Node {
	responseMsg := &ProposeResponse{
		CoreMessage: base.MakeCoreMessage(server.Address(), req.CoreMessage.From()),
		SessionId:   req.SessionId,
	}
	// fmt.Printf("handleProposeRequest %+v", req)
	newNode := server.copy()
	if req.N > newNode.n_p {
		newNode.n_p = req.N
		responseMsg.Ok = true
		responseMsg.N_p = req.N
		responseMsg.N_a = newNode.n_a
		responseMsg.V_a = newNode.v_a
	} else {
		responseMsg.Ok = false
		responseMsg.N_p = newNode.n_p
	}

	// rsp := []base.Message{
	// 	responseMsg,
	// }
	// rsp = append(rsp, server.Response...)

	newNode.SetSingleResponse(responseMsg)

	newNodes := []base.Node{
		newNode,
	}
	return newNodes
}

// one 1 new state, reject or accept
func (server *Server) handleAcceptRequest(req *AcceptRequest) []base.Node {
	responseMsg := AcceptResponse{
		CoreMessage: base.MakeCoreMessage(server.Address(), req.CoreMessage.From()),
		SessionId:   req.SessionId,
	}

	// serverAttr := (server.Attribute()).(ServerAttribute)
	n := req.N
	v := req.V
	newNode := server.copy()
	if req.N >= server.n_p {
		newNode.n_p = n
		newNode.n_a = n
		newNode.v_a = v
		responseMsg.Ok = true
		responseMsg.N_p = n
	} else {
		responseMsg.Ok = false
		responseMsg.N_p = n
	}
	// fmt.Printf("messages rsp %+v\n", responseMsg)
	newNode.SetSingleResponse(&responseMsg)
	// fmt.Printf("send AcceptResponse to %#v \n", responseMsg)
	newNodes := []base.Node{
		newNode,
	}

	return newNodes
}

// handle decided request
func (server *Server) handleDecideRequest(req *DecideRequest) []base.Node {
	newNodes := make([]base.Node, 0)

	newNode := server.copy()
	newNode.agreedValue = req.V
	// newNode.proposer.InitialValue = req.V
	newNodes = append(newNodes, newNode)
	// newNodes = append(newNodes, server)
	return newNodes
}

// handle propose response
func (server *Server) handleProposeResponse(res *ProposeResponse) []base.Node {
	i := 0
	for i = range server.peers {
		if server.peers[i] == res.From() {
			break
		}
	}
	newNodes := make([]base.Node, 0)
	// newNodes = append(newNodes, server)
	updatedServer := server.copy()

	// check session
	if res.SessionId != server.proposer.SessionId || server.proposer.Phase != Propose || server.proposer.Responses[i] {
		newNodes = append(newNodes, updatedServer)
		return newNodes
	}
	ok := res.Ok
	updatedServer.proposer.Responses[i] = true

	if ok {
		if res.N_a > updatedServer.proposer.N_a_max {
			updatedServer.proposer.N_a_max = res.N_a
			if res.V_a != nil {
				updatedServer.proposer.V = res.V_a
				updatedServer.proposer.InitialValue = res.V_a
			}
		}
		updatedServer.proposer.SuccessCount += 1
	}
	// if not ok just increase response count
	updatedServer.proposer.ResponseCount += 1
	// fmt.Printf("response count %+v \n", updatedServer.proposer.ResponseCount == 3)

	// Enter Accept phase if majority is reached
	majorityNeeded := len(updatedServer.peers)/2 + 1
	if updatedServer.proposer.ResponseCount >= majorityNeeded {
		// if majority response
		if updatedServer.proposer.SuccessCount >= majorityNeeded {
			//send the accept
			acceptNode := updatedServer.copy()

			rsp := make([]base.Message, 0)
			for _, peer := range updatedServer.peers {
				acceptMsg := &AcceptRequest{
					CoreMessage: base.MakeCoreMessage(updatedServer.Address(), peer),
					N:           acceptNode.proposer.N,
					V:           acceptNode.proposer.V,
					SessionId:   acceptNode.proposer.SessionId,
				}
				rsp = append(rsp, acceptMsg)
			}
			// fmt.Printf("messages %+v\n", rsp[0])
			//update the proposer to new phase
			acceptNode.proposer.Phase = Accept
			acceptNode.proposer.SuccessCount = 0
			acceptNode.proposer.ResponseCount = 0
			acceptNode.proposer.Responses = make([]bool, len(acceptNode.peers))

			acceptNode.CoreNode.SetResponse(rsp)
			newNodes = append(newNodes, acceptNode)
		}
		// Case 1b: Majority not OK, restart Propose
		//1c
		waitingNode := updatedServer.copy()
		waitingNode.proposer.Phase = Propose
		newNodes = append(newNodes, waitingNode)
		// }

	} else {
		// Case 2: Don't have majority responses, wait for more responses
		waitingNode := updatedServer.copy()
		waitingNode.proposer.Phase = Propose
		newNodes = append(newNodes, waitingNode)
	}

	return newNodes
}

// handle accept response
func (server *Server) handleAcceptResponse(res *AcceptResponse) []base.Node {
	i := 0
	for i = range server.peers {
		if server.peers[i] == res.From() {
			break
		}
	}
	newNodes := make([]base.Node, 0)
	// newNodes = append(newNodes, server)
	updatedServer := server.copy()
	if res.SessionId != server.proposer.SessionId || server.proposer.Phase != Accept || server.proposer.Responses[i] {
		// fmt.Println("reject")
		newNodes = append(newNodes, updatedServer)
		return newNodes
	}

	// server.proposer.Responses[i] = true
	// updatedServer := server.copy()
	updatedServer.proposer.Responses[i] = true

	if res.Ok {
		updatedServer.proposer.SuccessCount++
	}
	updatedServer.proposer.ResponseCount++

	majorityNeeded := len(updatedServer.peers)/2 + 1
	if updatedServer.proposer.ResponseCount >= majorityNeeded {
		// if majority response
		if updatedServer.proposer.SuccessCount >= majorityNeeded {
			//send the decided
			decideNode := updatedServer.copy()
			rsp := make([]base.Message, 0)
			for _, peer := range updatedServer.peers {
				acceptMsg := &DecideRequest{
					CoreMessage: base.MakeCoreMessage(updatedServer.Address(), peer),
					V:           decideNode.proposer.V,
					SessionId:   decideNode.proposer.SessionId,
				}
				rsp = append(rsp, acceptMsg)
			}
			//update the proposer to new phase
			decideNode.proposer.Phase = Decide
			decideNode.proposer.SuccessCount = 0
			decideNode.proposer.ResponseCount = 0
			decideNode.proposer.Responses = make([]bool, len(decideNode.peers))

			decideNode.CoreNode.SetResponse(rsp)

			newNodes = append(newNodes, decideNode)
		}
		// else {
		// Case 1b: Majority not OK and recieve all, restart Propose
		// if updatedServer.proposer.ResponseCount == len(server.peers) {
		// 	proposeNode := updatedServer.copy()
		// 	// proposeNode.proposer.N += 1000
		// 	newNodes = append(newNodes, proposeNode)
		// 	// proposeNode.StartPropose()
		// } else {
		// 	//1c
		// fmt.Printf("waaaaa %+v \n", updatedServer)
		waitingNode := updatedServer.copy()
		waitingNode.proposer.Phase = Accept
		newNodes = append(newNodes, waitingNode)

	} else {
		// Case 2: Don't have majority responses, wait for more responses
		waitingNode := updatedServer.copy()
		waitingNode.proposer.Phase = Accept
		newNodes = append(newNodes, waitingNode)
	}
	return newNodes
}

// handle decide response
func (server *Server) handleDecideResponse(res *DecideResponse) []base.Node {
	newNodes := make([]base.Node, 0)
	return newNodes
}

// To start a new round of Paxos.
func (server *Server) StartPropose() {
	//TODO: implement it
	// panic("implement me")
	// 	case 1: We have majority responses (ie. at least n/2+1 responses)
	//     case 1a: Majority OK
	//         case 1aa: all peers responded -> 1 node (proceed to accept)
	//         case 1ab: response(s) awaited -> 2 nodes (proceed to accept and wait for response(s))
	//     case 1b: Majority not OK
	//         (fill this -- as you mentioned: the nodes restart Propose)
	//  case 2: Don't have majority responses (ie. <= n/2 responses) -> 1 node (wait for responses)
	// newNodes := make([]base.Node, 0)
	if server.proposer.InitialValue == nil {
		return
	}
	// if server.proposer.N == 0 {
	// 	server.proposer.N = server.me + 1
	// } else {
	// 	server.proposer.N += 1000
	// }
	server.proposer.N = server.n_p + 1
	server.proposer.Phase = Propose
	server.proposer.ResponseCount = 0
	server.proposer.SuccessCount = 0
	server.proposer.SessionId++

	for i := range server.proposer.Responses {
		server.proposer.Responses[i] = false
	}
	if server.proposer.V == nil {
		server.proposer.V = server.proposer.InitialValue
	}

	// Create and send new ProposeRequest messages to each peer
	rsp := make([]base.Message, 0)
	for _, peer := range server.peers {
		proposeMsg := &ProposeRequest{
			CoreMessage: base.MakeCoreMessage(server.Address(), peer),
			N:           server.proposer.N,
			SessionId:   server.proposer.SessionId,
		}
		rsp = append(rsp, proposeMsg)
	}
	server.SetResponse(rsp)
}

// Returns a deep copy of server node
func (server *Server) copy() *Server {
	response := make([]bool, len(server.peers))
	for i, flag := range server.proposer.Responses {
		response[i] = flag
	}

	var copyServer Server
	copyServer.me = server.me
	// shallow copy is enough, assuming it won't change
	copyServer.peers = server.peers
	copyServer.n_a = server.n_a
	copyServer.n_p = server.n_p
	copyServer.v_a = server.v_a
	copyServer.agreedValue = server.agreedValue
	copyServer.proposer = Proposer{
		N:             server.proposer.N,
		Phase:         server.proposer.Phase,
		N_a_max:       server.proposer.N_a_max,
		V:             server.proposer.V,
		SuccessCount:  server.proposer.SuccessCount,
		ResponseCount: server.proposer.ResponseCount,
		Responses:     response,
		InitialValue:  server.proposer.InitialValue,
		SessionId:     server.proposer.SessionId,
	}

	// doesn't matter, timeout timer is state-less
	copyServer.timeout = server.timeout

	return &copyServer
}

func (server *Server) NextTimer() base.Timer {
	return server.timeout
}

// A TimeoutTimer tick simulates the situation where a proposal procedure times out.
// It will close the current Paxos round and start a new one if no consensus reached so far,
// i.e. the server after timer tick will reset and restart from the first phase if Paxos not decided.
// The timer will not be activated if an agreed value is set.
func (server *Server) TriggerTimer() []base.Node {
	if server.timeout == nil {
		return nil
	}
	newNodes := make([]base.Node, 0)
	if server.proposer.Phase == Propose && server.proposer.SuccessCount == len(server.peers) {
		acceptNode := server.copy()

		rsp := make([]base.Message, 0)
		for _, peer := range acceptNode.peers {
			acceptMsg := &AcceptRequest{
				CoreMessage: base.MakeCoreMessage(acceptNode.Address(), peer),
				N:           acceptNode.proposer.N,
				V:           acceptNode.proposer.V,
				SessionId:   acceptNode.proposer.SessionId,
			}
			rsp = append(rsp, acceptMsg)
		}
		// fmt.Printf("messages %+v\n", rsp[0])
		//update the proposer to new phase
		acceptNode.proposer.Phase = Accept
		acceptNode.proposer.SuccessCount = 0
		acceptNode.proposer.ResponseCount = 0
		acceptNode.proposer.Responses = make([]bool, len(acceptNode.peers))

		acceptNode.CoreNode.SetResponse(rsp)
		newNodes = append(newNodes, acceptNode)
	}
	// fmt.Printf("time out paxos %+v \n", server)
	subNode := server.copy()
	subNode.StartPropose()
	newNodes = append(newNodes, subNode)
	return newNodes //[]base.Node{subNode}
}

func (server *Server) Attribute() interface{} {
	return server.ServerAttribute
}

func (server *Server) Copy() base.Node {
	return server.copy()
}

func (server *Server) Hash() uint64 {
	return base.Hash("paxos", server.ServerAttribute)
}

func (server *Server) Equals(other base.Node) bool {
	otherServer, ok := other.(*Server)

	if !ok || server.me != otherServer.me ||
		server.n_p != otherServer.n_p || server.n_a != otherServer.n_a || server.v_a != otherServer.v_a ||
		(server.timeout == nil) != (otherServer.timeout == nil) {
		return false
	}

	if server.proposer.N != otherServer.proposer.N || server.proposer.V != otherServer.proposer.V ||
		server.proposer.N_a_max != otherServer.proposer.N_a_max || server.proposer.Phase != otherServer.proposer.Phase ||
		server.proposer.InitialValue != otherServer.proposer.InitialValue ||
		server.proposer.SuccessCount != otherServer.proposer.SuccessCount ||
		server.proposer.ResponseCount != otherServer.proposer.ResponseCount {
		return false
	}

	for i, response := range server.proposer.Responses {
		if response != otherServer.proposer.Responses[i] {
			return false
		}
	}

	return true
}

func (server *Server) Address() base.Address {
	return server.peers[server.me]
}
