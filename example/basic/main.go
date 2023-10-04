package main

import (
	"fmt"

	"github.com/markpassawat/go-pulse/pulse"
	"github.com/sirupsen/logrus"
)

func main() {
	p := pulse.NewPulse()
	a, err := p.Broadcast(&pulse.Asset{
		Symbol:    "ETH",
		Price:     4500,
		Timestamp: 1678912345,
	})
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(a.TxHash)
	// Output:
	// 044fd72daa83cd51f5ed10a9be8da73f4bf9b5ff167c1a8040e5db2ab15382c6

	b, err := p.BroadcastAndMonitor(&pulse.Asset{
		Symbol:    "ETH",
		Price:     4500,
		Timestamp: 1678912345,
	})
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(b)
	// Output:
	// &{131cc9cb7ee149a4a629d9fed28290026caf21bbb54cee8e661375ae9b124376 CONFIRMED Transaction has been processed and confirmed}

	c, err := p.MonitorStatus(a.TxHash)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(c)
	// Output:
	// &{044fd72daa83cd51f5ed10a9be8da73f4bf9b5ff167c1a8040e5db2ab15382c6 CONFIRMED Transaction has been processed and confirmed}

	d, err := p.MultipleBroadcast([]pulse.Asset{
		{
			Symbol:    "ETH",
			Price:     4500,
			Timestamp: 1678912345,
		},
		{
			Symbol:    "BTC",
			Price:     4500,
			Timestamp: 1678912345,
		},
	})
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(d)
	// Output:
	// [c76894b276ea228dffdbe0c166a1572a735e421feffa32401a3529857c799c9c 6787c66d4e9643230522a579376982f25ab4d56118d07560cff8ee8ee5281630]

	e := p.MultipleMonitorStatus(a.TxHash, b.TxHash)
	fmt.Println(e)
	// Output:
	// [{044fd72daa83cd51f5ed10a9be8da73f4bf9b5ff167c1a8040e5db2ab15382c6 CONFIRMED Transaction has been processed and confirmed}
	// {131cc9cb7ee149a4a629d9fed28290026caf21bbb54cee8e661375ae9b124376 CONFIRMED Transaction has been processed and confirmed}]

	f, err := p.MultipleBroadcastAndMonitor([]pulse.Asset{
		{
			Symbol:    "ETH",
			Price:     4500,
			Timestamp: 1678912345,
		},
		{
			Symbol:    "BTC",
			Price:     4500,
			Timestamp: 1678912345,
		},
	})
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(f)
	// Output:
	// [{a635ea045d86fe27319b491ac9d35e8227523f922401d2a66e47366b903003cb CONFIRMED Transaction has been processed and confirmed}
	// {ac143034403b45951506ed567297070d58c57d553fed9ccb4138929cde6a7174 CONFIRMED Transaction has been processed and confirmed}]
}
