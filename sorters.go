package main

import "sort"

func sortByCpuReq(nodes []*nodeTotals) {
	for _, n := range nodes {
		sort.Slice(n.pods, func(i, j int) bool {
			return n.pods[i].total.cpu.req.Cmp(n.pods[j].total.cpu.req) < 0
		})
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].total.cpu.req.Cmp(nodes[j].total.cpu.req) < 0
	})
}

func sortByCpuReqPercentage(nodes []*nodeTotals) {
	for _, n := range nodes {
		sort.Slice(n.pods, func(i, j int) bool {
			return n.pods[i].total.cpu.req.Cmp(n.pods[j].total.cpu.req) < 0
		})
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].total.cpu.ReqP() < nodes[j].total.cpu.ReqP()
	})
}

func sortByMemReq(nodes []*nodeTotals) {
	for _, n := range nodes {
		sort.Slice(n.pods, func(i, j int) bool {
			return n.pods[i].total.mem.req.Cmp(n.pods[j].total.mem.req) < 0
		})
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].total.mem.req.Cmp(nodes[j].total.mem.req) < 0
	})
}

func sortByMemReqPercentage(nodes []*nodeTotals) {
	for _, n := range nodes {
		sort.Slice(n.pods, func(i, j int) bool {
			return n.pods[i].total.mem.req.Cmp(n.pods[j].total.mem.req) < 0
		})
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].total.mem.ReqP() < nodes[j].total.mem.ReqP()
	})
}
