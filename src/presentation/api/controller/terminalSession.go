package apiController

import (
	"encoding/json"
	"log/slog"

	_ "github.com/goinfinite/os/src/domain/dto"
	_ "github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type TerminalSessionController struct {
	terminalSessionLiaison *liaison.TerminalSessionLiaison
	inputReader            tkPresentation.ApiRequestInputReader
}

func NewTerminalSessionController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *TerminalSessionController {
	return &TerminalSessionController{
		terminalSessionLiaison: liaison.NewTerminalSessionLiaison(
			persistentDbSvc, trailDbSvc,
		),
		inputReader: tkPresentation.ApiRequestInputReader{},
	}
}

type terminalSessionControlMessage struct {
	Type string `json:"type"`
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
}

type terminalSessionServerMessage struct {
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
}

const (
	terminalSessionMessageReady = "ready"
	terminalSessionMessageExit  = "exit"
)

var terminalSessionWsUpgrader = websocket.Upgrader{}

// ReadTerminalSessions	 godoc
// @Summary      ReadTerminalSessions
// @Description  List terminal sessions visible to the operator. Normal accounts see their own sessions. Super-admins see every account and may filter with accountId.
// @Tags         terminal-session
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id query  string  false  "TerminalSessionId"
// @Param        accountId query  uint  false  "AccountId (super-admin only)"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Success      200 {object} dto.ReadTerminalSessionsResponse
// @Router       /v1/terminal-sessions/ [get]
func (controller *TerminalSessionController) Read(echoContext echo.Context) error {
	requestData, requestParsingErr := controller.inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}

	return tkPresentation.LiaisonApiResponseEmitter(
		echoContext, controller.terminalSessionLiaison.Read(requestData),
	)
}

// CreateTerminalSession	 godoc
// @Summary      CreateTerminalSession
// @Description  Create a terminal session. The working directory must exist. The account cap is ten sessions.
// @Tags         terminal-session
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createTerminalSessionDto  body  dto.CreateTerminalSession  true  "command is optional."
// @Success      201 {object} entity.TerminalSession
// @Router       /v1/terminal-sessions/ [post]
func (controller *TerminalSessionController) Create(echoContext echo.Context) error {
	requestData, requestParsingErr := controller.inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}

	return tkPresentation.LiaisonApiResponseEmitter(
		echoContext, controller.terminalSessionLiaison.Create(requestData),
	)
}

// DeleteTerminalSession	 godoc
// @Summary      DeleteTerminalSession
// @Description  Kill a terminal session.
// @Tags         terminal-session
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id  path  string  true  "TerminalSessionId to delete."
// @Success      200 {object} object{} "TerminalSessionDeleted"
// @Router       /v1/terminal-sessions/{id}/ [delete]
func (controller *TerminalSessionController) Delete(echoContext echo.Context) error {
	requestData, requestParsingErr := controller.inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}

	return tkPresentation.LiaisonApiResponseEmitter(
		echoContext, controller.terminalSessionLiaison.Delete(requestData),
	)
}

func (controller *TerminalSessionController) applyControlMessage(
	attachHandle repository.TerminalSessionAttachHandle,
	payload []byte,
) {
	controlMessage := terminalSessionControlMessage{}
	err := json.Unmarshal(payload, &controlMessage)
	if err != nil {
		slog.Debug("ParseTerminalSessionControlMessageError", slog.String("err", err.Error()))
		return
	}

	if controlMessage.Type != "resize" {
		return
	}

	err = attachHandle.Resize(controlMessage.Cols, controlMessage.Rows)
	if err != nil {
		slog.Debug("ResizeTerminalSessionError", slog.String("err", err.Error()))
	}
}

func (controller *TerminalSessionController) pumpClientMessages(
	wsConn *websocket.Conn,
	attachHandle repository.TerminalSessionAttachHandle,
) {
	for {
		messageType, payload, err := wsConn.ReadMessage()
		if err != nil {
			_ = attachHandle.Close()
			return
		}

		switch messageType {
		case websocket.BinaryMessage:
			_, err = attachHandle.Write(payload)
			if err != nil {
				_ = attachHandle.Close()
				return
			}
		case websocket.TextMessage:
			controller.applyControlMessage(attachHandle, payload)
		}
	}
}

func (controller *TerminalSessionController) sendServerMessage(
	wsConn *websocket.Conn, messageType, message string,
) {
	serverMessage := terminalSessionServerMessage{
		Type: messageType, Message: message,
	}
	payload, err := json.Marshal(serverMessage)
	if err != nil {
		slog.Debug("MarshalTerminalSessionServerMessageError", slog.String("err", err.Error()))
		return
	}

	_ = wsConn.WriteMessage(websocket.TextMessage, payload)
}

func (controller *TerminalSessionController) pumpTerminalOutput(
	wsConn *websocket.Conn,
	attachHandle repository.TerminalSessionAttachHandle,
) {
	outputBuffer := make([]byte, 4096)
	for {
		bytesRead, readErr := attachHandle.Read(outputBuffer)
		if bytesRead > 0 {
			writeErr := wsConn.WriteMessage(
				websocket.BinaryMessage, outputBuffer[:bytesRead],
			)
			if writeErr != nil {
				return
			}
		}

		if readErr != nil {
			controller.sendServerMessage(wsConn, terminalSessionMessageExit, "")
			return
		}
	}
}

// AttachTerminalSession	 godoc
// @Summary      AttachTerminalSession
// @Description  Upgrade to WebSocket and attach to a terminal session. Binary frames carry terminal bytes in both directions. Text frames carry JSON control messages. The client sends {"type":"resize","cols":N,"rows":N}. The server sends {"type":"ready"}, {"type":"exit"} or {"type":"error"}.
// @Tags         terminal-session
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id  path  string  true  "TerminalSessionId to attach."
// @Success      101 {string} string "SwitchingProtocols"
// @Router       /v1/terminal-sessions/{id}/attach/ [get]
func (controller *TerminalSessionController) Attach(echoContext echo.Context) error {
	requestData, requestParsingErr := controller.inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}

	attachHandle, liaisonResponse := controller.terminalSessionLiaison.Attach(
		requestData,
	)
	if liaisonResponse.Status != tkPresentation.LiaisonResponseStatusSuccess {
		return tkPresentation.LiaisonApiResponseEmitter(echoContext, liaisonResponse)
	}

	// The nil CheckOrigin applies gorilla's checkSameOrigin, so a handshake
	// whose Origin host differs from the request Host is rejected.
	// nosemgrep: go.gorilla.security.audit.websocket-missing-origin-check.websocket-missing-origin-check
	wsConn, err := terminalSessionWsUpgrader.Upgrade(
		echoContext.Response(), echoContext.Request(), nil,
	)
	if err != nil {
		slog.Error("UpgradeTerminalSessionWsFailed", slog.String("err", err.Error()))
		_ = attachHandle.Close()
		return nil
	}
	defer func() { _ = wsConn.Close() }()
	defer func() { _ = attachHandle.Close() }()

	controller.sendServerMessage(wsConn, terminalSessionMessageReady, "")

	go controller.pumpClientMessages(wsConn, attachHandle)

	controller.pumpTerminalOutput(wsConn, attachHandle)

	return nil
}
