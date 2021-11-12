## Dynamic Merkle Tree

## 背景
我们常见的merkle-tree通常都是二叉树的结构，被经常用在一些如ipfs，p2p下载，blockchain等分布式去中心化的场景之下，通过这种数据结构可以实现高效、快速对比两个大的分片存储对象中有差异的部分

### 存储复杂度分析

如果merkle-tree使用在中心化的server端的时候，不得不考虑merkle-tree的hash节点的存储占用空间，比如我们管理1024个分片数据的时候，需要额外存储


```math
\sum_{i=0}^{10} 2^i = 2047
```
一共2047个节点空间

如果我们尝试使用四叉的merkle-tree

```math
\sum_{i=0}^{5} 4^i = 1365
```
可以看到，现在就只占用1365个节点空间了，相比之前的二叉结构的场景要节省约33%的开销，当然了，伴随而来的是搜索空间会有一个常数倍的下降，因为多叉数增多了嘛...

### 搜索复杂度分析

二叉merkle-tree的搜索复杂度
```math
2log_2N
```

四叉merkle-tree的搜索复杂度
```math
4log_4N
```

## 没有银弹
### 这种适合用在什么场景下？
海量服务器的配置管理中心平台，在服务端中，每台服务器的所有配置都存放在多叉的merkle-tree的数据结构中。
服务器的agent每次利用merkle-tree的特性能快速获取到服务端有变动修改的配置项，这样就不必拉取全量的配置项，大大减少与服务端的请求交互流量

与此同时，服务端要存储海量的服务器的所有配置，并要维护与服务器数量相等个的独立merkle-tree数据结构


在这里尝试提出并实现可以支持用户自定义N叉结构类型的merkle-tree

## 实现

### 实现参考
- [https://github.com/cbergoon/merkletree](https://github.com/cbergoon/merkletree)
### 支持的功能

- 构造MK-tree的实例
- 增加
- 修改

### Todo
- 持久化能力，支持特定二进制文件或者MongoDB存储
- 对比两个MK-tree的hash值