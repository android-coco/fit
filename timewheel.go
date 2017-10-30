package fit

import (
	"errors"
	// jkl_fmt "fmt"
	"sync"
	"time"
)

type TimeWheel struct {
	isOn             bool          // 启动/停止
	timeZoneDuration time.Duration // 卡槽/时钟区的跨度、粒度、周期时长
	wheelDuration    time.Duration // 时间轮的总周期 =  timeZoneDuration * timeZoneCount
	timeZoneCount    int           // 一个时间轮拥有的卡槽/时钟区的总数
	wheel            []*Iterator   // 卡槽/时钟区slice，存放iterator（timeout handler的容器） cap = timeZoneCount = maxIndex + 1
	ticker           *time.Ticker  // 循环的定时器
	lock             sync.Mutex    // 锁
	timezoneIndex    int           // 当前指向的卡槽/时区索引, 和ticker配合完成timewheel的轮询
	maxIndex         int           // 卡槽的最大索引，= timeZoneCount - 1
}

type Task interface {
	// 处理到期
	TimeWheelExpireFunc(key string)
}

// 时间轮卡槽迭代器，可以理解成某个车轮/时钟区的负责人，存储timeout handler,key为GUID
type Iterator struct {
	items map[string]*TimeOutHandler
}

// 超时处理对象
type TimeOutHandler struct {
	timeout time.Duration // 延迟时间
	index   int           // 触发时的卡槽索引
	turns   int           // 需要转动车轴的圈数
	task    Task          // 超时任务
}

const (
	// 时间轮一轮的周期时间 =  default_tick_duration * default_wheel_count = 1h
	default_timezone_duration = time.Second * 1 // 默认的卡槽/时钟区的跨度、粒度、周期时长
	default_timezone_count    = 360             // 默认时间轮拥有的卡槽/时钟区的总数
)

// initialize TimeWheel
func NewTimeWheel() *TimeWheel {
	tw := &TimeWheel{
		timeZoneDuration: default_timezone_duration,
		timeZoneCount:    default_timezone_count,
		timezoneIndex:    0,
		maxIndex:         default_timezone_count - 1,
		wheelDuration:    default_timezone_duration * default_timezone_count,
	}
	tw.wheel = tw.creatWheelIterators()
	return tw
}

// initialize wheel iterator slice
func (tw *TimeWheel) creatWheelIterators() []*Iterator {
	slice := make([]*Iterator, tw.timeZoneCount)
	for i := 0; i < tw.timeZoneCount; i++ {
		slice[i] = &Iterator{items: make(map[string]*TimeOutHandler)}
	}
	return slice
}

// start timewheel
func (tw *TimeWheel) Start() {
	tw.lock.Lock()
	defer tw.lock.Unlock()
	tw.isOn = true
	// JKWarning!  如果未调用Stop()， 调用NewTimeWheel重新生成timeWheel会怎样, tikcer会停止吗？ 待验证
	tw.ticker = time.NewTicker(tw.timeZoneDuration)
	go func() {
		for {
			select {
			case <-tw.ticker.C:
				// starting with index 1
				tw.timezoneIndex += 1
				// return to zero
				if tw.timezoneIndex == tw.maxIndex+1 {
					tw.timezoneIndex = 0
				}
				iterator := tw.wheel[tw.timezoneIndex]
				timeoutTasks := tw.fetchExpiredTasks(iterator)
				tw.notifyExpiredTaskHandler(timeoutTasks)
			default:
			}
		}
	}()
}

// 获取本轮指向的卡槽/时钟区中所有的超时任务
func (tw *TimeWheel) fetchExpiredTasks(iterator *Iterator) map[string]*TimeOutHandler {
	tw.lock.Lock()
	defer tw.lock.Unlock()
	handlers := make(map[string]*TimeOutHandler)
	for k, v := range iterator.items {
		if v.turns <= 0 {
			// 本轮超时的task handler，全部存储当前timeZoneIndex指向的卡槽/时钟区里
			handlers[k] = v
			delete(iterator.items, k)
		} else {
			// 下一轮甚至更后面的任务
			v.turns -= 1
		}
	}
	return handlers
}

