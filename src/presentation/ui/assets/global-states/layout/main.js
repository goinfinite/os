function devWsHotReload() {
  hotReloadWs = new WebSocket(
    "wss://" + document.location.host + "/dev/hot-reload"
  );
  hotReloadWs.onclose = () => {
    setTimeout(() => {
      location.reload();
    }, 2000);
  };
}

document.addEventListener("alpine:initializing", () => {
  Alpine.store("main", {
    displayScheduledTasksPopover: Alpine.$persist(false).as(
      "dash.displayScheduledTasksPopover"
    ),
    toggleScheduledTasksPopover() {
      this.displayScheduledTasksPopover = !this.displayScheduledTasksPopover;
    },
    refreshScheduledTasksPopover() {
      window.dispatchEvent(new CustomEvent("refresh:footer"));
      setTimeout(() => {
        this.displayScheduledTasksPopover = true;
      }, 1000);
    },
  });
});
