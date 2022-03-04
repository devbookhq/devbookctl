package message

import (
	"encoding/json"
	"fmt"

	fsitem "github.com/DevbookHQ/template-init/pkg/fs/item"
	"github.com/DevbookHQ/template-init/pkg/template"
)

type MessageType string
type IPCEvent string

const (
	// In Messages
	// Incoming messages from the unix socket.
	InMessageCreateDir  MessageType = "CreateDirectory"
	InMessageWriteFile  MessageType = "WriteFile"
	InMessageGetFile    MessageType = "GetFile"
	InMessageRemoveFile MessageType = "RemoveFile"
	InMessageListDir    MessageType = "ListDir"

	InMessageExecCmd         MessageType = "ExecCmd"
	InMessageKillCmd         MessageType = "KillCmd"
	InMessageListRunningCmds MessageType = "ListRunningCmds"
	//////////////////////////////////////////

	// Out Messages
	// Messages we send to the unix socket.
	OutMessageStatus MessageType = "Status"

	OutMessageFSEventCreate MessageType = "FSEventCreate"
	OutMessageFSEventWrite  MessageType = "FSEventWrite"
	OutMessageFSEventRemove MessageType = "FSEventRemove"

	OutMessageDirContent  MessageType = "DirContent"
	OutMessageFileContent MessageType = "FileContent"

	OutMessageCmdOut      MessageType = "CmdOut"
	OutMessageCmdExit     MessageType = "CmdExit"
	OutMessageRunningCmds MessageType = "RunningCmds"
	//////////////////////////////////////////
)

const (
	IPCEventMessage IPCEvent = "message"
)

type rawMessage struct {
	Type IPCEvent `json:"type"`
	// Type MessageType     `json:"type"`
	Data rawMessagePayload `json:"data"`
}

type rawMessagePayload struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Message struct {
	Type IPCEvent `json:"type"`
	// Type MessageType `json:"type"`
	Data MessagePayload `json:"data"`
}

// {
// 	type: "message",
// 	data: {
// 		type: "",
// 		payload: {
// 			command
// 		}
// 	}
// }

type MessagePayload struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

type MessagePayloadStatus struct {
	Status string `json:"status"`
}

type MessagePayloadFSEventCreate struct {
	Path string      `json:"path"`
	Type fsitem.Type `json:"type"`
}

type MessagePayloadFSEventWrite struct {
	Path string `json:"path"`
}

type MessagePayloadFSEventRemove struct {
	Path string      `json:"path"`
	Type fsitem.Type `json:"type"`
}

type MessagePayloadCreateDir struct {
	Path string `json:"path"`
}

type MessagePayloadGetFile struct {
	Path string `json:"path"`
}

type MessagePayloadWriteFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type MessagePayloadRemoveFile struct {
	Path string `json:"path"`
}

type MessagePayloadListDir struct {
	Path string `json:"path"`
}

type MessagePayloadDirContent struct {
	DirPath string        `json:"dirPath"`
	Content []fsitem.Item `json:"content"`
}

type MessagePayloadFileContent struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type MessagePayloadExecCmd struct {
	Cmd         string `json:"command"`
	ExecutionID string `json:"executionID"`
}

type MessagePayloadKillCmd struct {
	ExecutionID string `json:"executionID"`
}

type MessagePayloadListRunningCmds struct{}

type MessagePayloadCmdOut struct {
	ExecutionID string `json:"executionID"`
	Stdout      string `json:"stdout"`
	Stderr      string `json:"stderr"`
}

type MessagePayloadCmdExit struct {
	ExecutionID string `json:"executionID"`
	Err         string `json:"error"`
}

type MessagePayloadRunningCmds struct {
	TookAt int64               `json:"tookAt"`
	Cmds   []*template.Command `json:"commands"`
}

func (m *Message) UnmarshalJSON(bs []byte) error {
	// Check if marshalled struct's `Type` is one of the consts.
	// Return error if not.
	var msg rawMessage
	if err := json.Unmarshal(bs, &msg); err != nil {
		return err
	}

	// Here we unmarshall the message's payload field based on the message type.
	switch msg.Data.Type {
	case InMessageCreateDir:
		var d MessagePayloadCreateDir
		if err := json.Unmarshal(msg.Data.Payload, &d); err != nil {
			return fmt.Errorf("Failed to unmarshal `InMessageCreateDir`: %s", err)
		}
		m.Data.Payload = d
	case InMessageGetFile:
		var d MessagePayloadGetFile
		if err := json.Unmarshal(msg.Data.Payload, &d); err != nil {
			return fmt.Errorf("Failed to unmarshal `InMessageGetFile`: %s", err)
		}
		m.Data.Payload = d
	case InMessageWriteFile:
		var d MessagePayloadWriteFile
		if err := json.Unmarshal(msg.Data.Payload, &d); err != nil {
			return fmt.Errorf("Failed to unmarshal `InMessageWriteFile`: %s", err)
		}
		m.Data.Payload = d
	case InMessageRemoveFile:
		var d MessagePayloadRemoveFile
		if err := json.Unmarshal(msg.Data.Payload, &d); err != nil {
			return fmt.Errorf("Failed to unmarshal `InMessageRemoveFile`: %s", err)
		}
		m.Data.Payload = d
	case InMessageListDir:
		var d MessagePayloadListDir
		if err := json.Unmarshal(msg.Data.Payload, &d); err != nil {
			return fmt.Errorf("Failed to unmarshal `InMessageListDir`: %s", err)
		}
		m.Data.Payload = d
	case InMessageExecCmd:
		var d MessagePayloadExecCmd
		if err := json.Unmarshal(msg.Data.Payload, &d); err != nil {
			return fmt.Errorf("Failed to unmarshal `InMessageExecCmd`: %s", err)
		}
		m.Data.Payload = d
	case InMessageKillCmd:
		var d MessagePayloadKillCmd
		if err := json.Unmarshal(msg.Data.Payload, &d); err != nil {
			return fmt.Errorf("Failed to unmarshal `InMessageKillCmd`: %s", err)
		}
		m.Data.Payload = d
	case InMessageListRunningCmds:
		var d MessagePayloadListRunningCmds
		if err := json.Unmarshal(msg.Data.Payload, &d); err != nil {
			return fmt.Errorf("Failed to unmarshal `InMessageListRunningCmds`: %s", err)
		}
		m.Data.Payload = d
	default:
		return fmt.Errorf(
			"Unexpected `Message.Data.Type`. value=%s, msg=%s",
			msg.Data.Type,
			string(bs),
		)
	}
	m.Type = msg.Type
	m.Data.Type = msg.Data.Type

	return nil
}
