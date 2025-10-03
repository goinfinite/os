UiToolset.RegisterAlpineState(() => {
  Alpine.data("login", () => ({
    username: "",
    password: "",
    createSessionToken() {
      const shouldDisplayToast = false;
      UiToolset.JsonAjax(
        "POST",
        Infinite.OsApiBasePath + "/v1/auth/login/",
        {
          username: this.username,
          password: this.password,
        },
        shouldDisplayToast
      )
        .then((authResponse) => {
          Alpine.store("toast").displayToast("LoginSuccessful", "success");

          UiToolset.ToggleLoadingOverlay(true);
          document.cookie = `${Infinite.Envs.AccessTokenCookieKey}=${authResponse.tokenStr}; path=/; Secure; SameSite=Lax;`;
          window.location.href = document.baseURI + "overview/";
        })
        .catch((error) => {
          UiToolset.ToggleLoadingOverlay(false);
          Alpine.store("toast").displayToast(error.message, "danger");
        });
    },
    init() {
      document.cookie = `${Infinite.Envs.AccessTokenCookieKey}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
      const prefilledUsername = document.getElementById("prefilledUsername")?.value;
      const usernameBasicRegex = /^[\w\-]{2,64}$/;
      if (usernameBasicRegex.test(prefilledUsername)) {
        this.username = prefilledUsername;
      }
      
      const prefilledPassword = document.getElementById("prefilledPassword")?.value;
      const passwordBasicRegex = /^.{4,128}$/;
      if (passwordBasicRegex.test(prefilledPassword)) {
        this.password = prefilledPassword;
      }
    },
  }));
});
