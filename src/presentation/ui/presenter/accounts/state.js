UiToolset.RegisterAlpineState(() => {
  Alpine.data("accounts", () => ({
    // PrimaryState
    account: {},
    secureAccessPublicKey: {},
    resetPrimaryStates() {
      this.account = {
        id: "",
        groupId: "",
        username: "",
        password: "",
        isSuperAdmin: false,
        apiKey: "",
        secureAccessPublicKeys: [],
      };
      this.secureAccessPublicKey = {
        id: "",
        name: "",
        content: "",
      };
    },
    init() {
      this.resetPrimaryStates();
    },

    // ModalState
    isCreateAccountModalOpen: false,
    openCreateAccountModal() {
      this.resetPrimaryStates();

      this.isCreateAccountModalOpen = true;
    },
    closeCreateAccountModal() {
      this.isCreateAccountModalOpen = false;
    },
    isUpdateAccountModalOpen: false,
    openUpdateAccountModal(id, isSuperAdmin) {
      this.resetPrimaryStates();

      this.account.id = id;
      this.account.isSuperAdmin = isSuperAdmin;
      this.isUpdateAccountModalOpen = true;
    },
    closeUpdateAccountModal() {
      this.isUpdateAccountModalOpen = false;
    },
    isUpdateApiKeyModalOpen: false,
    openUpdateApiKeyModal(id, username) {
      this.resetPrimaryStates();

      this.account.id = id;
      this.account.username = username;
      this.isUpdateApiKeyModalOpen = true;
    },
    closeUpdateApiKeyModal() {
      this.isUpdateApiKeyModalOpen = false;
      this.account.apiKey = "";
    },
    updateApiKey() {
      const shouldDisplayToast = false;
      UiToolset.JsonAjax(
        "PUT",
        Infinite.OsApiBasePath + "/v1/account/",
        { id: this.account.id, shouldUpdateApiKey: true },
        shouldDisplayToast
      ).then((apiKey) => (this.account.apiKey = apiKey));
    },
    isSecureAccessPublicKeysModalOpen: false,
    openSecureAccessPublicKeysModal(id, username) {
      this.resetPrimaryStates();

      this.account.id = id;
      this.account.username = username;
      this.account.secureAccessPublicKeys = JSON.parse(
        document.getElementById("secureAccessPublicKeys_" + id).textContent
      );

      this.isSecureAccessPublicKeysModalOpen = true;
    },
    closeSecureAccessPublicKeysModal() {
      this.isSecureAccessPublicKeysModalOpen = false;
    },
    isCreateSecureAccessPublicKeyModalOpen: false,
    openCreateSecureAccessPublicKeyModal() {
      this.isCreateSecureAccessPublicKeyModalOpen = true;
    },
    closeCreateSecureAccessPublicKeyModal() {
      this.isCreateSecureAccessPublicKeyModalOpen = false;
    },
    isDeleteSecureAccessPublicKeyModalOpen: false,
    openDeleteSecureAccessPublicKeyModal(id, name) {
      this.secureAccessPublicKey.id = id;
      this.secureAccessPublicKey.name = name;
      this.isDeleteSecureAccessPublicKeyModalOpen = true;
    },
    closeDeleteSecureAccessPublicKeyModal() {
      this.isDeleteSecureAccessPublicKeyModalOpen = false;
    },
    deleteSecureAccessPublicKey() {
      htmx
        .ajax(
          "DELETE",
          Infinite.OsApiBasePath + `/v1/account/secure-access-public-key/${this.secureAccessPublicKey.id}/`,
          { swap: "none" }
        )
        .then(() => this.$dispatch("delete:secure-access-public-key"))
        .finally(() => this.closeDeleteSecureAccessPublicKeyModal());
    },
    isDeleteAccountModalOpen: false,
    openDeleteAccountModal(id, username) {
      this.resetPrimaryStates();

      this.account.id = id;
      this.account.username = username;
      this.isDeleteAccountModalOpen = true;
    },
    closeDeleteAccountModal() {
      this.resetPrimaryStates();

      this.isDeleteAccountModalOpen = false;
    },
    deleteAccount() {
      htmx
        .ajax("DELETE", Infinite.OsApiBasePath + `/v1/account/` + this.account.id, { swap: "none" })
        .then(() => this.$dispatch("delete:account"))
        .finally(() => this.closeDeleteAccountModal());
    },
  }));
});
