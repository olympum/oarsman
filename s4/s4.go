package s4

import (
	"bufio"
	"bytes"
	"github.com/huin/goserial"
	jww "github.com/spf13/jwalterweatherman"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	UsbRequest                        = "USB"   // Application starting communication’s
	WrResponse                        = "_WR_"  // Hardware Type
	ExitRequest                       = "EXIT"  // Application is exiting
	OkResponse                        = "OK"    // Packet Accepted
	ErrorResponse                     = "ERROR" // Unknown packet
	PingResponse                      = "PING"  // Ping
	ResetRequest                      = "RESET" // Request the rowing computer to reset
	ModelInformationRequest           = "IV?"   // Request Model Information
	ModelInformationResponse          = "IV"    // Current model information
	ReadMemoryRequest                 = "IR"    // Read a memory location
	ReadMemoryResponse                = "ID"    // Value from a memory location
	StrokeStartResponse               = "SS"    // Start of stroke
	StrokeEndResponse                 = "SE"    // End of stroke
	PulseCountResponse                = "P"     // Pulse Count in the last 25mS
	DisplaySetIntensityRequest        = "DI"    // Display: Set Intensity
	DisplaySetDistanceRequest         = "DD"    // Display: Set Distance
	WorkoutSetDistanceRequest         = "WSI"   // Define a distance workout
	WorkoutSetDurationRequest         = "WSU"   // Define a duration workout
	IntervalWorkoutSetDistanceRequest = "WII"   // Define an interval distance workout
	IntervalWorkoutSetDurationRequest = "WIU"   // Define an interval duration workout
	AddIntervalWorkoutRequest         = "WIN"   // Add/End an interval to a workout
)

type S4Interface interface {
	Run(workout *S4Workout)
	Exit()
}

type Packet struct {
	cmd  string
	data []byte
}

func (p Packet) Bytes() []byte {
	var b bytes.Buffer
	b.Write([]byte(p.cmd))
	if p.data != nil {
		b.Write(p.data)
	}
	b.Write([]byte("\n"))
	return b.Bytes()
}

const (
	Unset             = 0
	ResetWaitingPing  = 1
	ResetPingReceived = 2
	WorkoutStarted    = 3
	WorkoutCompleted  = 4
	WorkoutExited     = 5
)

const (
	Meters = "1"
)

type AtomicEvent struct {
	Time  int64
	Label string
	Value uint64
}

var EndAtomicEvent = AtomicEvent{}

type S4 struct {
	port       io.ReadWriteCloser
	scanner    *bufio.Scanner
	workout    *S4Workout
	aggregator Aggregator
	debug      bool
}

func findUsbSerialModem() string {
	contents, _ := ioutil.ReadDir("/dev")

	for _, f := range contents {
		if strings.Contains(f.Name(), "cu.usbmodem") {
			return "/dev/" + f.Name()
		}
	}

	return ""
}

func openPort() io.ReadWriteCloser {
	name := findUsbSerialModem()
	if len(name) == 0 {
		jww.FATAL.Println("S4 USB serial modem port not found")
		os.Exit(-1)
	}

	c := &goserial.Config{Name: name, Baud: 115200, CRLFTranslate: true}
	p, err := goserial.OpenPort(c)
	if err != nil {
		jww.FATAL.Println(err)
		os.Exit(-1)
	}

	return p
}

func NewS4(eventChannel chan<- AtomicEvent, aggregateEventChannel chan<- AggregateEvent, debug bool) S4Interface {
	p := openPort()
	aggregator := newAggregator(eventChannel, aggregateEventChannel)
	s4 := S4{port: p, scanner: bufio.NewScanner(p), aggregator: aggregator, debug: debug}
	return &s4
}

func (s4 *S4) write(p Packet) {
	n, err := s4.port.Write(p.Bytes())
	if err != nil {
		jww.FATAL.Println(err)
		os.Exit(-1)
	}
	if s4.debug {
		jww.DEBUG.Printf("written %s (%d+1 bytes)", strings.TrimRight(string(p.Bytes()), "\n"), n-1)
	}
	time.Sleep(25 * time.Millisecond) // yield per spec
}

func (s4 *S4) read() {
	for s4.scanner.Scan() {
		b := s4.scanner.Bytes()
		if len(b) > 0 {
			if s4.debug {
				jww.DEBUG.Printf("read %s (%d+1 bytes)", string(b), len(b))
			}
			s4.onPacketReceived(b)
			if s4.workout.state == WorkoutCompleted || s4.workout.state == WorkoutExited {
				return
			}
		}
	}

	if err := s4.scanner.Err(); err != nil {
		jww.FATAL.Println(err)
		os.Exit(-1)
	}
}

func (s4 *S4) Run(workout *S4Workout) {
	// send connection command and start listening
	s4.workout = workout
	s4.workout.state = Unset
	s4.write(Packet{cmd: UsbRequest})
	s4.read()
	s4.Exit()
}

func (s4 *S4) Exit() {
	if s4.workout.state != WorkoutExited {
		s4.write(Packet{cmd: ExitRequest})
		s4.workout.state = WorkoutExited
		s4.aggregator.consume(EndAtomicEvent)
	}
}

