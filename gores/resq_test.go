package gores

import (
	"fmt"
	"testing"
)

func PrintItem(item map[string]interface{}) error {
	var err error
	for k, v := range item {
		fmt.Printf("key: %s, value: %s\n", k, v)
	}
	return err
}

var (
	config = &Config{
		RedisURL:          "127.0.0.1:6379",
		RedisPassword:     "mypassword",
		BlpopMaxBlockTime: 1,
		MaxWorkers:        2,
		Queues:            []string{"TestJob", "TestItem"},
	}
	resq = NewResQ(config)
	args = map[string]interface{}{"id": 1}
	item = map[string]interface{}{
		"Name":              "TestItem",
		"Queue":             "TestItem",
		"Args":              args,
		"Enqueue_timestamp": resq.CurrentTime(),
		"Retry":             true,
		"Retry_every":       10,
	}
	tasks = map[string]interface{}{
		"TestItem": PrintItem,
	}
	stat = NewStat("TestItem", resq)
)

func TestResQPushPop(t *testing.T) {
	err := resq.push("TestItem", item)
	if err != nil {
		t.Errorf("ResQ Push returned ERROR")
	}

	ret1, err := resq.Pop("TestItem")
	if err != nil {
		t.Errorf("%s", err)
	}
	if val, _ := ret1["Name"]; val != "TestItem" {
		t.Errorf("ResQ Pop Value ERROR")
	}

	ret2, err := resq.Pop("TestItem")
	if err != nil {
		t.Errorf("%s", err)
	}
	if ret2 != nil {
		t.Errorf("ResQ Pop expected to return nil, but did not")
	}
}

func TestResQSize(t *testing.T) {
	err := resq.push("TestItem", item)
	if err != nil {
		t.Errorf("ResQ Push returned ERROR")
	}

	size, err := resq.Size("TestItem")
	if err != nil {
		t.Errorf("%s", err)
	}

	if size != 1 {
		t.Errorf("ResQ Size() expected to return 1, but returned %d", size)
	}

	resq.Pop("TestItem")
}

func TestStat(t *testing.T) {
	v := stat.Get()
	if v != 0 {
		t.Errorf("Stat Get() expected to return 0 but returned %d", v)
	}

	ok := stat.Incr()
	if ok == 0 {
		t.Errorf("Stat Incr() Error")
	}

	if stat.Get() != v+1 {
		t.Errorf("Stat Incr() did not increment")
	}

	ok = stat.Decr()
	if ok == 0 {
		t.Errorf("Stat Decr() error")
	}
	if stat.Get() != v {
		t.Errorf("Stat Decr() did not decrement")
	}

	ok = stat.Clear()
	if ok == 0 {
		t.Errorf("Stat Clear() error")
	}
}
