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
    refreshFooter() {
      htmx.ajax("GET", "/fragment/footer", {
        select: "#footer",
        target: "#footer",
        swap: "outerHTML transition:true",
      });
    },
    refreshScheduledTasksPopover() {
      this.refreshFooter();
      setTimeout(() => {
        this.displayScheduledTasksPopover = true;
      }, 1000);
    },
    navigateTo(path) {
      htmx.ajax("GET", path, {
        select: "#page-content",
        target: "#page-content",
        swap: "outerHTML transition:true",
      });
    },
  });
});
