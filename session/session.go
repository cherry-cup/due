/**
 * @Author: fuxiao
 * @Email: 576101059@qq.com
 * @Date: 2022/6/9 20:10
 * @Desc: TODO
 */

package session

import (
	"net"
	"sync"

	"github.com/dobyte/due/network"
)

type Session struct {
	rw     sync.RWMutex        // 读写锁
	uid    int64               // 用户ID
	conn   network.Conn        // 连接
	groups map[*Group]struct{} // 所在组
}

func NewSession() *Session {
	return &Session{}
}

// Init 初始化SESSION
func (s *Session) Init(conn network.Conn) {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.conn = conn
	s.groups = make(map[*Group]struct{})
}

// Reset 重置SESSION
func (s *Session) Reset() {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.uid = 0
	s.conn = nil
	s.groups = nil
}

// CID 获取连接ID
func (s *Session) CID() int64 {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.conn.ID()
}

// UID 获取用户ID
func (s *Session) UID() int64 {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.uid
}

// Bind 绑定用户ID
func (s *Session) Bind(uid int64) {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.uid = uid
	s.conn.Bind(uid)
	for group := range s.groups {
		group.addUserMap(uid, s)
	}
}

// Close 关闭会话
func (s *Session) Close(isForce ...bool) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	return s.conn.Close(isForce...)
}

// LocalIP 获取本地IP
func (s *Session) LocalIP() (string, error) {
	return s.conn.LocalIP()
}

// LocalAddr 获取本地地址
func (s *Session) LocalAddr() (net.Addr, error) {
	return s.conn.LocalAddr()
}

// RemoteIP 获取远端IP
func (s *Session) RemoteIP() (string, error) {
	return s.conn.RemoteIP()
}

// RemoteAddr 获取远端地址
func (s *Session) RemoteAddr() (net.Addr, error) {
	return s.conn.RemoteAddr()
}

// Send 发送消息（同步）
func (s *Session) Send(msg []byte, msgType ...int) error {
	return s.conn.Send(msg, msgType...)
}

// Push 发送消息（异步）
func (s *Session) Push(msg []byte, msgType ...int) error {
	return s.conn.Push(msg, msgType...)
}

// JoinGroup 加入群组
func (s *Session) JoinGroup(groups ...*Group) {
	s.rw.Lock()
	defer s.rw.Unlock()

	for i := range groups {
		group := groups[i]
		group.add(s)
		s.groups[group] = struct{}{}
	}
}

// 加入群组
func (s *Session) joinGroup(groups ...*Group) {
	s.rw.Lock()
	defer s.rw.Unlock()

	for i := range groups {
		group := groups[i]
		s.groups[group] = struct{}{}
	}
}

// QuitGroup 退出群组
func (s *Session) QuitGroup(groups ...*Group) {
	s.rw.Lock()
	defer s.rw.Unlock()

	for _, group := range groups {
		_ = group.remSession(s.CID())
		delete(s.groups, group)
	}
}

// 退出群组
func (s *Session) quitGroup(groups ...*Group) {
	s.rw.Lock()
	defer s.rw.Unlock()

	for _, group := range groups {
		delete(s.groups, group)
	}
}