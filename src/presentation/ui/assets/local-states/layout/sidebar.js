document.addEventListener("alpine:init", () => {
  Alpine.data("sidebar", () => ({
    clearSession(accessTokenCookieKey) {
      document.cookie =
        accessTokenCookieKey +
        "=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
      window.location.href = "/login/";
    },
    scrollToTop() {
      const menu = document.getElementById("sidebarMenu");
      menu.scrollTo({ top: 0, behavior: "smooth" });
    },
    toggleScrollButtonDisplay() {
      const menu = document.getElementById("sidebarMenu");
      const scrollButton = document.getElementById("scrollButton");
      if (menu.scrollTop > 0) {
        scrollButton.style.display = "flex";
      } else {
        scrollButton.style.display = "none";
      }
    },
  }));
});
