UiToolset.RegisterAlpineState(() => {
  Alpine.data("php", () => ({
    // PrimaryState
    vhostHostname: "",
    phpConfigs: {},
    resetPrimaryStates() {
      phpConfigsElement = document.getElementById("phpConfigs");
      if (!phpConfigsElement) {
        return;
      }
      this.phpConfigs = JSON.parse(phpConfigsElement.textContent);
    },
    init() {
      this.vhostHostname = document.getElementById("vhostHostname").value;
      this.resetPrimaryStates();
    },
    updateVhostHostname(selectedHostname) {
      this.vhostHostname = selectedHostname;
    },
    updatePhpConfigs() {
      UiToolset.JsonAjax(
        "PUT",
        Infinite.OsApiBasePath +
          "/v1/runtime/php/" +
          encodeURIComponent(this.vhostHostname) +
          "/",
        {
          version: this.phpConfigs.version.value,
          modules: this.phpConfigs.modules,
          settings: this.phpConfigs.settings,
        }
      ).then(() => this.$dispatch("refresh:runtimes-page-content"));
    },

    // AuxiliaryState
    selectedPhpVerticalTab: "modules",

    // ModalState
    isUpdatePhpVersionModalOpen: false,
    openUpdatePhpVersionModal() {
      this.isUpdatePhpVersionModalOpen = true;
    },
    closeUpdatePhpVersionModal() {
      this.isUpdatePhpVersionModalOpen = false;
    },
    updatePhpVersion() {
      this.closeUpdatePhpVersionModal();
      htmx
        .ajax(
          "PUT",
          Infinite.OsApiBasePath +
            "/v1/runtime/php/" +
            encodeURIComponent(this.vhostHostname) +
            "/",
          {
            swap: "none",
            values: { version: this.phpConfigs.version.value },
          }
        )
        .then(() => this.$dispatch("refresh:runtimes-page-content"));
    },
  }));
});
