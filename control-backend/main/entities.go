package main

//存一些简化版本的Ceph及其他实体的信息

// CephHost Ceph的Host，如master，node1等，有name和其包含的service
type CephHost struct {
	Hostname string
	Services []string
}

// CephMonitor 代表Ceph的monitor，有名称、排名、地址、活跃的Sessions（近10条数据）
type CephMonitor struct {
	Name         string
	Rank         int
	Address      string
	OpenSessions []int
}

// CephOSD 代表Ceph的OSD
type CephOSD struct {
	Address         string
	HostName        string //osd所在的host的name
	Id              int
	KbAll           int       // 全部空间，kb为单位
	KbUsed          int       // 已使用的空间。kb
	KbUsedData      int       //实际数据使用的kb
	KbUsedMeta      int       //元数据使用的kb
	NumPG           int       //PG（对象组）数量
	State           []string  //状态，一个元素一个状态
	WriteBytes      []float64 //写入数据量（B），近19条数据
	ReadBytes       []float64 //读出数据量（B），近19条数据
	ReadOperations  []float64 //平均读请求数 （/s），近19条数据
	WriteOperations []float64 //平均写请求数 （/s），近19条数据
}

type CephPool struct {
	Name string
	Size int // 副本个数

}
