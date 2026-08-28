package engine

import (
	"sync"
	"time"
)

type AuditEntry struct {
	Kind   string
	Ref    int64
	Action string
	At     time.Time
}

type auditLog struct {
	mu      sync.RWMutex
	entries []AuditEntry
}

func newAuditLog() *auditLog {
	return &auditLog{entries: make([]AuditEntry, 0, 16)}
}

func (l *auditLog) add(entry AuditEntry) {
	l.entries = append(l.entries, entry)
}

func (l *auditLog) snapshot() []AuditEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return append([]AuditEntry(nil), l.entries...)
}

type ReviewAudit struct{ log *auditLog }
type RunAudit struct{ log *auditLog }
type MeasurementAudit struct{ log *auditLog }
type SampleAudit struct{ log *auditLog }
type AlertAudit struct{ log *auditLog }
type ReportAudit struct{ log *auditLog }

func NewReviewAudit() *ReviewAudit           { return &ReviewAudit{log: newAuditLog()} }
func NewRunAudit() *RunAudit                 { return &RunAudit{log: newAuditLog()} }
func NewMeasurementAudit() *MeasurementAudit { return &MeasurementAudit{log: newAuditLog()} }
func NewSampleAudit() *SampleAudit           { return &SampleAudit{log: newAuditLog()} }
func NewAlertAudit() *AlertAudit             { return &AlertAudit{log: newAuditLog()} }
func NewReportAudit() *ReportAudit           { return &ReportAudit{log: newAuditLog()} }

func (a *ReviewAudit) Add(ref int64, action string, at time.Time) {
	a.log.add(AuditEntry{"review", ref, action, at})
}
func (a *RunAudit) Add(ref int64, action string, at time.Time) {
	a.log.add(AuditEntry{"run", ref, action, at})
}
func (a *MeasurementAudit) Add(ref int64, action string, at time.Time) {
	a.log.add(AuditEntry{"measurement", ref, action, at})
}
func (a *SampleAudit) Add(ref int64, action string, at time.Time) {
	a.log.add(AuditEntry{"sample", ref, action, at})
}
func (a *AlertAudit) Add(ref int64, action string, at time.Time) {
	a.log.add(AuditEntry{"alert", ref, action, at})
}
func (a *ReportAudit) Add(ref int64, action string, at time.Time) {
	a.log.add(AuditEntry{"report", ref, action, at})
}

func (a *ReviewAudit) Snapshot() []AuditEntry      { return a.log.snapshot() }
func (a *RunAudit) Snapshot() []AuditEntry         { return a.log.snapshot() }
func (a *MeasurementAudit) Snapshot() []AuditEntry { return a.log.snapshot() }
func (a *SampleAudit) Snapshot() []AuditEntry      { return a.log.snapshot() }
func (a *AlertAudit) Snapshot() []AuditEntry       { return a.log.snapshot() }
func (a *ReportAudit) Snapshot() []AuditEntry      { return a.log.snapshot() }
