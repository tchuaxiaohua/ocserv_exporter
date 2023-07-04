#### 一、概述

---

关于ocserv 的介绍及安装这里不再赘述，可以参考文章[ocserv 部署](https://wiki.tbchip.com/pages/ldap/oserv/)，这边开发关于`ocserv`的监控客户端，主要是为了监控该服务的用户连接使用情况。

#### 二、项目背景

---

监控指标说明:

* **`ocserv 端口`**: 服务存活性监控
* **`ocserv 客户端在线统计`**: vpn用户连接数统计
* **`ocserv 客户端连接详情`**：vpn用户连接状态

软件版本：

* **`GO版本`**:1.20
* **`ocserv 版本`**: 1.1.1

#### 三、项目介绍

---

> 源码地址[ocserv_exporter](https://github.com/tchuaxiaohua/ocserv_exporter.git)

##### 3.1 Prometheus接入

~~~sh
  - job_name: "ocserv-exporter"
    scrape_interval: 60s
    metrics_path: '/metrics'
    static_configs:
    - targets: ['172.100.10.10:18086']
      labels:
        appname: "ocserv-exporter"
~~~

##### 3.2 告警规则

> 告警规则按需配置

~~~sh
- alert: ocserv down
    expr: sum(ocserv_status) by(instance) == 1 
    for: 5m
    labels:
      severity: critical
    annotations:
       summary: '主机{{ $labels.instance }}，ocserv 连接异常！'
  - alert: ocserv client 带宽使用详情
    expr: sum(ocserv_client_info{bandwidth="receive"}) by(instance,hostname,id) / 1024 > 10
    for: 10m
    labels:
      severity: warring
    annotations:
       summary: 'vpn 主机{{ $labels.instance }}，用户{{ $labels.instance }} 下行带宽使用超过 10MB/sec,已经持续10分钟,当前值: {{ printf "%.2f" $value }} MB/sec'
~~~

