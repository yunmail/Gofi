package environment

import (
	"context"
	"flag"
	"fmt"
	"gofi/ent"
	//import sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//version ,will be replaced at compile time by [-ldflags="-X 'gofi/environment.Version=vX.X.X'"]
var version string = "UNKOWN VERSION"

const (
	//DefaultPort default port to listen Gofi监听的默认端口号
	DefaultPort = "8080"
	portUsage   = "port to expose web services"
	ipUsage     = "server side ip for web client to request,default is lan ip"
)

//Environment 上下文对象
type Environment struct {
	Version           string
	Port              string
	DatabaseName      string
	AppName           string
	ServerAddress     string
	ServerIP          string //ServerIP server side ip for web client to request,default is lan ip
	WorkDir           string
	DefaultStorageDir string
	CustomStorageDir  string
	LogDir            string
	DatabaseFilePath  string
	OrmClient         *ent.Client
	configuration     *ent.Configuration
}

var instance = new(Environment)
var isFlagBind = false

func init() {
	bindFlags()
}

func bindFlags() {
	if !isFlagBind {
		flag.StringVar(&instance.Port, "port", DefaultPort, portUsage)
		flag.StringVar(&instance.Port, "p", DefaultPort, portUsage+" (shorthand)")
		flag.StringVar(&instance.ServerIP, "ip", "", ipUsage)
		isFlagBind = true
	}
}

//InitContext 初始化Context,只能初始化一次
func InitContext() {
	flag.Parse()
	instance.Version = version
	instance.AppName = "gofi"
	instance.WorkDir = instance.getWorkDirectoryPath()
	instance.DatabaseName = instance.AppName + ".db"
	instance.DatabaseFilePath = filepath.Join(instance.WorkDir, instance.DatabaseName)
	instance.DefaultStorageDir = filepath.Join(instance.WorkDir, "storage")
	instance.LogDir = filepath.Join(instance.WorkDir, "log")

	// if ip is empty, obtain lan ip to instead.
	if instance.ServerIP == "" || !CheckIP(instance.ServerIP) {
		instance.ServerIP = instance.GetLanIP()
	}
	instance.ServerAddress = instance.ServerIP + ":" + instance.Port
	instance.OrmClient = instance.initDatabase()
	instance.configuration = instance.queryConfiguration()
	instance.CustomStorageDir = instance.configuration.CustomStoragePath
}

//CheckIP 校验IP是否有效
func CheckIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

//Get 返回当前Context实例
func Get() *Environment {
	return instance
}

//GetConfiguration 获取当前设置项
func (environment *Environment) GetConfiguration() *ent.Configuration {

	if environment.configuration != nil {
		// 动态字段赋值
		environment.configuration.Version = environment.Version
		environment.configuration.AppPath = environment.WorkDir
		environment.configuration.DefStoragePath = environment.DefaultStorageDir
	}

	return environment.configuration
}

func (environment *Environment) SetConfiguration(configuration *ent.Configuration) {
	if configuration == nil {
		return
	}

	// 动态字段赋值
	environment.configuration.Version = environment.Version
	environment.configuration.AppPath = environment.WorkDir
	environment.configuration.DefStoragePath = environment.DefaultStorageDir

	environment.configuration = configuration
}

//GetStorageDir 获取当前仓储目录
func (environment *Environment) GetStorageDir() string {
	if len(environment.CustomStorageDir) == 0 {
		return environment.DefaultStorageDir
	}
	return environment.CustomStorageDir
}

func (environment *Environment) initDatabase() *ent.Client {
	// connect to database
	client, err := ent.Open("sqlite3", "file:"+environment.DatabaseFilePath+"?cache=shared&_fk=1")
	if err != nil {
		logrus.Println(err)
		panic("failed to connect database")
	}

	if environment.IsTestEnvironment() {
		logrus.Info("on environment,skip database migrate")
	} else {
		// fixme https://github.com/facebookincubator/ent/pull/221# sqlite3 panic bug
		// migrate database
		if err := client.Schema.Create(context.Background(), ); err != nil {
			logrus.Fatalf("failed creating schema resources: %v", err)
		}
	}

	return client
}

func (environment *Environment) queryConfiguration() *ent.Configuration {
	var configuration *ent.Configuration
	var err error
	//obtain first record
	configuration, err = environment.OrmClient.Configuration.Query().Only(context.Background())

	if err != nil {
		//create new record if there is no record exist
		configuration, err = environment.OrmClient.Configuration.
			Create().
			SetInitialized(false).
			SetDatabaseFilePath(environment.DatabaseFilePath).
			SetLogDirectoryPath(environment.LogDir).
			SetThemeStyle("light").
			SetThemeColor("#1890FF").
			SetNavMode("top").
			SetCreated(time.Time{}).
			SetUpdated(time.Time{}).
			Save(context.Background())

		if err != nil {
			logrus.Error(err)
		}
	}

	// 动态字段赋值
	configuration.Version = environment.Version
	configuration.AppPath = environment.WorkDir
	configuration.DefStoragePath = environment.DefaultStorageDir

	return configuration
}

//IsTestEnvironment 当前是否测试环境
func (environment *Environment) IsTestEnvironment() bool {
	for _, value := range os.Args {
		if strings.Contains(value, "-test.v") {
			return true
		}
	}
	return false
}

//getWorkDirectoryPath 获取工作目录
func (environment *Environment) getWorkDirectoryPath() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

//GetLanIP 返回本地ip
func (environment *Environment) GetLanIP() string {
	addresses, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logrus.Infof("print all ip address: %v\n\t", addresses)

	for _, address := range addresses {
		ipNet, ok := address.(*net.IPNet)

		if !ok || ipNet.IP.IsLoopback() || ipNet.IP.To4() == nil {
			continue
		}

		// 当前ip属于私有地址,直接返回
		if isIpBelongToPrivateIpNet(ipNet.IP) {
			return ipNet.IP.To4().String()
		}
	}

	return "127.0.0.1"
}

// 某个ip是否属于私有网段
func isIpBelongToPrivateIpNet(ip net.IP) bool {
	for _, ipNet := range getInternalIpNetArray() {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

// 返回私有网段切片
func getInternalIpNetArray() []*net.IPNet {
	var ipNetArrays []*net.IPNet

	for _, ip := range []string{"192.168.0.0/16", "172.16.0.0/12", "10.0.0.0/8"} {
		_, ipNet, _ := net.ParseCIDR(ip)
		ipNetArrays = append(ipNetArrays, ipNet)
	}

	return ipNetArrays
}
