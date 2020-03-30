package xormlog

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"xorm.io/builder"
	xormlog "xorm.io/xorm/log"
)

type LogCtx struct {
	logger   *logrus.Logger
	showSQL  bool
	asyncSQL bool
	sqlCtx   chan sqlCtx
	done     chan struct{}
	lock     sync.Mutex
}

type sqlCtx struct {
	fullSql     string
	executeTime time.Duration
	err         error
}

var xormlog2logrusLevel = map[xormlog.LogLevel]logrus.Level{
	xormlog.LOG_DEBUG:   logrus.DebugLevel,
	xormlog.LOG_INFO:    logrus.InfoLevel,
	xormlog.LOG_WARNING: logrus.WarnLevel,
	xormlog.LOG_ERR:     logrus.ErrorLevel,
}

var logrus2xormlogLevel = map[logrus.Level]xormlog.LogLevel{
	logrus.DebugLevel: xormlog.LOG_DEBUG,
	logrus.InfoLevel:  xormlog.LOG_INFO,
	logrus.WarnLevel:  xormlog.LOG_WARNING,
	logrus.ErrorLevel: xormlog.LOG_ERR,
}

var (
	_ xormlog.ContextLogger = &LogCtx{}
)

func NewLogCtx(l *logrus.Logger) *LogCtx {
	return &LogCtx{logger: l, sqlCtx: make(chan sqlCtx, 10), done: make(chan struct{}, 1)}
}

func (l *LogCtx) AsyncShowSQL(show bool) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.asyncSQL = show
	if l.asyncSQL {
		l.recvSql()
		return
	}
	// 否则关闭通道
	close(l.sqlCtx)
}

func (l *LogCtx) recvSql() {
	go func() {
		for {
			select {
			// 一直阻塞,等待接收sql
			case s, ok := <-l.sqlCtx:
				// 如果通道关闭了，则退出协程
				if !ok {
					return
				}
				l.logger.Infof("[SQL] %v - %v", s.fullSql, s.executeTime)
			case <-l.done:
				return
			}
		}
	}()
}

func (l *LogCtx) Close() {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.asyncSQL {
		l.done <- struct{}{}
	}
}

func (l *LogCtx) BeforeSQL(ctx xormlog.LogContext) {}

func (l *LogCtx) AfterSQL(ctx xormlog.LogContext) {
	if l.asyncSQL {
		go func() {
			fullSqlStr, err := builder.ConvertToBoundSQL(ctx.SQL, ctx.Args)
			l.sqlCtx <- sqlCtx{
				fullSql:     fullSqlStr,
				executeTime: ctx.ExecuteTime,
				err:         err,
			}
		}()
		return
	}
	fullSqlStr, err := builder.ConvertToBoundSQL(ctx.SQL, ctx.Args)
	if err != nil {
		l.logger.Errorf("[SQL] %v %v - %v", ctx.SQL, ctx.Args, ctx.ExecuteTime)
	} else {
		l.logger.Infof("[SQL] %s - %v", fullSqlStr, ctx.ExecuteTime)
	}
}

func (l *LogCtx) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}

func (l *LogCtx) Errorf(format string, v ...interface{}) {
	l.logger.Errorf(format, v...)
}

func (l *LogCtx) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

func (l *LogCtx) Warnf(format string, v ...interface{}) {
	l.logger.Warnf(format, v...)
}

func (l *LogCtx) Level() xormlog.LogLevel {
	return logrus2xormlogLevel[l.logger.GetLevel()]
}

func (l *LogCtx) SetLevel(lv xormlog.LogLevel) {
	l.logger.SetLevel(xormlog2logrusLevel[lv])
}

func (l *LogCtx) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}

func (l *LogCtx) IsShowSQL() bool {
	return l.showSQL
}
