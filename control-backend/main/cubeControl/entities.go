package cubeControl

//存一些简化版本的Ceph及其他实体的信息

// CephHost Ceph的Host，如master，node1等，有name和其包含的service
type CephHost struct {
	Hostname string   `json:"Hostname,omitempty"`
	Services []string `json:"Services,omitempty"`
}

// CephMonitor 代表Ceph的monitor，有名称、排名、地址、活跃的Sessions（近10条数据）
type CephMonitor struct {
	Name         string `json:"Name,omitempty"`
	Rank         int    `json:"Rank"`
	Address      string `json:"Address,omitempty"`
	OpenSessions []int  `json:"OpenSessions,omitempty"`
}

// CephOSD 代表Ceph的OSD
type CephOSD struct {
	Address         string    `json:"Address,omitempty"`
	HostName        string    `json:"HostName,omitempty"` //osd所在的host的name
	Uid             int       `json:"Uid"`
	KbAll           int       `json:"KbAll"`                     // 全部空间，kb为单位
	KbUsed          int       `json:"KbUsed"`                    // 已使用的空间。kb
	KbUsedData      int       `json:"KbUsedData"`                //实际数据使用的kb
	KbUsedMeta      int       `json:"KbUsedMeta"`                //元数据使用的kb
	NumPG           int       `json:"NumPG"`                     //PG（对象组）数量
	State           []string  `json:"State,omitempty"`           //状态，一个元素一个状态
	WriteBytes      []float64 `json:"WriteBytes,omitempty"`      //写入数据量（B），近19条数据
	ReadBytes       []float64 `json:"ReadBytes,omitempty"`       //读出数据量（B），近19条数据
	ReadOperations  []float64 `json:"ReadOperations,omitempty"`  //平均读请求数 （/s），近19条数据
	WriteOperations []float64 `json:"WriteOperations,omitempty"` //平均写请求数 （/s），近19条数据
}

// CephPool 代表ceph的储存池
type CephPool struct {
	Name       string `json:"Name,omitempty"`
	Replica    int    `json:"Replica"`              // 副本个数
	PG         int    `json:"PG"`                   // PG数
	CreateTime string `json:"CreateTime,omitempty"` //创建时间
}

// CephPerformance 包含Ceph的总体状态信息和性能信息
// 只含一份，需刷新
type CephPerformance struct {
	ReadBytesPerSec       int      `json:"ReadBytesPerSec,omitempty"`       // 每秒读的bytes
	ReadOperationsPerSec  int      `json:"ReadOperationsPerSec,omitempty"`  // 每秒读操作数
	WriteBytesPerSec      int      `json:"WriteBytesPerSec,omitempty"`      // 每秒写的bytes
	WriteOperationPerSec  int      `json:"WriteOperationPerSec,omitempty"`  // 每秒写操作数
	RecoveringBytesPerSec int      `json:"RecoveringBytesPerSec,omitempty"` // 每秒恢复数据流量
	TotalBytes            int      `json:"TotalBytes,omitempty"`            // 集群可用总容量(bytes)
	TotalUsedBytes        int      `json:"TotalUsedBytes,omitempty"`        // 集群已占用容(bytes)
	HealthStatus          string   `json:"HealthStatus,omitempty"`          // 健康状态总体，如HEALTH_WARN
	HealthStatusDetailed  []string `json:"HealthStatusDetailed,omitempty"`  // 健康状态事件（可能有多个）详细说明，如xxx has recently crashed
	HostNum               int      `json:"HostNum,omitempty"`               // Host（存储服务节点）的个数
	MonitorNum            int      `json:"MonitorNum,omitempty"`            // 就绪的monitor数量
	OSDReadyNum           int      `json:"OSDReadyNum,omitempty"`           // 就绪的osd数量
	OSDNotReadyNum        int      `json:"OSDNotReadyNum,omitempty"`        // 未就绪的osd数量
	ObjectReplicatedNum   int      `json:"ObjectReplicatedNum,omitempty"`   // 储存的(包含副本的)对象总数
	ObjectNum             int      `json:"ObjectNum,omitempty"`             // 储存的独立对象总数
	ObjectDegradedNum     int      `json:"ObjectDegradedNum,omitempty"`     // 处于降级状态的对象总数
	ObjectMisplacedNum    int      `json:"ObjectMisplacedNum,omitempty"`    // 处于未归置状态的对象总数
	ObjectNotFoundNum     int      `json:"ObjectNotFoundNum,omitempty"`     // 处于丢失状态的对象总数
	PoolNum               int      `json:"PoolNum,omitempty"`               // 存储池总数
}

type CephOSDBucket struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	MaxObjects string `json:"max_object"`
	MaxSize    string `json:"max_size"`
	CreateTime string `json:"time"`
}
