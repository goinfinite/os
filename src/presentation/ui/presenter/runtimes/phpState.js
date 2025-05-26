document.addEventListener("alpine:init", () => {
  Alpine.data("php", () => ({
    // Primary States
    phpConfigs: {},
    resetPrimaryStates() {
      phpConfigsElement = document.getElementById("phpConfigs");
      if (!phpConfigsElement) {
        return;
      }
      this.phpConfigs = JSON.parse(phpConfigsElement.textContent);
    },
    init() {
      this.resetPrimaryStates();
    },
    updatePhpConfigs() {
      Infinite.JsonAjax(
        "PUT",
        "/api/v1/runtime/php/" + this.vhostHostname + "/",
        {
          version: this.phpConfigs.version.value,
          modules: this.phpConfigs.modules,
          settings: this.phpConfigs.settings,
        }
      ).then(() => this.$dispatch("refresh:runtimes-page-content"));
    },

    // Auxiliary States
    selectedPhpVerticalTab: "modules",

    // Modal States
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
        .ajax("PUT", "/api/v1/runtime/php/" + this.vhostHostname + "/", {
          swap: "none",
          values: { version: this.phpConfigs.version.value },
        })
        .then(() => this.$dispatch("refresh:runtimes-page-content"));
    },
  }));
});
