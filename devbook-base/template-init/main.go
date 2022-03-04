package main

import (
	"time"

	"github.com/DevbookHQ/template-init/pkg/environment"
	"github.com/DevbookHQ/template-init/pkg/fs"
	"github.com/DevbookHQ/template-init/pkg/log"
	"github.com/DevbookHQ/template-init/pkg/runner"
	msg "github.com/DevbookHQ/template-init/pkg/runner/message"
	"github.com/DevbookHQ/template-init/pkg/template"
)

var (
	rconn runner.Connection
	templ *template.Template

	createDir       = make(chan msg.Message)
	getFile         = make(chan msg.Message)
	writeFile       = make(chan msg.Message)
	removeFile      = make(chan msg.Message)
	listDir         = make(chan msg.Message)
	execCmd         = make(chan msg.Message)
	killCmd         = make(chan msg.Message)
	listRunningCmds = make(chan msg.Message)
)

func watchRunnerMessages() {
	// The main dispatch bus. We listen to new messages and based on their type
	// we perform the correct action.
	for {
		// We can safely type cast to a specific MessageData type because each channel
		// is subscribed to a single InMessageType topic which has corresponding MessageData type.
		select {
		case lc := <-listRunningCmds:
			log.Log(
				msg.WithMessageData(lc.Data),
			).Info("Received `ListRunningCmds`")
			// The `ListCmds` message has an empty payload so we don't need to type cast it
			// or create a variable for it like with other messages.
			cmds := templ.RunningCommands()
			msgPayload := msg.MessagePayload{
				Type: msg.OutMessageRunningCmds,
				Payload: msg.MessagePayloadRunningCmds{
					TookAt: time.Now().UnixMilli(),
					Cmds:   cmds,
				},
			}
			rconn.Send(msgPayload)
		case kc := <-killCmd:
			log.Log(
				msg.WithMessageData(kc.Data),
			).Info("Received `KillCmd`")
			d := kc.Data.Payload.(msg.MessagePayloadKillCmd)
			templ.KillCommand(d.ExecutionID)
		case ec := <-execCmd:
			log.Log(
				msg.WithMessageData(ec.Data),
			).Info("Received `ExecCmd`")
			d := ec.Data.Payload.(msg.MessagePayloadExecCmd)

			done := make(chan error)
			stdout := make(chan string)
			stderr := make(chan string)

			go func() {
				for {
					select {
					case outStr := <-stdout:
						msgPayload := msg.MessagePayload{
							Type: msg.OutMessageCmdOut,
							Payload: msg.MessagePayloadCmdOut{
								ExecutionID: d.ExecutionID,
								Stdout:      outStr,
								Stderr:      "",
							},
						}
						rconn.Send(msgPayload)
					case errStr := <-stderr:
						msgPayload := msg.MessagePayload{
							Type: msg.OutMessageCmdOut,
							Payload: msg.MessagePayloadCmdOut{
								ExecutionID: d.ExecutionID,
								Stdout:      "",
								Stderr:      errStr,
							},
						}
						rconn.Send(msgPayload)
					case err := <-done:
						var errStr string
						if err != nil {
							log.Log(
								log.WithFields(log.Fields{
									"cmd": d.Cmd,
								}),
								log.WithError(err),
							).Error("Command exited with error")
							errStr = err.Error()
						} else {
							log.Log(
								log.WithFields(log.Fields{
									"cmd": d.Cmd,
								}),
							).Info("Command exited ok")
						}

						msgPayload := msg.MessagePayload{
							Type: msg.OutMessageCmdExit,
							Payload: msg.MessagePayloadCmdExit{
								ExecutionID: d.ExecutionID,
								Err:         errStr,
							},
						}
						rconn.Send(msgPayload)
						return
					}
				}
			}()
			go templ.ExecCommand(d.ExecutionID, d.Cmd, stdout, stderr, done)
			// TODO: We should send ack message back to Runner.
		case cd := <-createDir:
			log.Log(
				msg.WithMessageData(cd.Data),
			).Info("Received `CreateDir`")

			d := cd.Data.Payload.(msg.MessagePayloadCreateDir)
			if err := templ.Filesystem.CreateDir(d.Path); err != nil {
				log.Log(
					log.WithField("path", d.Path),
					log.WithError(err),
				).Error("Failed to create dir")
			}
			// TODO: We should send ack message back to Runner.
		case gf := <-getFile:
			log.Log(
				msg.WithMessageData(gf.Data),
			).Info("Received `GetFile`")

			d := gf.Data.Payload.(msg.MessagePayloadGetFile)
			content, err := templ.Filesystem.GetFile(d.Path)
			if err != nil {
				log.Log(
					log.WithField("path", d.Path),
					log.WithError(err),
				).Error("Failed to get file")
				continue
			}

			msgPayload := msg.MessagePayload{
				Type: msg.OutMessageFileContent,
				Payload: msg.MessagePayloadFileContent{
					Path:    d.Path,
					Content: string(content),
				},
			}
			rconn.Send(msgPayload)
			// TODO: We should send ack message back to Runner.
		case wf := <-writeFile:
			log.Log(
				msg.WithMessageData(wf.Data),
			).Info("Received `WriteFile`")

			d := wf.Data.Payload.(msg.MessagePayloadWriteFile)
			if err := templ.Filesystem.WriteFile(d.Path, d.Content); err != nil {
				log.Log(
					log.WithField("path", d.Path),
					log.WithError(err),
				).Error("Failed to write file")
			}
			// TODO: We should send ack message back to Runner.
		case rf := <-removeFile:
			log.Log(
				msg.WithMessageData(rf.Data),
			).Info("Received `RemoveFile`")

			d := rf.Data.Payload.(msg.MessagePayloadRemoveFile)
			if err := templ.Filesystem.RemoveFileOrDir(d.Path); err != nil {
				log.Log(
					log.WithField("path", d.Path),
					log.WithError(err),
				).Error("Failed to remove file")
			}
			// TODO: We should send ack message back to Runner.
		case ld := <-listDir:
			log.Log(
				msg.WithMessageData(ld.Data),
			).Info("Received `ListDir`")

			d := ld.Data.Payload.(msg.MessagePayloadListDir)

			content, err := templ.Filesystem.ListDir(d.Path)
			if err != nil {
				log.Log(
					log.WithField("path", d.Path),
					log.WithError(err),
				).Error("Failed to list dir")
				continue
			}

			msgPayload := msg.MessagePayload{
				Type: msg.OutMessageDirContent,
				Payload: msg.MessagePayloadDirContent{
					DirPath: d.Path,
					Content: content,
				},
			}
			rconn.Send(msgPayload)
			// TODO: We should send ack message back to Runner.
			// TODO: We should send ack message back to Runner.
		}
	}
}

