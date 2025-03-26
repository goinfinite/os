document.addEventListener("alpine:init", () => {
  Alpine.data("mappings", () => ({
    // Primary states
    virtualHost: {},
    get vhostHostnameWithTrailingSlash() {
      return this.virtualHost.hostname + "/";
    },
    mapping: {},
    resetPrimaryStates() {
      this.virtualHost = {
        hostname: "",
        type: "top-level",
        rootDirectory: "",
        parentHostname: "",
      };
      this.mapping = {
        id: 0,
        path: "",
        matchPattern: "begins-with",
        targetType: "url",
        targetValue: "",
        targetHttpResponseCode: "",
      };
    },
    init() {
      this.resetPrimaryStates();
    },

    // Auxiliary states
    isAdvancedSettingsClosed: true,
    isCreateMappingFromVhost: false,
    resetAuxiliaryStates() {
      this.isAdvancedSettingsClosed = true;
      this.isCreateMappingFromVhost = false;
    },
    get shouldDisableCreateVhostSubmitButton() {
      if (this.virtualHost.type == "alias") {
        return (
          this.virtualHost.hostname == "" ||
          this.virtualHost.parentHostname == ""
        );
      }

      return this.virtualHost.hostname == "";
    },
    get shouldDisableCreateMappingSubmitButton() {
      const isResponseCodeType = this.mapping.targetType == "response-code";
      const isTargetHttpResponseCodeRequired =
        isResponseCodeType || this.mapping.targetType == "inline-html";
      if (
        isTargetHttpResponseCodeRequired &&
        this.mapping.targetHttpResponseCode == ""
      ) {
        return true;
      }

      const isTargetValueRequired =
        !isResponseCodeType && this.mapping.targetType != "static-files";
      if (isTargetValueRequired && this.mapping.targetValue == "") {
        return true;
      }

      return this.virtualHost.hostname == "";
    },

    // Modal states
    isCreateVhostModalOpen: false,
    openCreateVhostModal() {
      this.resetPrimaryStates();

      this.isCreateVhostModalOpen = true;
    },
    closeCreateVhostModal() {
      this.isCreateVhostModalOpen = false;
    },
    isDeleteVhostModalOpen: false,
    openDeleteVhostModal(vhostHostname) {
      this.resetPrimaryStates();

      this.virtualHost.hostname = vhostHostname;
      this.isDeleteVhostModalOpen = true;
    },
    closeDeleteVhostModal() {
      this.isDeleteVhostModalOpen = false;
    },
    deleteVhostElement() {
      htmx
        .ajax("DELETE", "/api/v1/vhosts/" + this.virtualHost.hostname + "/", {
          swap: "none",
        })
        .finally(() => {
          this.closeDeleteVhostModal();
        });
    },
    isCreateMappingModalOpen: false,
    openCreateMappingModal() {
      this.resetPrimaryStates();
      this.resetAuxiliaryStates();

      this.isCreateMappingModalOpen = true;
    },
    isCreateMappingFromVhostModalOpen: false,
    openCreateMappingFromVhostModal(vhostHostname) {
      this.resetPrimaryStates();
      this.resetAuxiliaryStates();

      this.virtualHost.hostname = vhostHostname;
      this.isCreateMappingFromVhostModalOpen = true;
      this.isCreateMappingFromVhost = true;
    },
    closeCreateMappingModal() {
      this.isCreateMappingModalOpen = false;
      this.isCreateMappingFromVhostModalOpen = false;
    },
    isDeleteMappingModalOpen: false,
    openDeleteMappingModal(mappingId, mappingPath) {
      this.resetPrimaryStates();

      this.mapping.id = mappingId;
      this.mapping.path = mappingPath;
      this.isDeleteMappingModalOpen = true;
    },
    closeDeleteMappingModal() {
      this.isDeleteMappingModalOpen = false;
    },
    deleteMappingElement() {
      htmx
        .ajax("DELETE", "/api/v1/vhosts/mapping/" + this.mapping.id + "/", {
          swap: "none",
        })
        .finally(() => {
          this.closeDeleteMappingModal();
        });
    },
  }));
});
