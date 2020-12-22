# ec2-network-monitor

This application exposes all of the metrics that `ethtool` provides as
Prometheus metrics.

In the context of Kubernetes, run this application as a `DaemonSet` so that one
runs on every node of your cluster, configured with `hostNetworking: true`.
This will allow you to capture the metrics from the EC2 host directly.

For more information about what `ethtool` exposes see [Monitoring EC2 Network
Performance
ENA](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/monitoring-network-performance-ena.html).

## Install

```
helm repo add ec2nm https://jfreeland.github.io/ec2-network-monitor
helm install -n kube-system ec2-network-monitor ec2nm/ec2-network-monitor
```

## Configuration

| Parameter   | Description                              | Default                                     |
| ----------- | ---------------------------------------- | ------------------------------------------- |
| annotations | Annotations to set on the DaemonSet pods | Configured for standard Prometheus scraping |
