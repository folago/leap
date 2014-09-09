//https://github.com/leapmotion/leapjs/blob/master/PROTOCOL.md

package leap

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	glm "github.com/folago/googlmath"
	"github.com/gorilla/websocket"
)

func init() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

const (
	Hovering string = "hovering"
	None     string = "none"
	Touching string = "touching"
)

var (
	WSPort string = "6437"
	WSHost string = "127.0.0.1"
	WSURL  string = "ws://127.0.0.1:6437/v1.json"
)

type Pointable struct {
	//Valid                 bool    `json: "valid"`
	Tool                  bool    `json:"tool"`
	ID                    int     `json:"id"`
	HandId                int     `json:"handId"`
	Length                float32 `json:"length"`
	Width                 float32 `json:"width"`
	Direction             Vector  `json:"direction"`
	TipPosition           Vector  `json:"tipPosition"`
	StabilizedTipPosition Vector  `json:"stabilizedTipPosition"`
	TipVelocity           Vector  `json:"tipVelocity"`
	TouchZone             string  `json:"touchZone"`
	TouchDistance         float32 `json:"touchDistance"`
	TimeVisible           float32 `json:"timeVisible"`
}

type Finger struct {
	ID int `json: "id"`
	*Pointable
}

type Tool struct {
	ID int `json: "id"`
	*Pointable
}

type Hand struct {
	Direction              Vector      `json:"direction"` //: array of floats (vector)
	Fingers                []Pointable `json:"fingers"`
	Tools                  []Pointable `json:"tools"`
	Pointables             []Pointable `json:"pointables"`
	ID                     int         `json:"id"`           //: integer
	PalmNormal             Vector      `json:"palmNormal"`   //: array of floats (vector)
	PalmPosition           Vector      `json:"palmPosition"` //: array of floats (vector)
	PalmVelocity           Vector      `json:"palmVelocity"` //: array of floats (vector)
	SphereCenter           Vector      `json:"sphereCenter"` //: array of floats (vector)
	SphereRadius           float32     `json:"sphereRadius"` //: float
	StabilizedPalmPosition Vector      `json:"stabilizedPalmPosition"`
	R                      Matrix      `json:"r"` //: nested arrays of floats (matrix)
	S                      float32     `json:"s"` //: float
	T                      Vector      `json:"t"` //: array of floats (vector)
	PinchStrength          float32     `json:"pinchStrength"`
	GrabStrength           float32     `json:"grabStrength"`
	Confidence             float32     `json:"confidence"`
	Type                   string      `json:"type"`
	//Valid                  bool        `json: "valid"`
}

type Gesture struct {
	ID            int           `json:"id"`
	State         string        `json:"state"`
	Type          string        `json:"type"`
	Duration      time.Duration `json:"duration"` //in microseconds
	HandIds       []int         `json:"handIds"`
	PointableIds  []int         `json:"pointableIds"`
	Speed         float32       `json:"speed"`
	Radius        float32       `json:"radius"`
	Progress      float32       `json:"progress"`
	Center        Vector        `json:"center"`
	Normal        Vector        `json:"normal"`
	StartPosition Vector        `json:"startPosition"`
	Position      Vector        `json:"position"`
	Direction     Vector        `json:"direction"`
}

type InteractionBox struct {
	//Valid bool `json: "valid"`
	//The center of the InteractionBox in device coordinates (millimeters). This point is equidistant from all sides of the box
	Center Vector `json:"center"`
	Size   Vector `json:"size"`
}

//homebrew
func (ib *InteractionBox) Norm(p Vector) Vector {
	vec := Vector{}
	vec.X = (p.X + ib.Center.X) / (ib.Size.X * 2)
	vec.Y = (p.Y + ib.Center.Y) / (ib.Size.Y * 2)
	vec.Z = (p.Z + ib.Center.Z) / (ib.Size.Z * 2)

	//vec := p.Add(ib.Center.Vector3)
	//vec := p.Add(ib.Center.Vector3)
	//vec.X /= 2 * ib.Size.X
	//vec.Y /= 2 * ib.Size.Y
	//vec.Z /= 2 * ib.Size.Z
	//return Vector{vec}
	return vec
}

//from the JS API
func (ib *InteractionBox) NormalizePoint(p Vector, clamp bool) Vector {
	vec := Vector{}
	vec.X = (p.X-ib.Center.X)/ib.Size.X + 0.5
	vec.Y = (p.Y-ib.Center.Y)/ib.Size.Y + 0.5
	vec.Z = (p.Z-ib.Center.Z)/ib.Size.Z + 0.5

	if clamp {
		vec.X = glm.Clampf(vec.X, 0, 1)
		vec.Y = glm.Clampf(vec.Y, 0, 1)
		vec.Z = glm.Clampf(vec.Z, 0, 1)
	}
	return vec
}

func (ib *InteractionBox) DenormalizePoint(p Vector) Vector {
	v := Vector{}
	v.X = (p.X-0.5)*ib.Size.X + ib.Center.X
	v.Y = (p.Y-0.5)*ib.Size.Y + ib.Center.Y
	v.Z = (p.Z-0.5)*ib.Size.Z + ib.Center.Z
	return v
}

