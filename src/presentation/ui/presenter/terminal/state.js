UiToolset.RegisterAlpineState(() => {
  class SessionTerminalManager {
    constructor(onConnectionStatusChange, onSessionExit) {
      this.onConnectionStatusChange = onConnectionStatusChange;
      this.onSessionExit = onSessionExit;
      this.sessionTerminals = new Map();
      this.mountRetryLimit = 25;
      this.mountRetryDelayMs = 200;
      this.reconnectDelayMs = 2000;
      this.terminalFontName = "JetBrains Mono";
      this.terminalFontFamily = `"${this.terminalFontName}", monospace`;
      this.terminalFontSize = 15;
      this.fontWaitTimeoutMs = 3000;
    }

    resize(sessionId) {
      const sessionTerminal = this.sessionTerminals.get(sessionId);
      if (!sessionTerminal?.fitAddon) {
        return;
      }

      const terminalContainer = sessionTerminal.term.element?.parentElement;
      if (
        !terminalContainer ||
        terminalContainer.offsetWidth === 0 ||
        terminalContainer.offsetHeight === 0
      ) {
        return;
      }

      sessionTerminal.fitAddon.fit();

      if (sessionTerminal.socket?.readyState !== WebSocket.OPEN) {
        return;
      }
      if (
        sessionTerminal.lastSentCols === sessionTerminal.term.cols &&
        sessionTerminal.lastSentRows === sessionTerminal.term.rows
      ) {
        return;
      }

      sessionTerminal.lastSentCols = sessionTerminal.term.cols;
      sessionTerminal.lastSentRows = sessionTerminal.term.rows;
      sessionTerminal.socket.send(
        JSON.stringify({
          type: "resize",
          cols: sessionTerminal.term.cols,
          rows: sessionTerminal.term.rows,
        }),
      );
    }

    connect(sessionId) {
      const sessionTerminal = this.sessionTerminals.get(sessionId);
      if (!sessionTerminal) {
        return;
      }

      sessionTerminal.lastSentCols = 0;
      sessionTerminal.lastSentRows = 0;

      const attachUrl = new URL(
        `api/v1/terminal-sessions/${sessionId}/attach/`,
        document.baseURI,
      );
      attachUrl.protocol = attachUrl.protocol === "https:" ? "wss:" : "ws:";

      const socket = new WebSocket(attachUrl.toString());
      socket.binaryType = "arraybuffer";
      sessionTerminal.socket = socket;

      socket.onopen = () => {
        this.onConnectionStatusChange(sessionId, true);
        this.resize(sessionId);
      };
      socket.onmessage = (event) => {
        if (typeof event.data !== "string") {
          sessionTerminal.term.write(new Uint8Array(event.data));
          return;
        }

        const serverMessage = JSON.parse(event.data);
        if (serverMessage.type === "exit") {
          sessionTerminal.isExited = true;
          this.onSessionExit(sessionId);
        }
      };
      socket.onclose = () => {
        this.onConnectionStatusChange(sessionId, false);
        if (sessionTerminal.isExited) {
          return;
        }
        sessionTerminal.reconnectTimer = setTimeout(
          () => this.connect(sessionId),
          this.reconnectDelayMs,
        );
      };
      socket.onerror = () => socket.close();
    }

    async waitForTerminalFont() {
      const fontStylesheet = document.getElementById(
        "terminal-font-stylesheet",
      );
      const stylesheetReady = new Promise((resolve) => {
        if (!fontStylesheet || fontStylesheet.sheet) {
          resolve();
          return;
        }
        fontStylesheet.addEventListener("load", resolve, { once: true });
        fontStylesheet.addEventListener("error", resolve, { once: true });
      });

      const fontReady = stylesheetReady.then(() =>
        document.fonts.load(
          `${this.terminalFontSize}px "${this.terminalFontName}"`,
        ),
      );

      let fontWaitTimer;
      const fontWaitTimeout = new Promise((resolve) => {
        fontWaitTimer = setTimeout(resolve, this.fontWaitTimeoutMs);
      });

      try {
        await Promise.race([fontReady, fontWaitTimeout]);
      } catch (error) {
        console.error(`LoadTerminalFontError: ${error}`);
      } finally {
        clearTimeout(fontWaitTimer);
      }
    }

    async mount(sessionId) {
      const sessionTerminal = this.sessionTerminals.get(sessionId);
      if (!sessionTerminal || sessionTerminal.term) {
        return;
      }

      const terminalContainer = document.getElementById(
        `terminal-${sessionId}`,
      );
      if (!terminalContainer) {
        return;
      }

      if (typeof Terminal === "undefined" || typeof FitAddon === "undefined") {
        sessionTerminal.mountAttempts += 1;
        if (sessionTerminal.mountAttempts > this.mountRetryLimit) {
          Alpine.store("toast").displayToast(
            "TerminalRendererUnavailable",
            "danger",
          );
          return;
        }
        sessionTerminal.mountRetryTimer = setTimeout(
          () => this.mount(sessionId),
          this.mountRetryDelayMs,
        );
        return;
      }

      await this.waitForTerminalFont();

      if (
        this.sessionTerminals.get(sessionId) !== sessionTerminal ||
        sessionTerminal.term
      ) {
        return;
      }

      const term = new Terminal({
        cursorBlink: true,
        fontFamily: this.terminalFontFamily,
        fontSize: this.terminalFontSize,
        lineHeight: 1,
        theme: { background: "#041118" },
      });
      const fitAddon = new FitAddon.FitAddon();
      term.loadAddon(fitAddon);
      term.open(terminalContainer);
      fitAddon.fit();

      term.onData((terminalData) => {
        if (!sessionTerminal.socket) {
          return;
        }
        if (sessionTerminal.socket.readyState !== WebSocket.OPEN) {
          return;
        }
        sessionTerminal.socket.send(new TextEncoder().encode(terminalData));
      });

      sessionTerminal.term = term;
      sessionTerminal.fitAddon = fitAddon;
      sessionTerminal.resizeObserver = new ResizeObserver(() =>
        this.resize(sessionId),
      );
      sessionTerminal.resizeObserver.observe(terminalContainer);
      this.connect(sessionId);
    }

    open(sessionId) {
      if (this.sessionTerminals.has(sessionId)) {
        return;
      }

      this.sessionTerminals.set(sessionId, {
        term: null,
        fitAddon: null,
        socket: null,
        reconnectTimer: null,
        mountRetryTimer: null,
        mountAttempts: 0,
        resizeObserver: null,
        lastSentCols: 0,
        lastSentRows: 0,
        isExited: false,
      });

      void this.mount(sessionId);
    }

    close(sessionId) {
      const sessionTerminal = this.sessionTerminals.get(sessionId);
      if (!sessionTerminal) {
        return;
      }

      if (sessionTerminal.mountRetryTimer) {
        clearTimeout(sessionTerminal.mountRetryTimer);
      }
      if (sessionTerminal.reconnectTimer) {
        clearTimeout(sessionTerminal.reconnectTimer);
      }
      if (sessionTerminal.resizeObserver) {
        sessionTerminal.resizeObserver.disconnect();
      }
      if (sessionTerminal.socket) {
        sessionTerminal.socket.onclose = null;
        sessionTerminal.socket.close();
      }
      if (sessionTerminal.term) {
        sessionTerminal.term.dispose();
      }

      this.sessionTerminals.delete(sessionId);
    }

    disposeAll() {
      for (const sessionId of this.sessionTerminals.keys()) {
        this.close(sessionId);
      }
    }
  }

  Alpine.data("terminal", () => ({
    // PrimaryState
    sessions: [],
    openTabs: [],
    activeTabId: "",
    renamingSessionId: "",
    renameInputValue: "",
    isCreateSessionFormExpanded: false,
    isCreateSessionLoading: false,
    createSessionForm: {
      name: "",
      runAsUsername: "nobody",
      workingDir: "/app",
      command: "",
    },
    accountHomeDir: "",

    // DerivedState
    resolveSessionLabel(session) {
      if (session.name) {
        return session.name;
      }
      return `${session.accountUsername}@${session.id.slice(0, 4)}`;
    },

    isSessionAttached(sessionId) {
      return this.openTabs.some((tab) => tab.sessionId === sessionId);
    },

    isTabConnected(sessionId) {
      const openTab = this.openTabs.find((tab) => tab.sessionId === sessionId);
      return openTab?.isConnected ?? false;
    },

    resetCreateSessionForm() {
      this.createSessionForm.name = "";
      this.createSessionForm.runAsUsername = "nobody";
      this.createSessionForm.workingDir = "/app";
      this.createSessionForm.command = "";
    },

    applyRunAsUsernameDefault() {
      if (this.createSessionForm.runAsUsername === "nobody") {
        this.createSessionForm.workingDir = "/app";
        return;
      }
      this.createSessionForm.workingDir = this.accountHomeDir;
    },

    expandCreateSessionForm() {
      this.isCreateSessionFormExpanded = true;
    },

    collapseCreateSessionForm() {
      this.isCreateSessionFormExpanded = false;
      this.resetCreateSessionForm();
    },

    selectTab(sessionId) {
      this.activeTabId = sessionId;
      this.$nextTick(() => this.sessionTerminalManager.resize(sessionId));
    },

    closeTab(sessionId) {
      const tabIndex = this.openTabs.findIndex(
        (tab) => tab.sessionId === sessionId,
      );
      if (tabIndex === -1) {
        return;
      }

      this.openTabs.splice(tabIndex, 1);
      if (this.activeTabId === sessionId) {
        this.activeTabId =
          this.openTabs.length > 0 ? this.openTabs[0].sessionId : "";
      }
      this.sessionTerminalManager.close(sessionId);
    },

    openTab(session) {
      if (this.isSessionAttached(session.id)) {
        this.selectTab(session.id);
        return;
      }

      this.openTabs.push({ sessionId: session.id, isConnected: false });
      this.activeTabId = session.id;

      this.$nextTick(() => this.sessionTerminalManager.open(session.id));
    },

    async loadSessions() {
      try {
        const response = await fetch(
          `${Infinite.OsApiBasePath}/v1/terminal-sessions/`,
          {
            method: "GET",
            headers: { Accept: "application/json" },
          },
        );
        if (!response.ok) {
          throw new Error(`BadHttpResponseCode: ${response.status}`);
        }

        const jsonResponse = await response.json();
        this.sessions = jsonResponse.body.terminalSessions;

        for (const tab of [...this.openTabs]) {
          const sessionStillExists = this.sessions.some(
            (session) => session.id === tab.sessionId,
          );
          if (!sessionStillExists) {
            this.closeTab(tab.sessionId);
          }
        }
      } catch (error) {
        console.error(`ReadTerminalSessionsError: ${error}`);
        Alpine.store("toast").displayToast(
          "ReadTerminalSessionsError",
          "danger",
        );
      }
    },

    isInsideSessionsModal() {
      return this.$el.closest("#terminal-sessions-modal-body") !== null;
    },

    async createSession() {
      if (this.isCreateSessionLoading) {
        return;
      }
      this.isCreateSessionLoading = true;
      try {
        const requestBody = {
          workingDir: this.createSessionForm.workingDir,
          runAsUsername: this.createSessionForm.runAsUsername,
        };
        if (this.createSessionForm.name.length > 0) {
          requestBody.name = this.createSessionForm.name;
        }
        if (this.createSessionForm.command.length > 0) {
          requestBody.command = this.createSessionForm.command;
        }

        const response = await fetch(
          `${Infinite.OsApiBasePath}/v1/terminal-sessions/`,
          {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(requestBody),
          },
        );
        const jsonResponse = await response.json().catch(() => ({}));
        if (!response.ok) {
          throw new Error(jsonResponse.body || "CreateTerminalSessionFailed");
        }

        const createdSession = jsonResponse.body;
        this.collapseCreateSessionForm();
        await this.loadSessions();
        window.dispatchEvent(new Event("update:terminal-session"));
        if (this.isInsideSessionsModal()) {
          this.openTab(createdSession);
          return;
        }
        this.$store.main.openTerminalSessionsModal(createdSession.id);
      } catch (error) {
        Alpine.store("toast").displayToast(error.message, "danger");
      } finally {
        this.isCreateSessionLoading = false;
      }
    },

    async createSessionWithDefaults() {
      this.resetCreateSessionForm();
      await this.createSession();
    },

    async killSession(sessionId) {
      try {
        const response = await fetch(
          `${Infinite.OsApiBasePath}/v1/terminal-sessions/${sessionId}/`,
          { method: "DELETE" },
        );
        if (!response.ok) {
          throw new Error("DeleteTerminalSessionFailed");
        }

        this.closeTab(sessionId);
        await this.loadSessions();
        window.dispatchEvent(new Event("update:terminal-session"));
      } catch (error) {
        Alpine.store("toast").displayToast(error.message, "danger");
      }
    },

    startRename(sessionId, currentName) {
      this.renamingSessionId = sessionId;
      this.renameInputValue = currentName;
    },

    cancelRename() {
      this.renamingSessionId = "";
      this.renameInputValue = "";
    },

    async submitRename(sessionId) {
      try {
        const response = await fetch(
          `${Infinite.OsApiBasePath}/v1/terminal-sessions/${sessionId}/`,
          {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              id: sessionId,
              name: this.renameInputValue,
            }),
          },
        );
        const jsonResponse = await response.json().catch(() => ({}));
        if (!response.ok) {
          throw new Error(jsonResponse.body || "UpdateTerminalSessionFailed");
        }

        this.cancelRename();
        await this.loadSessions();
        window.dispatchEvent(new Event("update:terminal-session"));
      } catch (error) {
        Alpine.store("toast").displayToast(error.message, "danger");
      }
    },

    openTerminalSession(sessionId) {
      this.$store.main.openTerminalSessionsModal(sessionId);
    },

    async openPendingSession() {
      await this.loadSessions();

      if (!this.isInsideSessionsModal()) {
        return;
      }

      const pendingSessionId = this.$store.main.pendingTerminalSessionId;
      if (!pendingSessionId) {
        return;
      }
      this.$store.main.pendingTerminalSessionId = "";

      const session = this.sessions.find(
        (candidateSession) => candidateSession.id === pendingSessionId,
      );
      if (!session) {
        return;
      }
      this.openTab(session);
    },

    sessionTerminalManager: null,
    focusPendingSessionHandler: null,
    refreshSessionsHandler: null,

    init() {
      this.accountHomeDir = this.$el.dataset.accountHomeDir ?? "";
      this.sessionTerminalManager = new SessionTerminalManager(
        (sessionId, isConnected) => {
          const openTab = this.openTabs.find(
            (tab) => tab.sessionId === sessionId,
          );
          if (!openTab) {
            return;
          }
          openTab.isConnected = isConnected;
        },
        (sessionId) => {
          this.closeTab(sessionId);
          void this.loadSessions();
          window.dispatchEvent(new Event("update:terminal-session"));
        },
      );

      this.focusPendingSessionHandler = () => this.openPendingSession();
      this.refreshSessionsHandler = () => this.loadSessions();
      window.addEventListener(
        "focus:terminal-session",
        this.focusPendingSessionHandler,
      );
      window.addEventListener(
        "update:terminal-session",
        this.refreshSessionsHandler,
      );
      void this.openPendingSession();
    },

    destroy() {
      if (this.focusPendingSessionHandler) {
        window.removeEventListener(
          "focus:terminal-session",
          this.focusPendingSessionHandler,
        );
      }
      if (this.refreshSessionsHandler) {
        window.removeEventListener(
          "update:terminal-session",
          this.refreshSessionsHandler,
        );
      }
      this.sessionTerminalManager.disposeAll();
      this.openTabs = [];
    },
  }));
});
