const terminalRendererRetryLimit = 25;
const terminalRendererRetryDelayMs = 200;
const terminalReconnectDelayMs = 2000;

UiToolset.RegisterAlpineState(() => {
  Alpine.data("terminal", () => ({
    // PrimaryState
    sessions: [],
    workingDir: "/app",
    command: "",
    openTabs: [],
    activeTabId: "",
    isLoading: false,

    // DerivedState
    isSessionAttached(sessionId) {
      return this.openTabs.some((tab) => tab.sessionId === sessionId);
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
      } catch (error) {
        console.error(`ReadTerminalSessionsError: ${error}`);
        Alpine.store("toast").displayToast(
          "ReadTerminalSessionsError",
          "danger",
        );
      }
    },

    async createSession() {
      this.isLoading = true;
      try {
        const requestBody = { workingDir: this.workingDir };
        if (this.command.length > 0) {
          requestBody.command = this.command;
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

        this.command = "";
        await this.loadSessions();
        this.openTab(jsonResponse.body);
      } catch (error) {
        Alpine.store("toast").displayToast(error.message, "danger");
      } finally {
        this.isLoading = false;
      }
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
      } catch (error) {
        Alpine.store("toast").displayToast(error.message, "danger");
      }
    },

    openTab(session) {
      if (this.isSessionAttached(session.id)) {
        this.activeTabId = session.id;
        return;
      }

      const tab = {
        sessionId: session.id,
        title: `${session.accountUsername}@${session.id.slice(0, 4)}`,
        term: null,
        fitAddon: null,
        socket: null,
        isConnected: false,
        reconnectTimer: null,
        mountAttempts: 0,
      };
      this.openTabs.push(tab);
      this.activeTabId = session.id;

      this.$nextTick(() => this.mountTerminal(tab));
    },

    mountTerminal(tab) {
      const terminalContainer = document.getElementById(
        `terminal-${tab.sessionId}`,
      );
      if (!terminalContainer) {
        return;
      }

      if (typeof Terminal === "undefined" || typeof FitAddon === "undefined") {
        tab.mountAttempts += 1;
        if (tab.mountAttempts > terminalRendererRetryLimit) {
          Alpine.store("toast").displayToast(
            "TerminalRendererUnavailable",
            "danger",
          );
          return;
        }
        setTimeout(() => this.mountTerminal(tab), terminalRendererRetryDelayMs);
        return;
      }

      const term = new Terminal({
        cursorBlink: true,
        fontFamily: "monospace",
        theme: { background: "#041118" },
      });
      const fitAddon = new FitAddon.FitAddon();
      term.loadAddon(fitAddon);
      term.open(terminalContainer);
      fitAddon.fit();

      term.onData((terminalData) => {
        if (!tab.socket || tab.socket.readyState !== WebSocket.OPEN) {
          return;
        }
        tab.socket.send(new TextEncoder().encode(terminalData));
      });

      tab.term = term;
      tab.fitAddon = fitAddon;
      this.connectTab(tab);
    },

    connectTab(tab) {
      const attachUrl = new URL(
        `api/v1/terminal-sessions/${tab.sessionId}/attach/`,
        document.baseURI,
      );
      attachUrl.protocol = attachUrl.protocol === "https:" ? "wss:" : "ws:";

      const socket = new WebSocket(attachUrl.toString());
      socket.binaryType = "arraybuffer";
      tab.socket = socket;

      socket.onopen = () => {
        tab.isConnected = true;
        this.sendResize(tab);
      };
      socket.onmessage = (event) => {
        if (typeof event.data === "string") {
          return;
        }
        tab.term.write(new Uint8Array(event.data));
      };
      socket.onclose = () => {
        tab.isConnected = false;
        tab.reconnectTimer = setTimeout(
          () => this.connectTab(tab),
          terminalReconnectDelayMs,
        );
      };
      socket.onerror = () => socket.close();
    },

    sendResize(tab) {
      if (!tab.socket || tab.socket.readyState !== WebSocket.OPEN) {
        return;
      }

      tab.fitAddon.fit();
      tab.socket.send(
        JSON.stringify({
          type: "resize",
          cols: tab.term.cols,
          rows: tab.term.rows,
        }),
      );
    },

    selectTab(sessionId) {
      this.activeTabId = sessionId;
      this.$nextTick(() => {
        const tab = this.openTabs.find((item) => item.sessionId === sessionId);
        if (!tab?.fitAddon) {
          return;
        }
        this.sendResize(tab);
      });
    },

    disposeTab(tab) {
      if (tab.reconnectTimer) {
        clearTimeout(tab.reconnectTimer);
      }
      if (tab.socket) {
        tab.socket.onclose = null;
        tab.socket.close();
      }
      if (tab.term) {
        tab.term.dispose();
      }
    },

    closeTab(sessionId) {
      const tabIndex = this.openTabs.findIndex(
        (tab) => tab.sessionId === sessionId,
      );
      if (tabIndex === -1) {
        return;
      }

      this.disposeTab(this.openTabs[tabIndex]);
      this.openTabs.splice(tabIndex, 1);
      if (this.activeTabId === sessionId) {
        this.activeTabId =
          this.openTabs.length > 0 ? this.openTabs[0].sessionId : "";
      }
    },

    init() {
      this.loadSessions();
    },

    destroy() {
      for (const tab of this.openTabs) {
        this.disposeTab(tab);
      }
      this.openTabs = [];
    },
  }));
});