type Frame struct {
	CurrentFrameRate float32        `json:"currentFrameRate"`
	Gestures         []Gesture      `json:"gestures"`
	Hands            []Hand         `json:"hands"`
	Pointables       []Pointable    `json:"pintables"`
	Fingers          []Pointable    `json:"fingers"`
	Tools            []Pointable    `json:"tools"`
	ID               int            `json:"id"`
	InteractionBox   InteractionBox `json:"interactionBox"`
	R                Matrix         `json:"r"`
	S                float32        `json:"s"`
	T                Vector         `json:"t"`
	// microseconds elapsed since the Leap started
	Timestamp time.Duration `json:"timestamp"`
	//Valid     bool          `json: "valid"`
}

/*
func (f *Frame) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, f)
	if err != nil {
		switch t := err.(type) {
		case *json.UnmarshalFieldError:
			log.Println(t)
		case *json.UnmarshalTypeError:
			log.Println(t)
		case *json.UnsupportedTypeError:
			log.Println(t)
		case *json.UnsupportedValueError:
			log.Println(t)
		case *json.SyntaxError:
			log.Println(t)
		case *json.InvalidUnmarshalError:
			log.Println(t)
		}
		return err
	}
	f.Timestamp *= time.Microsecond
	return nil
}
*/
type Vector struct {
	glm.Vector3
}

func (v *Vector) UnmarshalJSON(data []byte) error {
	vp := [3]float32{}
	err := json.Unmarshal(data, &vp)
	if err != nil {
		switch t := err.(type) {
		case *json.UnmarshalFieldError:
			log.Println(t)
		case *json.UnmarshalTypeError:
			log.Println(t)
		case *json.UnsupportedTypeError:
			log.Println(t)
		case *json.UnsupportedValueError:
			log.Println(t)
		case *json.SyntaxError:
			log.Println(t)
		case *json.InvalidUnmarshalError:
			log.Println(t)
		}
		return err
	}
	v.X, v.Y, v.Z = vp[0], vp[1], vp[2]
	return nil

}

type Matrix struct {
	glm.Matrix3
}

func (m *Matrix) UnmarshalJSON(data []byte) error {
	mp := [3][3]float32{}
	err := json.Unmarshal(data, &mp)
	if err != nil {
		switch t := err.(type) {
		case *json.UnmarshalFieldError:
			log.Println(t)
		case *json.UnmarshalTypeError:
			log.Println(t)
		case *json.UnsupportedTypeError:
			log.Println(t)
		case *json.UnsupportedValueError:
			log.Println(t)
		case *json.SyntaxError:
			log.Println(t)
		case *json.InvalidUnmarshalError:
			log.Println(t)
		}
		return err
	}
	m.M11, m.M12, m.M13 = mp[0][0], mp[0][1], mp[0][2]
	m.M21, m.M22, m.M23 = mp[1][0], mp[1][1], mp[1][2]
	m.M31, m.M32, m.M33 = mp[2][0], mp[2][1], mp[2][2]
	return nil

}

type Device struct {
	Conn        *websocket.Conn
	Frames, Bin chan *Frame
	quit        chan struct{}
}

func (dev *Device) Close() {
	close(dev.quit)
}

//enables or disables the gestures form the leapmotion
func (dev *Device) GestEnable(enable bool) error {
	msg := struct {
		EnableGestures bool `json:"enableGestures"`
	}{enable}
	err := dev.Conn.WriteJSON(msg)
	//err := dev.Conn.WriteMessage(websocket.TextMessage, []byte("{enableGestures: true}"))
	if err != nil {
		return err
	}
	return nil
}

//Dial try to connects to the websocket and if successful will start a goroutine
//to decode frames from the websocket and send them on the Device channel
func Dial(url string) (*Device, error) {
	//origin default to localhost
	//if origin == "" {
	//	origin = "http://localhost/"
	//}
	//ws, err := websocket.Dial(url, "", origin)
	//if err != nil {
	//	return nil, err
	//}
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		if err == websocket.ErrBadHandshake {
			return nil, fmt.Errorf("%v\n Response: %+v\n", err, resp)
		}
	}
	dev := &Device{
		Conn:   conn,
		Frames: make(chan *Frame, 2),
		Bin:    make(chan *Frame, 100),
		quit:   make(chan struct{})}
	go listener(dev)
	return dev, nil
}

//this is the decoder of the frames, here are also fixed some quirks
//as the timestamps being in microseconds
func listener(dev *Device) {
	//dec := json.NewDecoder(dev.Conn)
	var f *Frame
	for {
		select {
		case <-dev.quit:
			close(dev.Frames)
			dev.Conn.Close()
			return
		case f = <-dev.Bin:
		default: //new frame
			f = &Frame{}
		}
		//if err := dec.Decode(f); err == io.EOF { //shutdown
		//	dev.Close()
		//} else if err != nil {
		//	log.Panic("Error decoding from websocket ", err)
		//}

		err := dev.Conn.ReadJSON(f)
		//_, p, err := dev.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error receiving %+v\n", err)
		}
		//log.Println(string(p))

		//fix the timestamps
		f.Timestamp *= time.Microsecond
		for _, g := range f.Gestures {
			g.Duration *= time.Microsecond
		}
		//try to giva a frame to the clinet
		select {
		case dev.Frames <- f:
		default: //if not ready skip this frame
			log.Println("client not ready, skip frame")
			select {
			case dev.Bin <- f:
			default:
			}
		}
	}
}
