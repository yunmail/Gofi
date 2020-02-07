package controllers

import (
	"context"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"gofi/env"
	"gofi/environment"
	"gofi/i18n"
	"gofi/util"
	"path/filepath"
)

//UpdateSetting 更新设置
func UpdateSetting(ctx iris.Context) {
	// 初始化完成且处于Preview环境,不允许更改设置项
	if env.IsPreview() && environment.Get().GetConfiguration().Initialized {
		_, _ = ctx.JSON(NewResource().Fail().Message(i18n.Translate(i18n.OperationNotAllowedInPreviewMode)).Build())
		return
	}

	configuration := environment.Get().GetConfiguration()

	// 用客户端给定的Configuration覆盖数据库持久化的Configuration
	// 避免Body为空的时候ReadJson报错,导致后续不能默认初始化，这里用ContentLength做下判断
	if err := ctx.ReadJSON(configuration); ctx.GetContentLength() != 0 && err != nil {
		logrus.Error(err)
		_, _ = ctx.JSON(NewResource().Fail().Build())
	}
	updateBuilder := configuration.Update().SetInitialized(true)

	path := filepath.Clean(configuration.CustomStoragePath)
	workDir := environment.Get().WorkDir
	defaultStorageDir := environment.Get().DefaultStorageDir

	// 是否使用默认地址
	useDefaultDir := path == "" || path == defaultStorageDir

	// 如果使用默认仓库路径
	if useDefaultDir {
		path = defaultStorageDir
	}

	logrus.Printf("工作目录是%v \n", workDir)
	logrus.Printf("dir目录是%v \n", path)

	// 判断给定的目录是否存在
	if !util.FileExist(path) {
		_, _ = ctx.JSON(NewResource().Fail().Message(i18n.Translate(i18n.DirIsNotExist, path)))
		return
	}

	// 判断给定的路径是否是目录
	if !util.IsDirectory(path) {
		_, _ = ctx.JSON(NewResource().Fail().Message(i18n.Translate(i18n.IsNotDir, path)))
		return
	}

	// 如果文件夹不存在，创建文件夹
	util.MkdirIfNotExist(defaultStorageDir)

	ormClient := environment.Get().OrmClient

	// 使用事务
	tx, err := ormClient.Tx(context.Background())

	if err != nil {
		logrus.Errorf("starting a transaction: %v", err)
		_, _ = ctx.JSON(NewResource().Fail().Build())
	}

	// 持久化到数据库
	configuration, _ = updateBuilder.
		SetCustomStoragePath(path).
		Save(context.Background())

	// 更新单例属性
	environment.Get().CustomStorageDir = path

	// 路径合法，初始化成功，持久化该路径。
	logrus.Infof("use default path %s, setup success", path)

	err = tx.Commit()

	if err != nil {
		logrus.Errorf("commit a transaction: %v", err)
		_, _ = ctx.JSON(NewResource().Fail().Build())
	}

	environment.Get().SetConfiguration(configuration)

	GetConfiguration(ctx)
}

//Setup 初始化
func Setup(ctx iris.Context) {
	// 已经初始化过
	if environment.Get().GetConfiguration().Initialized {
		_, _ = ctx.JSON(NewResource().Fail().Message(i18n.Translate(i18n.GofiIsAlreadyInitialized)).Build())
		return
	}

	UpdateSetting(ctx)
}

//GetConfiguration 获取设置项
func GetConfiguration(ctx iris.Context) {
	configuration := environment.Get().GetConfiguration()
	_, _ = ctx.JSON(NewResource().Payload(configuration).Build())
}
