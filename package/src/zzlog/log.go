package zzlog

type Log_t struct {
}

//踪迹日志
func (p *Log_t) Trace() {

}

//调试日志
func (p *Log_t) Debug() {

}

//报告日志
func (p *Log_t) Info() {
}

//公告日志
func (p *Log_t) Notice() {
}

//警告日志
func (p *Log_t) Warning() {

}

//错误日志
func (p *Log_t) Error() {

}

//临界日志
func (p *Log_t) Crit() {
}

//不可用日志
func (p *Log_t) Emerg() {
}
