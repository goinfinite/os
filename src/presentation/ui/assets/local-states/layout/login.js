document.addEventListener("alpine:init", () => {
  Alpine.data("login", () => ({
    username: "",
    password: "",
    createSessionToken() {
      const shouldDisplayToast = false;
      Infinite.JsonAjax(
        "POST",
        "/api/v1/auth/login/",
        {
          username: this.username,
          password: this.password,
        },
        shouldDisplayToast
      )
        .then((authResponse) => {
          Alpine.store("toast").displayToast("LoginSuccessful", "success");
          document.cookie = `${Infinite.Envs.AccessTokenCookieKey}=${authResponse.tokenStr}; path=/`;
          window.location.href = "/overview/";
        })
        .catch((error) =>
          Alpine.store("toast").displayToast(error.message, "danger")
        );
    },
    init() {
      document.cookie = `${Infinite.Envs.AccessTokenCookieKey}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
    },
  }));
});
