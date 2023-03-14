package main

import (
	"fmt"

	"github.com/dolfly/webshark/pkg/sharkd/shark"
)

func main() {

	sharkcli := shark.NewSharkdClient("/tmp/sharkd.sock")
	{
		r, err := sharkcli.Load(shark.LoadParam{
			File: "/tmp/test.pcapng",
		})
		fmt.Printf("%+v,%v\n==\n", r, err)
	}
	// {
	// 	r, err := sharkcli.Analyse()
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }
	// {
	// 	r, err := sharkcli.Bye()
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }

	// {
	// 	r, err := sharkcli.Check(&shark.CheckParam{})
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }

	// {
	// 	r, err := sharkcli.Complete(&shark.CompleteParam{})
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }
	// {
	// 	r, err := sharkcli.Download(&shark.DownloadParam{})
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }
	// {
	// 	r, err := sharkcli.DumpConf(&shark.DumpConfParam{})
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }
	// {
	// 	r, err := sharkcli.Follow(&shark.FollowParam{
	// 		Follow: "UDP",
	// 		Filter: "udp.stream eq 0",
	// 	})
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }
	// {
	// 	r, err := sharkcli.Frame(&shark.FrameParam{
	// 		Frame: 1,
	// 	})
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }
	// {
	// 	r, err := sharkcli.Frames(&shark.FramesParam{})
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }
	// {
	// 	r, err := sharkcli.Info()
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }
	// {
	// 	r, err := sharkcli.Intervals(&shark.IntervalsParam{})
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }
	// {
	// 	r, err := sharkcli.IOGraph(&shark.IOGraphParam{
	// 		Graph0:  "packets",
	// 		Filter0: "frame.number <= 100",
	// 	})
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }

	// {
	// 	r, err := sharkcli.SetComment(&shark.SetCommentParam{
	// 		Frame:   1,
	// 		Comment: "test",
	// 	})
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }
	// {
	// 	r, err := sharkcli.SetConf(&shark.SetConfParam{
	// 		Name:  "tcp.desegment_tcp_streams",
	// 		Value: false,
	// 	})
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }
	// {
	// 	r, err := sharkcli.DumpConf(&shark.DumpConfParam{
	// 		Pref: "tcp.desegment_tcp_streams",
	// 	})
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }
	// {
	// 	r, err := sharkcli.Status()
	// 	fmt.Printf("%+v,%v\n==\n", r, err)
	// }
	{
		r, err := sharkcli.Tap(&shark.TapParam{
			Tap0: "conv:Ethernet",
		})
		fmt.Printf("%+v,%v\n==\n", r, err)
	}
}
