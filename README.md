`kubectl top` shows only actual CPU/memory usage (sorted randomly). This tool complements that by aggregating all resource requests and limits across all pods and sorts the resulting list by either cpu or mem requests.

Default output is just pending (unscheduled) pods and nodes, `-pods` lists the pods for each node as well. The sorting criteria can uses absolute values by default, `-rel` uses percentages relative to node capacity instead.

`kubectl config current-context` is used by default, but `-context` and `-config` are provided as well to override.

```
$ bin/ktop
                                             POD:	  CPU(REQ)	  CPU(LIM)	  MEM(REQ)	  MEM(LIM)
                       namebuddy-d56cf8ddb-49jdx:	         1	         1	       1Gi	       1Gi
                       rbacbuddy-cd5d6db84-c9hvk:	         1	         1	      48Mi	      48Mi
                  resourcebuddy-847b679555-p48lh:	         1	         1	      48Mi	      48Mi
                                prometheus-k8s-1:	     8010m	    32010m	   71730Mi	  102450Mi
                                         PENDING:	    11010m	    35010m	   72850Mi	  103570Mi

                                            NODE:	  CPU(REQ)	  CPU(LIM)	  MEM(REQ)	  MEM(LIM)
   gke-tierstaging-central-redis-1-31b1db81-j47g:	1930m (12%)	2700m (17%)	3830Mi (4%)	4574Mi (5%)
   gke-tierstaging-central-redis-1-1e8d628e-7qvc:	3930m (25%)	4700m (30%)	5878Mi (6%)	6622Mi (7%)
 gke-tierstaging-central-default-1-35cc0318-pnhp:	7170m (91%)	10350m (131%)	11176Mi (42%)	15774Mi (59%)
 gke-tierstaging-central-default-1-35cc0318-clnf:	7300m (92%)	8650m (109%)	7816Mi (29%)	8766Mi (33%)
 gke-tierstaging-central-default-1-51c06d36-jvd7:	7353m (93%)	8798m (111%)	6408Mi (24%)	7608Mi (29%)
 gke-tierstaging-central-default-1-51c06d36-t8f6:	7400m (94%)	8850m (112%)	5328Mi (20%)	6328Mi (24%)
 gke-tierstaging-central-default-1-35cc0318-bdtl:	7500m (95%)	10650m (135%)	7872Mi (30%)	8968Mi (34%)
 gke-tierstaging-central-default-1-51c06d36-96jq:	7500m (95%)	11650m (147%)	5824Mi (22%)	6920Mi (26%)
 gke-tierstaging-central-default-1-35cc0318-qjml:	7558m (96%)	12808m (162%)	15624744Ki (57%)	17293864Ki (64%)
 gke-tierstaging-central-default-1-51c06d36-9vjt:	7700m (97%)	10650m (135%)	15522Mi (58%)	16472Mi (62%)
gke-tierstaging-central-high-mem-2-2dbec7af-9265:	15500m (98%)	27650m (174%)	8630Mi (9%)	11116Mi (12%)
gke-tierstaging-central-high-mem-2-18985ed3-s8m6:	15580m (98%)	16650m (105%)	25810Mi (27%)	26760Mi (28%)
```

```
$ bin/ktop -mem -pods
                                                            POD:	  CPU(REQ)	  CPU(LIM)	  MEM(REQ)	  MEM(LIM)
                                      rbacbuddy-cd5d6db84-c9hvk:	         1	         1	      48Mi	      48Mi
                                 resourcebuddy-847b679555-p48lh:	         1	         1	      48Mi	      48Mi
                                      namebuddy-d56cf8ddb-49jdx:	         1	         1	       1Gi	       1Gi
                                               prometheus-k8s-1:	     8010m	    32010m	   71730Mi	  102450Mi
                                                        PENDING:	    11010m	    35010m	   72850Mi	  103570Mi

                                                           NODE:	  CPU(REQ)	  CPU(LIM)	  MEM(REQ)	  MEM(LIM)
---------------------------------------------------------------------------------------------------------------------------------
                  gke-tierstaging-central-redis-1-31b1db81-j47g:	1930m (12%)	2700m (17%)	3830Mi (4%)	4574Mi (5%)
---------------------------------------------------------------------------------------------------------------------------------
       kube-proxy-gke-tierstaging-central-redis-1-31b1db81-j47g:	      100m	         0	         0	         0
                                              trace-proxy-cfgzj:	         0	         0	         0	         0
                                              calico-node-nk885:	      120m	         0	         0	         0
                                         conntrack-bumper-g7k22:	         0	         0	         0	         0
                                  redis-hostvm-configurer-lqnjv:	       10m	      500m	      16Mi	     128Mi
                                            node-exporter-tfjpv:	      100m	      200m	      30Mi	      50Mi
                                      fluentd-gcp-v2.0.10-9zlw9:	      100m	         0	     200Mi	     300Mi
                                        fluentd-forwarder-c4pr7:	      500m	         1	    1536Mi	       2Gi
                                                 dd-agent-5ktmj:	         1	         1	       2Gi	       2Gi
---------------------------------------------------------------------------------------------------------------------------------
                gke-tierstaging-central-default-1-51c06d36-t8f6:	7400m (94%)	8850m (112%)	5328Mi (20%)	6328Mi (24%)
---------------------------------------------------------------------------------------------------------------------------------
                                              calico-node-d2vwl:	      120m	         0	         0	         0
     kube-proxy-gke-tierstaging-central-default-1-51c06d36-t8f6:	      100m	         0	         0	         0
                                         conntrack-bumper-x5wbj:	         0	         0	         0	         0
                                              trace-proxy-qp4lv:	         0	         0	         0	         0
                                            ip-masq-agent-btwzw:	       10m	         0	      16Mi	         0
                                            node-exporter-hxx7z:	      100m	      200m	      30Mi	      50Mi
                                 resourcebuddy-598949798f-zp6jn:	         1	         1	      48Mi	      48Mi
```


## TODO

* [ ] there is sometimes a significant discrepancy between what kubectl node describe shows and ktop
* [ ] add usage metrics https://kubernetes.io/docs/tasks/debug-application-cluster/core-metrics-pipeline/
* [ ] support reverse sort order
* [ ] some filtering will likely be needed for larger clusters
  * [x] allow filtering to only nodes and pods from a given namespace
  * [x] allow filtering to nodes named matching a pattern
