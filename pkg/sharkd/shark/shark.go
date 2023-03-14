package shark

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type request struct {
	Version string `json:"jsonrpc"`
	Id      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

type response struct {
	Version string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Result  interface{} `json:"result"`
	Error   any         `json:"error"`
}

type SharkdClient struct {
	sockpath  string
	connected bool
	conn      net.Conn
}

// https://wiki.wireshark.org/sharkd-JSON-RPC

func NewSharkdClient(sockpath string) *SharkdClient {
	return &SharkdClient{
		sockpath:  sockpath,
		connected: false,
	}
}

type AnalyseResult struct {
	Frames    int64    `json:"frames"`
	Protocols []string `json:"protocols"`
	First     float64  `json:"first,omitempty"`
	Last      float64  `json:"last,omitempty"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#analyse
func (c *SharkdClient) Analyse() (r *AnalyseResult, err error) {
	r = &AnalyseResult{}

	res, err := c.send(1, "analyse", nil)
	if err != nil {
		return
	}
	rsp := response{
		Result: &AnalyseResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*AnalyseResult), nil
}

type ByeResult struct {
	Status string `json:"status"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#bye

func (c *SharkdClient) Bye() (r *ByeResult, err error) {
	res, err := c.send(3, "bye", nil)
	if err != nil {
		return
	}
	rsp := response{
		Result: &ByeResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*ByeResult), nil
}

type CheckParam struct {
	Field  string `json:"field,omitempty"`
	Filter string `json:"filter,omitempty"`
}

type CheckResult struct {
	Status string `json:"status"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#check

func (c *SharkdClient) Check(param *CheckParam) (r *CheckResult, err error) {
	res, err := c.send(4, "check", param)
	if err != nil {
		return
	}
	rsp := response{
		Result: &CheckResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*CheckResult), nil
}

type CompleteParam struct {
	Field string `json:"field,omitempty"`
	Pref  string `json:"pref,omitempty"`
}

type Field struct {
	Reference string `json:"f"`
	Type      string `json:"t"`
	Name      string `json:"n"`
}

type Preference struct {
	Name string `json:"f"`
	Desc string `json:"d"`
}

type CompleteResult struct {
	Fields []Field      `json:"field,omitempty"`
	Prefs  []Preference `json:"pref,omitempty"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#complete
func (c *SharkdClient) Complete(param *CompleteParam) (r *CompleteResult, err error) {
	res, err := c.send(5, "complete", param)
	if err != nil {
		return
	}
	rsp := response{
		Result: &CompleteResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*CompleteResult), nil
}

type DownloadParam struct {
	Token string `json:"token"`
}

type DownloadResult struct {
	File string `json:"file"`
	Mime string `json:"mime"`
	Data string `json:"data"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#download

func (c *SharkdClient) Download(param *DownloadParam) (r *DownloadResult, err error) {
	res, err := c.send(6, "download", param)
	if err != nil {
		return
	}
	rsp := response{
		Result: &DownloadResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*DownloadResult), nil
}

type DumpConfParam struct {
	Pref string `json:"pref,omitempty"`
}

type DumpConfResult struct {
	Prefs map[string]struct {
		Binary   int    `json:"b,omitempty"` //0 - not set, 1 - set
		Desc     string `json:"d,omitempty"`
		DropDown []struct {
			Value    int64  `json:"v,omitempty"`
			Selected int    `json:"s,omitempty"`
			Desc     string `json:"d,omitempty"`
		} `json:"e,omitempty"`
		Range  string        `json:"r,omitempty"`
		String string        `json:"s,omitempty"`
		Table  []interface{} `json:"t,omitempty"`
	}
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#dumpconf
func (c *SharkdClient) DumpConf(param *DumpConfParam) (r *DumpConfResult, err error) {
	res, err := c.send(7, "dumpconf", param)
	if err != nil {
		return
	}
	rsp := response{
		Result: &DumpConfResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*DumpConfResult), nil
}

type FollowParam struct {
	Follow string `json:"follow"`
	Filter string `json:"filter"`
}

type FollowResult struct {
	SHost    string `json:"shost"`
	SPort    string `json:"sport"`
	SBytes   int64  `json:"sbytes"`
	CHost    string `json:"chost"`
	CPort    string `json:"cport"`
	CBytes   int64  `json:"cbytes"`
	Payloads []struct {
		N int64  `json:"n"`
		D string `json:"d"`
		S int    `json:"s"`
	} `json:"payloads"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#follow

func (c *SharkdClient) Follow(param *FollowParam) (r *FollowResult, err error) {
	res, err := c.send(8, "follow", param)
	if err != nil {
		return
	}
	rsp := response{
		Result: &FollowResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*FollowResult), nil
}

type FrameParam struct {
	Frame   int  `json:"frame"`
	Proto   bool `json:"proto,omitempty"`
	Ref     bool `json:"ref_frame,omitempty"`
	Pre     bool `json:"prev_frame,omitempty"`
	Columns bool `json:"columns,omitempty"`
	Color   bool `json:"color,omitempty"`
	Bytes   bool `json:"bytes,omitempty"`
	Hidden  bool `json:"hidden,omitempty"`
}

type FrameResult map[string]interface{}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#frame
func (c *SharkdClient) Frame(param *FrameParam) (r *FrameResult, err error) {
	res, err := c.send(8, "frame", param)
	if err != nil {
		return
	}
	rsp := response{
		Result: &FrameResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*FrameResult), nil
}

type FramesParam struct {
	Filter string `json:"filter,omitempty"`
	Skip   int    `json:"skip,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Refs   string `json:"refs,omitempty"`
}
type FramesResult []struct {
	Columns []string `json:"c,omitempty"`
	Num     int64    `json:"num,omitempty"`
	Bg      string   `json:"bg,omitempty"`
	Fg      string   `json:"fg,omitempty"`
}

// TODO:具体实现
// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#frames
func (c *SharkdClient) Frames(param *FramesParam) (r *FramesResult, err error) {
	res, err := c.send(9, "frames", param)
	if err != nil {
		return
	}
	rsp := response{
		Result: &FramesResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*FramesResult), nil
}

type InfoResult struct {
	Columns []struct {
		Name   string `json:"name"`
		Format string `json:"format"`
	} `json:"columns,omitempty"`
	Stats []struct {
		Name string `json:"name,omitempty"`
		Tap  string `json:"tap,omitempty"`
	} `json:"stats,omitempty"`
	Ftypes  []string `json:"ftypes,omitempty"`
	Version string   `json:"version,omitempty"`
	Nstat   []struct {
		Name string `json:"name,omitempty"`
		Tap  string `json:"tap,omitempty"`
	} `json:"nstat,omitempty"`
	Convs []struct {
		Name string `json:"name,omitempty"`
		Tap  string `json:"tap,omitempty"`
	} `json:"convs,omitempty"`
	Seqa []struct {
		Name string `json:"name,omitempty"`
		Tap  string `json:"tap,omitempty"`
	} `json:"seqa,omitempty"`
	Taps []struct {
		Name string `json:"name,omitempty"`
		Tap  string `json:"tap,omitempty"`
	} `json:"taps,omitempty"`
	Eo []struct {
		Name string `json:"name,omitempty"`
		Tap  string `json:"tap,omitempty"`
	} `json:"eo,omitempty"`
	Srt []struct {
		Name string `json:"name,omitempty"`
		Tap  string `json:"tap,omitempty"`
	} `json:"srt,omitempty"`
	Rtd []struct {
		Name string `json:"name,omitempty"`
		Tap  string `json:"tap,omitempty"`
	} `json:"rtd,omitempty"`
	Follow []struct {
		Name string `json:"name,omitempty"`
		Tap  string `json:"tap,omitempty"`
	} `json:"follow,omitempty"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#info
func (c *SharkdClient) Info() (r *InfoResult, err error) {
	res, err := c.send(10, "info", nil)
	if err != nil {
		return
	}
	rsp := response{
		Result: &InfoResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*InfoResult), nil
}

type IntervalsParam struct {
	Interval int    `json:"interval,omitempty"`
	Filter   string `json:"filter,omitempty"`
}

type IntervalsResult struct {
	Intervals []interface{} `json:"intervals"`
	Last      int64         `json:"last"`
	Frames    int64         `json:"frames"`
	Bytes     int64         `json:"bytes"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#intervals
func (c *SharkdClient) Intervals(param *IntervalsParam) (r *IntervalsResult, err error) {
	res, err := c.send(12, "intervals", param)
	if err != nil {
		return
	}
	rsp := response{
		Result: &IntervalsResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*IntervalsResult), nil
}

type IOGraphParam struct {
	Interval int    `json:"interval,omitempty"`
	Filter   string `json:"filter,omitempty"`
	Graph0   string `json:"graph0"`
	Filter0  string `json:"filter0,omitempty"`
	Graph1   string `json:"graph1,omitempty"`
	Filter1  string `json:"filter1,omitempty"`
	Graph2   string `json:"graph2,omitempty"`
	Filter2  string `json:"filter2,omitempty"`
	Graph3   string `json:"graph3,omitempty"`
	Filter3  string `json:"filter3,omitempty"`
	Graph4   string `json:"graph4,omitempty"`
	Filter4  string `json:"filter4,omitempty"`
	Graph5   string `json:"graph5,omitempty"`
	Filter5  string `json:"filter5,omitempty"`
	Graph6   string `json:"graph6,omitempty"`
	Filter6  string `json:"filter6,omitempty"`
	Graph7   string `json:"graph7,omitempty"`
	Filter7  string `json:"filter7,omitempty"`
	Graph8   string `json:"graph8,omitempty"`
	Filter8  string `json:"filter8,omitempty"`
	Graph9   string `json:"graph9,omitempty"`
	Filter9  string `json:"filter9,omitempty"`
}

type IOGraphResult struct {
	IOGraph []interface{} `json:"iograph"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#iograph
func (c *SharkdClient) IOGraph(param *IOGraphParam) (r *IOGraphResult, err error) {
	res, err := c.send(13, "iograph", param)
	if err != nil {
		return
	}
	rsp := response{
		Result: &IOGraphResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*IOGraphResult), nil
}

type LoadParam struct {
	File string `json:"file"`
}

type LoadResult struct {
	Status string `json:"status"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#load
func (c *SharkdClient) Load(param LoadParam) (r *LoadResult, err error) {
	r = &LoadResult{}
	res, err := c.send(14, "load", param)
	if err != nil {
		return
	}
	rsp := response{
		Result: &LoadResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*LoadResult), nil
}

type SetCommentParam struct {
	Frame   int64  `json:"frame"`
	Comment string `json:"comment,omitempty"`
}

type SetCommentResult struct {
	Status string `json:"status"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#setcomment
func (c *SharkdClient) SetComment(param *SetCommentParam) (r *SetCommentResult, err error) {
	res, err := c.send(15, "setcomment", param)
	if err != nil {
		return
	}
	rsp := response{
		Result: &SetCommentResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*SetCommentResult), nil
}

type SetConfParam struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type SetConfResult struct {
	Status string `json:"status"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#setconf
func (c *SharkdClient) SetConf(param *SetConfParam) (r *SetConfResult, err error) {
	res, err := c.send(16, "setconf", param)
	if err != nil {
		return
	}
	rsp := response{
		Result: &SetConfResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*SetConfResult), nil
}

type StatusResult struct {
	Frames   int64   `json:"frames"`
	Duration float64 `json:"duration"`
	Filename string  `json:"filename"`
	Filesize int64   `json:"filesize"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#status
func (c *SharkdClient) Status() (r *StatusResult, err error) {
	res, err := c.send(17, "status", nil)
	if err != nil {
		return
	}
	rsp := response{
		Result: &StatusResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*StatusResult), nil
}

type TapParam struct {
	Tap0  string `json:"tap0"`
	Tap1  string `json:"tap1,omitempty"`
	Tap2  string `json:"tap2,omitempty"`
	Tap3  string `json:"tap3,omitempty"`
	Tap4  string `json:"tap4,omitempty"`
	Tap5  string `json:"tap5,omitempty"`
	Tap6  string `json:"tap6,omitempty"`
	Tap7  string `json:"tap7,omitempty"`
	Tap8  string `json:"tap8,omitempty"`
	Tap9  string `json:"tap9,omitempty"`
	Tap10 string `json:"tap10,omitempty"`
	Tap11 string `json:"tap11,omitempty"`
	Tap12 string `json:"tap12,omitempty"`
	Tap13 string `json:"tap13,omitempty"`
	Tap14 string `json:"tap14,omitempty"`
	Tap15 string `json:"tap15,omitempty"`
}

type TapResult struct {
	Taps []struct {
		Tap     string        `json:"tap"`
		Type    string        `json:"type"`
		Details []interface{} `json:"details"`
	} `json:"taps"`
}

// https://wiki.wireshark.org/sharkd-JSON-RPC-Request-Syntax#tap
func (c *SharkdClient) Tap(param *TapParam) (r *TapResult, err error) {
	res, err := c.send(18, "tap", param)
	if err != nil {
		return
	}
	rsp := response{
		Result: &TapResult{},
	}
	err = json.Unmarshal(res, &rsp)
	if err != nil {
		return
	}
	return rsp.Result.(*TapResult), nil
}
func (c *SharkdClient) connect() error {
	conn, err := net.DialTimeout("unix", c.sockpath, 2*time.Second)
	if err != nil {
		return err
	}
	c.conn = conn
	c.connected = true
	return nil
}

func (c *SharkdClient) send(id uint64, method string, params interface{}) (res []byte, err error) {
	if !c.connected {
		err = c.connect()
		if err != nil {
			return
		}
	}
	data, err := json.Marshal(request{
		Version: "2.0",
		Id:      id,
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return
	}
	fmt.Printf("send(%d): %s \n", len(data), string(data))
	c.conn.Write([]byte(string(data) + "\n"))
	reader := bufio.NewReaderSize(c.conn, 10*1024*1024)
	res, _, err = reader.ReadLine()
	fmt.Printf("recv(%d): %s \n", len(res), string(res))
	return
}

func (c *SharkdClient) Close() (err error) {
	if c.connected {
		err = c.conn.Close()
		if err == nil {
			c.connected = false
		}
	}
	return
}
