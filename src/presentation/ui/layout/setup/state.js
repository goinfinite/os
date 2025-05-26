document.addEventListener("alpine:init", () => {
  Alpine.data("setup", () => ({
    username: "",
    password: "",
    setupInfiniteOsAndLogin() {
      const shouldDisplayToast = false;
      Infinite.JsonAjax(
        "POST",
        "/api/v1/setup/",
        {
          username: this.username,
          password: this.password,
        },
        shouldDisplayToast
      )
        .then(() => {
          Infinite.JsonAjax(
            "POST",
            "/api/v1/auth/login/",
            {
              username: this.username,
              password: this.password,
            },
            shouldDisplayToast
          ).then((authResponse) => {
            if (!authResponse.tokenStr) {
              return Alpine.store("toast").displayToast(
                error.message,
                "danger"
              );
            }

            Alpine.store("toast").displayToast("LoginSuccessful", "success");
            document.cookie = `${Infinite.Envs.AccessTokenCookieKey}=${authResponse.tokenStr}; path=/`;
            window.location.href = "/overview/";
          });
        })
        .catch((error) =>
          Alpine.store("toast").displayToast(error.message, "danger")
        );
    },
  }));
});
