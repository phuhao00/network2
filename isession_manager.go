package network

type ISessionManager[k comparable, v ISession] interface {
	AddSession(id k, session v) bool
	DelSession(id k)
	GeSession(id k) v
	CountSessions() int
	Clear()
	RangeSessions(func(id k, session v) bool)
}
