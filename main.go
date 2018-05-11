package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

const PENDING = "PENDING"

var pending = &nodeTotals{name: PENDING}

var (
	flags            = flag.NewFlagSet("", flag.ExitOnError)
	kubeConfig       = flags.String("config", "$HOME/.kube/config", "kubeconfig file")
	kubeContext      = flags.String("context", "", "cluster context")
	showPods         = flags.Bool("pods", false, "print pod details for each node")
	nodeCapacity     = flags.Bool("cap", false, "print node capacity details (only)")
	sortByMemory     = flags.Bool("mem", false, "sort by memory rather than cpu (requests)")
	sortByPercentage = flags.Bool("rel", false, "sort by percentage relative to node capacity")
	forNamespace     = flags.String("ns", "", "filter pods/nodes down to specific namespace")
	nodesNamed       = flags.String("nodes", "", "filter pods/nodes down to nodes matching the pattern")
)

func getSorter() func([]*nodeTotals) {
	if *sortByMemory {
		if *sortByPercentage {
			return sortByMemReqPercentage
		}
		return sortByMemReq
	}
	if *sortByPercentage {
		return sortByCpuReqPercentage
	}
	return sortByCpuReq
}

func main() {
	flags.Parse(os.Args[1:])

	if *kubeContext == "" {
		cc, err := exec.Command("kubectl", "config", "current-context").Output()
		if err != nil {
			panic(err)
		}
		*kubeContext = string(bytes.TrimSpace(cc))
	}
	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: os.ExpandEnv(*kubeConfig)},
		&clientcmd.ConfigOverrides{CurrentContext: *kubeContext},
	).ClientConfig()
	if err != nil {
		panic(err)
	}

	k8s, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	if *nodesNamed != "" {
		*nodesNamed = "*" + strings.Trim(*nodesNamed, "*") + "*"
	}

	// Collect node info
	totals := map[string]*nodeTotals{PENDING: pending}
	max := 0
	nodes, err := k8s.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, node := range nodes.Items {
		nn := node.Name
		if *nodesNamed != "" {
			ok, err := filepath.Match(*nodesNamed, node.Name)
			if err != nil {
				panic(err)
			}
			if !ok {
				continue
			}
		}
		if max < len(nn) {
			max = len(nn)
		}
		totals[node.Name] = &nodeTotals{
			name:        node.Name,
			allocatable: node.Status.Allocatable,
			capacity:    node.Status.Capacity,
		}
	}
	if *nodeCapacity {
		sorted := nodeList(totals)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].allocatable.Cpu().Cmp(*sorted[j].allocatable.Cpu()) < 0
		})
		format := "%" + fmt.Sprintf("%d.%d", max, max) + "s:\t%10s\t%10s\t%10s\t%10s\n"
		fmt.Printf(format, "NODE", "CPU(ALLOC)", "CPU(CAP)", "MEM(ALLOC)", "MEM(CAP)")
		for _, n := range sorted {
			fmt.Printf(format, n.name,
				n.allocatable.Cpu(), n.capacity.Cpu(),
				n.allocatable.Memory(), n.capacity.Memory())
		}
		return
	}

	// Collect pod totals
	pods, err := k8s.CoreV1().Pods(*forNamespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, pod := range pods.Items {
		nn := pod.Spec.NodeName
		if nn == "" {
			nn = PENDING
		}
		n := totals[nn]
		if n == nil {
			continue
		}
		if *showPods && max < len(pod.Name) {
			max = len(pod.Name)
		}
		p := getPodTotals(&pod)
		(&n.total).Add(&p.total, &n.allocatable)
		n.pods = append(n.pods, p)
	}

	// Sort totals
	sorted := nodeList(totals)
	getSorter()(sorted)
	getSorter()([]*nodeTotals{pending})
	format := "%" + fmt.Sprintf("%d.%d", max, max) + "s:\t%10s\t%10s\t%10s\t%10s"
	divider := strings.Repeat("-", max+66)

	// Output results
	if len(pending.pods) > 0 {
		fmt.Printf(format+"\n", "POD", "CPU(REQ)", "CPU(LIM)", "MEM(REQ)", "MEM(LIM)")
		for _, p := range pending.pods {
			fmt.Println(p.Format(format))
		}
		fmt.Println(pending.Format(format))
		fmt.Println()
	}
	fmt.Printf(format+"\n", "NODE", "CPU(REQ)", "CPU(LIM)", "MEM(REQ)", "MEM(LIM)")
	if *showPods {
		fmt.Println(divider)
	}
	for _, n := range sorted {
		fmt.Println(n.Format(format))
		if !*showPods {
			continue
		}
		fmt.Println(divider)
		for _, p := range n.pods {
			fmt.Println(p.Format(format))
		}
		fmt.Println(divider)
	}
}
