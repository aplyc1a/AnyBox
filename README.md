# AnyBox

跨平台工具集,持续开发中! 测试程序获得：output/*

## 0x00 编译环境准备

```shell
go env -w GOPROXY=https://proxy.golang.com.cn,direct
go env -w GO111MODULE=on
apt-get install libssl-dev zlib1g-dev libzstd-dev  autoconf automake libtool make gcc pkg-config
apt-get install binutils-mingw-w64*          gcc-mingw-w64*          g++-mingw-w64*          mingw-w64*
#apt-get install binutils-aarch64-linux-gnu   gcc-aarch64-linux-gnu   g++-aarch64-linux-gnu
#apt-get install binutils-arm-linux-gnueabihf gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf  cpp-arm-linux-gnueabihf gcc-arm-linux-gnueabihf gccgo-arm-linux-gnueabihf

# yara support
git clone https://github.com/VirusTotal/yara
cd yara
./bootstrap.sh
#./configure --host=x86_64-w64-mingw32 --enable-static --disable-shared
#./configure --enable-static --disable-shared
make

```

## 0x01 工具清单

| Name     | Description     | Target     |
| -------- | -------- | -------- |
| pudu | Windows及Linux平台失陷主机上的信息收集、打包、搜集 | windows/linux |
| jihe | 主机端上常见的攻击痕迹场景的分析 | linux |

## 0x02 pudu-普渡

pudu提供Windows及Linux平台失陷主机上的信息收集、打包、搜集的功能。（Information collection on the attacked host.）

### HELP

```text
Usage of pudu:
===============
Common parameters:
  --help/-h                 : Show these help message.
  --verb/-v                 : Enable verbose mode.
  --lang                    : Language for output (en or cn).
  --version/-V              : Show version.
  --quiet/-q                : Enable quiet mode.
  --logfile/-log            : Path to log file.
---------------
Classic parameters:
  --fast                    : Do Fast collecting.(default:normal mode).
  --all                     : Do the complete collecting.(default:normal mode).
  --pack                    : Pack the data after collecting.(Also:-p/--packing/--packaging).
  --packonly                : Just pack the collecting data.(Also:-po/--packingonly/--packagingonly).
  --net-mon                 : Record the system network session per 1s.
  --fd-mon                  : Record the system process-fd per 1s.
  --cache                   : Caching evtx/statmap/syslog/browserhistory into sqlitedb.(Also -c/--cachedb/--caching).
---------------
Advanced parameters:
  --target <filename,...>   : Set the file targets.(Also:-ts/--targets).
  --ignore <filename,...>   : Set the file/dirs to be ignored.(Also:-ign/-its/--ignore-target/--ignore-targets)
  --hashmap                 : Do the `target` hashmapping!
  --md5sum                  : Do the `target` md5sum-hashmapping!
  --sha256sum               : Do the `target` sha256sum-hashmapping!
  --sha1sum                 : Do the `target` sha1sum-hashmapping!
  --copy                    : Copy the `target` files based on your provided!
  --yara <rulesfile>        : Scans the `target` used yarascan!
  --grep <string>           : Search for PATTERNS in each `target`.
  --grep-rl <string>        : Search for PATTERNS in each `target`.
  --grep-rn <string>        : Search for PATTERNS in each `target`.

```

### EXAMPLE

#### *1 收集*

面对失陷主机信息收集的场景，主要有三种工作模式：快速、正常、完全。下面是具体实例：


```shell
# 快速
./pudu --fast

# 普通
./pudu

# 完全
./pudu --all

```

存在不少辅助参数在收集时可以按照需求开启：

```text
-v        :获取更多输出，能够看到程序执行过程中的verbose信息。
-q        :安静输出，少说废话。
-log      :记录执行过程所有输出到日志。
--pack    :收集完毕后及时进行打包。打包后的文件放在.cache目录下，名如gair-%TIMESTAMP%.tar.gz。
--packonly:不进行收集，只打包.cache目录下本机的数据，通常用于各方面收集完毕后统一打包。打包后的文件放在.cache目录下，名为gair-xxxxxx.tar.gz。

如：./pudu --fast -v --pack -log record.log
```

#### *2 搜集*

按照不同的条件收集主机上感兴趣的数据

这种模式下。要求指定感兴趣的目标（--target）及待忽略的目标（--ignore）。目标可以是一个也可以是多个，多个目标以','分隔，还可以通过添加@符号从文件中按行读取目标清单。

指定目标后，我们可以进行YaraScan（Yara规则扫描）、Grep（文件内容匹配检索）、HashMap（哈希计算）、Copy(自定义文件提取)，下面是部分示例：

*YaraScan:*

```shell

./pudu --yara <rulesfile> --target "/data/web/uploads/,/home/appuser,@evil_list.txt"

```

*Grep:*

```shell

#查找所有带有请求网站内容的文件
./pudu --grep-rl "curl http://" --target "evil1,evil2,@evil_list.txt"

#查找所有匹配到密码成功的行
.\pudu.exe --grep "/index.php?s=captch" --target "E:\\wwwroot"
#查找所有匹配到密码成功的行(带上行号)
.\pudu.exe --grep-rn "/index.php?s=captch" --target "E:\\wwwroot"

```

*Copy:*

```shell

# 批量收集自定义文件。相关文件收集到前期收集到的.cache，主机数据路径下。
./pudu --copy --target "~/infected,/tmp,/etc" --ignore "/tmp/ssh-xxxxx,/etc/shadow,/etc/shadow-"

```


*HashMap:*

```shell

# 根据提供的目标清单，遍历目录下文件或指定文件并进行哈希计算。
./pudu --hashmap --target "/tmp"
./pudu --sha1sum --target "/tmp"
./pudu --sha256sum --target "/tmp"
./pudu --md5sum --target "/tmp"

```

## 0x03 jihe-几何

当前实现了Linux下主机端上常见的攻击痕迹场景的分析。

### HELP

```text
Usage of jihe:
===============
Common parameters:
  --help/-h: Show these help message
  --verbose/-v: Enable verbose mode
  --lang: Language for output (en or cn)
  --version/-V: Show version
  --quiet/-q: Enable quiet mode
  --logfile/-log: Path to log file 
---------------
Specific parameters:
  --import: Import the collect-data.(also -i)
  --importonly: Only import the collect-data.(also -io,--import-only)
  --mid: Provide the machine id.
  --cache: Make the cache.(also -mc/-makecache)
  --recache: Only rebuild the cache.(also -rc/-rco)
  --dropcache: Remove the cache db.(also -dc)
  --cacheonly: Only make the cache.(also -co)
  --dropcacheonly: Only remove the cache db.(also -dco)

```

### EXAMPLE

*1 导入分析*

导入打包的数据，解压到.cache目录并立刻执行分析操作。可以通过verbose参数控制输出信息的等级。支持导入目录，会自动获取该目录下的所有主机数据，可用于迁移。但要求必须设置--mid区分分析目标才能进一步分析。

```shell
./jihe --import gair_20240629002935.tar.gz

./jihe --import .cache_old --mid 0fd37b449c1b4dffc82f15c1b4757d48
```

仅仅导入打包的数据，但不执行分析。导入的数据放置在.cache目录下，以设备ID命名。

```shell
./jihe --importonly gair_20240629002935.tar.gz
```

当.cache目录下有且仅有一个目标时，--mid参数可以省略。当分析目标已经导入后，后续无需再使用--import进行重复导入。换句话说，当处于现网分析时，可以直接运行如下命令：

```shell
./jihe -v
```
