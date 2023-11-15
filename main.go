package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Ttype represents a meaninglessness of our life

type Ttype struct {
	id         int
	cT         time.Time // время создания
	fT         time.Time // время выполнения
	taskRESULT []byte
}
type Result struct {
	tt  Ttype
	err error
}

func taskCreater(ctx context.Context) <-chan Ttype {
	tch := make(chan Ttype, 10)
	go func() {
		defer close(tch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				tt := Ttype{cT: time.Now(), id: int(time.Now().UnixNano())}
				tch <- tt // передаем таск на выполнение
			}
		}

	}()
	return tch
}

func removeTrailingZeros(n string) int {
	for n[len(n)-1] == '0' {
		n = n[:len(n)-1]
	}
	num, _ := strconv.Atoi(n)
	return num
}

func taskWorker(chachacha <-chan Ttype) <-chan Result {
	rc := make(chan Result)
	go func() {
		defer close(rc)
		for a := range chachacha {
			res := Result{
				tt: a,
			}
			if removeTrailingZeros(strconv.Itoa(a.cT.Nanosecond()))%2 > 0 {
				res.err = fmt.Errorf("task id %d time %s, error %s", a.id, a.cT, a.taskRESULT)
			} else {
				res.tt.taskRESULT = []byte("task has been successed")
			}
			a.fT = time.Now()
			rc <- res
			time.Sleep(time.Millisecond * 150)
		}
	}()

	return rc
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(3*time.Second, cancel)
	done := []Ttype{}
	errs := []error{}
	result := taskWorker(taskCreater(ctx))
	for res := range result {
		if res.err != nil {
			errs = append(errs, res.err)
		} else {
			done = append(done, res.tt)
		}
	}

	println("Errors:")
	for _, r := range errs {
		println(r.Error())
	}

	println("Done tasks:")
	for _, r := range done {
		println(r.id)
	}
}
