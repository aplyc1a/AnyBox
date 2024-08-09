package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aplyc1a/gair-framework/cache"
	"github.com/aplyc1a/gair-framework/capability"
	"github.com/aplyc1a/gair-framework/clt"
	"github.com/aplyc1a/gair-framework/config"
	"github.com/aplyc1a/gair-framework/yara"
	"github.com/aplyc1a/utils/crypto"
	"github.com/aplyc1a/utils/fs"
	"github.com/aplyc1a/utils/log"
)

var (
	Fast      bool
	All       bool
	NoStatMap bool
	Pack      bool
	PackOnly  bool

	IsCustom             bool
	Cachedb              bool
	CustomTarget         string
	CustomIgnoredTargets string

	Filemap   bool
	Md5Sum    bool
	Sha256Sum bool
	Sha1Sum   bool
	Hashes    bool

	AddCustomCopy bool

	YaraScan string

	GrepText   string
	GrepRLText string
	GrepRNText string

	MonitoringNetSession bool
	MonitoringFdSession  bool

	Help bool
)

var innername string = "pudu"

func init() {
	// Set Mode
	flag.BoolVar(&Fast, "fast", false, "")
	flag.BoolVar(&Fast, "Fast", false, "")
	flag.BoolVar(&All, "all", false, "")
	flag.BoolVar(&All, "All", false, "")

	// Packing
	flag.BoolVar(&Pack, "packaging", false, "")
	flag.BoolVar(&Pack, "packing", false, "")
	flag.BoolVar(&Pack, "p", false, "")
	flag.BoolVar(&Pack, "pack", false, "")
	// PackOnly
	flag.BoolVar(&PackOnly, "packonly", false, "")
	flag.BoolVar(&PackOnly, "po", false, "")
	flag.BoolVar(&PackOnly, "packingonly", false, "")
	flag.BoolVar(&PackOnly, "packagingonly", false, "")

	flag.BoolVar(&Cachedb, "cache", false, "")
	flag.BoolVar(&Cachedb, "cachedb", false, "")
	flag.BoolVar(&Cachedb, "caching", false, "")
	flag.BoolVar(&Cachedb, "c", false, "")

	// --targets：设置目标对象。直接指定或从文件读取
	flag.StringVar(&CustomTarget, "target", "", "")
	flag.StringVar(&CustomTarget, "targets", "", "")
	flag.StringVar(&CustomTarget, "ts", "", "")
	flag.StringVar(&CustomIgnoredTargets, "ignore", "", "")
	flag.StringVar(&CustomIgnoredTargets, "ignore-target", "", "")
	flag.StringVar(&CustomIgnoredTargets, "ignore-targets", "", "")
	flag.StringVar(&CustomIgnoredTargets, "its", "", "")
	flag.StringVar(&CustomIgnoredTargets, "ign", "", "")

	// --statmap
	flag.BoolVar(&NoStatMap, "NoStatmap", false, "")
	flag.BoolVar(&NoStatMap, "nostatmap", false, "")
	flag.BoolVar(&NoStatMap, "nostat", false, "")
	flag.BoolVar(&Filemap, "statmap", false, "")

	// hash
	flag.BoolVar(&Md5Sum, "md5sum", false, "File to compute md5 hash for.")
	flag.BoolVar(&Sha256Sum, "sha256sum", false, "File to compute sha256 hash for.")
	flag.BoolVar(&Sha1Sum, "sha1sum", false, "File to compute sha1 hash for.")
	flag.BoolVar(&Hashes, "hashes", false, "File to compute hash(md5/sha1/sha256) for.")

	// --copy filename
	flag.BoolVar(&AddCustomCopy, "copy", false, "")

	// yara
	flag.StringVar(&YaraScan, "yarascan", "", "Runs with the yara!")

	// search
	flag.StringVar(&GrepText, "g", "", "String to search for in file.")
	flag.StringVar(&GrepText, "grep", "", "String to search for in file.")
	flag.StringVar(&GrepText, "grep-r", "", "String to search for in file.")
	flag.StringVar(&GrepText, "search", "", "String to search for in file.")
	flag.StringVar(&GrepRLText, "grep-rl", "", "String to search for in file.")
	flag.StringVar(&GrepRNText, "grep-rn", "", "String to search for in file.")

	// monitoring
	flag.BoolVar(&MonitoringNetSession, "net-monitoring", false, "")
	flag.BoolVar(&MonitoringNetSession, "net-mon", false, "")
	flag.BoolVar(&MonitoringFdSession, "fd-monitoring", false, "")
	flag.BoolVar(&MonitoringFdSession, "fd-mon", false, "")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", innername)
		// 创建一个切片存储参数的名称和描述
		commonOptions := []struct {
			name  string
			usage string
		}{
			{"--help/-h", "Show these help message."},
			{"--verb/-v", "Enable verbose mode."},
			{"--lang", "Language for output (en or cn)."},
			{"--version/-V", "Show version."},
			{"--quiet/-q", "Enable quiet mode."},
			{"--logfile/-log", "Path to log file."},
		}
		fmt.Fprintf(flag.CommandLine.Output(), "===============\n")
		fmt.Println("Common parameters:")
		// 按自定义顺序输出参数的帮助信息
		for _, option := range commonOptions {
			fmt.Fprintf(flag.CommandLine.Output(), "  %-26s: %s\n", option.name, option.usage)
		}

		cltOptions := []struct {
			name  string
			usage string
		}{
			{"--fast", "Do Fast collecting.(default:normal mode)."},
			{"--all", "Do the complete collecting.(default:normal mode)."},
			{"--all", "Do the complete collecting.(default:normal mode)."},
			{"--pack", "Pack the data after collecting.(Also:-p/--packing/--packaging)."},
			{"--packonly", "Just pack the collecting data.(Also:-po/--packingonly/--packagingonly)."},
			{"--net-mon", "Record the system network session per 1s."},
			{"--fd-mon", "Record the system process-fd per 1s."},
			{"--cache", "Caching evtx/statmap/syslog/browserhistory into sqlitedb.(Also -c/--cachedb/--caching)."},
		}

		fmt.Fprintf(flag.CommandLine.Output(), "---------------\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Classic parameters:\n")
		// 按自定义顺序输出参数的帮助信息
		for _, option := range cltOptions {
			fmt.Fprintf(flag.CommandLine.Output(), "  %-26s: %s\n", option.name, option.usage)
		}

		cltOptions = []struct {
			name  string
			usage string
		}{
			{"--target <filename,...>", "Set the file targets.(Also:-ts/--targets)."},
			{"--ignore <filename,...>", "Set the file/dirs to be ignored.(Also:-ign/-its/--ignore-target/--ignore-targets)"},

			{"--hashmap", "Do the `target` hashmapping!"},
			{"--md5sum", "Do the `target` md5sum-hashmapping!"},
			{"--sha256sum", "Do the `target` sha256sum-hashmapping!"},
			{"--sha1sum", "Do the `target` sha1sum-hashmapping!"},

			{"--copy", "Copy the `target` files based on your provided!"},

			{"--yara <rulesfile>", "Scans the `target` used yarascan!"},

			{"--grep <string>", "Search for PATTERNS in each `target`."},
			{"--grep-rl <string>", "Search for PATTERNS in each `target`."},
			{"--grep-rn <string>", "Search for PATTERNS in each `target`."},
		}
		fmt.Fprintf(flag.CommandLine.Output(), "---------------\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Advanced parameters:\n")

		// 按自定义顺序输出参数的帮助信息
		for _, option := range cltOptions {
			fmt.Fprintf(flag.CommandLine.Output(), "  %-26s: %s\n", option.name, option.usage)
		}
	}

	// 解析命令行参数
	flag.Parse()

	if config.Cfg.Help {
		flag.Usage()
		os.Exit(1)
	}

	if config.Cfg.Version {
		//config.ShowLogo()
		fmt.Printf("%s\n", config.Version)
		os.Exit(1)
	}
	checkParams()
	config.LoadTranslations()
	config.CheckCommonsRequiredParams()
}

