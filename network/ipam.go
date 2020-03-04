package network

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"net"
	"os"
	"path"
	"strings"
)

const ipamDefaultAllocatorPath = "/var/run/mydocker/network/ipam/subnet.json"

type IPAM struct {
	SubnetAllocatorPath string
	Subnets             *map[string]string
}

var ipAllocator = &IPAM{
	SubnetAllocatorPath: ipamDefaultAllocatorPath,
}

/*
加载数据
就是把 /var/run/mydocker/network/ipam/subnet.json 里面的数据加载到 IPAM 的 Subnets 中，方便操作
*/
func (ipam *IPAM) load() error {
	// 如果该文件不存在 返回
	if _, err := os.Stat(ipam.SubnetAllocatorPath); err != nil {
		log.Errorf("load error err:%v", err)
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	// 打开该文件
	subnetConfigFile, err := os.Open(ipam.SubnetAllocatorPath)
	defer subnetConfigFile.Close()
	if err != nil {
		return err
	}
	// 将该文件的内容读到subnetJson中
	subnetJson := make([]byte, 2000)
	n, err := subnetConfigFile.Read(subnetJson)
	if err != nil {
		return err
	}

	//log.Printf("n:%d\n", n)
	//log.Println(subnetJson)

	// 将subnetJson中内容加载到ipam.Subnets中
	err = json.Unmarshal(subnetJson[:n], ipam.Subnets)
	if err != nil {
		log.Errorf("Error dump allocation info, %v", err)
		return err
	}
	return nil
}

/*
持久化数据 dump
就是把 IPAM 中的 Subnets 持久化到 IPAM 的 SubnetAllocatorPath（/var/run/mydocker/network/ipam/subnet.json）
因为 Subnets 里面保存着容器网络的地址分配情况
*/
func (ipam *IPAM) dump() error {
	// ipam.SubnetAllocatorPath 文件夹与文件分离开
	ipamConfigFileDir, _ := path.Split(ipam.SubnetAllocatorPath)
	if _, err := os.Stat(ipamConfigFileDir); err != nil {
		// 如果不存在 就逐级生成文件夹
		if os.IsNotExist(err) {
			os.MkdirAll(ipamConfigFileDir, 0644)
		} else {
			return err
		}
	}
	// O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	// O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	// O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file when opened.
	// 如果该文件不存在 就创建一个
	subnetConfigFile, err := os.OpenFile(ipam.SubnetAllocatorPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	defer subnetConfigFile.Close()
	if err != nil {
		return err
	}
	// 将ipam.Subnets的内容持久化
	// 也就是将所有网络的分配情况保存到该文件中
	ipamConfigJson, err := json.Marshal(ipam.Subnets)
	if err != nil {
		return err
	}

	_, err = subnetConfigFile.Write(ipamConfigJson)
	if err != nil {
		return err
	}

	return nil
}

/*
地址分配：
1. 先从 /var/run/mydocker/network/ipam/subnet.json 中加载数据到 ipam 的 Subnets，如果该文件不存在，Subnets是一个空 map，里面什么网络信息都没有。
2. 根据 bitmap 分配 ip
3. 将已经有数据的 Subnets 持久化到 /var/run/mydocker/network/ipam/subnet.json 中
*/
func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP, err error) {
	// 存放网段中地址分配信息的数组
	// 无论ipamDefaultAllocatorPath是否存在都先new一个
	ipam.Subnets = &map[string]string{}

	// 从文件中加载已经分配的网段信息
	err = ipam.load()
	if err != nil {
		log.Errorf("Error dump allocation info, %v", err)
	}
	// 得到网络号
	_, subnet, _ = net.ParseCIDR(subnet.String())

	//log.Printf("Allocate subnet:%s, ipam.Subnets:%v\n", subnet, ipam.Subnets)
	// one表示前缀的个数 size表示ip地址的个数 ipv4==>size=32
	one, size := subnet.Mask.Size()

	//log.Printf("Allocate one:%d, size:%d\n", one, size)

	// 如果该网络还不在ipam.Subnets中, 则初始化一个
	// 那怎么知道该网络有多少个ip地址呢 size-one就表示主机号占的位数 2的(size-one)方就有多少个主机ip
	if _, exist := (*ipam.Subnets)[subnet.String()]; !exist {
		(*ipam.Subnets)[subnet.String()] = strings.Repeat("0", 1<<uint8(size-one))
	}

	//log.Printf("Allocate one:%s\n", (*ipam.Subnets)[subnet.String()])

	for c := range (*ipam.Subnets)[subnet.String()] {
		// 如果第c个ip没有被分配 则分配
		if (*ipam.Subnets)[subnet.String()][c] == '0' {
			ipalloc := []byte((*ipam.Subnets)[subnet.String()])
			ipalloc[c] = '1'
			(*ipam.Subnets)[subnet.String()] = string(ipalloc)

			// 查一下c 用32位如何表示
			ip = subnet.IP
			for t := uint(4); t > 0; t -= 1 {
				[]byte(ip)[4-t] += uint8(c >> ((t - 1) * 8))
			}
			ip[3] += 1
			break
		}
	}
	// 持久化数据
	ipam.dump()
	return
}

func (ipam *IPAM) Release(subnet *net.IPNet, ipaddr *net.IP) error {
	ipam.Subnets = &map[string]string{}

	_, subnet, _ = net.ParseCIDR(subnet.String())

	err := ipam.load()
	if err != nil {
		log.Errorf("Error dump allocation info, %v", err)
	}

	c := 0
	releaseIP := ipaddr.To4()
	releaseIP[3] -= 1
	for t := uint(4); t > 0; t -= 1 {
		c += int(releaseIP[t-1]-subnet.IP[t-1]) << ((4 - t) * 8)
	}

	ipalloc := []byte((*ipam.Subnets)[subnet.String()])
	ipalloc[c] = '0'
	(*ipam.Subnets)[subnet.String()] = string(ipalloc)

	ipam.dump()
	return nil
}
