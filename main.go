package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"rnv-mmq/commands"
	"rnv-mmq/database"
	"rnv-mmq/middlewares"
	"rnv-mmq/models"
	"rnv-mmq/providers"
	"rnv-mmq/routes/apiRoute"
	"rnv-mmq/routes/webRoute"
	"rnv-mmq/settings"
	"rnv-mmq/tools"
	"rnv-mmq/types"
	"rnv-mmq/wrongs"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

// 启动守护进程
func launchDaemon(title string) {
	if syscall.Getppid() == 1 {
		if err := os.Chdir("./"); err != nil {
			panic(err)
		}
		// syscall.Umask(0) // TODO TEST
		return
	}

	currentDir := os.Getenv("PWD")
	logDir := fmt.Sprintf("%s/logs", currentDir)
	if !tools.NewFileSystem(logDir).IsExist() {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			fmt.Println("创建日志目录错误：" + err.Error())
			return
		}
	}
	filenameOnToday := fmt.Sprintf("%s/logs/%s.log", currentDir, time.Now().Format(types.FormatDate))
	fp, err := os.OpenFile(filenameOnToday, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = fp.Close()
	}()
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Stdout = fp
	cmd.Stderr = fp
	cmd.Stdin = nil
	// cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // TODO TEST
	if err = cmd.Start(); err != nil {
		panic(err)
	}

	startLogContent := fmt.Sprintf("%s 程序启动成功 [进程号%d] 启动于：%s\r\n", title, cmd.Process.Pid, time.Now().Format(types.FormatDatetime))
	tools.StdoutSuccess("程序启动成功", "").EchoLine(fmt.Sprintf("[进程号%d] 启动于：%s", cmd.Process.Pid, time.Now().Format(types.FormatDatetime)))

	_, _ = fp.WriteString("\r\n\r\n\r\n>>>>>>>>>>>>>>>>>>>>\r\n" + startLogContent)
	os.Exit(0)
}

type (
	// WebServiceCommand web-service 服务
	WebServiceCommand struct{}

	// WebServiceCommandSetting web-service 服务配置
	WebServiceCommandSetting struct {
		WebAddr           string
		TcpServerEnable   bool
		TcpServerAddr     string
		TcpClientEnable   bool
		TcpClientAddr     string
		RabbitMqEnable    bool
		RabbitMqAddr      string
		RabbitMqUsername  string
		RabbitMqPassword  string
		RabbitMqVhost     string
		RabbitMqQueueName string
	}
)

var autoMigrateDst = []any{
	models.UserModel{},    // 用户
	models.TaskModel{},    // 任务
	models.TaskLogModel{}, // 任务日志
	models.FileModel{},    // 文件
}

// bootAutoMigrate 初始化数据库迁移
func bootAutoMigrate(dst ...interface{}) {
	if errAutoMigrate := database.
		NewGormLauncher().
		SetDbDriver("").
		GetConn("").
		AutoMigrate(dst...); errAutoMigrate != nil {
		fmt.Println("数据库迁移错误:", errAutoMigrate.Error())
		os.Exit(1)
	}
}