func checkParams() {
	IsCustom = false
	if Fast && All {
		flag.Usage()
		os.Exit(1)
	}

	if MonitoringNetSession && MonitoringFdSession {
		flag.Usage()
		os.Exit(1)
	}

	if GrepRLText != "" || GrepRNText != "" || GrepText != "" {
		if CustomTarget == "" {
			flag.Usage()
			os.Exit(1)
		}
	}

}

func run() {
	isAdmin := fs.IsPriviledge()
	clt.CheckSpace()
	clt.InitEnv()

	log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "BaicInfo"}))
	clt.CltBasic()

	if PackOnly {
		log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "Packing"}))
		clt.Packing()
		os.Exit(1)
	}

	if Fast && isAdmin {

		clt.QuickStart()

		if Pack {
			log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "Packing"}))
			clt.Packing()
		}
	} else if All && isAdmin {
		log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "statmap"}))
		clt.FullStart()
		if Cachedb {
			configStr := config.Mid + ".db"
			cacheInstance := cache.NewInstance()
			cacheInstance.InitEnv(config.TargetCwd, config.Mid)
			cacheInstance.Start(configStr)
		}

		if Pack {
			log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "Packing"}))
			clt.Packing()
		}
	} else {
		target := []string{}
		ignoredTarget := []string{}

		if CustomIgnoredTargets != "" {
			IsCustom = true
			if strings.Contains(CustomIgnoredTargets, ",") {
				ignoredTarget = append(ignoredTarget, strings.Split(CustomIgnoredTargets, ",")...)
			} else {
				ignoredTarget = append(ignoredTarget, CustomIgnoredTargets)
			}
		}
		if CustomTarget != "" {
			IsCustom = true
			if strings.Contains(CustomTarget, ",") {
				target = append(target, strings.Split(CustomTarget, ",")...)
			} else {
				target = append(target, CustomTarget)
			}

			if Filemap && !NoStatMap {
				IsCustom = true
				HASHES := crypto.SHA256_BITS
				log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "DrawStatMapX"}))
				clt.DrawStatMapC(clt.Custom, target, ignoredTarget, "statmap-custom.txt", HASHES)
				os.Exit(1)
			}

			if Hashes {
				IsCustom = true
				HASHES := crypto.MD5_BITS & crypto.SHA1_BITS & crypto.SHA256_BITS
				log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "DrawHashMap"}))
				clt.DrawHashMapC(clt.Custom, target, ignoredTarget, HASHES)
			}

			if Md5Sum {
				IsCustom = true
				HASHES := crypto.MD5_BITS
				log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "DrawHashMap"}))
				clt.DrawHashMapC(clt.Custom, target, ignoredTarget, HASHES)
			}

			if Sha256Sum {
				IsCustom = true
				HASHES := crypto.SHA256_BITS
				log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "DrawHashMap"}))
				clt.DrawHashMapC(clt.Custom, target, ignoredTarget, HASHES)
			}

			if Sha1Sum {
				IsCustom = true
				HASHES := crypto.SHA1_BITS
				log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "DrawHashMap"}))
				clt.DrawHashMapC(clt.Custom, target, ignoredTarget, HASHES)
			}

			if AddCustomCopy {
				log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "CustomCopy"}))
				clt.CopyCustomFiles(target)
			}

			if YaraScan != "" {
				IsCustom = true
				if fs.IsFileExist(YaraScan) {
					log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "YaraScan"}))
					yara.RunAll(YaraScan, target, ignoredTarget)
				}
			}

			if GrepText != "" {
				IsCustom = true
				if !fs.IsFileExist(CustomTarget) {
					log.Critical("%s not exist!\n", CustomTarget)
				}
				clt.SearchFileLikeGrepR(GrepText, target, ignoredTarget)
			}

			if GrepRLText != "" {
				IsCustom = true
				if !fs.IsFileExist(CustomTarget) {
					log.Critical("%s not exist!\n", CustomTarget)
				}
				clt.SearchFileLikeGrepRL(GrepRLText, target, ignoredTarget)
			}

			if GrepRNText != "" {
				IsCustom = true
				if !fs.IsFileExist(CustomTarget) {
					log.Critical("%s not exist!\n", CustomTarget)
				}
				clt.SearchFileLikeGrepRN(GrepRNText, target, ignoredTarget)
			}

		}

		if MonitoringNetSession {
			IsCustom = true
			log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "NetworkHunting"}))
			stopChan := make(chan bool)
			go capability.HuntNetProc(stopChan)
			go func() {
				time.Sleep(1 * time.Hour)
				close(stopChan)
			}()

			go func() {
				log.Info(config.Translate("QuitHunting", nil))
			}()
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				if scanner.Text() == "q" {
					close(stopChan)
					break
				}
			}
			msg := ""
			if config.OS == "windows" {
				msg = fmt.Sprintf("%s\\", config.TargetCwd)
			} else {
				msg = fmt.Sprintf("`%s/", config.TargetCwd)
			}

			log.Info(config.Translate("HuntingResultMsg", map[string]interface{}{"Item": msg}))
			os.Exit(1)
		}

		if MonitoringFdSession {
			IsCustom = true
			log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "FileActivityHunting"}))
			stopChan := make(chan bool)
			capability.HuntFdProc(stopChan)
			go func() {
				time.Sleep(1 * time.Hour)
				close(stopChan)
			}()

			go func() {
				log.Info(config.Translate("QuitHunting", nil))
			}()
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				if scanner.Text() == "q" {
					close(stopChan)
					break
				}
			}
			msg := ""
			if config.OS == "windows" {
				msg = fmt.Sprintf("%s\\", config.TargetCwd)
			} else {
				msg = fmt.Sprintf("`%s/", config.TargetCwd)
			}

			log.Info(config.Translate("HuntingResultMsg", map[string]interface{}{"Item": msg}))
			os.Exit(1)
		}

		if !IsCustom {
			log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "NormalStart"}))
			clt.NormalStart()
			if Cachedb {
				configStr := config.Mid + ".db"
				cacheInstance := cache.NewInstance()
				cacheInstance.InitEnv(config.TargetCwd, config.Mid)
				cacheInstance.Start(configStr)
			}

		}

		if Pack {
			log.Info(config.Translate("WorkingOn", map[string]interface{}{"Item": "Packing"}))
			clt.Packing()
		}
	}
}

func main() {
	config.InitLog()
	defer config.CloseLog()

	run()
}
