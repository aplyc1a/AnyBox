package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aplyc1a/gair-framework/cache"
	"github.com/aplyc1a/gair-framework/config"
	"github.com/aplyc1a/gair-framework/ir"
	"github.com/aplyc1a/utils/fs"
	"github.com/aplyc1a/utils/log"
	"github.com/aplyc1a/utils/mid"
)

var (
	Import     string = config.CACHEDIR
	ImportOnly string = config.CACHEDIR
	Mid        string = ""
	innername  string = "jihe"
	LiveMode   bool   = false
)

var Cache bool = false
var ReCache bool = false
var DropCache bool = false
var CacheOnly bool = false
var DropCacheOnly bool = false
var Debug bool = false

func init() {

	// 新增应用程序特有的命令行参数
	flag.StringVar(&Import, "import", config.CACHEDIR, "")
	flag.StringVar(&Import, "i", config.CACHEDIR, "")
	flag.StringVar(&ImportOnly, "import-only", config.CACHEDIR, "")
	flag.StringVar(&ImportOnly, "importonly", config.CACHEDIR, "")
	flag.StringVar(&ImportOnly, "io", config.CACHEDIR, "")
	flag.StringVar(&Mid, "deviceid", "", "")
	flag.StringVar(&Mid, "mid", "", "")
	flag.StringVar(&Mid, "machineid", "", "")

	flag.BoolVar(&Cache, "cache", false, "")
	flag.BoolVar(&Cache, "makecache", false, "")
	flag.BoolVar(&Cache, "mc", false, "")

	flag.BoolVar(&ReCache, "recache", false, "")
	flag.BoolVar(&ReCache, "rc", false, "")
	flag.BoolVar(&ReCache, "rco", false, "")

	flag.BoolVar(&DropCache, "dropcache", false, "")
	flag.BoolVar(&DropCache, "dc", false, "")
	//flag.BoolVar(&LiveMode, "livemode", false, "")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", innername)
		// 创建一个切片存储参数的名称和描述
		commonOptions := []struct {
			name  string
			usage string
		}{
			{"--help/-h", "Show these help message"},
			{"--verbose/-v", "Enable verbose mode"},
			{"--lang", "Language for output (en or cn)"},
			{"--version/-V", "Show version"},
			{"--quiet/-q", "Enable quiet mode"},
			{"--logfile/-log", "Path to log file "},
		}
		fmt.Fprintf(flag.CommandLine.Output(), "===============\n")
		fmt.Println("Common parameters:")
		// 按自定义顺序输出参数的帮助信息
		for _, option := range commonOptions {
			fmt.Fprintf(flag.CommandLine.Output(), "  %s: %s\n", option.name, option.usage)
		}

		cltOptions := []struct {
			name  string
			usage string
		}{
			{"--import", "Import the collect-data.(also -i)"},
			{"--importonly", "Only import the collect-data.(also -io,--import-only)"},
			{"--mid", "Provide the machine id."},
			{"--cache", "Make the cache.(also -mc/-makecache)"},
			{"--recache", "Only rebuild the cache.(also -rc/-rco)"},
			{"--dropcache", "Remove the cache db.(also -dc)"},
			{"--cacheonly", "Only make the cache.(also -co)"},
			{"--dropcacheonly", "Only remove the cache db.(also -dco)"},
			//{"--livemode", "Set the work mode."},
		}
		fmt.Fprintf(flag.CommandLine.Output(), "---------------\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Specific parameters:\n")
		// 按自定义顺序输出参数的帮助信息
		for _, option := range cltOptions {
			fmt.Fprintf(flag.CommandLine.Output(), "  %s: %s\n", option.name, option.usage)
		}
	}

	// 解析命令行参数
	flag.Parse()

	if config.Cfg.Help {
		flag.Usage()
		os.Exit(1)
	}

	if config.Cfg.Version {
		fmt.Printf("%s\n", config.Version)
		os.Exit(1)
	}

	if config.Cfg.Debug {
		Debug = true
	}

	config.LoadTranslations()
	config.CheckCommonsRequiredParams()
	cwd, err := os.Getwd()
	if err != nil {
		log.Error(config.Translate("GetwdFailed", map[string]interface{}{"Err": err}))
		return
	}
	config.Cwd = cwd
	if Mid != "" {
		config.Mid = Mid
	}

}