// notify timeOut task handler, call TimeWheelExpireFunc()
func (tw *TimeWheel) notifyExpiredTaskHandler(timeoutTasks map[string]*TimeOutHandler) {
	for k, handler := range timeoutTasks {
		handler.task.TimeWheelExpireFunc(k)
	}
}

// stop timewheel
func (tw *TimeWheel) Stop() {
	tw.isOn = false
	tw.ticker.Stop()
}

// add task to the timewheel, the task will be scheduled at the suitable timezone
func (tw *TimeWheel) AddTask(task Task, timeout time.Duration) (string, error) {
	if task == nil {
		return "", errors.New("TimeWheel AddTask() Fail! Reason: task can't be nil")
	} else if timeout <= 0 {
		return "", errors.New("TimeWheel AddTask() Fail! Reason: timeout must be greater than 0")
	}
	timeoutHandler := &TimeOutHandler{
		timeout: timeout,
		task:    task,
	}
	sid := GUID()
	tw.scheduleTimeoutTask(timeoutHandler, sid)
	return sid, nil
}

// remove task with taskid
func (tw *TimeWheel) RemoveTask(taskId string) (*TimeOutHandler, error) {
	if taskId == "" {
		return nil, errors.New("RemoveTask Fail! Reason: taskId is invalid")
	}
	tw.lock.Lock()
	// 此处的查找算法可以再优化

	var dup *TimeOutHandler
	for _, iterator := range tw.wheel {
		if v, ok := iterator.items[taskId]; ok {
			delete(iterator.items, taskId)
			dup = v
			break
		}
	}
	tw.lock.Unlock()
	if dup != nil {
		return dup, nil
	} else {
		return nil, errors.New("RemoveTask Fail! Reason: can't find the value for taskId")
	}
}

func (tw *TimeWheel) UpdateTicker(taskId string) error {
	handler, err := tw.RemoveTask(taskId)
	if err != nil {
		return err
	} else {
		tw.scheduleTimeoutTask(handler, taskId)
		return nil
	}
}

// Core algorithm！timewheel will schedule the timeout task to the suitable timezone
func (tw *TimeWheel) scheduleTimeoutTask(timeoutTask *TimeOutHandler, tid string) {
	tw.lock.Lock()
	defer tw.lock.Unlock()

	// 不足一个卡槽粒度（间隔），算一个卡槽粒度值
	if timeoutTask.timeout < tw.timeZoneDuration {
		timeoutTask.timeout = tw.timeZoneDuration
	}

	// 剩余的时间轮圈数
	remainingTurns := timeoutTask.timeout / tw.wheelDuration
	// 按时间轮周期算，剩多少时长
	remainingRoundDelay := timeoutTask.timeout % tw.wheelDuration
	if remainingRoundDelay == 0 {
		remainingTurns -= 1
	}

	// 按卡槽周期/粒度算，剩多少卡槽
	relativetimezoneIndex := remainingRoundDelay / tw.timeZoneDuration
	// 按卡槽周期/粒度算，还剩多少碎片时长
	remainingTickerDelay := remainingRoundDelay % tw.timeZoneDuration
	// 有剩余，算入下一个卡槽  (timeoutTask.timeout = tw.timeZoneDuration 后应该不会存在 != 0, 待验证)
	if remainingTickerDelay != 0 {
		relativetimezoneIndex += 1
	}

	// 计算实际剩余的圈数和最终卡槽index = 当前的index + 剩余的相对卡槽数, 然后需要对最大卡槽数取余
	stopIndex := tw.timezoneIndex + int(relativetimezoneIndex)
	if stopIndex >= tw.timeZoneCount {
		stopIndex %= tw.timeZoneCount
		timeoutTask.turns = int(remainingTurns) + stopIndex/tw.timeZoneCount
	} else {
		timeoutTask.turns = int(remainingTurns)
	}

	timeoutTask.index = stopIndex
	iterator := tw.wheel[stopIndex]
	iterator.items[tid] = timeoutTask
	tw.wheel[stopIndex] = iterator
}
