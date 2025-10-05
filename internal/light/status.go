package light

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const (
	readBufferSize = 1024 // in bytes
)

type StatusResponse struct {
	Method *string `json:"method"`
	Env    *string `json:"env,omitempty"`
	Result *struct {
		Mac     *string `json:"mac,omitempty"`
		RSSI    *int    `json:"rssi,omitempty"`
		State   *bool   `json:"state,omitempty"`
		SceneID *int    `json:"sceneId,omitempty"`
		Temp    *int    `json:"temp,omitempty"`
		Dimming *int    `json:"dimming,omitempty"`
	} `json:"result,omitempty"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func CheckStatus(conn net.PacketConn, lightIP net.Addr, timeout time.Duration) (*StatusResponse, error) {
	message := `{"method": "getPilot"}`
	_, err := conn.WriteTo([]byte(message), lightIP)
	if err != nil {
		return nil, err
	}
	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, readBufferSize)
	n, addr, err := conn.ReadFrom(buffer)
	if err != nil {
		return nil, err
	}
	if addr.String() != lightIP.String() {
		return nil, fmt.Errorf("response from unexpected address: %s", addr.String())
	}
	var response StatusResponse
	err = json.Unmarshal(buffer[:n], &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
