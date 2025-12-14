package netx

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"slices"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/nikmy/algo/syncx"
	"github.com/nikmy/algo/testx/faulty"
	"github.com/nikmy/algo/testx/synctest"
)

type NetworkConfig struct {
	Topology       Topology
	SendLossProb   float64
	RecvLossProb   float64
	NodeFaultProb  float64
	ReorderingProb float64
	DuplicateProb  float64
}

type AssertableRemoteCall struct {
	caller string
	target string
	method string
	input  proto.Message
	output proto.Message
	outErr error
}

func RunStress(t *testing.T, cfg NetworkConfig, checks ...Check) {
	nodes := cfg.Topology.NetNodes

	n := &network{
		adjacent: cfg.Topology.Adjacent,
		nodesMap: make(map[string]nodeWrapper, len(nodes)),
		allNodes: make([]nodeWrapper, len(nodes)),
		nodeChan: make(map[string]faultyChannel, len(nodes)),
	}

	nodeCtrl := faulty.NewController(t, 42)
	nodeCtrl.SetFaultProbability(cfg.NodeFaultProb)

	world := faulty.NewController(t, 7)
	world.SetFaultProbability(0.7)

	for _, node := range nodes {
		mock := nodeWrapper{
			wrld: world,
			Node: node,
			net:  n,
		}

		n.nodesMap[node.Addr()] = mock
		n.nodeChan[node.Addr()] = newFaultyChannel(t, cfg)
		n.allNodes = append(n.allNodes, mock)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	world.Parallel(false)
	allReady := syncx.NewBarrier(len(n.allNodes))
	for i := range n.allNodes {
		node := n.allNodes[i]
		world.Go(func() {
			node.Init()
			allReady.Pass()
			node.Work(ctx, node)
		})
		world.Go(func() {
			for ctx.Err() != nil {
				for _, c := range checks {
					if !c.Valid() {
						world.Fail(c.Report())
					}
				}
				world.Yield()
			}
		})
	}

	world.Wait()
}

type network struct {
	adjacent map[string][]string
	allNodes []nodeWrapper
	nodesMap map[string]nodeWrapper
	nodeChan map[string]faultyChannel
}

func (n *network) makeOpsWithContext(ctx context.Context, ctrl *faulty.Controller) []synctest.Operation {
	ops := make([]synctest.Operation, 0, len(n.allNodes))
	allReady := syncx.NewBarrier(len(n.allNodes)*3 + 1)

	for _, node := range n.allNodes {
		ops = append(ops, synctest.Operation{
			Actors: 1,
			Runner: func() {
				node.Init()
				allReady.Pass()
				node.Work(ctx, node)
			},
		})
		ops = append(ops, synctest.Operation{
			Actors: 1,
			Runner: func() {
				if ctrl.Fault() {
					node.Fail()
				}
			},
		})
	}

	ops = append(ops, synctest.Operation{
		Runner: nil,
	})

	return ops
}

func (n *network) send(rpc AssertableRemoteCall) error {
	channel, ok := n.nodeChan[rpc.target]
	if !ok || !slices.Contains(n.adjacent[rpc.caller], rpc.target) {
		return errors.New("node is unreachable")
	}

	err := channel.send(rpc)

	return err
}

func (n *network) pull(receiver string) *AssertableRemoteCall {
	return n.nodeChan[receiver].pull()
}

type nodeWrapper struct {
	Node

	net  *network
	fail *faulty.Controller
	wrld *faulty.Controller

	calls faultyChannel
	dials faultyChannel
}

func (n nodeWrapper) Dial(_ context.Context, peer string, method string, arg proto.Message, result proto.Message) error {
	// double check here is for
	// avoiding
	if n.wrld.Fault() {
		n.wrld.Yield()
		if n.wrld.Fault() {
			n.wrld.Yield()
		}
	}
	return n.net.send(AssertableRemoteCall{
		caller: n.Addr(),
		target: peer,
		method: method,
		input:  arg,
		output: result,
	})
}

func (n nodeWrapper) Broadcast(ctx context.Context, peers []Node, method string, message proto.Message) ([]Node, error) {
	var (
		successes []Node
		failures  []error
	)

	var ok Ok
	for _, peer := range peers {
		err := n.Dial(ctx, peer.Addr(), method, message, &ok)
		if err != nil {
			failures = append(failures, err)
		} else {
			successes = append(successes, peer)
		}

		if ctx.Err() != nil {
			break
		}
	}

	return successes, errors.Join(failures...)
}
