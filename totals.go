package main

import (
	"fmt"

	"gopkg.in/inf.v0"
	v1 "k8s.io/api/core/v1"
	res "k8s.io/apimachinery/pkg/api/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type resource struct {
	req  res.Quantity
	lim  res.Quantity
	cap  *res.Quantity
	reqP *int64 // req/cap percentage
	limP *int64 // lim/cap percentage
}

func (r *resource) Add(r2 *resource, cap *res.Quantity) {
	r.req.Add(r2.req)
	r.lim.Add(r2.lim)
	r.cap = cap
}

func (r *resource) LimP() int64 {
	if r.cap == nil || r.cap.IsZero() {
		return 0
	}
	if r.limP != nil {
		return *r.limP
	}
	r.limP = new(int64)
	*r.limP = percent(&r.lim, r.cap)
	return *r.limP
}

func (r *resource) ReqP() int64 {
	if r.cap == nil || r.cap.IsZero() {
		return 0
	}
	if r.reqP != nil {
		return *r.reqP
	}
	r.reqP = new(int64)
	*r.reqP = percent(&r.req, r.cap)
	return *r.reqP
}

func (r *resource) Lim() string {
	if r.cap == nil || r.cap.IsZero() {
		return r.lim.String()
	}
	return fmt.Sprintf("%s (%d%%)", &r.lim, r.LimP())
}

func (r *resource) Req() string {
	if r.cap == nil || r.cap.IsZero() {
		return r.req.String()
	}
	return fmt.Sprintf("%s (%d%%)", &r.req, r.ReqP())
}

func percent(val, cap *res.Quantity) int64 {
	var perc inf.Dec
	perc.Add(&perc, val.AsDec())
	perc.Mul(&perc, inf.NewDec(100, 0))
	perc.QuoRound(&perc, cap.AsDec(), 0, inf.RoundHalfEven)
	return perc.UnscaledBig().Int64()
}

type total struct {
	cpu resource
	mem resource
}

func (t *total) Add(t2 *total, cap *v1.ResourceList) {
	var cpu, mem *res.Quantity
	if cap != nil {
		cpu = cap.Cpu()
		mem = cap.Memory()
	}
	(&t.cpu).Add(&t2.cpu, cpu)
	(&t.mem).Add(&t2.mem, mem)
}

func (t *total) Format(format string, name string) string {
	return fmt.Sprintf(format,
		name,
		t.cpu.Req(),
		t.cpu.Lim(),
		t.mem.Req(),
		t.mem.Lim(),
	)
}

type podTotals struct {
	name string
	total
}

func (p *podTotals) Format(format string) string {
	return p.total.Format(format, p.name)
}

type nodeTotals struct {
	name        string
	pods        []*podTotals
	allocatable v1.ResourceList
	capacity    v1.ResourceList
	total
}

func (n *nodeTotals) Format(format string) string {
	return n.total.Format(format, n.name)
}

func getContainerTotals(c *v1.Container) *total {
	return &total{
		cpu: resource{
			req: *c.Resources.Requests.Cpu(),
			lim: *c.Resources.Limits.Cpu(),
		},
		mem: resource{
			req: *c.Resources.Requests.Memory(),
			lim: *c.Resources.Limits.Memory(),
		},
	}
}

func getPodTotals(pod *v1.Pod) *podTotals {
	t := &podTotals{name: pod.Name}
	for _, c := range pod.Spec.Containers {
		t.Add(getContainerTotals(&c), nil)
	}
	return t
}

func nodeList(nodes map[string]*nodeTotals) (list []*nodeTotals) {
	for _, n := range nodes {
		if n == pending || len(n.pods) == 0 {
			continue
		}
		list = append(list, n)
	}
	return list
}
