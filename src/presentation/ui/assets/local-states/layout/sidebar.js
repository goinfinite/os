document.addEventListener("alpine:init", () => {
  Alpine.data("sidebar", () => ({
    isSidebarCollapsed: Alpine.$persist(true).as(
      "osDashboard.isSidebarCollapsed"
    ),
    currentPath: window.location.pathname,
    sidebarItems: [
      { Label: "Overview", Icon: "ph-speedometer", Link: "/overview/" },
      { Label: "Accounts", Icon: "ph-users-three", Link: "/accounts/" },
      { Label: "Databases", Icon: "ph-database", Link: "/databases/" },
      { Label: "Runtime", Icon: "ph-code", Link: "/runtimes/" },
      { Label: "Cron Jobs", Icon: "ph-clock", Link: "/crons/" },
      { Label: "File Manager", Icon: "ph-files", Link: "/file-manager/" },
      { Label: "Mappings", Icon: "ph-graph", Link: "/mappings/" },
      { Label: "SSL Certificates", Icon: "ph-lock", Link: "/ssls/" },
      { Label: "Marketplace", Icon: "ph-storefront", Link: "/marketplace/" },
    ],
    clearSession() {
      document.cookie = `${Infinite.Envs.AccessTokenCookieKey}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
      window.location.href = "/login/";
    },
    scrollToTop() {
      const menu = document.getElementById("sidebarMenu");
      menu.scrollTo({ top: 0, behavior: "smooth" });
    },
    toggleScrollButtonDisplay() {
      const menu = document.getElementById("sidebarMenu");
      const scrollButton = document.getElementById("scrollButton");

      let scrollButtonDisplay = "none";
      if (menu.scrollTop > 0) {
        scrollButtonDisplay = "flex";
      }
      scrollButton.style.display = scrollButtonDisplay;
    },
    isActivePath(path) {
      return this.currentPath.startsWith(path);
    },
    updateCurrentPath(path) {
      this.currentPath = path;
    },
  }));
});

document.addEventListener("htmx:afterSettle", function (event) {
  // Update Alpine.js state when HTMX completes a request
  const sidebar = document.getElementById("sidebar")?.__x;
  if (sidebar) {
    sidebar.$data.currentPath = window.location.pathname;
  }

  // Update document title based on the active menu item
  const currentPath = window.location.pathname;
  const menuItems = sidebar?.$data.sidebarItems || [
    { Label: "Overview", Link: "/overview/" },
    { Label: "Accounts", Link: "/accounts/" },
    { Label: "Databases", Link: "/databases/" },
    { Label: "Runtime", Link: "/runtimes/" },
    { Label: "Cron Jobs", Link: "/crons/" },
    { Label: "File Manager", Link: "/file-manager/" },
    { Label: "Mappings", Link: "/mappings/" },
    { Label: "SSL Certificates", Link: "/ssls/" },
    { Label: "Marketplace", Link: "/marketplace/" },
  ];

  for (const item of menuItems) {
    if (currentPath.startsWith(item.Link)) {
      document.title = `${item.Label} - Infinite OS`;
      break;
    }
  }
});
