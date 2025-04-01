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
    isCreateMappingFromVirtualHost: false,
    resetAuxiliaryStates() {
      this.isAdvancedSettingsClosed = true;
      this.isCreateMappingFromVirtualHost = false;
    },
    get shouldDisableCreateVirtualHostSubmitButton() {
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
    isCreateVirtualHostModalOpen: false,
    openCreateVirtualHostModal() {
      this.resetPrimaryStates();

      this.isCreateVirtualHostModalOpen = true;
    },
    closeCreateVirtualHostModal() {
      this.isCreateVirtualHostModalOpen = false;
    },
    isDeleteVirtualHostModalOpen: false,
    openDeleteVirtualHostModal(vhostHostname) {
      this.resetPrimaryStates();

      this.virtualHost.hostname = vhostHostname;
      this.isDeleteVirtualHostModalOpen = true;
    },
    closeDeleteVirtualHostModal() {
      this.isDeleteVirtualHostModalOpen = false;
    },
    deleteVirtualHostElement() {
      this.closeDeleteVirtualHostModal();
      htmx
        .ajax("DELETE", "/api/v1/vhosts/" + this.virtualHost.hostname + "/", {
          swap: "none",
        })
        .then(() => this.$dispatch("refresh:mappings-table"));
    },
    isCreateMappingModalOpen: false,
    openCreateMappingModal() {
      this.resetPrimaryStates();
      this.resetAuxiliaryStates();

      this.isCreateMappingModalOpen = true;
    },
    isCreateMappingFromVirtualHostModalOpen: false,
    openCreateMappingFromVirtualHostModal(vhostHostname) {
      this.resetPrimaryStates();
      this.resetAuxiliaryStates();

      this.virtualHost.hostname = vhostHostname;
      this.isCreateMappingFromVirtualHostModalOpen = true;
      this.isCreateMappingFromVirtualHost = true;
    },
    closeCreateMappingModal() {
      this.isCreateMappingModalOpen = false;
      this.isCreateMappingFromVirtualHostModalOpen = false;
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
      this.closeDeleteMappingModal();
      htmx
        .ajax("DELETE", "/api/v1/vhosts/mapping/" + this.mapping.id + "/", {
          swap: "none",
        })
        .then(() => this.$dispatch("refresh:mappings-table"));
    },
  }));
});
