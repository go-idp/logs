package manager

import (
	"sort"
	"sync"
	"time"

	"github.com/go-zoox/core-utils/safe"
)

type Manager interface {
	Create(id string)
	Update(id string, message string)
	Delete(id string)
	//
	Status() Status
	//
	IsRunning(id string) bool
}

type Task struct {
	ID         string     `json:"id"`
	Size       safe.Int64 `json:"size"`
	StartedAt  *time.Time `json:"started_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

type Status struct {
	Count  StatusCount  `json:"count"`
	Detail StatusDetail `json:"detail"`
}

type StatusCount struct {
	Total    int64 `json:"total"`
	Running  int64 `json:"running"`
	Finished int64 `json:"finished"`
}

type StatusDetail struct {
	Runnings  []*Task `json:"runnings"`
	Finisheds []*Task `json:"finisheds"`
}

type manager struct {
	runnings  *safe.Map[string, *Task]
	finisheds *safe.Map[string, *Task]

	status Status

	sync.RWMutex
}

func New() Manager {
	runnings := safe.NewMap[string, *Task](func(mc *safe.MapConfig) {
		mc.Capacity = 30
	})

	finisheds := safe.NewMap[string, *Task](func(mc *safe.MapConfig) {
		mc.Capacity = 100
	})

	return &manager{
		runnings:  runnings,
		finisheds: finisheds,
	}
}

func (m *manager) Create(id string) {
	m.Lock()
	defer m.Unlock()

	ins := &Task{
		ID:        id,
		StartedAt: now(),
	}

	m.runnings.Set(id, ins)
	m.status.Count.Total++
	m.status.Count.Running++
}

func (m *manager) Update(id string, message string) {
	m.Lock()
	defer m.Unlock()

	ins := m.runnings.Get(id)
	if ins == nil {
		return
	}

	ins.Size.Set(ins.Size.Get() + int64(len(message)))
	ins.UpdatedAt = now()

	m.runnings.Set(id, ins)
}

func (m *manager) Delete(id string) {
	m.Lock()
	defer m.Unlock()

	ins := m.runnings.Get(id)
	if ins == nil {
		return
	}

	ins.FinishedAt = now()

	m.runnings.Del(id)
	m.finisheds.Set(id, ins)
	m.status.Count.Running--
	m.status.Count.Finished++
}

func (m *manager) Status() Status {
	m.RLock()
	defer m.RUnlock()

	m.status.Detail.Runnings = m.GetRunnings()
	m.status.Detail.Finisheds = m.GetFinisheds()

	return m.status
}

func (m *manager) GetRunnings() (tasks []*Task) {
	tasks = make([]*Task, 0)

	m.runnings.ForEach(func(key string, task *Task) bool {
		tasks = append(tasks, task)
		return false
	})

	sort.Slice(tasks, func(i, j int) bool {
		return (*tasks[i].StartedAt).After(*(tasks[j].StartedAt))
	})

	return
}

func (m *manager) GetFinisheds() (tasks []*Task) {
	tasks = make([]*Task, 0)

	m.finisheds.ForEach(func(key string, task *Task) bool {
		tasks = append(tasks, task)
		return false
	})

	sort.Slice(tasks, func(i, j int) bool {
		return (*tasks[i].FinishedAt).After(*(tasks[j].FinishedAt))
	})

	return
}

func (m *manager) IsRunning(id string) bool {
	m.RLock()
	defer m.RUnlock()

	if task := m.runnings.Get(id); task != nil {
		return true
	}

	return false
}

func now() *time.Time {
	now := time.Now()
	return &now
}
