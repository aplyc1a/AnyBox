package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aplyc1a/gair-framework/config"
	"github.com/aplyc1a/gair-framework/winevtx"
	"github.com/aplyc1a/utils/fs"
)

var (
	Import    string
	EventID   string
	EventType string
	KeyWord   string
	//IP        bool
	EXE bool
)

var (
	version      string = "0.0.1"
	innername    string = "evtxdump"
	default_path string = `C:\Windows\System32\winevt\Logs`
)

func init() {

	// 新增应用程序特有的命令行参数
	flag.StringVar(&Import, "import", default_path, "")
	flag.StringVar(&Import, "i", default_path, "")

	flag.StringVar(&EventID, "id", "", "")
	//flag.StringVar(&EventType, "service", "", "")
	//flag.StringVar(&EventType, "srv", "", "")
	flag.StringVar(&KeyWord, "keyword", "", "")

	//flag.BoolVar(&IP, "ip", false, "")
	flag.BoolVar(&EXE, "exe", false, "")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", innername)
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
			{"--import", "select the path of *.evtx."},
			{"--id", "select the event id."},
			//{"--srv", "select the event type."},
			{"--keyword", "select the record matched the provide keyword."},
			//{"--ip", "select the record matched ip."},
			{"--exe", "select the record matched exe.(equals --keyword .exe)"},
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
		fmt.Printf("%s\n", version)
		os.Exit(1)
	}

	config.LoadTranslations()
	config.CheckCommonsRequiredParams()

}

func CheckRequiredParams() {
	if EXE && KeyWord != "" {
		fmt.Printf("Option Conflicts between --exe and --keyword.\n")
		os.Exit(1)
	}
	if EventID != "" {
		if !fs.IsNumeric(EventID) {
			fmt.Printf("--id should provided eventid, like 4624, 4648.\n")
			os.Exit(1)
		}
	}
}

func main() {
	var err error
	now := time.Now()
	base := now.Format("20060102150405")
	if !filepath.IsAbs(Import) {
		Import, err = filepath.Abs(Import)
		if err != nil {
			fmt.Printf("err:%v\n", err)
		}
		if Import == default_path {
			if !fs.CheckAdminWindows() {
				fmt.Printf("Please run %s with administrator priviledge!\n", os.Args[0])
				os.Exit(1)
			}
		}
	}

	ret, err := fs.IsDirectory(Import)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		os.Exit(1)
	}

	var files []string
	if ret {
		files, err = fs.GetDirChildFiles(Import)
		if err != nil {
			fmt.Printf("err:%v\n", err)
			os.Exit(1)
		}
	} else {
		files = append(files, filepath.Base(Import))
	}
	if EXE {
		KeyWord = ".exe"
	}
	if len(files) > 1 {
		os.MkdirAll(base, 0644)
		os.Chdir(base)
		for _, f := range files {
			filename := filepath.Base(f[:len(f)-len(filepath.Ext(f))])
			json := filename + ".json"
			csv := filename + ".csv"
			winevtx.ParsedEvtxToJson(fs.Join(Import, f), json, KeyWord, EventID)
			ret, err := fs.IsFileEmpty(json)
			if err != nil {
				fmt.Printf("err:%v\n", err)
				e := os.Remove(json)
				if e != nil {
					fmt.Printf("e:%v\n", e)
				}
				continue
			}
			if ret {
				e := os.Remove(json)
				if e != nil {
					fmt.Printf("e:%v\n", e)
				}
				continue
			}
			fs.ParseJsonFileToCSV(json, csv, true)
			e := os.Remove(json)
			if e != nil {
				fmt.Printf("e:%v\n", e)
			}
			fmt.Printf("..%s --> %s\n", f, fs.Join(base, csv))
		}

		/*
			if IP {
				fmt.Printf("============================================\n")
				matched1, err := fs.Grep_rn("./", fs.RMTIPv4Pattern)
				if err != nil {
					fmt.Printf("err:%v\n", err)
				}
				matched2, err := fs.Grep_rn("./", fs.RMTIPv6Pattern)
				if err != nil {
					fmt.Printf("err:%v\n", err)
				}
				for _, x := range matched1 {
					fmt.Printf("%v\n", x)
				}
				for _, x := range matched2 {
					fmt.Printf("%v\n", x)
				}
			}
		*/

	} else {
		filename := filepath.Base(files[0][:len(files[0])-len(filepath.Ext(files[0]))])
		json := filename + ".json"
		csv := filename + ".csv"
		winevtx.ParsedEvtxToJson(Import, json, KeyWord, EventID)
		ret, err := fs.IsFileEmpty(json)
		if err != nil {
			e := os.Remove(json)
			if e != nil {
				fmt.Printf("e:%v\n", e)
			}
			fmt.Printf("err:%v\n", err)
			os.Exit(1)
		}
		if ret {
			e := os.Remove(json)
			if e != nil {
				fmt.Printf("e:%v\n", e)
			}
			fmt.Printf("err:%s is empty\n", json)
			os.Exit(1)
		}
		fs.ParseJsonFileToCSV(json, csv, true)
		e := os.Remove(json)
		if e != nil {
			fmt.Printf("e:%v\n", e)
		}
		fmt.Printf("...%s --> %s\n", Import, csv)
		/*
			if IP {
				fmt.Printf("============================================\n")
				matched1, err := fs.Grep_rn(csv, fs.RMTIPv4Pattern)
				if err != nil {
					fmt.Printf("err:%v\n", err)
				}
				matched2, err := fs.Grep_rn(csv, fs.RMTIPv6Pattern)
				if err != nil {
					fmt.Printf("err:%v\n", err)
				}
				for _, x := range matched1 {
					fmt.Printf("%v\n", x)
				}
				for _, x := range matched2 {
					fmt.Printf("%v\n", x)
				}
			}
		*/
	}

}
