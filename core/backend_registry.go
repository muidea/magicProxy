package core

import (
	"log"
	"sync"
	"time"
)

type commandAction int

const (
	putData   commandAction = iota // 存放数据
	fetchData                      // 获取数据
	heartbeat                      // 心跳数据
	end                            // 停止
)

type commandData struct {
	action commandAction
	value  interface{}
	result chan<- interface{} //单向Channel
}

type putInData struct {
	backend Backend
}

type putInResult struct {
	result bool
}

type fetchOutData struct {
}

type fetchOutResult struct {
	result Backend
	err    error
}

type endData struct {
}

type endResult struct {
	result bool
}

// CommandDataChannel commandData channel
type CommandDataChannel chan *commandData

//Registry backend registry
type Registry struct {
	joinGroup   *sync.WaitGroup
	runningFlag bool
	backendList []Backend
	idleList    []Backend

	commandDataChannel CommandDataChannel

	config *Config
}

// NewRegistry new registry
func NewRegistry(cfg *Config, join *sync.WaitGroup) *Registry {
	registry := &Registry{config: cfg, joinGroup: join}
	registry.backendList = []Backend{}
	registry.idleList = []Backend{}

	go registry.run()

	return registry
}

// Close close registry
func (s *Registry) Close() {
	s.end()
}

// FetchOut fetchout backend
func (s *Registry) FetchOut() (ret Backend, err error) {
	reply := make(chan interface{})

	param := &commandData{action: fetchData, result: reply}

	s.commandDataChannel <- param

	result := (<-reply).(*fetchOutResult)

	ret = result.result
	err = result.err

	return
}

// PutIn putin backend
func (s *Registry) PutIn(val Backend) {
	param := &commandData{action: putData, value: &putInData{backend: val}, result: nil}

	s.commandDataChannel <- param
}

func (s *Registry) heartbeat() {
	if s.commandDataChannel != nil {
		param := &commandData{action: heartbeat}

		s.commandDataChannel <- param
	}
}

func (s *Registry) end() {
	reply := make(chan interface{})

	param := &commandData{action: end, result: reply}

	s.commandDataChannel <- param

	<-reply
}

func (s *Registry) checkTimeOut() {
	s.joinGroup.Add(1)
	defer s.joinGroup.Done()

	timeOutTimer := time.NewTicker(1 * time.Second)
	for s.runningFlag {
		select {
		case <-timeOutTimer.C:
			s.heartbeat()
		}
	}

	log.Printf("checkTimeout finish.")
}

func (s *Registry) run() {
	s.runningFlag = true
	s.commandDataChannel = make(CommandDataChannel)

	go s.checkTimeOut()

	s.joinGroup.Add(1)
	defer s.joinGroup.Done()

	for s.runningFlag {
		var commandData *commandData
		select {
		case data, ok := <-s.commandDataChannel:
			if ok {
				commandData = data
			}
		}

		if commandData == nil {
			continue
		}

		switch commandData.action {
		case putData:
			val := commandData.value.(*putInData)
			ret := s.putInInternal(val)
			if ret != nil && commandData.result != nil {
				commandData.result <- ret
			}
		case fetchData:
			val := commandData.value.(*fetchOutData)
			ret := s.fetchOutInternal(val)
			if ret != nil && commandData.result != nil {
				commandData.result <- ret
			}
		case heartbeat:
			s.heartbeatInternal()
		case end:
			ret := s.endInternal()
			if ret != nil && commandData.result != nil {
				commandData.result <- ret
			}
		default:
		}
	}

	close(s.commandDataChannel)
	s.commandDataChannel = nil

	log.Printf("registry finished")
}

func (s *Registry) fetchOutInternal(val *fetchOutData) (ret *fetchOutResult) {
	if len(s.idleList) == 0 {
		backend, err := s.getBackend(s.config)

		ret = &fetchOutResult{result: backend, err: err}

		if err == nil {
			s.backendList = append(s.backendList, backend)
		}

		return
	}

	ret = &fetchOutResult{result: s.idleList[0], err: nil}

	s.backendList = append(s.backendList, s.idleList[0])
	s.idleList = s.idleList[1:]

	return
}

func (s *Registry) putInInternal(val *putInData) (ret *putInResult) {
	newList := []Backend{}
	for _, cur := range s.backendList {
		if cur != val.backend {
			newList = append(newList, cur)
		}
	}
	s.backendList = newList
	s.idleList = append(s.idleList, val.backend)

	ret = &putInResult{result: true}

	return
}

func (s *Registry) heartbeatInternal() {
	newList := []Backend{}
	for _, val := range s.backendList {
		err := val.Ping()
		if err == nil {
			newList = append(newList, val)
		} else {
			log.Printf("ping using backend failed, err:%s", err.Error())
		}
	}
	s.backendList = newList

	newIdle := []Backend{}
	for _, val := range s.idleList {
		err := val.Ping()
		if err == nil {
			newIdle = append(newIdle, val)
		} else {
			log.Printf("ping idle backend failed, err:%s", err.Error())
		}
	}
	s.idleList = newIdle
}

func (s *Registry) endInternal() (ret *endResult) {
	newList := []Backend{}
	for _, val := range s.backendList {
		val.Close()
	}
	s.backendList = newList

	for _, val := range s.idleList {
		val.Close()
	}
	s.idleList = newList

	s.runningFlag = false

	ret = &endResult{result: true}

	return
}

func (s *Registry) getBackend(cfg *Config) (ret Backend, err error) {
	co := new(backConn)
	err = co.Connect(cfg.MySQLURI, cfg.MySQLUser, cfg.MySQLPswd)
	if err != nil {
		log.Printf("connect backend failed, url:%s, %s", cfg.MySQLURI, err.Error())
		co.Close()
		return nil, err
	}

	return co, nil
}
