package server

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  watch_manager
 * @Version: 1.0.0
 * @Date: 2020/2/20 11:50 上午
 */

type WatchManager struct {
	WatchTable  map[string][]*Protolcol
	Watch2Paths map[*Protolcol][]string
}

type WatcherEvent struct {
	Type  EventType
	State State
	Path  string
}

type Event struct {
	Type  int32
	State int32
	Path  string
}

func (s *WatchManager) AddWatch(path string, protocol *Protolcol) {
	list, ok := s.WatchTable[path]
	if !ok {
		list = make([]*Protolcol, 0)
		s.WatchTable[path] = list
	}
	list = append(list, protocol)

	paths, ok := s.Watch2Paths[protocol]
	if !ok {
		paths = make([]string, 0)
		s.Watch2Paths[protocol] = paths
	}
	paths = append(paths, path)
}

func (s *WatchManager) RemoveWatch(protocol *Protolcol) {
	paths, ok := s.Watch2Paths[protocol]
	if !ok {
		return
	}
	for _, path := range paths {
		protocols := s.WatchTable[path]
		for key, value := range protocols {
			if value == protocol {
				protocols = append(protocols[:key], protocols[key+1:]...)
			}
		}
	}
}

func (s *WatchManager) TriggerWatch(path string, eventType EventType) {
	event := &WatcherEvent{
		Type:  eventType,
		State: StateConnected,
		Path:  path,
	}
	watchers, ok := s.WatchTable[path]
	if !ok || len(watchers) == 0 {
		return
	}
	for _, watch := range watchers {
		for key, value := range s.Watch2Paths[watch] {
			if value == path {
				s.Watch2Paths[watch] = append(s.Watch2Paths[watch][:key], s.Watch2Paths[watch][key+1:]...)
			}
		}
	}
	for _, watch := range watchers {
		watch.Process(event)
	}
}
