package filewathcher

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-fsnotify/fsnotify"
)

// Event represents a single file system notification.
type Event struct {
	Name string // Relative path to the file or directory.
	Op   Op     // File operation that triggered the event.
}

// Op describes a set of file operations.
type Op uint32

// These are the generalized file operations that can trigger a notification.
const (
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

func (o Op) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if o&Create == Create {
		buffer.WriteString("|CREATE")
	}
	if o&Remove == Remove {
		buffer.WriteString("|REMOVE")
	}
	if o&Write == Write {
		buffer.WriteString("|WRITE")
	}
	if o&Rename == Rename {
		buffer.WriteString("|RENAME")
	}
	if o&Chmod == Chmod {
		buffer.WriteString("|CHMOD")
	}

	// If buffer remains empty, return no event names
	if buffer.Len() == 0 {
		return ""
	}

	// Return a list of event names, with leading pipe character stripped
	return fmt.Sprintf("%s", buffer.String()[1:])
}

// String returns a string representation of the event in the form
// "file: REMOVE|WRITE|..."
func (e Event) String() string {

	// Return a list of event names, with leading pipe character stripped
	return fmt.Sprintf("%q: %s", e.Name, e.Op.String())
}

//默认的文件监制器
var defaultWather *fsnotify.Watcher
var pathsLock sync.RWMutex

//监控关闭通知channel
var quitChan = make(chan int)

//缓存的监控处理信息 目录-->事件-->处理函数
var wathPaths = make(map[string]map[Op]func())

func init() {
	//文件监控
	var err error
	defaultWather, err = fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
		panic("创建默认监控器失败")
	}
}

//增加监控目录
func AddWathPath(filePath string, op Op, handler func()) error {

	pathsLock.Lock()
	defer pathsLock.Unlock()
	//增加事件处理函数
	events, ok := wathPaths[filePath]
	if !ok {
		//目录未监控过则设置目录监控
		err := defaultWather.Add(filePath)
		if err != nil {
			fmt.Println("增加监控目录失败,path:", filePath)
			return err
		}
		events = make(map[Op]func())
	}

	_, ok = events[op]
	if ok {
		return errors.New("目录的此事件已有处理函数")
	}

	fmt.Printf("增加监控目录%s的%s事件\n", filePath, op.String())

	events[op] = handler
	wathPaths[filePath] = events

	return nil
}

//移除监控目录
func RemoveWathPath(filePath string, op Op) error {
	pathsLock.Lock()
	defer pathsLock.Unlock()

	events, ok := wathPaths[filePath]
	if !ok {
		//目录未在监控列表，则直接返回成功
		fmt.Println("filePath:", filePath, "不在监控目录中")
		return nil
	}

	_, ok = events[op]
	if ok {
		//删除监控事件
		delete(events, op)
		fmt.Printf("删除目录%s的%s事件\n", filePath, op.String())
		fmt.Println("len(events):", len(events))
		fmt.Println("events:", events)
		//目录下的监控事件为空，则删除目录的监控
		if len(events) == 0 {
			err := defaultWather.Remove(filePath)
			if err != nil {
				fmt.Println("删除监控目录失败", filePath)
				return errors.New("删除监控目录失败")
			}
			//从监控目录中删除
			delete(wathPaths, filePath)
			return nil
		}
	}
	wathPaths[filePath] = events
	fmt.Println("删除事件后wathPaths:", wathPaths)

	return nil
}

//停止监控器
func Stop() {
	close(quitChan)
}

//关闭监控器
func Close() {
	close(quitChan)
	err := defaultWather.Close()
	if err != nil {
		fmt.Println("关闭文件监控器失败")
	} else {
		fmt.Println("关闭文件监控器成功")
	}
}

//启动监控器,带参数true为启动后台，不带则阻塞
func Run(bFlag ...bool) {
	if len(bFlag) >= 1 && bFlag[0] == true {
		go doEvent()
	} else {
		doEvent()
	}
}

func doEvent() {
	for {
		select {
		case event := <-defaultWather.Events:
			//根据不同的目录进行处理
			proWatherEvent(event)

		case err := <-defaultWather.Errors:
			fmt.Println("error:", err)
		case <-quitChan:
			fmt.Println("接收到退出信息")
			return
		}
	}
}

//事件处理器
func proWatherEvent(event fsnotify.Event) {
	fmt.Println("proWatherEvent event:", event.String())
	pathsLock.RLock()
	defer pathsLock.RUnlock()

	curEvent := Event{}
	curEvent.Name = event.Name
	curEvent.Op = Op(event.Op)

	for k, v := range wathPaths {
		if strings.HasSuffix(event.Name, k) {
			for ek, ev := range v {
				if ek&curEvent.Op == Create {
					go ev()
				}
				if ek&curEvent.Op == Write {
					go ev()
				}

				if ek&curEvent.Op == Remove {
					go ev()
				}

				if ek&curEvent.Op == Rename {
					go ev()
				}

				if ek&curEvent.Op == Chmod {
					go ev()
				}
			}
		}
	}
}
