package main

import (
	"CubeUniverse/universalFuncs"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
	"strings"
	"time"
)

// ceph-api有关的交互内容
// 工具函数不需要被外部调用 只调用功能函数即可

var cephToken string

// const cephApiBase = "https://ceph-dashboard-in-cluster.rook-ceph.svc.cluster.local:8443/"
const cephApiBase = "https://192.168.79.11:30701/" //TODO: 调试后改回来

// GetCephHosts 获取Ceph的Hosts状态
func GetCephHosts() ([]CephHost, error) {
	req, _ := http.NewRequest("GET", cephApiBase+"api/host", nil)
	resJson, err := SendHttpsForJson(req, true)
	if err != nil {
		return nil, errors.New("访问cephAPI出错：" + err.Error())
	}
	var hosts []CephHost
	for i := 0; i < len(resJson.MustArray()); i++ {
		var host CephHost
		hostJson := resJson.GetIndex(i)
		host.Hostname = hostJson.Get("hostname").MustString()
		servicesJson, err := hostJson.Get("services").Array()
		if err != nil {
			return nil, err
		}
		host.Services = []string{}
		for j := 0; j < len(servicesJson); j++ {
			host.Services = append(host.Services, hostJson.Get("services").GetIndex(j).Get("type").MustString())
		}
		hosts = append(hosts, host)
	}
	return hosts, err
}

// GetCephMonitor 获取Ceph的Monitor状态
func GetCephMonitor() (inQuorumMonitor []CephMonitor, outQuorumMonitor []CephMonitor, errr error) {
	req, _ := http.NewRequest("GET", cephApiBase+"api/monitor", nil)
	resJson, err := SendHttpsForJson(req, true)
	if err != nil {
		return nil, nil, errors.New("访问cephAPI出错：" + err.Error())
	}
	inQuorumJson := resJson.Get("in_quorum")
	var inQuorumMons []CephMonitor
	for i := 0; i < len(inQuorumJson.MustArray()); i++ {
		mon := CephMonitor{}
		monJson := inQuorumJson.GetIndex(i)
		mon.Address = monJson.Get("addr").MustString()
		mon.Name = monJson.Get("name").MustString()
		mon.Rank = monJson.Get("rank").MustInt()
		sessionsJson := monJson.Get("stats").Get("num_sessions")
		for j := 0; j < len(sessionsJson.MustArray()); j++ {
			session := sessionsJson.GetIndex(j).GetIndex(1).MustInt()
			mon.OpenSessions = append(mon.OpenSessions, session)
		}
		inQuorumMons = append(inQuorumMons, mon)
	}
	outQuorumJson := resJson.Get("out_quorum")
	var outQuorumMons []CephMonitor
	for i := 0; i < len(outQuorumJson.MustArray()); i++ {
		mon := CephMonitor{}
		monJson := outQuorumJson.GetIndex(i)
		mon.Address = monJson.Get("addr").MustString()
		mon.Name = monJson.Get("name").MustString()
		mon.Rank = monJson.Get("rank").MustInt()
		sessionsJson := monJson.Get("stats").Get("num_sessions")
		for j := 0; j < len(sessionsJson.MustArray()); j++ {
			session := sessionsJson.GetIndex(j).GetIndex(1).MustInt()
			mon.OpenSessions = append(mon.OpenSessions, session)
		}
		outQuorumMons = append(outQuorumMons, mon)
	}
	return inQuorumMons, outQuorumMons, err
}

// GetCephOSD 获取Ceph的OSD的状态
func GetCephOSD() ([]CephOSD, error) {
	req, _ := http.NewRequest("GET", cephApiBase+"api/osd", nil)
	resJson, err := SendHttpsForJson(req, true)
	if err != nil {
		return nil, err
	}
	var osds []CephOSD
	for i := 0; i < len(resJson.MustArray()); i++ {
		osd := CephOSD{}
		osdJson := resJson.GetIndex(i)
		osd.Address = osdJson.Get("public_addr").MustString()
		osd.Id = osdJson.Get("id").MustInt()
		osd.HostName = osdJson.Get("host").Get("name").MustString()
		osdStatJson := osdJson.Get("osd_stats")
		osd.KbAll = osdStatJson.Get("kb").MustInt()
		osd.KbUsed = osdStatJson.Get("kb_used").MustInt()
		osd.KbUsedData = osdStatJson.Get("kb_used_data").MustInt()
		osd.KbUsedMeta = osdStatJson.Get("kb_used_meta").MustInt()
		osd.NumPG = osdStatJson.Get("num_pgs").MustInt()
		osd.State = osdJson.Get("state").MustStringArray()
		statHistoryJson := osdJson.Get("stats_history")
		opInByte := statHistoryJson.Get("op_in_bytes")
		opOutByte := statHistoryJson.Get("op_out_bytes")
		opR := statHistoryJson.Get("op_r")
		opW := statHistoryJson.Get("op_w")
		for i := 0; i < len(opInByte.MustArray()); i++ {
			osd.WriteBytes = append(osd.WriteBytes, opInByte.GetIndex(i).GetIndex(1).MustFloat64())
		}
		for i := 0; i < len(opOutByte.MustArray()); i++ {
			osd.ReadBytes = append(osd.ReadBytes, opOutByte.GetIndex(i).GetIndex(1).MustFloat64())
		}
		for i := 0; i < len(opR.MustArray()); i++ {
			osd.ReadOperations = append(osd.ReadOperations, opR.GetIndex(i).GetIndex(1).MustFloat64())
		}
		for i := 0; i < len(opW.MustArray()); i++ {
			osd.WriteOperations = append(osd.WriteOperations, opR.GetIndex(i).GetIndex(1).MustFloat64())
		}
		osds = append(osds, osd)
	}
	return osds, nil
}