func initEnv() {

	//设置Cwd目录，已提前设置好
	log.Debug("Set Cwd=%s\n", config.Cwd)
	//设置Mid，已提前设置好
	if config.Mid != "" {
		log.Debug("Set Mid=%s\n", config.Mid)
	}

	//设置CacheCwd目录
	err := os.Chdir(config.CACHEDIR)
	if err != nil {
		log.Critical("Failed to change working directory to %s:%v", config.CACHEDIR, err)
		os.Exit(1)
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Error(config.Translate("GetwdFailed", map[string]interface{}{"Err": err}))
		return
	}
	config.CacheCwd = cwd
	log.Debug("Set CacheCwd=%s\n", config.CacheCwd)

	//如果指定了导入，那就导入，如果导入的是单数据，Mid直接指向它
	if ImportOnly != "" {
		importFile(ImportOnly, config.CACHEDIR)
		os.Exit(1)
	}
	if Import != "" {
		importFile(Import, config.CACHEDIR)
	}

	MidList, err := mid.GetDeviceIDsFromDir(config.CacheCwd)
	if err != nil {
		log.Critical("Failed to fetch any data-set by Mid\n")
		os.Exit(1)
	}

	if len(MidList) == 0 {
		//未在目录下发现Mid
		// not mid dataset found
		log.Critical("No dataset found, where is it!\n")
		os.Exit(1)
	} else if len(MidList) == 1 {
		if Mid != "" {
			//目录下的Mid文件包含给定的Mid substr
			if strings.Contains(MidList[0], Mid) {
				Mid = MidList[0]
				log.Debug("Reset Mid=%s\n", Mid)
			} else {
				//目录下的Mid文件包含给定的Mid substr,提供了错误的mid
				// provided a wrong mid
				log.Critical("Wrong Mid provided! Please try again\n")
				os.Exit(1)
			}
		} else {
			//Mid未设置，直接设置为目录下的Mid
			Mid = MidList[0]
			log.Debug("Set Mid=%s\n", Mid)
		}
	} else {
		log.Warning("[!] %d datasets detected.\n", len(MidList))
		//当目录下存在多个Mid时，
		//如果给定的Mid不为空
		if Mid != "" {
			MidMatched := fs.FilterStringsContainingSubstring(MidList, Mid)
			//如果给定的Mid在目录下的Mid中搜不到
			if len(MidMatched) == 0 {
				// provided a wrong mid
				log.Critical("Wrong Mid provided! Please try again!\n")
				os.Exit(1)
			} else if len(MidMatched) > 1 {
				//如果给定的Mid在目录下的Mid中搜到了多个
				log.Critical("Matched %d times! please try again!\n", len(MidMatched))
				os.Exit(1)

			} else {
				//如果给定的Mid在目录下的Mid中搜到了1个
				if Mid == MidMatched[0] {
					Mid = MidMatched[0]
				} else {
					Mid = MidMatched[0]
					log.Debug("Reset Mid=%s\n", Mid)
				}
			}
		} else {
			// please provide mid
			log.Critical("Use `--mid` to choose one!\n")
			os.Exit(1)
		}
	}

	err = os.Chdir(Mid)
	if err != nil {
		log.Critical("Failed to change working directory to %s:%v", Mid, err)
		os.Exit(1)
	}
	config.TargetCwd, err = os.Getwd()
	if err != nil {
		log.Error(config.Translate("GetwdFailed", map[string]interface{}{"Err": err}))
		return
	}
	log.Debug("Set TargetCwd=%s\n", config.TargetCwd)
	return
}

func main() {
	config.InitLog()
	defer config.CloseLogFile()

	log.Debug(config.Translate("WorkingOn", map[string]interface{}{"Item": "initEnv"}))
	initEnv()

	configStr := Mid + ".db"
	log.Debug("configStr:%s\n", configStr)
	log.Blank()

	if Cache || CacheOnly {
		log.Debug(config.Translate("WorkingOn", map[string]interface{}{"Item": "Caching"}))

		cacheInstance := cache.NewInstance()
		cacheInstance.InitEnv(config.TargetCwd, Mid)
		cacheInstance.Start(configStr)
		if CacheOnly {
			os.Exit(1)
		}
		log.Blank()
	}

	if DropCache || DropCacheOnly {
		cacheInstance := cache.NewInstance()
		cacheInstance.InitEnv(config.TargetCwd, Mid)
		cacheInstance.RemoveDB()
		if DropCacheOnly {
			os.Exit(1)
		}
		log.Blank()
	}

	if ReCache {
		cacheInstance := cache.NewInstance()
		cacheInstance.InitEnv(config.TargetCwd, Mid)
		cacheInstance.RemoveDB()
		cacheInstance.Start(configStr)
		os.Exit(1)
		log.Blank()
	}

	//search

	irInstance := ir.NewInstance()
	irInstance.InitEnv(config.TargetCwd, Mid)
	irInstance.Start()
	//irInstance.ShowEvent()
	//irInstance.ExportReport()
	irInstance.Close()

	log.Info(config.Translate("ProgramFinished", nil))
}

// 文件导入函数
func importFile(path, cacheDir string) (string, error) {
	a, _ := filepath.Abs(path)
	b, _ := filepath.Abs(cacheDir)
	if a == b {
		return cacheDir, nil
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("路径不存在: %s", path)
	}

	if info.IsDir() {
		if mid.IsDeviceIDValid(filepath.Base(path)) {
			fmt.Printf("Copy :%s %s \n", path, cacheDir)
			err = fs.CopyDir(path, cacheDir)
			if err != nil {
				fmt.Printf("copy %s failed\n", path)
				return "", err
			}
			Mid = filepath.Base(path)
		} else {
			fList, err := fs.GetDirChildFiles(path)
			if err != nil {
				return "", err
			}
			for _, f := range fList {
				if b, err := fs.IsDirectory(fs.Join(path, f)); !b || err != nil {
					continue
				}
				if mid.IsDeviceIDValid(f) {
					fmt.Printf("Copy :%s %s \n", fs.Join(path, f), fs.Join(cacheDir, f))
					err = fs.CopyDir(fs.Join(path, f), fs.Join(cacheDir, f))
					if err != nil {
						fmt.Printf("copy %s failed\n", fs.Join(path, f))
					}
				}
			}
		}

		return cacheDir, nil
	}

	if fs.IsTarGzFile(path) {
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			err = os.Mkdir(cacheDir, os.ModePerm)
			if err != nil {
				return "", err
			}
		}

		err := fs.UncompressTarGz(path, cacheDir)
		if err != nil {
			return "", err
		}
		fmt.Printf("Copy and uncompressed successfully! \n")
		return cacheDir, nil
	}

	return "", fmt.Errorf("不支持的文件类型")
}
