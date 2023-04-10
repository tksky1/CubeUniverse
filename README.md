![](https://pic.imgdb.cn/item/6401b226f144a0100783a222.png)

# CubeUniverse: 基于集群的云原生分布式海量数据存储平台

## 简介

CubeUniverse海量数据存储平台可以做到海量数据高效存储,能够存储结构化、半/非结构化数据，也能够以文件、块、 对象存储，集群高可用、自动弹性伸缩，同时可以基于该存储系统实现有状态容器化应用的自动或半自动灾备。配合该存储后，应用在容器化后可以像传统应用一样使用该存储系统完成数据或文件存储，并保持容器本身灵活及高可用性特征。

CubeUniverse基于成熟的分布式存储系统ceph进行设计，将ceph与新一代容器编排系统kubernetes融合，通过CubeUniverse多项组件疏通ceph与kubernetes的整合与交互，实现云原生、一键部署的便捷特性；

CubeUniverse提供统一的存储平台、快速部署、统一接口、组件监控、自动扩缩容等完善的存储维护体系，为企业用户打造更便捷的存储环境，助力DevOps发展。

## 部署

CubeUniverse平台的部署使用云原生方式，操作简单方便，可以一键部署。  
在满足下述`部署条件`的kubernetes集群上，  在CubeUniverse目录并使用`root`用户执行指令：  

```shell
kubectl create -f deployment/CubeUniverse.yml && watch -n 0.5 kubectl get pod -n cubeuniverse
```
即可完成部署并持续显示组件部署情况。之后CubeUniverse组件会自动下载和部署，并视网速在30分钟内完成整个集群的构建。  

### 部署条件

集群至少包括三个工作节点和一个主节点；

集群的每个工作节点应安装有一块**没有写入数据和文件系统**的卷或磁盘，其大小不小于10GB；

Linux**内核版本4.7以上**；**Kubernetes版本1.20以上**；关闭SELINUX和swap；

集群已经正确部署、所有节点运行正常无污点，且已安装网络插件，所有节点能够访问互联网；

国内Dockerhub连接较慢，建议配置镜像，必要时使用代理。

## 接口

TODO

## 控制面板

TODO

## 文档

TODO