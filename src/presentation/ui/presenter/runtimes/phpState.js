UiToolset.RegisterAlpineState(() => {
  Alpine.data("php", () => ({
    // PrimaryState
    vhostHostname: "",
    phpConfigs: {},
    initialPhpModuleStatuses: {},
    resetPrimaryStates() {
      const phpConfigsElement = document.getElementById("phpConfigs");
      if (!phpConfigsElement) {
        return;
      }
      this.phpConfigs = JSON.parse(phpConfigsElement.textContent);
      this.initialPhpModuleStatuses = {};
      for (const module of this.phpConfigs.modules ?? []) {
        this.initialPhpModuleStatuses[module.name] = module.status;
      }
    },
    init() {
      this.vhostHostname = document.getElementById("vhostHostname").value;
      this.resetPrimaryStates();
    },
    updateVhostHostname(selectedHostname) {
      this.vhostHostname = selectedHostname;
    },
    displayModuleParsingFailures(parsingFailures) {
      const failureMessages = parsingFailures.map((failure) => {
        const moduleName =
          failure.name || "module #" + (Number(failure.index) + 1);
        return moduleName + ": " + failure.reason;
      });
      Alpine.store("toast").displayToast(
        "PhpModulesParsingFailed: " + failureMessages.join(", "),
        "danger"
      );
    },
    refreshRuntimesPage() {
      this.$dispatch("refresh:runtimes-page-content");
    },
    updatePhpConfigs() {
      const changedModules = (this.phpConfigs.modules ?? []).filter(
        (module) =>
          this.initialPhpModuleStatuses[module.name] !== module.status
      );
      const requestBody = {
        version: this.phpConfigs.version.value,
        settings: this.phpConfigs.settings,
      };
      if (changedModules.length > 0) {
        requestBody.modules = changedModules;
      }

      UiToolset.JsonAjax(
        "PUT",
        Infinite.OsApiBasePath +
          "/v1/runtime/php/" +
          encodeURIComponent(this.vhostHostname) +
          "/",
        requestBody
      ).then((responseBody) => {
        const parsingFailures =
          responseBody?.failedModulesWithParsingErrors ?? [];
        if (parsingFailures.length > 0) {
          this.displayModuleParsingFailures(parsingFailures);
        }

        if (responseBody?.taskId === undefined) {
          this.refreshRuntimesPage();
          return;
        }

        Alpine.store("toast").displayToast("PhpUpdateScheduled", "success");
        this.$store.main.refreshScheduledTasksPopover();
        this.refreshRuntimesPage();
      });
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
        .then(() => this.refreshRuntimesPage());
    },
  }));
});