// <-----------工具函数，不需要外部调用------------>

// SendHttpsRequest 工具函数，根据ceph要求发送https请求
func SendHttpsRequest(request *http.Request, withToken bool) (*http.Response, error) {
	if withToken && cephToken == "" {
		err := GetCephToken()
		if err != nil {
			return nil, errors.New("获取CephToken失败！" + err.Error())
		}
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	request.Header.Add("Accept", "application/vnd.ceph.api.v1.0+json")
	request.Header.Add("Content-Type", "application/json")
	if withToken {
		request.Header.Add("Authorization", "Bearer "+cephToken)
	}
	client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if strings.Contains(response.Status, "401") {
		err := GetCephToken()
		if err != nil {
			return nil, errors.New("获取CephToken失败！" + err.Error())
		}
		response, err := SendHttpsRequest(request, withToken)
		if err != nil {
			return nil, err
		}
		return response, err
	}
	return response, err
}

// ParseResponseString 工具函数，将返回体转换成map[string]string
func ParseResponseString(response *http.Response) (map[string]string, error) {
	midResult, err := ParseResponseMap(response)
	result := make(map[string]string)
	for key, value := range midResult {
		result[key] = fmt.Sprintf("%v", value)
	}
	return result, err
}

// ParseResponseMap 工具函数，将返回体转换成map[string]interface{}
func ParseResponseMap(response *http.Response) (map[string]interface{}, error) {
	body, err := io.ReadAll(response.Body)
	body = bytes.ReplaceAll(body, []byte("\n"), []byte(""))
	result := make(map[string]interface{})
	if err == nil {
		err = json.Unmarshal(body, &result)
	}
	return result, err
}

// SetCubeUniverseAccount 工具类，在初始化时设置ceph-api的账号
func SetCubeUniverseAccount() error {
	// 进入tool-box pod创建cubeUniverse账号
	selector := labels.SelectorFromSet(map[string]string{"app": "rook-ceph-tools"})
	toolBoxPods, err := clientSet.CoreV1().Pods("rook-ceph").List(context.TODO(),
		v1.ListOptions{LabelSelector: selector.String()})
	if len(toolBoxPods.Items) == 0 {
		return errors.New(err.Error() + " 获取toolbox pod失败！")
	}
	toolBoxPod := toolBoxPods.Items[0]
	config, _ := rest.InClusterConfig()
	outstd, outerr, err := universalFuncs.ExecInPod(config, "rook-ceph", toolBoxPod.Name,
		"echo \"cubeuniverse\" >> p.txt && ceph dashboard ac-user-create cubeuniverse -i p.txt administrator")
	if outerr != "" || err != nil {
		log.Print(err.Error() + " " + outstd + " 在toolbox执行指令失败！")
	}
	return err
}

// GetCephToken 工具函数，用于get CephAPI的token
func GetCephToken() error {
	// 申请token
	req, err := http.NewRequest("POST", cephApiBase+"api/auth",
		strings.NewReader("{\"username\": \"cubeuniverse\", \"password\": \"cubeuniverse\"}"))
	if err != nil {
		return err
	}
	response, err := SendHttpsRequest(req, false)
	if err != nil {
		return err
	}
	parsedResponse, err := ParseResponseString(response)
	if err != nil {
		return err
	}
	cephToken = parsedResponse["token"]
	return nil
}

// SendHttpsForJson 工具函数，发送https请求并将返回体转为Json格式
func SendHttpsForJson(request *http.Request, withToken bool) (*simplejson.Json, error) {
	res, err := SendHttpsRequest(request, true)
	if err != nil {
		return nil, err
	}
	body, _ := io.ReadAll(res.Body)
	resJson, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}
	if !strings.Contains(res.Status, "OK") {
		err = errors.New("CephAPI请求未正确返回：" + res.Status + ", url：" + request.URL.String())
	}
	return resJson, err
}
