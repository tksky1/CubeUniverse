package cubeControl

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

// CephPool 代表ceph的储存池
type CephPool struct {
	Name       string
	Replica    int    // 副本个数
	PG         int    // PG数
	CreateTime string //创建时间
}

// CephPerformance 包含Ceph的总体状态信息和性能信息
// 只含一份，需刷新
type CephPerformance struct {
	ReadBytesPerSec       int      // 每秒读的bytes
	ReadOperationsPerSec  int      // 每秒读操作数
	WriteBytesPerSec      int      // 每秒写的bytes
	WriteOperationPerSec  int      // 每秒写操作数
	RecoveringBytesPerSec int      // 每秒恢复数据流量
	TotalBytes            int      // 集群可用总容量(bytes)
	TotalUsedBytes        int      // 集群已占用容(bytes)
	HealthStatus          string   // 健康状态总体，如HEALTH_WARN
	HealthStatusDetailed  []string // 健康状态事件（可能有多个）详细说明，如xxx has recently crashed
	HostNum               int      // Host（存储服务节点）的个数
	MonitorNum            int      // 就绪的monitor数量
	OSDReadyNum           int      // 就绪的osd数量
	OSDNotReadyNum        int      // 未就绪的osd数量
	ObjectReplicatedNum   int      // 储存的(包含副本的)对象总数
	ObjectNum             int      // 储存的独立对象总数
	ObjectDegradedNum     int      // 处于降级状态的对象总数
	ObjectMisplacedNum    int      // 处于未归置状态的对象总数
	ObjectNotFoundNum     int      // 处于丢失状态的对象总数
	PoolNum               int      // 存储池总数
}
