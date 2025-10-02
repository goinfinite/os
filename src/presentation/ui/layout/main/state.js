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
// "UiToolset.RegisterAlpineState" is not used here on purpose:
// 1. Registration is done at the "alpine:initializing" event instead of "init".
// 2. $store is globally accessible, so it doesn't require registration on page reload.
// 3. The mainLayout is not reloaded during page transitions.
document.addEventListener("alpine:initializing", () => {
  Alpine.store("main", {
    // RoutingState
    activeRoute: String(document.location.pathname),
    isActiveRoute(path) {
      return this.activeRoute.startsWith(path);
    },
    navigateTo(path) {
      this.activeRoute = path;

      let baseUri = document.baseURI;
      if (baseUri.endsWith("/")) {
        baseUri = baseUri.slice(0, -1);
      }
      const newPath = baseUri + path;
      htmx.ajax("GET", newPath, {
        source: "#htmx-routing-attributes-element",
        select: "#page-content",
        target: "#page-content",
        swap: "outerHTML transition:true",
      });
    },
    clearUserSession() {
      document.cookie = `${Infinite.Envs.AccessTokenCookieKey}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
      window.location.href = document.baseURI + "login/";
    },
    init() {
      window.addEventListener("popstate", () => {
        this.activeRoute = String(document.location.pathname);
      });
    },

    // FooterState
    refreshFooter() {
      htmx
        .ajax("GET", document.baseURI + "fragment/footer", {
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
