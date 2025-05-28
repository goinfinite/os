Infinite.RegisterAlpineState(sidebarAlpineState);

function sidebarAlpineState() {
  Alpine.data("sidebar", () => ({
    isSidebarCollapsed: Alpine.$persist(false).as(
      "osDashboard.isSidebarCollapsed"
    ),
    clearSession() {
      document.cookie = `${Infinite.Envs.AccessTokenCookieKey}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
      window.location.href = "/login/";
    },
  }));
}
