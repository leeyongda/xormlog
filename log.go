package xormlog

import (
	"github.com/sirupsen/logrus"
	"xorm.io/builder"
	xormlog "xorm.io/xorm/log"
)

type LogCtx struct {
	logger  *logrus.Logger
	showSQL bool
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

func NewLogCtx(l *logrus.Logger) *LogCtx {
	return &LogCtx{logger: l}
}

func (l *LogCtx) BeforeSQL(ctx xormlog.LogContext) {}

func (l *LogCtx) AfterSQL(ctx xormlog.LogContext) {
	fullSqlStr, err := builder.ConvertToBoundSQL(ctx.SQL, ctx.Args)
	if err != nil {
		l.logger.Errorf("[SQL] %v %v - %v", ctx.SQL, ctx.Args, ctx.ExecuteTime)
	} else {
		l.logger.Infof("[SQL] %v - %v", fullSqlStr, ctx.ExecuteTime)
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
