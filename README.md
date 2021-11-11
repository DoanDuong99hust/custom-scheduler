#Custom Scheduler Base on Nodes Status

### Running and testing the Custom Scheduler
1. Open 3 terminals
2. Terminal 1: kubectl proxy
3. Terminal 2: 
   - go build . (from within the scheduler/ folder)
   - ./scheduler.
4. Terminal 3: 
   - kubectl create -f deployments/testcustom.yaml 
   - kubectl get pods -o wide (to see which node its been scheduled on; should be the 'best node' identified by the custom scheduler. See logs in Terminal 2 to verify)
## Implement and modify from
https://github.com/meeramurali/Custom-Kubernetes-Scheduler
## References
[1]  Kubernetes tutorial. Url: https://www.tutorialspoint.com/kubernetes/index.htm

[2]  Kubernetes 101: Pods, Nodes, Containers, and Clusters. Url: https://medium.com/google-cloud/kubernetes-101-pods-nodes-containers-and-clusters-c1509e409e16

[3]  Kubernetes concepts. Url: https://kubernetes.io/docs/concepts/

[4]  Sysbench workload. Url: https://wiki.gentoo.org/wiki/Sysbench#Using_the_memory_workload

[5]  Prometheus. Url: https://prometheus.io/

[6]  Node Exporter. Url: https://prometheus.io/docs/guides/node-exporter/

[7]  A Deep Dive into Kubernetes Metrics. Url: https://blog.freshtracks.io/a-deep-dive-into-kubernetes-metrics-part-2-c869581e9f29

[8]  The Kubernetes Authors. Managing Compute Resources for Containers. url: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/ (accessed: 06.10.2019)

[9]  Eduar Tua. Scheduler Algorithm in Kubernetes. Url: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-scheduling/scheduler_algorithm.md (accessed: 06.10.2019)

[10]  Eduar Tua. The Kubernetes Scheduler. Url: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-scheduling/scheduler.md. (accessed: 06.10.2019)

[11]  4 cool Kubernetes tools for mastering clusters. Url: https://www.infoworld.com/article/3196250/4-cool-kubernetes-tools-for-mastering-clusters.html

[12]  GopherCon 2016: Kelsey Hightower - Building a custom Kubernetes scheduler. Url: https://www.youtube.com/watch?v=IYcL0Un1io0

[13]  Hightower Toy Scheduler. Url: https://github.com/kelseyhightower/scheduler

[14] Google Kubernetes Engine. Url: https://cloud.google.com/kubernetes-engine/
