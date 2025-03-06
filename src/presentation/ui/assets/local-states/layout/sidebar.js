document.addEventListener("alpine:init", () => {
  Alpine.data("sidebar", () => ({
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
  }));
});
