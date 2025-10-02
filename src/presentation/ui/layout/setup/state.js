UiToolset.RegisterAlpineState(() => {
  Alpine.data("setup", () => ({
    username: "",
    password: "",
    setupInfiniteOsAndLogin() {
      const shouldDisplayToast = false;
      UiToolset.JsonAjax(
        "POST",
        Infinite.OsApiBasePath + "/v1/setup/",
        {
          username: this.username,
          password: this.password,
        },
        shouldDisplayToast
      )
        .then(() => {
          UiToolset.JsonAjax(
            "POST",
            Infinite.OsApiBasePath + "/v1/auth/login/",
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
            document.cookie = `${Infinite.Envs.AccessTokenCookieKey}=${authResponse.tokenStr}; path=/; Secure; SameSite=Lax;`;
            window.location.href = document.baseURI + "overview/";
          });
        })
        .catch((error) =>
          Alpine.store("toast").displayToast(error.message, "danger")
        );
    },
  }));
});