func (s4 *S4) onPacketReceived(b []byte) {
	// responses can start with:
	// _ : _WR_
	// O : OK
	// E : ERROR
	// P : PING, P
	// S : SS, SE
	c := b[0]
	switch c {
	case '_':
		s4.wRHandler(b)
	case 'I':
		s4.informationHandler(b)
	case 'O':
		s4.oKHandler()
	case 'E':
		s4.errorHandler()
	case 'P':
		s4.pingHandler(b)
	case 'S':
		s4.strokeHandler(b)
	default:
		jww.INFO.Printf("Unrecognized packet: %s", string(b))
	}
}

func (s4 *S4) wRHandler(b []byte) {
	s := string(b)
	if s == "_WR_" {
		s4.write(Packet{cmd: ModelInformationRequest})
	} else {
		jww.INFO.Printf("Unknown WaterRower init command %s\n", s)
	}
}

func (s4 *S4) readMemoryRequest(address string, size string) {
	cmd := ReadMemoryRequest + size
	data := []byte(address)
	s4.write(Packet{cmd: cmd, data: data})
}

func (s4 *S4) oKHandler() {
	s4.aggregator.consume(AtomicEvent{
		Time:  millis(),
		Label: "okay",
		Value: 0})
}

func (s4 *S4) errorHandler() {
	if s4.workout.state == ResetPingReceived {
		s4.write(Packet{cmd: ResetRequest})
		s4.workout.state = ResetWaitingPing
	}
}

func (s4 *S4) pingHandler(b []byte) {
	c := b[1]
	switch c {
	case 'I': // PING
		if s4.workout.state == ResetWaitingPing {
			s4.workout.state = ResetPingReceived
			for e := s4.workout.workoutPackets.Front(); e != nil; e = e.Next() {
				s4.write(e.Value.(Packet))
			}
		}
		s4.aggregator.consume(AtomicEvent{
			Time:  millis(),
			Label: "ping",
			Value: 0})
	default: // P
		// The WaterRower is now provided with a toothed wheel and optical
		// detector such that a pulse train containing 57 pulses per
		// revolution can be recorded for the purposes of paddle speed
		// measurement
		pulses := string(b[1:3])
		value, _ := strconv.ParseUint(pulses, 16, 8)
		s4.aggregator.consume(AtomicEvent{
			Time:  millis(),
			Label: "pulses_per_25ms",
			Value: value})
	}
}

type MemoryEntry struct {
	label string
	size  string
	base  int
}

var g_memorymap = map[string]MemoryEntry{
	"055": MemoryEntry{"total_distance_meters", "D", 16},
	"1A9": MemoryEntry{"stroke_rate", "S", 16},
	"088": MemoryEntry{"watts", "D", 16},
	"08A": MemoryEntry{"calories", "T", 16},
	"148": MemoryEntry{"speed_cm_s", "D", 16},
	"1A0": MemoryEntry{"heart_rate", "D", 16}}

func (s4 *S4) strokeHandler(b []byte) {
	c := b[1]
	switch c {
	case 'S': // SS
		if s4.workout.state == ResetPingReceived {
			s4.workout.state = WorkoutStarted
			// these are the things we want captured from the S4
			for address, mmap := range g_memorymap {
				s4.readMemoryRequest(address, mmap.size)
			}
		}
		s4.aggregator.consume(AtomicEvent{
			Time:  millis(),
			Label: "stroke_start",
			Value: 1})
	case 'E': // SE
		s4.aggregator.consume(AtomicEvent{
			Time:  millis(),
			Label: "stroke_end",
			Value: 0})
	}
}

func (s4 *S4) informationHandler(b []byte) {
	c1 := b[1]
	switch c1 {
	case 'V': // version
		// e.g. IV40210
		msg := string(b)
		jww.INFO.Printf("WaterRower S%s %s.%s\n", msg[2:3], msg[3:5], msg[5:7])
		model, _ := strconv.ParseInt(msg[2:3], 0, 0)  // 4
		fwHigh, _ := strconv.ParseInt(msg[3:5], 0, 0) // 2
		fwLow, _ := strconv.ParseInt(msg[5:7], 0, 0)  // 10
		if model != 4 {
			jww.INFO.Println("not an S4 monitor")
		}
		if fwHigh != 2 {
			jww.INFO.Println("unsupported major S4 firmware version")
		}
		if fwLow != 10 {
			jww.INFO.Println("unsupported minor S4 firmware version")
		}

		// we are ready to start workout
		s4.workout.state = ResetWaitingPing
		s4.write(Packet{cmd: ResetRequest})

	case 'D': // memory value
		size := b[2]
		address := string(b[3:6])

		var l int
		switch size {
		case 'S':
			l = 1
		case 'D':
			l = 2
		case 'T':
			l = 3
		}
		v, err := strconv.ParseUint(string(b[6:(6+2*l)]), 16, 8*l)
		if err == nil {
			s4.aggregator.consume(AtomicEvent{
				Time:  millis(),
				Label: g_memorymap[address].label,
				Value: v})
			// we re-request the data
			if s4.workout.state == WorkoutStarted {
				s4.readMemoryRequest(address, string(size))
			}
		} else {
			jww.INFO.Println("error parsing int: ", err)
		}
	}
}

func millis() int64 {
	// we operate at 25ms resolution, so Unix() is too coarse
	// we use a syscall directly to avoid time parsing costs
	var tv syscall.Timeval
	syscall.Gettimeofday(&tv)
	millis := (int64(tv.Sec)*1e3 + int64(tv.Usec)/1e3)
	return millis
}
