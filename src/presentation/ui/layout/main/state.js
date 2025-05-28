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
    // RoutingState
    activeRoute: String(document.location.pathname),
    isActiveRoute(path) {
      return this.activeRoute.startsWith(path);
    },
    navigateTo(path) {
      this.activeRoute = path;
      htmx.ajax("GET", path, {
        source: "#htmx-routing-attributes-element",
        select: "#page-content",
        target: "#page-content",
        swap: "outerHTML transition:true",
      });
    },
    clearUserSession() {
      document.cookie = `${Infinite.Envs.AccessTokenCookieKey}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
      window.location.href = "/login/";
    },
    init() {
      window.addEventListener("popstate", () => {
        this.activeRoute = String(document.location.pathname);
      });
    },

    // FooterState
    refreshFooter() {
      htmx
        .ajax("GET", "/fragment/footer", {
          select: "#footer",
          target: "#footer",
          swap: "outerHTML transition:true",
        })
        .catch(() => {
          console.error("FooterRefreshFailed");
          this.clearUserSession();
        });
    },
    // - ScheduledTasksState
    displayScheduledTasksPopover: Alpine.$persist(false).as(
      "osDash.displayScheduledTasksPopover"
    ),
    toggleScheduledTasksPopover() {
      this.displayScheduledTasksPopover = !this.displayScheduledTasksPopover;
    },
    refreshScheduledTasksPopover() {
      this.refreshFooter();
      setTimeout(() => {
        this.displayScheduledTasksPopover = true;
      }, 1000);
    },
  });
});
