# /deployment: 用云原生的方式部署CubeUniverse

---

本目录存放CubeUniverse需部署在k8s中的deployment和其他资源。  
使用deployment、job等云原生的方式部署，可以真正实现和k8s进行整合，充分利用k8s各个部分松耦合的特性，并使CubeUniverse相关组件和api-server充分交互，发挥出强大的功能，并使得应用可以一键部署。

要部署CubeUniverse，请在本目录使用k8s提供的命令：
```shell
kubectl create -f CubeUniverse.yml
```
之后CubeUniverse组件将自动开始自检和安装，整个存储平台的搭建应当视网速在三十分钟左右自动全部完成。