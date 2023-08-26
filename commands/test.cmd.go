package commands

import (
	"context"
	"fmt"
	"time"

	"rnv-mmq/database"
	"rnv-mmq/tools"
	"rnv-mmq/wrongs"

	"github.com/gosuri/uiprogress"
	uuid "github.com/satori/go.uuid"
	"github.com/schollz/progressbar/v3"
)

// TestCommand 测试用
type TestCommand struct{}

// NewTestCommand 构造函数
func NewTestCommand() *TestCommand {
	return &TestCommand{}
}

func (receiver TestCommand) uuid() []string {
	c := make(chan string)
	go func(c chan string) {
		uuidStr := uuid.NewV4().String()
		c <- uuidStr
	}(c)
	go tools.NewTimer(5).Ticker()
	return []string{<-c}
}

func (receiver TestCommand) ls() []string {
	_, res := (&tools.Cmd{}).Process("ls", "-la")
	return []string{res}
}

func (receiver TestCommand) redis() []string {
	if _, err := database.NewRedis(0).SetValue("test", "AAA", 15*time.Minute); err != nil {
		wrongs.ThrowForbidden(err.Error())
	}

	for i := 0; i < 100000; i++ {
		if val, err := database.NewRedis(0).GetValue("test"); err != nil {
			wrongs.ThrowForbidden(err.Error())
		} else {
			fmt.Println(i, val)
		}
	}

	return []string{""}
}

func (receiver TestCommand) uiProcesses() []string {
	bar := progressbar.Default(100)
	for i := 0; i < 100; i++ {
		err := bar.Add(1)
		if err != nil {
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}

	return []string{}
}

func (receiver TestCommand) uiProcesses2() []string {
	uiprogress.Start()
	bar := uiprogress.AddBar(100).AppendCompleted().PrependElapsed()
	for i := 0; i < 100; i++ {
		bar.Incr()
		time.Sleep(50 * time.Millisecond)
	}
	uiprogress.Stop()

	return []string{}
}

func worker(ctx context.Context, done chan bool) {
	// 模拟一个需要执行5秒的任务
	time.Sleep(5 * time.Second)

	select {
	case <-ctx.Done():
		fmt.Println("Worker stopped due to timeout")
	default:
		fmt.Println("Worker completed successfully")
	}

	done <- true
}

func (receiver TestCommand) timeout() []string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	done := make(chan bool, 1)
	go worker(ctx, done)

	select {
	case <-done:
		fmt.Println("All workers completed")
	case <-ctx.Done():
		fmt.Println("Program stopped due to timeout")
	}

	return []string{}
}

type (
	SyncTasks struct {
		Name                string  `gorm:"type:varchar(50);not null;default:'';comment:名称;"`
		Status              uint8   `gorm:"type:tinyint unsigned;not null;default:2;comment:状态;"`
		ParagraphUniqueCode string  `gorm:"type:char(4);not null;default:'';comment:段编码;"`
		Remark              string  `gorm:"type:varchar(255);not null;default:'';comment:备注;"`
		Project             string  `gorm:"type:varchar(50);default '';not null;comment:项目;"`
		RequestUrl          string  `gorm:"type:varchar(100);default '';not null;comment:请求url;"`
		RequestMethod       string  `gorm:"type:varchar(20);default '';not null;comment:请求method;"`
		RequestContent      *string `gorm:"type:longtext;comment:请求内容;"`
		ResponseContent     *string `gorm:"type:longtext;comment:响应内容;"`
		BatchCode           string  `gorm:"type:char(36);default:'';not null;comment:批次号;"`
	}
	RequestContent struct {
		UpdateEquipments []UpdateEquipment `json:"update_equipments"`
	}

	UpdateEquipment struct {
		Status                       string  `json:"status"`
		SerialNumber                 string  `json:"serial_number"`
		InstalledAt                  string  `json:"installed_at"`
		UpdatedAt                    string  `json:"updated_at"`
		UniqueCode                   string  `json:"unique_code"`
		WorkshopUniqueCode           string  `json:"workshop_unique_code"`
		StationUniqueCode            string  `json:"station_unique_code"`
		CenterUniqueCode             string  `json:"center_unique_code"`
		CrossingUniqueCode           string  `json:"crossing_unique_code"`
		LineUniqueCode               string  `json:"line_unique_code"`
		InstallLocationUniqueCode    string  `json:"install_location_unique_code"`
		OutdoorInstallLocationUuid   string  `json:"outdoor_install_location_uuid"`
		OutdoorInstallLocationExtend *string `json:"outdoor_install_location_extend"`
		SceneWorkAreaUniqueCode      string  `json:"scene_work_area_unique_code"`
	}

	Person struct {
		Name string
		Age  uint64
		Rank uint64
	}
)

// Handle 执行命令
func (receiver TestCommand) Handle(params []string) []string {
	switch params[0] {
	case "uuid":
		return receiver.uuid()
	case "ls":
		return receiver.ls()
	case "redis":
		return receiver.redis()
	case "uiProcesses":
		return receiver.uiProcesses()
	case "uiProcesses2":
		return receiver.uiProcesses2()
	case "timeout":
		return receiver.timeout()
	default:
		return []string{"没有找到命令"}
	}
}