// 执行command命令
func launchCommand(commandName string, commandParams, tmp []string, daemon bool) {
	var (
		commandResults []string
		commandSetting WebServiceCommandSetting
		appSetting     *ini.File = settings.NewSetting().App
		appName        string    = appSetting.Section("app").Key("name").String()
	)

	commandSetting = WebServiceCommandSetting{
		WebAddr:           appSetting.Section("web-service").Key("addr").MustString(":8091"),
		TcpServerEnable:   appSetting.Section("tcp-server-service").Key("enable").MustBool(false),
		TcpServerAddr:     appSetting.Section("tcp-server-service").Key("addr").MustString("0.0.0.0:8092"),
		TcpClientEnable:   appSetting.Section("tcp-client-service").Key("enable").MustBool(false),
		TcpClientAddr:     appSetting.Section("tcp-client-service").Key("addr").MustString("127.0.0.1:8080"),
		RabbitMqEnable:    appSetting.Section("rabbit-mq-service").Key("enable").MustBool(false),
		RabbitMqAddr:      appSetting.Section("rabbit-mq-service").Key("addr").MustString("127.0.0.1:5672"),
		RabbitMqUsername:  appSetting.Section("rabbit-mq-service").Key("username").MustString(""),
		RabbitMqPassword:  appSetting.Section("rabbit-mq-service").Key("password").MustString(""),
		RabbitMqVhost:     appSetting.Section("rabbit-mq-service").Key("vhost").MustString(""),
		RabbitMqQueueName: appSetting.Section("rabbit-mq-service").Key("queue-name").MustString(""),
	}

	if daemon {
		launchDaemon(fmt.Sprintf("%s %s", appName, commandName))
	}

	switch commandName {
	default:
		// println(tools.StdoutWrong(fmt.Sprintf("【执行失败】%s", strings.Join(tmp, " ")), "").GetContentAndNext("没有找到命令"))
		// os.Exit(-1)
		// 启动web-service服务

		// 模型
		bootAutoMigrate(autoMigrateDst...)

		// 创建TCP服务端
		if commandSetting.TcpServerEnable {
			go providers.TcpServerHandler(commandSetting.TcpServerAddr)
		}

		// 创建TCP客户端
		if commandSetting.TcpClientEnable {
			go providers.TcpClientHandler(commandSetting.TcpClientAddr)
		}

		// 创建RabbitMQ监听
		if commandSetting.RabbitMqEnable {
			go providers.RabbitMqHandler(
				commandSetting.RabbitMqUsername,
				commandSetting.RabbitMqPassword,
				commandSetting.RabbitMqAddr,
				commandSetting.RabbitMqVhost,
				commandSetting.RabbitMqQueueName,
			)
		}

		engine := gin.Default()                                     // 启动服务引擎
		engine.Use(middlewares.Cors())                              // 允许跨域
		engine.Use(middlewares.TimeoutMiddleware(time.Second * 60)) // 超时处理
		engine.Use(wrongs.RecoverHandler)                           // 异常处理
		webRoute.Router{}.Register(engine)                          // 加载web路由
		apiRoute.RouterHandle{}.Register(engine)                    // 加载api路由

		// 绑定web-socket路由
		engine.GET("/ws", func(ctx *gin.Context) {
			providers.WebsocketHandler(ctx)
		})

		// 启动web-service服务
		if err := engine.Run(commandSetting.WebAddr); err != nil {
			log.Fatalf("[web-service-error] [启动web服务错误] %v", err)
		}

		log.Printf("[app] [程序停止] [%s]", time.Now().Format(types.FormatDatetime))
	case "test":
		commandResults = commands.NewTestCommand().Handle(commandParams)
	}

	fmt.Println(tools.StdoutDebug(fmt.Sprintf("【执行完成 %s】 「%s」 ", tools.NewTime().SetTimeNowAdd8Hour().ToDateTimeString(), strings.Join(tmp, " ")), "").GetContentAndNext(strings.Join(commandResults, " ")))
	os.Exit(0)
}

// main 程序入口
func main() {

	// params explain:
	// t=command 执行命令行
	// t='web-service' 启动web服务 可选项-port=8080 端口号、-daemon=false 守护进程
	var (
		title, commandName string
		commandParams, tmp []string
		daemon             bool
	)
	flag.StringVar(&title, "t", "", "命令终端参数")
	flag.BoolVar(&daemon, "daemon", false, "是否启动守护进程")
	flag.Parse()

	commandName = ""
	commandParams = make([]string, 0)

	if title != "" {
		tmp = strings.Split(title, " ")
		commandName = tmp[0]
		commandParams = tmp[1:]
	}

	launchCommand(commandName, commandParams, tmp, daemon)
}