func sendFSEventToRunner(ev fs.Event) {
	var msgPayload msg.MessagePayload
	switch e := ev.(type) {
	case *fs.CreateEvent:
		msgPayload = msg.MessagePayload{
			Type: msg.OutMessageFSEventCreate,
			Payload: msg.MessagePayloadFSEventCreate{
				Path: e.Path(),
				Type: e.Type(),
			},
		}
	case *fs.WriteEvent:
		msgPayload = msg.MessagePayload{
			Type: msg.OutMessageFSEventWrite,
			Payload: msg.MessagePayloadFSEventWrite{
				Path: e.Path(),
			},
		}
	case *fs.RemoveEvent:
		msgPayload = msg.MessagePayload{
			Type: msg.OutMessageFSEventRemove,
			Payload: msg.MessagePayloadFSEventRemove{
				Path: e.Path(),
				Type: e.Type(),
			},
		}
	default:
		log.Log(
			fs.WithEvent(ev),
		).Error("Cannot send even to Runner - received unexpected filesystem event")
		return
	}
	rconn.Send(msgPayload)
}

func watchTemplateFSEvents() {
	for {
		select {
		case ev := <-templ.Filesystem.WatcherEvents:
			log.Log(
				fs.WithEvent(ev),
			).Info("Filesystem event")
			sendFSEventToRunner(ev)
		case err := <-templ.Filesystem.WatcherErrors:
			log.Log(
				log.WithError(err),
			).Error("Filesystem error")
		}
	}
}

