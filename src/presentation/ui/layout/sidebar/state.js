document.addEventListener("alpine:init", () => {
  Alpine.data("sidebar", () => ({
    isSidebarCollapsed: Alpine.$persist(true).as(
      "osDashboard.isSidebarCollapsed"
    ),
    clearSession() {
      document.cookie = `${Infinite.Envs.AccessTokenCookieKey}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
      window.location.href = "/login/";
    },
    isActivePath(path) {
      return window.location.pathname.startsWith(path);
    },
  }));
});
