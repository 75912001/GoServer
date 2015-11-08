package zzlog

type Log struct {
}

//踪迹日志
func (p *Log) Trace() {

}

//调试日志
func (p *Log) Debug() {

}

//报告日志
func (p *Log) Info() {
}

//公告日志
func (p *Log) Notice() {
}

//警告日志
func (p *Log) Warning() {

}

//错误日志
func (p *Log) Error() {

}

//临界日志
func (p *Log) Crit() {
}

//不可用日志
func (p *Log) Emerg() {
}
