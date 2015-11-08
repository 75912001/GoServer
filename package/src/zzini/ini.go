//读取ini文件
//#开始的是注释行,不读取
//以下是ini文件的例子
/*
[server]
ip=192.168.8.101
port=9988

[common]
#child fd max value, def:20000
max_fd_num=1000
#tcp listen number, def:1024
#listen_num=1024
*/

//使用方法
/*
var ini ZZIni
ini.Path = "xxx.ini"
ini.Load()
ip := ini.Get_val_def("server", "ip", "")
*/

package zzini

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//加载文件
func (p *ZZIni) Load() (err error) {
	p.init()

	file, err := os.Open(p.Path)

	if nil != err {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var section string

	for {
		line, err := reader.ReadString('\n')
		if nil != err {
			break
		}
		line = strings.TrimSpace(line)
		switch {
		case 0 == len(line):
		case '#' == line[0]:
		//匹配[xxx]然后存储
		case '[' == line[0] && ']' == line[len(line)-1]:
			section = line[1 : len(line)-1]
		default:
			symbolIndex := strings.IndexAny(line, "=")
			if -1 == symbolIndex {
				break
			}
			key := line[0:symbolIndex]
			value := line[symbolIndex+1:]
			fmt.Println(section, key, value)
			p.load(section, key, value)
		}
	}

	return err
}

//获取对应的值
func (p *ZZIni) Get(section string, key string, defaultValue string) (value string) {
	sectionValue, valid := p.sectionMap[section]
	if valid {
		keyValue, valid := sectionValue[key]
		if valid {
			return keyValue
		}
	}
	return defaultValue
}

type KEY_MAP map[string]string
type SECTION_MAP map[string]KEY_MAP

//ini文件
type ZZIni struct {
	Path       string      //文件的路径
	sectionMap SECTION_MAP //存取配置文件
}

//加载文件到内存中
func (p *ZZIni) load(section string, key string, value string) {
	_, valid := p.sectionMap[section]
	if valid {
		p.sectionMap[section][key] = value
	} else {
		keyMap := make(KEY_MAP)
		keyMap[key] = value
		p.sectionMap[section] = keyMap
	}
}

//初始化
func (p *ZZIni) init() {
	p.sectionMap = make(SECTION_MAP)
}