func main() {
	log.Info("tinit main")

	// Get the current environment.
	log.Info("Parsing environment")
	env, err := environment.New()
	if err != nil {
		log.Log(
			log.WithError(err),
		).Fatal("Failed to initialize environment")
	}
	log.Log(
		environment.WithEnv(env),
	).Info("Environment parsed")

	// Create template from the environment.
	log.Log(
		environment.WithEnv(env),
	).Info("Deriving template from environment")
	t, err := template.FromEnvironment(env)
	if err != nil {
		log.Log(
			log.WithError(err),
			environment.WithEnv(env),
		).Fatal("Failed to create template based on environment")
	}
	log.Log(
		environment.WithEnv(env),
		template.WithTemplate(t),
	).Info("Derived template from environment")
	templ = t
	defer templ.Stop()

	// Prepare a communication channel to Runner via provided unix socket.
	rconn = runner.NewConnection(env.RUNNER_SOCKET_PATH)
	defer rconn.Close()

	// Subscribe to messages from runner.
	rconn.Subscribe(msg.InMessageCreateDir, createDir)
	rconn.Subscribe(msg.InMessageWriteFile, writeFile)
	rconn.Subscribe(msg.InMessageRemoveFile, removeFile)
	rconn.Subscribe(msg.InMessageListDir, listDir)
	rconn.Subscribe(msg.InMessageGetFile, getFile)
	rconn.Subscribe(msg.InMessageExecCmd, execCmd)
	rconn.Subscribe(msg.InMessageKillCmd, killCmd)
	rconn.Subscribe(msg.InMessageListRunningCmds, listRunningCmds)
	log.Info("Ready to receive messages from Runner")

	log.Log(
		template.WithTemplate(templ),
	).Info("Starting template")
	templ.Start()
	log.Log(
		template.WithTemplate(templ),
	).Info("Started template")

	// Connect to Runner.
	log.Info("Connecting to the IPC socket server")
	go rconn.DialAndRead()

	for {
		select {
		case err := <-rconn.Err:
			f := log.Fields{}
			switch e := err.(type) {
			case *runner.RunnerMessageParseError:
				f["errType"] = "RunnerMessageParseError:"
				f["offset"] = e.Offset
				f["errByte"] = e.ErrByte
				f["errString"] = string(e.ErrByte)
			case *runner.ReadError:
				f["errType"] = "ReadError"
				f["socketPath"] = e.SocketPath
				f["errByte"] = e.ErrByte
			case *runner.ReadEOF:
				f["errType"] = "ReadEOF"
				f["socketPath"] = e.SocketPath
			case *runner.WriteError:
				f["errType"] = "WriteError"
				f["socketPath"] = e.SocketPath
				f["offset"] = e.Offset
				f["errByte"] = e.ErrByte
			case *runner.ConnError:
				f["errType"] = "ConnError"
				f["socketPath"] = e.SocketPath
			default:
				f["errType"] = "builtin"
			}

			log.Log(
				log.WithError(err),
				log.WithFields(f),
			).Error("Runner connection error")
		case <-rconn.Ready:
			log.Info("Connected to the IPC socket server")
			go watchTemplateFSEvents()
			go watchRunnerMessages()
			// Announce Runner that the environment is ready.
			log.Info("Sending 'Ready' to Runner")
			rconn.SendReady()
			log.Info("Sent 'Ready' to Runner")
			//case err := <-templ.Done:
			//	if err != nil {
			//		log.Log(
			//			log.WithError(err),
			//			template.WithTemplate(templ),
			//		).Error("Template finished with error")
			//	} else {
			//		log.Log(
			//			template.WithTemplate(templ),
			//		).Info("Template finished without error")
			//	}
		}
	}

	// Block the main goroutine forever. We don't want to exit even when the template has problems.
	// This is useful for example when a user broke a template manually inside the web app.
	// We want to let user fix the error they made and keep working.
	// TODO: We don't have any system that tries to restart the template when user actually fixes the error.
	c := make(chan struct{})
	<-c
}
