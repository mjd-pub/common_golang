package logs

import (
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Loggers struct {
	LogType string
}

func Log(msg string, data interface{}) {
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}

	logger.Info(msg, zap.Any("log", data))
}

func InnerLog(msg string, data map[string]interface{}) {
	logger, err := NewInnerLogger()
	var zapData []zap.Field
	if err != nil {
		panic(err)
	}

	if data != nil {
		for k, v := range data {
			zapData = append(zapData, zap.Any(k, v))
		}

		logger.Info(msg, zapData...)
	} else {
		logger.Info(msg, zap.String("log", "日志为空"))
	}

}

//自定义logger
//官方log包与zap对接
//loggers结构体内，根据logType设置模板
func (l Loggers) Print(v ...interface{}) {
	var cfg zap.Config
	cfg = zap.NewProductionConfig()
	errFile, _ := filepath.Abs("./log/error.log")
	cfg.OutputPaths = []string{
		errFile,
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := cfg.Build()
	res := l.setTemplate(v...)
	logger.Error("[Error]", zap.Any(l.LogType, res))
}

/*配置模板*/
func (l Loggers) setTemplate(v ...interface{}) map[string]interface{} {
	switch l.LogType {
	case "mysqlError":
		return map[string]interface{}{
			"file":    v[1],
			"errCode": v[2].(*mysql.MySQLError).Number,
			"errInfo": v[2].(*mysql.MySQLError).Message,
			"debug":   v[2].(*mysql.MySQLError).Error(),
		}
	default:
		return map[string]interface{}{}
	}
}

func Error(data interface{}) {
	logger, _ := NewLogger()
	logger.Error("[Error]", zap.Any("request-error", data))
}

func NewLogger() (*zap.Logger, error) {
	var cfg zap.Config
	errFile, _ := filepath.Abs("./log/error.log")
	env := viper.GetString("params.env")
	//if gin.Mode() == "debug" || gin.Mode() == "test" {
	if env == "dev" || env == "test" {
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.Encoding = "json"
		cfg.EncoderConfig = zap.NewProductionEncoderConfig()
		errFile, _ := filepath.Abs("./log/error.log")
		cfg.ErrorOutputPaths = []string{
			errFile,
		}
		cfg.InitialFields = map[string]interface{}{
			"ginMode": env,
		}
	} else {
		cfg = zap.NewProductionConfig()
		errFile, _ := filepath.Abs("./log/error.log")
		cfg.ErrorOutputPaths = []string{
			errFile,
		}
		cfg.InitialFields = map[string]interface{}{
			"ginMode": env,
		}
	}

	cfg.OutputPaths = []string{
		errFile,
	}

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfg.Build()
}

func NewInnerLogger() (*zap.Logger, error) {
	var cfg zap.Config

	//if gin.Mode() == "debug" || gin.Mode() == "test" {
	env := viper.GetString("params.env")
	if env == "dev" || env == "test" {
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.Encoding = "json"
		cfg.EncoderConfig = zap.NewProductionEncoderConfig()
		errFile, _ := filepath.Abs("./log/error.log")
		cfg.ErrorOutputPaths = []string{
			errFile,
		}
		cfg.InitialFields = map[string]interface{}{
			"ginMode": env,
		}
	} else {
		cfg = zap.NewProductionConfig()
		errFile, _ := filepath.Abs("./log/error.log")
		cfg.ErrorOutputPaths = []string{
			errFile,
		}
		cfg.InitialFields = map[string]interface{}{
			"ginMode": env,
		}
	}

	timeDec := time.Now()
	logName := timeDec.Format("20060102_15") + ".log"
	fileDay := timeDec.Format("20060102") + "/"
	file, _ := filepath.Abs("./log/" + fileDay)

	_, err := os.Stat(file)

	logPath := file + "/" + logName
	if err != nil && os.IsNotExist(err) {
		//创建天文件夹
		err := os.MkdirAll(file, 0755)

		if err != nil {
			panic("Create Log Dir Fail:" + err.Error())
		}
		_, err = os.Stat(logPath)

		if err != nil && os.IsNotExist(err) {
			_, err = os.Create(logPath)
			if err != nil {
				panic("Create Log Dir Fail:" + err.Error())
			}
		}
	}
	cfg.OutputPaths = []string{
		logPath,
	}

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfg.Build()
}

//自动生成日志
func generateFileLogForTime(rootFilePath string) string {
	timeDec := time.Now()
	fileDay := timeDec.Format("20060102")
	file, _ := filepath.Abs(rootFilePath + fileDay)
	_, err := os.Stat(file)

	var logFile string

	if err != nil && os.IsNotExist(err) {
		//创建天文件夹
		err := os.MkdirAll(file, 0755)

		if err != nil {
			panic("Create Log Dir Fail:" + err.Error())
		}
	}

	logFile = file + "/" + timeDec.Format("2006-01-02-15") + ".log"
	_, err = os.Stat(logFile)

	if err != nil && os.IsNotExist(err) {
		_, err = os.Create(logFile)
		if err != nil {
			panic("Create Log Dir Fail:" + err.Error())
		}
	}

	return logFile

}

//自动生成inner日志

//删除日志
//@params days int 几天前
func DelRemove(days int) {
	//计算前一天日志
	d := time.Now()
	beforeDay := d.AddDate(0, 0, -days)
	beforeLog := beforeDay.Format("20060102")
	log.Printf("需要删除%s当天及之前的日志", beforeLog)

	//获取全部文件夹
	dirpath, err := filepath.Abs("./log")

	if err != nil {
		panic(err)
	}
	files, err := ioutil.ReadDir(dirpath)

	if err != nil {
		panic(err)
	}

	var delLogDir []string
	for _, file := range files {
		if file.IsDir() && beforeLog >= file.Name() {
			delLogDir = append(delLogDir, dirpath+"/"+file.Name())
			err := os.RemoveAll(dirpath + "/" + file.Name())
			/*if os.IsNotExist(err) {
				log.Fatal("NotExist")
			}*/
			if err != nil {
				panic(err)
			}
		}
	}
	if len(delLogDir) > 0 {
		log.Println("已删除日志文件夹:", delLogDir)
	}
}
