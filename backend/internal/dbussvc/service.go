package dbussvc

import (
	"log"
	"sync"

	"github.com/godbus/dbus/v5"
)

const (
	ServiceName   = "net.dgkim.SendToLinux"
	ObjectPath    = "/net/dgkim/SendToLinux"
	InterfaceName = "net.dgkim.SendToLinux"
)

type RecentItem struct {
	ID    string
	Type  string
	Value string
	Size  uint32
}

type statusState struct {
	url     string
	port    uint32
	running bool
	qrPath  string
}

type Service struct {
	conn *dbus.Conn

	statusMu sync.RWMutex
	status   statusState

	recentMu sync.RWMutex
	recent   []RecentItem
}

func New(conn *dbus.Conn) *Service {
	return &Service{conn: conn}
}

func (s *Service) GetStatus() (string, uint32, bool, *dbus.Error) {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return s.status.url, s.status.port, s.status.running, nil
}

func (s *Service) GetQrPath() (string, *dbus.Error) {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return s.status.qrPath, nil
}

func (s *Service) GetRecentItems(limit uint32) ([]RecentItem, *dbus.Error) {
	s.recentMu.RLock()
	defer s.recentMu.RUnlock()

	if limit == 0 || int(limit) >= len(s.recent) {
		out := make([]RecentItem, len(s.recent))
		copy(out, s.recent)
		return out, nil
	}
	out := make([]RecentItem, limit)
	copy(out, s.recent[:limit])
	return out, nil
}

func (s *Service) EmitTestSignal() error {
	value := "test"
	return s.conn.Emit(dbus.ObjectPath(ObjectPath), InterfaceName+".ItemReceived", "test-0", "text", value, uint32(len(value)))
}

func (s *Service) EmitItemReceived(item RecentItem) {
	if err := s.conn.Emit(dbus.ObjectPath(ObjectPath), InterfaceName+".ItemReceived", item.ID, item.Type, item.Value, item.Size); err != nil {
		log.Printf("emit ItemReceived: %v", err)
	}
}

func (s *Service) AddRecent(item RecentItem) {
	const maxRecent = 50
	s.recentMu.Lock()
	defer s.recentMu.Unlock()
	s.recent = append([]RecentItem{item}, s.recent...)
	if len(s.recent) > maxRecent {
		s.recent = s.recent[:maxRecent]
	}
}

func (s *Service) SetStatus(url string, port uint32, running bool) {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	s.status.url = url
	s.status.port = port
	s.status.running = running
}

func (s *Service) SetQrPath(path string) {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	s.status.qrPath = path
}
