package bgp

import "net/netip"

type BgpPeer struct {
	LocalNode    *BgpNode
	RemoteNode   *BgpNode
	LocalPrefix  netip.Prefix
	RemotePrefix netip.Prefix
	LocalIface   string
	RemoteIface  string
}

type BgpTopology struct {
	peerMap map[string][]BgpPeer
}

func (t *BgpTopology) GetPeers(node string) []BgpPeer {
	return t.peerMap[node]
}

// Build BgpTopology with BgpPeers, edges are treated as undirected
type BgpTopologyBuilder struct {
	peers []BgpPeer
}

func (b *BgpTopologyBuilder) AddPeer(peer BgpPeer) {
	b.peers = append(b.peers, peer)
}

func (b *BgpTopologyBuilder) Build() *BgpTopology {
	t := &BgpTopology{
		peerMap: make(map[string][]BgpPeer),
	}

	for _, v := range b.peers {
		localName := v.LocalNode.name
		remoteName := v.RemoteNode.name

		t.peerMap[localName] = append(t.peerMap[localName], v)

		// Reverse
		revPeer := BgpPeer{
			LocalNode:    v.RemoteNode,
			RemoteNode:   v.LocalNode,
			LocalPrefix:  v.RemotePrefix,
			RemotePrefix: v.LocalPrefix,
			LocalIface:   v.RemoteIface,
			RemoteIface:  v.LocalIface,
		}

		t.peerMap[remoteName] = append(t.peerMap[remoteName], revPeer)
	}

	return t
}
