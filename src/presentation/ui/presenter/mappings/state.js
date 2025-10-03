UiToolset.RegisterAlpineState(() => {
  Alpine.data("mappings", () => ({
    // PrimaryState
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
        isWildcard: false,
        isPrimary: false,
      };
      this.mapping = {
        id: 0,
        path: "",
        matchPattern: "begins-with",
        targetType: "url",
        targetValue: "",
        targetHttpResponseCode: "",
        shouldUpgradeInsecureRequests: "false",
        mappingSecurityRuleId: "",
      };
    },
    init() {
      this.resetPrimaryStates();
    },

    // AuxiliaryState
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

    // ModalState
    isCreateVirtualHostModalOpen: false,
    openCreateVirtualHostModal() {
      this.resetPrimaryStates();

      this.isCreateVirtualHostModalOpen = true;
    },
    closeCreateVirtualHostModal() {
      this.isCreateVirtualHostModalOpen = false;
    },
    isUpdateVirtualHostModalOpen: false,
    openUpdateVirtualHostModal(vhostHostname) {
      this.resetPrimaryStates();

      const vhostEntity = JSON.parse(
        document.getElementById("vhostEntity_" + vhostHostname).textContent
      );
      this.virtualHost.hostname = vhostEntity.hostname;
      this.virtualHost.isWildcard = vhostEntity.isWildcard;

      this.isUpdateVirtualHostModalOpen = true;
    },
    closeUpdateVirtualHostModal() {
      this.isUpdateVirtualHostModalOpen = false;
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
        .ajax("DELETE", Infinite.OsApiBasePath + "/v1/vhost/" + this.virtualHost.hostname + "/", {
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
    isUpdateMappingModalOpen: false,
    openUpdateMappingModal(mappingId) {
      this.resetPrimaryStates();

      this.mapping = JSON.parse(
        document.getElementById("mappingEntity_" + mappingId).textContent
      );
      this.virtualHost.hostname = this.mapping.hostname;

      this.isUpdateMappingModalOpen = true;
    },
    closeUpdateMappingModal() {
      this.isUpdateMappingModalOpen = false;
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
        .ajax("DELETE", Infinite.OsApiBasePath + "/v1/vhost/mapping/" + this.mapping.id + "/", {
          swap: "none",
        })
        .then(() => this.$dispatch("refresh:mappings-table"));
    },
  }));
});
