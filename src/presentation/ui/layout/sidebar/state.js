document.addEventListener("alpine:init", () => {
  Alpine.data("sidebar", () => ({
    isSidebarCollapsed: Alpine.$persist(true).as(
      "osDashboard.isSidebarCollapsed"
    ),
    activeItemPath: "/overview/",
    init() {
      this.activeItemPath = window.location.pathname;
    },
    isActivePath(path) {
      return String(this.activeItemPath).startsWith(path);
    },
    clearSession() {
      document.cookie = `${Infinite.Envs.AccessTokenCookieKey}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
      window.location.href = "/login/";
    },
  }));
});
